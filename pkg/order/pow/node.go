package pow

import (
	"context"
	"crypto/sha256"
	"fmt"
	"github.com/cbergoon/merkletree"
	"github.com/ethereum/go-ethereum/event"
	"github.com/meshplus/bitxhub-core/agency"
	"github.com/meshplus/bitxhub-kit/types"
	"github.com/meshplus/bitxhub-model/pb"
	"github.com/meshplus/bitxhub/pkg/order"
	"github.com/meshplus/bitxhub/pkg/order/etcdraft"
	raftproto "github.com/meshplus/bitxhub/pkg/order/etcdraft/proto"
	"github.com/meshplus/bitxhub/pkg/order/mempool"
	"github.com/meshplus/bitxhub/pkg/peermgr"
	"github.com/sirupsen/logrus"
	"math/big"
	"strconv"
	"sync"
	"time"
)

type Node struct {
	ID               uint64
	commitC          chan *pb.CommitEvent         // block channel
	logger           logrus.FieldLogger           // logger
	mempool          mempool.MemPool              // transaction pool
	proposeC         chan *raftproto.RequestBatch // proposed listenReadyBlock, input channel
	stateC           chan *mempool.ChainState
	txCache          *mempool.TxCache // cache the transactions received from api
	batchMgr         *etcdraft.BatchTimer
	lastExec         uint64               // the index of the last-applied block
	packSize         int                  // maximum number of transaction packages
	blockTick        time.Duration        // block packed period
	peerMgr          peermgr.PeerManager  // network manager
	difficulty       uint64               //pow's difficultly
	getChainMetaFunc func() *pb.ChainMeta // current chain meta

	ctx    context.Context
	cancel context.CancelFunc
	sync.RWMutex
}

func (n *Node) GetPendingTxByHash(hash *types.Hash) pb.Transaction {
	return n.mempool.GetTransaction(hash)
}

func (n *Node) Start() error {
	go n.txCache.ListenEvent()
	go n.run()
	return nil
}

func (n *Node) Stop() {
	n.cancel()
}

func (n *Node) GetPendingNonceByAccount(account string) uint64 {
	return n.mempool.GetPendingNonceByAccount(account)
}

func (n *Node) DelNode(delID uint64) error {
	return nil
}

func (n *Node) Prepare(tx pb.Transaction) error {
	if err := n.Ready(); err != nil {
		return err
	}
	n.txCache.RecvTxC <- tx
	return nil
}

func (n *Node) Commit() chan *pb.CommitEvent {
	return n.commitC
}

func (n *Node) Step(msg []byte) error {
	return nil
}

func (n *Node) Ready() error {
	return nil
}

// ReportState get the latest block.
func (n *Node) ReportState(height uint64, blockHash *types.Hash, txHashList []*types.Hash) {
	state := &mempool.ChainState{
		Height:     height,
		BlockHash:  blockHash,
		TxHashList: txHashList,
	}
	n.stateC <- state
}

// Quorum needn't Quorum in pow.
func (n *Node) Quorum() uint64 {
	return 1
}

func (n *Node) SubscribeTxEvent(ch chan<- pb.Transactions) event.Subscription {
	return n.mempool.SubscribeTxEvent(ch)
}

func init() {
	agency.RegisterOrderConstructor("pow", NewNode)
}

func NewNode(opts ...agency.ConfigOption) (agency.Order, error) {
	// TODO: node sync(including block_sync and tx_sync).
	var options []order.Option
	for i, _ := range opts {
		options = append(options, opts[i].(order.Option))
	}

	config, err := order.GenerateConfig(options...)
	if err != nil {
		return nil, fmt.Errorf("generate config: %w", err)
	}
	if err != nil {
		return nil, fmt.Errorf("new leveldb: %w", err)
	}
	ctx, cancel := context.WithCancel(context.Background())

	batchTimeout, memConfig, err := generatePowConfig(config.RepoRoot)
	mempoolConf := &mempool.Config{
		ID:              config.ID,
		ChainHeight:     config.Applied,
		Logger:          config.Logger,
		StoragePath:     config.StoragePath,
		GetAccountNonce: config.GetAccountNonce,

		BatchSize:      memConfig.BatchSize,
		PoolSize:       memConfig.PoolSize,
		TxSliceSize:    memConfig.TxSliceSize,
		TxSliceTimeout: memConfig.TxSliceTimeout,
	}
	batchC := make(chan *raftproto.RequestBatch, 1024)
	mempoolInst, err := mempool.NewMempool(mempoolConf)
	if err != nil {
		return nil, fmt.Errorf("create mempool instance: %w", err)
	}
	txCache := mempool.NewTxCache(mempoolConf.TxSliceTimeout, mempoolConf.TxSliceSize, config.Logger)
	batchTimerMgr := etcdraft.NewTimer(batchTimeout, config.Logger)

	//TODO: read the difficulty from last block.
	lastblocknum := config.GetChainMetaFunc().Height
	lastblock, err := config.GetBlockByHeight(lastblocknum)
	if err != nil {
		fmt.Errorf("read lastblock from ledger: %w", err)
	}
	difficulty := CalcDifficulty(lastblock)
	powNode := &Node{
		ID:               config.ID,
		commitC:          make(chan *pb.CommitEvent, 1024),
		stateC:           make(chan *mempool.ChainState),
		lastExec:         config.Applied,
		mempool:          mempoolInst,
		txCache:          txCache,
		batchMgr:         batchTimerMgr,
		peerMgr:          config.PeerMgr,
		proposeC:         batchC,
		logger:           config.Logger,
		ctx:              ctx,
		cancel:           cancel,
		difficulty:       difficulty,
		getChainMetaFunc: config.GetChainMetaFunc,
	}
	powNode.logger.Infof("POW lastExec = %d", powNode.lastExec)
	powNode.logger.Infof("POW batch timeout = %v", batchTimeout)
	return powNode, nil
}

