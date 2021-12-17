package app

import (
	"github.com/meshplus/bitxhub-core/governance"
	orderPeerMgr "github.com/meshplus/bitxhub-core/peer-mgr"
	"github.com/meshplus/bitxhub-kit/types"
	"github.com/meshplus/bitxhub-model/pb"
	"github.com/meshplus/bitxhub/internal/model/events"
	"github.com/meshplus/bitxhub/internal/repo"
	"github.com/sirupsen/logrus"
	"sync/atomic"
	"time"
)

func (bxh *BitXHub) start() {
	bxh.currentBlockHash = bxh.Ledger.GetChainMeta().BlockHash
	bxh.mockledger = &orderBlockLedger{
		slice:      make([]*pb.CommitEvent, 0, size),
		leadgerMap: make(map[uint64]*pb.CommitEvent),
	}
	bxh.delayCh = make(chan *pb.CommitEvent, 1000)

	go bxh.listenEvent()

	lastTime := atomic.LoadInt64(&timpStamp)
	lastCounter := atomic.LoadInt64(&counter)
	go func() {
		ticker := time.NewTicker(2 * time.Second)
		for {
			select {
			case <-bxh.Ctx.Done():
				return
			case <-ticker.C:
				currentTime := atomic.LoadInt64(&timpStamp)
				currentCounter := atomic.LoadInt64(&counter)
				if currentTime != lastTime {
					time := (currentTime - lastTime) / int64(time.Second)
					c := float64(currentCounter-lastCounter) / float64(time)
					if c != 0 {
						bxh.logger.WithFields(logrus.Fields{
							"tps":            c,
							"lastCounter":    lastCounter,
							"currentCounter": currentCounter,
							"currentTime":    currentTime,
							"lastTime":       lastTime,
						}).Infof("append to tpsList")

						lastCounter = currentCounter
						lastTime = currentTime
						if tpsList == nil {
							tpsList = make([]float64, 0)
						}

						tpsList = append(tpsList, c)
					}
				}
			}
		}
	}()

	go func() {
		//var index uint64
		for {
			select {
			case commitEvent := <-bxh.Order.Commit():
				bxh.logger.WithFields(logrus.Fields{
					"height": commitEvent.Block.BlockHeader.Number,
					"count":  len(commitEvent.Block.Transactions.Transactions),
				}).Info("Generated block")
				//bxh.BlockExecutor.ExecuteBlock(commitEvent)

				// store in slice
				bxh.mockledger.slice = append(bxh.mockledger.slice, commitEvent)
				txCount := len(commitEvent.LocalList)
				atomic.AddInt64(&counter, int64(txCount))
				atomic.StoreInt64(&timpStamp, commitEvent.Block.BlockHeader.Timestamp)

				// record blockHash
				commitEvent.Block.BlockHeader.ParentHash = bxh.currentBlockHash
				bxh.currentBlockHash = commitEvent.Block.BlockHeader.Hash()
				commitEvent.Block.BlockHash = bxh.currentBlockHash

				txHashList := make([]*types.Hash, 0)
				go bxh.Order.ReportState(commitEvent.Block.BlockHeader.Number, commitEvent.Block.BlockHash, txHashList)

				bxh.delayCh <- commitEvent
				//// store in map
				//bxh.mockledger.leadgerMap[index]=commitEvent
				//index++
			case <-bxh.Ctx.Done():
				return
			}
		}
	}()

	go func() {
		ticker := time.NewTicker(200 * time.Second)
		for {
			select {
			case <-bxh.Ctx.Done():
				return
			case <-ticker.C:
				bxh.calTPSorDelay(tpsList)
				bxh.calTPSorDelay(delayList)
			}
		}
	}()

	// calculate tx delay
	go func() {
		for {
			select {
			case <-bxh.Ctx.Done():
				return
			case ev := <-bxh.delayCh:
				err := bxh.handleTxDelay(ev)
				if err != nil {
					bxh.logger.Errorf("handle cal delay err: %s", err)
					return
				}
			}
		}
	}()
}

func (bxh *BitXHub) listenEvent() {
	blockCh := make(chan events.ExecutedEvent)
	orderMsgCh := make(chan orderPeerMgr.OrderMessageEvent)
	nodeCh := make(chan events.NodeEvent)
	configCh := make(chan *repo.Repo)

	blockSub := bxh.BlockExecutor.SubscribeBlockEvent(blockCh)
	orderMsgSub := bxh.PeerMgr.SubscribeOrderMessage(orderMsgCh)
	nodeSub := bxh.BlockExecutor.SubscribeNodeEvent(nodeCh)
	configSub := bxh.repo.SubscribeConfigChange(configCh)

	defer blockSub.Unsubscribe()
	defer orderMsgSub.Unsubscribe()
	defer nodeSub.Unsubscribe()
	defer configSub.Unsubscribe()

	for {
		select {
		case ev := <-blockCh:
			go bxh.Order.ReportState(ev.Block.BlockHeader.Number, ev.Block.BlockHash, ev.TxHashList)
			go bxh.Router.PutBlockAndMeta(ev.Block, ev.InterchainMeta)
		case ev := <-orderMsgCh:
			go func() {
				if err := bxh.Order.Step(ev.Data); err != nil {
					bxh.logger.Error(err)
				}
			}()
		case ev := <-nodeCh:
			switch ev.NodeEventType {
			case governance.EventLogout:
				go func() {
					if err := bxh.Order.Ready(); err != nil {
						bxh.logger.Error(err)
						return
					}
					if err := bxh.Order.DelNode(ev.NodeId); err != nil {
						bxh.logger.Error(err)
					}
				}()
			}
		case config := <-configCh:
			bxh.ReConfig(config)
		case <-bxh.Ctx.Done():
			return
		}
	}
}
