package main

//
// Created by YangZhi
// 			4/9/2019
//

import (
	//_ "container/heap"
	//_ "fmt"

	"log"
	"server/platform/util"
	"sync/atomic"
)

//
type TimerCallback func(timerID uint32, dt int32, args interface{}) bool

// https://www.ibm.com/developerworks/cn/linux/l-cn-timers/
// https://blog.csdn.net/yueguanghaidao/article/details/46290539
// https://github.com/cloudwu/skynet/blob/master/skynet-src/skynet_timer.c

// ScopedTimer 基于最小堆(最小生成树)实现的线程局部定时器
// ScopedTimer 提供给线程内部使用，所以是安全的
type ScopedTimer interface {
	// 定时器所属协程ID
	ThreadID() uint32
	// 返回所在Worker
	GetWorker() Worker
	// 返回其中TimeWheel
	GetTimeWheel() TimeWheel
	// 不指定回调和timerID
	CreateTimer(delay, interval int32, args interface{}) uint32
	// 不指定回调指定timerID
	CreateTimerWithID(timerID uint32, delay, interval int32, args interface{}) uint32
	// 指定回调和timerID
	CreateTimerWithIDCB(timerID uint32, delay, interval int32, handler TimerCallback, args interface{}) uint32
	//  指定回调不指定timerID
	CreateTimerWithCB(delay, interval int32, handler TimerCallback, args interface{}) uint32
	// 撤销定时器
	RemoveTimer(timerID uint32)
	// 撤销所有
	RemoveTimers()
	// 轮询定时回调 默认 handler(timerID, dt, args)
	// 如果创建定时器时指定了回调函数handler
	// 则执行handler(timerID, dt, args)回调，否则执行update(timerID, dt, args)回调
	Poll(pid uint32, update TimerCallback) bool
}

// timerEvent 定时器事件
type timerEvent struct {
	// 定时器ID
	timerID uint32
	// 延迟执行等待(s)
	//delay int32
	// 执行间隔时间(s)
	interval int32
	// 上次开始执行时间
	last Timestamp
	// 下次开始执行时间
	//expr Timestamp
	// 回调函数
	handler TimerCallback
	// 回调参数
	args interface{}
}

// scopedTimer 基于最小堆(最小生成树)实现的定时器
type scopedTimer struct {
	x        uint32          // 用于自动生成timerID
	pid      uint32          // 定时器所属goroutine
	timers   Orderedmap      // 排序 map[timestamp] = timer
	timerIDs map[uint32]bool // 保存要删除的timerID集合
	w        Worker          // 定时器所在worker
	t        TimeWheel       // 定时器的时间轮
}

// NewScopedTimer 创建
// NewScopedTimer w 指定Worker
// NewScopedTimer t 指定TimeWheel
func NewScopedTimer(pid uint32, w Worker, t TimeWheel) ScopedTimer {
	if util.GoroutineID() != pid {
		log.Fatalln("NewScopedTimer")
	}
	return &scopedTimer{
		pid:      pid,
		w:        w,
		t:        t,
		timers:   *newOrderedmap(),
		timerIDs: map[uint32]bool{}}
}

// ThreadID 定时器所属协程ID
func (s *scopedTimer) ThreadID() uint32 {
	return s.pid
}

// GetWorker 返回所在Worker
func (s *scopedTimer) GetWorker() Worker {
	return s.w
}

// GetTimeWheel 返回其中TimeWheel
func (s *scopedTimer) GetTimeWheel() TimeWheel {
	return s.t
}

// RemoveTimer 撤销定时器
func (s *scopedTimer) RemoveTimer(timerID uint32) {
	s.addRemoves(timerID)
}

// addRemoves 添加到撤销表
func (s *scopedTimer) addRemoves(timerID uint32) {
	// 线程安全
	s.AssertInThread(util.GoroutineID())
	if timerID != 0 {
		s.timerIDs[timerID] = true
	}
}

// isRemoveID 判断是否撤销
func (s *scopedTimer) isRemoveID(timerID uint32) bool {
	if _, ok := s.timerIDs[timerID]; ok {
		// timerID 在删除表中则移除 timerIDs = append(timerIDs[:i], timerIDs[i+1:]...)
		delete(s.timerIDs, timerID)
		return true
	}
	return false
}

// RemoveTimers 撤销所有
func (s *scopedTimer) RemoveTimers() {

}

