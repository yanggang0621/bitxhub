package mempool

import (
	"math"
	"sync"
	"time"

	"github.com/google/btree"
	"github.com/meshplus/bitxhub-model/pb"
	raftproto "github.com/meshplus/bitxhub/pkg/order/etcdraft/proto"
	"github.com/sirupsen/logrus"
)

type mempoolImpl struct {
	localID     uint64
	batchSize   uint64
	txSliceSize uint64
	batchSeqNo  uint64 // track the sequence number of block
	poolSize    uint64
	logger      logrus.FieldLogger
	txStore     *transactionStore // store all transactions info
}

func newMempoolImpl(config *Config) *mempoolImpl {
	mpi := &mempoolImpl{
		localID:     config.ID,
		batchSeqNo:  config.ChainHeight,
		logger:      config.Logger,
		txSliceSize: config.TxSliceSize,
	}
	mpi.txStore = newTransactionStore()
	if config.BatchSize == 0 {
		mpi.batchSize = DefaultBatchSize
	} else {
		mpi.batchSize = config.BatchSize
	}
	if config.PoolSize == 0 {
		mpi.poolSize = DefaultPoolSize
	} else {
		mpi.poolSize = config.PoolSize
	}
	if config.TxSliceSize == 0 {
		mpi.txSliceSize = DefaultTxSetSize
	} else {
		mpi.txSliceSize = config.TxSliceSize
	}
	mpi.logger.Infof("MemPool batch size = %d", mpi.batchSize)
	mpi.logger.Infof("MemPool tx slice size = %d", mpi.batchSize)
	mpi.logger.Infof("MemPool batch seqNo = %d", mpi.batchSeqNo)
	mpi.logger.Infof("MemPool pool size = %d", mpi.poolSize)
	return mpi
}

func (mpi *mempoolImpl) ProcessTransactions(txs []*pb.Transaction, isLeader, isLocal bool) *raftproto.RequestBatch {
	validTxs := make(map[string][]*pb.Transaction)
	for _, tx := range txs {
		// check the sequence number of tx
		txAccount := tx.Account()
		currentSeqNo := mpi.txStore.nonceCache.getPendingNonce(txAccount)
		if tx.Nonce < currentSeqNo {
			mpi.logger.Warningf("Account %s, current sequence number is %d, required %d", txAccount, tx.Nonce, currentSeqNo+1)
			continue
		}
		// check the existence of hash of this tx
		txHash := tx.TransactionHash.String()
		if txPointer := mpi.txStore.txHashMap[txHash]; txPointer != nil {
			mpi.logger.Warningf("Tx [account: %s, nonce: %d, hash: %s] already received", txAccount, tx.Nonce, txHash)
			continue
		}
		_, ok := validTxs[txAccount]
		if !ok {
			validTxs[txAccount] = make([]*pb.Transaction, 0)
		}
		validTxs[txAccount] = append(validTxs[txAccount], tx)
	}

	// Process all the new transaction and merge any errors into the original slice
	dirtyAccounts := mpi.txStore.insertTxs(validTxs, isLocal)

	// send tx to mempool store
	mpi.processDirtyAccount(dirtyAccounts)

	// generator batch by block size
	if isLeader && mpi.txStore.priorityNonBatchSize >= mpi.batchSize {
		batch, err := mpi.generateBlock()
		if err != nil {
			mpi.logger.Errorf("Generator batch failed")
			return nil
		}
		return batch
	}
	return nil
}

func (mpi *mempoolImpl) processDirtyAccount(dirtyAccounts map[string]bool) {
	for account := range dirtyAccounts {
		if list, ok := mpi.txStore.allTxs[account]; ok {
			// search for related sequential txs in allTxs
			// and add these txs into priorityIndex and parkingLotIndex
			pendingNonce := mpi.txStore.nonceCache.getPendingNonce(account)
			readyTxs, nonReadyTxs, nextDemandNonce := list.filterReady(pendingNonce)
			mpi.txStore.nonceCache.setPendingNonce(account, nextDemandNonce)

			// insert ready txs into priorityIndex and set ttlIndex for these txs.
			for _, tx := range readyTxs {
				if !mpi.txStore.priorityIndex.data.Has(makeTimeoutKey(account, tx)) {
					mpi.txStore.priorityIndex.insertByTimeoutKey(account, tx)
					if list.items[tx.Nonce].local {
						// no need to rebroadcast tx from other nodes to reduce network overhead
						mpi.txStore.ttlIndex.insertByTtlKey(account, tx.Nonce, tx.Timestamp)
					}
				}
			}
			mpi.txStore.updateEarliestTimestamp()
			mpi.txStore.priorityNonBatchSize = mpi.txStore.priorityNonBatchSize + uint64(len(readyTxs))

			// inset non-ready txs into parkingLotIndex.
			for _, tx := range nonReadyTxs {
				mpi.txStore.parkingLotIndex.insertByOrderedQueueKey(account, tx)
			}
		}
	}
}

