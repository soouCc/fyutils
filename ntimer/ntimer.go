package ntimer

import (
	"fmt"
	"github.com/liu-junyong/go-logger/logger"
	"runtime/debug"
	"time"
)

func NewTimer() *NTimer {
	nt := new(NTimer)
	nt.init()

	return nt
}

type NTimer struct {
	C         chan string
	running   bool
	tm        time.Duration
	fn        func(args ...interface{})
	args      []interface{}
	StartTime int64 // 计时器开始时间
	setTime float64 // 计时器开始时间
}

func (nt *NTimer) init() {
	nt.C = make(chan string, 10)
	nt.running = false
}

func (nt *NTimer)Remaining()int{
	r:=int(nt.setTime)-nt.RunTime()
	return r
}

//获取生于时间
func (nt *NTimer)LeftTime()int64  {
	left:=int64(nt.tm.Seconds())-(time.Now().Unix() - nt.StartTime)
	logger.Debug(int64(nt.tm.Seconds()))
	logger.Debug(time.Now().Unix() - nt.StartTime)
	if left<0{
		left = 0
	}
	return left
}


func (nt *NTimer) RunTime() int {
	now1:=time.Now().Unix()
	t:=int(now1 - nt.StartTime)
	return t
}

func (nt *NTimer) GetTime() float64 {
	return nt.setTime
}

func (nt *NTimer) IsRunning() bool {
	return nt.running
}

func (nt *NTimer) Close() {
	if nt.running {
		nt.C <- "close"
	} else {
		close(nt.C)
	}
}

func (nt *NTimer) Set(sec float64, fn func(args ...interface{}), args ...interface{}) {
	tm := int(sec * 1000)
	nt.tm = time.Duration(time.Millisecond * time.Duration(tm))
	nt.setTime = sec
	logger.Info(fmt.Sprintf("ntier Set nt.tm is %f ", sec))
	nt.args = args
	nt.fn = fn
}

func (nt *NTimer) SetTime(sec float64) {
	tm := int(sec * 1000)
	nt.tm = time.Duration(time.Millisecond * time.Duration(tm))
	nt.setTime = sec
	logger.Info(fmt.Sprintf("ntier SetTime nt.tm is %d ", nt.tm))
}

func (nt *NTimer) TimerStop() int {
	time.Sleep(10 * time.Nanosecond)
	if nt.running {
		nt.C <- "stop"
		logger.Info(fmt.Sprintf("TimerStop stop ------------------- "))
		return int(time.Now().Unix() - nt.StartTime)
	}
	return 0
}

func (nt *NTimer) TimerRun() {
	logger.Info("ntimer TimerRun nt.running = ", nt.running)
	if nt.running {
		nt.C <- "do"
	}
}

func (nt *NTimer) TimerStart() {
	nt.StartTime = time.Now().Unix()
	fn := func() {
		defer func() {
			if err := recover(); err != nil {
				logger.Error(fmt.Sprintf("CacheError, err:%v ", err))
				logger.Error(string(debug.Stack()))
			}
		}()
		nt.running = true
		select {
		case sig := <-nt.C:
			nt.running = false
			logger.Info(fmt.Sprintf("ntier TimerStart case sig nt =%t %#v  sig = %s ", nt.running, nt, sig))
			switch sig {
			case "stop":
				break
			case "do":
				nt.fn(nt.args...)
			case "close":
				close(nt.C)
			}
		case <-time.After(nt.tm):
			nt.running = false
			logger.Info(fmt.Sprintf("ntier 执行 nt.tm is %d nt = %#v ", nt.tm, nt))
			nt.fn(nt.args...)
		}
	}
	go fn()
}