// CreateTimerWithID 不带回调带ID
func (s *scopedTimer) CreateTimerWithID(timerID uint32, delay, interval int32, args interface{}) uint32 {
	return s.createTimer(timerID, delay, interval, nil, args)
}

// CreateTimer 不带回调不带ID
func (s *scopedTimer) CreateTimer(delay, interval int32, args interface{}) uint32 {
	return s.CreateTimerWithID(atomic.AddUint32(&s.x, 1), delay, interval, args)
}

// CreateTimerWithIDCB 带回调带ID
func (s *scopedTimer) CreateTimerWithIDCB(timerID uint32, delay, interval int32, handler TimerCallback, args interface{}) uint32 {
	return s.createTimer(timerID, delay, interval, handler, args)
}

// CreateTimerWithCB 带回调不带ID
func (s *scopedTimer) CreateTimerWithCB(delay, interval int32, handler TimerCallback, args interface{}) uint32 {
	return s.CreateTimerWithIDCB(atomic.AddUint32(&s.x, 1), delay, interval, handler, args)
}

// compare 比较大小
func compare(a, b interface{}) bool {
	return a.(Timestamp).Greater(b.(Timestamp))
}

// AssertInThread 安全断言
func (s *scopedTimer) AssertInThread(pollingThreadID uint32) {
	if s.pid != pollingThreadID {
		log.Fatalf(" scopedTimer::assertInThread\n")
	}
}

// createTimer 带回调带ID
func (s *scopedTimer) createTimer(timerID uint32, delay, interval int32, handler TimerCallback, args interface{}) uint32 {
	// 线程安全
	s.AssertInThread(util.GoroutineID())
	// 创建 timer
	timer := &timerEvent{timerID: timerID, interval: interval, last: TimeNowMilliSec(), handler: handler, args: args}
	// 放在 map[timestamp] = timer 中，并对 timestamp 进行关键字排序
	s.timers.insert(TimeAdd(timer.last, delay), timer, compare)
	// 打印调试
	// s.Keys()
	// 返回定时器ID
	return timerID
}

// Keys 从栈顶节点开始打印
func (s *scopedTimer) Keys() {
	i := 0
	for elem := s.timers.front(); elem != nil; elem = elem.Next() {
		key := elem.Value.(*Pair).key.(Timestamp)
		val := elem.Value.(*Pair).val.(*timerEvent)
		log.Printf("--- *** ScopedTimer[%d:%v] = %d", i, key.SinceUnixEpoch(), val.timerID)
		i++
	}
}

// Poll 定时器轮询 true定时器已空 false定时器不空
func (s *scopedTimer) Poll(pid uint32, update TimerCallback) bool {
	// 线程安全
	s.AssertInThread(pid)
	if s.timers.empty() {
		return true
	}
	// 进入循环
	for {
		//log.Printf("--- *** ScopedTimer:: Poll %s...", CreateToken())
		now := TimeNowMilliSec()
		// 取出栈顶Timestamp
		k, v := s.timers.top()
		ts := k.(Timestamp)
		t := v.(*timerEvent)
		if ts.Greater(now) {
			return false
		}
		// 先移除
		s.timers.pop()
		// 判断是否撤销
		if s.isRemoveID(t.timerID) {
			// 删除
		} else if t.handler != nil { // 先执行handler回调如果有的话
			// 执行handler回调 handler(timerID, elapsed, args)
			if t.handler(t.timerID, TimeDiff(now, t.last), t.args) {
				// 下次开始执行时间，从当前handler执行之后开始算
				if t.interval > 0 {
					t.last = now
					// 再次添加到有序表
					s.timers.insert(TimeNowMilliSec().Add(t.interval), t, compare)
				} else {
					// 不再需要则销毁
				}
			} else {
				// 不再需要则销毁
			} // 否则执行update回调如果有的话
		} else if update != nil {
			// 执行update回调 update(timerID, elapsed, args)
			if update(t.timerID, TimeDiff(now, t.last), t.args) {
				// 下次开始执行时间，从当前update执行之后开始算
				if t.interval > 0 {
					t.last = now
					// 再次添加到有序表
					s.timers.insert(TimeNowMilliSec().Add(t.interval), t, compare)
				} else {
					// 不再需要则销毁
				}
			} else {
				// 不再需要则销毁
			}
		} else {
			// 不再需要则销毁
		}
		// 容器已空则返回
		if s.timers.empty() {
			return true
		}
	}
}