// getBlock fetches next block of transactions for consensus,
// batchedTx are all txs sent to consensus but were not committed yet, mempool should filter out such txs.
func (mpi *mempoolImpl) generateBlock() (*raftproto.RequestBatch, error) {
	// tx which has lower timestamp will be observed first in priority index iterator.
	// and if first seen tx's nonce isn't the required nonce for the account,
	// it will be stored in skip DS first.
	mpi.logger.Debugf("Length of non-batched transactions: %d", mpi.txStore.priorityNonBatchSize)
	var batchSize uint64
	if poolLen := mpi.txStore.priorityNonBatchSize; poolLen > mpi.batchSize {
		batchSize = mpi.batchSize
	} else {
		batchSize = mpi.txStore.priorityNonBatchSize
	}

	skippedTxs := make(map[orderedIndexKey]bool)
	result := make([]orderedIndexKey, 0, mpi.batchSize)
	mpi.txStore.priorityIndex.data.Ascend(func(a btree.Item) bool {
		tx := a.(*orderedTimeoutKey)
		// if tx has existed in bathedTxs, ignore this tx
		if _, ok := mpi.txStore.batchedTxs[orderedIndexKey{tx.account, tx.nonce}]; ok {
			return true
		}
		txSeq := tx.nonce
		commitNonce := mpi.txStore.nonceCache.getCommitNonce(tx.account)
		var seenPrevious bool
		if txSeq >= 1 {
			_, seenPrevious = mpi.txStore.batchedTxs[orderedIndexKey{account: tx.account, nonce: txSeq - 1}]
		}
		if txSeq == 3 {
			mpi.logger.Infof("seenPrevious %v and commitNonce is %d", seenPrevious, commitNonce)
			mpi.logger.Infof("batched txs is %v", mpi.txStore.batchedTxs)
		}
		// include transaction if it's "next" for given account or
		// we've already sent its ancestor to Consensus
		if seenPrevious || (txSeq == commitNonce) {
			ptr := orderedIndexKey{account: tx.account, nonce: txSeq}
			mpi.txStore.batchedTxs[ptr] = true
			result = append(result, ptr)
			if uint64(len(result)) == batchSize {
				return false
			}

			// check if we can now include some txs that were skipped before for given account
			skippedTxn := orderedIndexKey{account: tx.account, nonce: txSeq + 1}
			for {
				if _, ok := skippedTxs[skippedTxn]; !ok {
					break
				}
				mpi.txStore.batchedTxs[skippedTxn] = true
				result = append(result, skippedTxn)
				if uint64(len(result)) == batchSize {
					return false
				}
				skippedTxn.nonce++
			}
		} else {
			skippedTxs[orderedIndexKey{tx.account, txSeq}] = true
		}
		return true
	})

	if len(result) == 0 && mpi.txStore.priorityNonBatchSize > 0 {
		mpi.logger.Error("===== NOTE!!! Leader generate a batch with 0 txs")
		mpi.txStore.priorityNonBatchSize = 0
		return nil, nil
	}

	// convert transaction pointers to real values
	txList := make([]*pb.Transaction, len(result))
	for i, v := range result {
		rawTransaction := mpi.txStore.getTxByOrderKey(v.account, v.nonce)
		txList[i] = rawTransaction
	}
	mpi.batchSeqNo++
	batchSeqNo := mpi.batchSeqNo
	batch := &raftproto.RequestBatch{
		TxList: txList,
		Height: batchSeqNo,
	}
	if mpi.txStore.priorityNonBatchSize >= uint64(len(txList)) {
		mpi.txStore.priorityNonBatchSize = mpi.txStore.priorityNonBatchSize - uint64(len(txList))
	}
	mpi.logger.Debugf("Leader generate a batch with %d txs, which height is %d, and now there are %d pending txs.", len(txList), batchSeqNo, mpi.txStore.priorityNonBatchSize)
	return batch, nil
}

