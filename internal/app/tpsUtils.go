package app

import (
	"fmt"
	"github.com/meshplus/bitxhub-model/pb"
	"github.com/sirupsen/logrus"
	"strconv"
	"time"
)

const size = 10000

var timpStamp int64
var counter int64
var tpsList []float64
var delayList []float64

type orderBlockLedger struct {
	slice      []*pb.CommitEvent
	leadgerMap map[uint64]*pb.CommitEvent
}

func (bxh *BitXHub) calTPSorDelay(list []float64) {
	skip := len(list) / 8
	begin := skip
	end := len(list) - skip
	var sum float64
	for index := begin; index < end; index++ {
		sum += list[index]
	}

	if len(list) == len(tpsList) {
		bxh.logger.Infof("the average tps is: %f, sum is %f, begin is %d, end is %d ", sum/float64(end-begin), sum, begin, end)
	} else {
		bxh.logger.Infof("the average delay is: %f, sum is %f, begin is %d, end is %d ", sum/float64(end-begin), sum, begin, end)
	}
}

func (bxh *BitXHub) handleTxDelay(ev *pb.CommitEvent) error {
	endTime := ev.Block.BlockHeader.Timestamp
	var delay float64
	for _, tx := range ev.Block.Transactions.Transactions {
		var startTime int64
		startTime = tx.GetTimeStamp()
		delay += float64(endTime-startTime) / float64(time.Millisecond)
	}
	cnt := len(ev.LocalList)

	if delayList == nil {
		delayList = make([]float64, 0)
	}
	bxh.logger.WithFields(logrus.Fields{
		"delay": fmt.Sprintf("%s%s", strconv.FormatFloat(delay/float64(cnt), 'f', 5, 64), "ms"),
	}).Infof("append to delayList")

	delayList = append(delayList, delay/float64(cnt))
	return nil
}