func (n *Node) run() {
	var (
		abort = make(chan struct{})
	)

	for {
		select {
		case <-n.ctx.Done():
			n.Stop()

		case txSet := <-n.txCache.TxSetC:
			// start batch timer when this node receives the first transaction
			if !n.batchMgr.IsBatchTimerActive() {
				n.batchMgr.StartBatchTimer()
			}
			if batch := n.mempool.ProcessTransactions(txSet.Transactions, true, true); batch != nil {
				n.batchMgr.StopBatchTimer()
				n.proposeC <- batch
			}

		case <-n.batchMgr.BatchTimeoutEvent():
			n.batchMgr.StopBatchTimer()
			n.logger.Debug("Batch timer expired, try to create a batch")
			if n.mempool.HasPendingRequest() {
				if batch := n.mempool.GenerateBlock(); batch != nil {
					n.postProposal(batch)
				}
			} else {
				n.logger.Debug("The length of priorityIndex is 0, skip the batch timer")
			}

		case state := <-n.stateC:
			if state.Height%10 == 0 {
				n.logger.WithFields(logrus.Fields{
					"height": state.Height,
					"hash":   state.BlockHash.String(),
				}).Info("Report checkpoint")
			}
			n.mempool.CommitTransactions(state)

		case proposal := <-n.proposeC:
			n.logger.WithFields(logrus.Fields{
				"proposal_height": proposal.Height,
				"tx_count":        len(proposal.TxList.Transactions),
			}).Debugf("Receive proposal from mempool")

			if proposal.Height != n.lastExec+1 {
				n.logger.Warningf("Expects to execute seq=%d, but get seq=%d, ignore it", n.lastExec+1, proposal.Height)
				return
			}
			n.logger.Infof("======== Call execute, height=%d", proposal.Height)
			mineBlock := &pb.Block{
				BlockHeader: &pb.BlockHeader{
					Version:    []byte("1.0.0"),
					Number:     proposal.Height,
					Timestamp:  time.Now().Unix(),
					Difficulty: n.difficulty,
					BlockNonce: 0,
				},
				Transactions: proposal.TxList,
			}

			localList := make([]bool, len(proposal.TxList.Transactions))
			for i := 0; i < len(proposal.TxList.Transactions); i++ {
				localList[i] = true
			}
			l2Root, err := n.buildTxMerkleTree(mineBlock.Transactions.Transactions)
			if err != nil {
				panic(err)
			}
			mineBlock.BlockHeader.TxRoot = l2Root
			mineBlock.BlockHeader.ParentHash = n.getChainMetaFunc().BlockHash
			mineBlock.PowBlockHash = mineBlock.Hash()
			n.mine(mineBlock, localList, abort, 0)
			break
		}
	}

}

// start mine
func (n *Node) mine(block *pb.Block, localList []bool, abort chan struct{}, nonce uint64) {
search:
	for {
		select {
		case <-abort:
			// Mining terminated, update stats and abort
			n.logger.Infof("======== POW nonce search aborted.start next term mine========")

			// TODO: start next mining.

			break
		default:
			diff := big.NewInt(int64(block.BlockHeader.Difficulty))
			target := new(big.Int).Div(two256, diff)
			result := CalcBlockNonce(block.PowBlockHash, nonce)
			if new(big.Int).SetBytes(result).Cmp(target) <= 0 {
				n.logger.WithFields(logrus.Fields{
					"pow nonce found and reported": nonce,
				})
				executeEvent := &pb.CommitEvent{
					Block:     block,
					LocalList: localList,
				}
				n.commitC <- executeEvent
				n.lastExec++
				n.logger.WithFields(logrus.Fields{
					"new block height = ": n.lastExec,
				})
				break search
			}
			nonce++
		}
	}
}

func (n *Node) postProposal(batch *raftproto.RequestBatch) {
	n.proposeC <- batch
	n.batchMgr.StartBatchTimer()
}

func CalcDifficulty(block *pb.Block) uint64 {
	// TODO: Calculate the difficulty.
	return 1
}

func CalcBlockNonce(blockHash *types.Hash, blockNonce uint64) []byte {
	//TODO: Calculate the PoW value of this BlockNonce.
	var hashdata = blockHash.String() + strconv.Itoa(int(blockNonce))
	var hash = sha256.New()
	hash.Write([]byte(hashdata))
	hashed := hash.Sum(nil)
	return hashed
}

// buildTxMerkleTree return the blockHeaader.TxRoot
func (n *Node) buildTxMerkleTree(txs []pb.Transaction) (*types.Hash, error) {
	// TODO: implement buildTxMerkleTree.
	return nil, nil
}

func calcMerkleRoot(contents []merkletree.Content) (*types.Hash, error) {
	if len(contents) == 0 {
		return &types.Hash{}, nil
	}

	tree, err := merkletree.NewTree(contents)
	if err != nil {
		return nil, err
	}

	return types.NewHash(tree.MerkleRoot()), nil
}

// validate the nonce and tx
func (n *Node) verifyPow(block *pb.Block) bool {
	return true
}

//TODO: handle the fork and follow the longest chain.