// processCommitTransactions removes the transactions in ready.
func (mpi *mempoolImpl) processCommitTransactions(state *ChainState) {
	dirtyAccounts := make(map[string]bool)
	// update current cached commit nonce for account
	updateAccounts := make(map[string]uint64)
	// update current cached commit nonce for account
	for _, txHash := range state.TxHashList {
		strHash := txHash.String()
		txPointer := mpi.txStore.txHashMap[strHash]
		txPointer, ok := mpi.txStore.txHashMap[strHash]
		if !ok {
			mpi.logger.Warningf("Remove transaction %s failed, Can't find it from txHashMap", strHash)
			continue
		}
		preCommitNonce := mpi.txStore.nonceCache.getCommitNonce(txPointer.account)
		newCommitNonce := txPointer.nonce + 1
		if preCommitNonce < newCommitNonce {
			mpi.txStore.nonceCache.setCommitNonce(txPointer.account, newCommitNonce)
			// Note!!! updating pendingNonce to commitNonce for the restart node
			pendingNonce := mpi.txStore.nonceCache.getPendingNonce(txPointer.account)
			if pendingNonce < newCommitNonce {
				updateAccounts[txPointer.account] = newCommitNonce
				mpi.txStore.nonceCache.setPendingNonce(txPointer.account, newCommitNonce)
			}
		}
		delete(mpi.txStore.txHashMap, strHash)
		delete(mpi.txStore.batchedTxs, *txPointer)
		dirtyAccounts[txPointer.account] = true
	}
	// clean related txs info in cache
	for account := range dirtyAccounts {
		commitNonce := mpi.txStore.nonceCache.getCommitNonce(account)
		if list, ok := mpi.txStore.allTxs[account]; ok {
			// remove all previous seq number txs for this account.
			removedTxs := list.forward(commitNonce)
			// remove index smaller than commitNonce delete index.
			var wg sync.WaitGroup
			wg.Add(4)
			go func(ready map[string][]*pb.Transaction) {
				defer wg.Done()
				list.index.removeBySortedNonceKey(removedTxs)
			}(removedTxs)
			go func(ready map[string][]*pb.Transaction) {
				defer wg.Done()
				mpi.txStore.priorityIndex.removeByTimeoutKey(removedTxs)
			}(removedTxs)
			go func(ready map[string][]*pb.Transaction) {
				defer wg.Done()
				mpi.txStore.ttlIndex.removeByTtlKey(removedTxs)
				mpi.txStore.updateEarliestTimestamp()
			}(removedTxs)
			go func(ready map[string][]*pb.Transaction) {
				defer wg.Done()
				mpi.txStore.parkingLotIndex.removeByOrderedQueueKey(removedTxs)
			}(removedTxs)
			wg.Wait()
		}
	}
	readyNum := uint64(mpi.txStore.priorityIndex.size())
	// set priorityNonBatchSize to min(nonBatchedTxs, readyNum),
	if mpi.txStore.priorityNonBatchSize > readyNum {
		mpi.txStore.priorityNonBatchSize = readyNum
	}
	for account, pendingNonce := range updateAccounts {
		mpi.logger.Debugf("Account %s update its pendingNonce to %d by commitNonce", account, pendingNonce)
	}
	mpi.logger.Debugf("Replica %d removes batches in mempool, and now there are %d non-batched txs,"+
		"priority len: %d, parkingLot len: %d, batchedTx len: %d, txHashMap len: %d", mpi.localID, mpi.txStore.priorityNonBatchSize,
		mpi.txStore.priorityIndex.size(), mpi.txStore.parkingLotIndex.size(), len(mpi.txStore.batchedTxs), len(mpi.txStore.txHashMap))
}

func (mpi *mempoolImpl) GetTimeoutTransactions(rebroadcastDuration time.Duration) [][]*pb.Transaction {
	// all the tx whose live time is less than lowBoundTime should be rebroadcast
	mpi.logger.Debugf("------- Start gathering timeout txs, ttl index len is %d -----------", mpi.txStore.ttlIndex.index.Len())
	currentTime := time.Now().UnixNano()
	if currentTime < mpi.txStore.earliestTimestamp+rebroadcastDuration.Nanoseconds() {
		// if the latest incoming tx has not exceeded the timeout limit, then none will be timeout
		return [][]*pb.Transaction{}
	}

	txList := make([]*pb.Transaction, 0)
	timeoutItems := make([]*sortedTtlKey, 0)
	mpi.txStore.ttlIndex.index.Ascend(func(i btree.Item) bool {
		item := i.(*sortedTtlKey)
		if item.liveTime > math.MaxInt64 {
			// TODO(tyx): if this tx has rebroadcast many times and exceeded a final limit,
			// it is expired and will be removed from mempool
			return true
		}
		// if this tx has not exceeded the rebroadcast duration, break iteration
		timeoutTime := item.liveTime + rebroadcastDuration.Nanoseconds()
		txMap, ok := mpi.txStore.allTxs[item.account]
		if !ok || currentTime < timeoutTime {
			return false
		}
		tx := txMap.items[item.nonce].tx
		txList = append(txList, tx)
		timeoutItems = append(timeoutItems, item)
		return true
	})
	if len(txList) == 0 {
		return [][]*pb.Transaction{}
	}
	for _, item := range timeoutItems {
		// update the liveTime of timeout txs
		item.liveTime = currentTime
		mpi.txStore.ttlIndex.items[makeAccountNonceKey(item.account, item.nonce)] = currentTime
	}
	// shard txList into fixed size in case txList is too large to broadcast one time
	return shardTxList(txList, mpi.txSliceSize)
}

func shardTxList(txs []*pb.Transaction, batchLen uint64) [][]*pb.Transaction {
	begin := uint64(0)
	end := uint64(len(txs)) - 1
	totalLen := uint64(len(txs))

	// shape txs to fixed size in case totalLen is too large
	batchNums := totalLen / batchLen
	if totalLen%batchLen != 0 {
		batchNums++
	}
	shardedLists := make([][]*pb.Transaction, 0, batchNums)
	for i := uint64(0); i <= batchNums; i++ {
		actualLen := batchLen
		if end-begin+1 < batchLen {
			actualLen = end - begin + 1
		}
		if actualLen == 0 {
			continue
		}
		shardedList := make([]*pb.Transaction, actualLen)
		for j := uint64(0); j < batchLen && begin <= end; j++ {
			shardedList[j] = txs[begin]
			begin++
		}
		shardedLists = append(shardedLists, shardedList)
	}
	return shardedLists
}
