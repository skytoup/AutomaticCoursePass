package modify

import (
	"autoCourse/config"
	"autoCourse/db"
	"autoCourse/util"
	"time"
)

// 定时检测，进行修改
func TickerCheckModify(checkWaitCountTime string) {
	go func() {
		dur, err := time.ParseDuration(checkWaitCountTime)
		util.ErrExit(err)
		t := time.NewTicker(dur) // 定时检测修改
		CheckWaitingCount()
		for {
			<-t.C
			CheckWaitingCount()
		}
	}()
}

// 检查等待数，大于一定量，开始修改
func CheckWaitingCount() {
	wc := db.GetAcountWaitCount()
	if wc >= config.MinModifyAccountCount { // 达到一定人数，开始修改
		TryModify()
	}
}
