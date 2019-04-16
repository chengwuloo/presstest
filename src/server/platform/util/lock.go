package util

//
// Created by YangZhi
// 			4/9/2019
//

import (
	"runtime"
	"sync/atomic"
)

//
type SpinLock struct {
	atm, cur uint64
}

//
func NewPxMutex() *SpinLock {
	return &SpinLock{atm: 0, cur: 1}
}

//
func (s *SpinLock) Wait() {
	idx := atomic.AddUint64(&s.atm, 1)
	for idx != s.cur {
		runtime.Gosched()
	}
}

//
func (s *SpinLock) Signal() {
	s.cur++
}

///////////////////////////////////////////////////////////////////////////////////////

//
const (
	Idle   = 100
	Locked = 101
)

//
type AxMutex struct {
	stat int32
}

//
func NewAxMutex() *AxMutex {
	return &AxMutex{stat: int32(Idle)}
}

//
func (s *AxMutex) Lock() {
	for s.tryLock() == false {
		runtime.Gosched()
	}
}

//
func (s *AxMutex) Unlock() {
	s.tryUnLock()
}

//
func (s *AxMutex) tryLock() bool {
	return atomic.CompareAndSwapInt32(&s.stat, Idle, Locked)
}

//
func (s *AxMutex) tryUnLock() bool {
	return atomic.CompareAndSwapInt32(&s.stat, Locked, Idle)
}
