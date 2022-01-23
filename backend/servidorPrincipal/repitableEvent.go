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
	EventLookUp map[uint]RepitableEvent
}

func (this *RepitableEventHandler) RegisterRepitableEvent(key uint, callBack func(), interval time.Duration) {
	this.EventLookUp[key] = makeRepitableEvent(interval, callBack)
}
func (this *RepitableEventHandler) ChangeCallBack(key uint, callBack func()) {
	this.EventLookUp[key].ticker.Stop()

	oldTimeInterval := this.EventLookUp[key].TimeInterval
	this.EventLookUp[key] = makeRepitableEvent(oldTimeInterval, callBack)
}
func (this *RepitableEventHandler) ChangeInterval(key uint, interval time.Duration) {
	this.EventLookUp[key].ticker.Stop()

	oldCallBack := this.EventLookUp[key].CallBack
	this.EventLookUp[key] = makeRepitableEvent(interval, oldCallBack)
}
func (this *RepitableEventHandler) StopAll() {
	for _, event := range this.EventLookUp {
		event.ticker.Stop()
	}
}
func (this *RepitableEventHandler) Stop(key uint) {
	this.EventLookUp[key].ticker.Stop()
}

func makeRepitableEventHandler() RepitableEventHandler {
	return RepitableEventHandler{make(map[uint]RepitableEvent)}
}
