package main

import (
	"time"
)

type RepitableEvent struct {
	TimeInterval time.Duration
	CallBack     func()
	ticker       *time.Ticker
}

func makeRepitableEvent(interval time.Duration, callBack func()) RepitableEvent {
	ticker := time.NewTicker(interval)
	go func() {
		for {
			select {
			case <-ticker.C:
				callBack()
			}
		}
	}()
	return RepitableEvent{interval, callBack, ticker}
}

type RepitableEventHandler struct {
	eventLookUp map[uint]RepitableEvent
}

func (this *RepitableEventHandler) RegisterRepitableEvent(key uint, callBack func(), interval time.Duration) {
	this.eventLookUp[key] = makeRepitableEvent(interval, callBack)
}
func (this *RepitableEventHandler) ChangeCallBack(key uint, callBack func()) {
	this.eventLookUp[key].ticker.Stop()

	oldTimeInterval := this.eventLookUp[key].TimeInterval
	this.eventLookUp[key] = makeRepitableEvent(oldTimeInterval, callBack)
}
func (this *RepitableEventHandler) ChangeInterval(key uint, interval time.Duration) {
	this.eventLookUp[key].ticker.Stop()

	oldCallBack := this.eventLookUp[key].CallBack
	this.eventLookUp[key] = makeRepitableEvent(interval, oldCallBack)
}
func (this *RepitableEventHandler) StopAll() {
	for _, event := range this.eventLookUp {
		event.ticker.Stop()
	}
}

func makeRepitableEventHandler() RepitableEventHandler {
	return RepitableEventHandler{make(map[uint]RepitableEvent)}
}
