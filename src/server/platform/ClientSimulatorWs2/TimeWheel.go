package main

//
// Created by YangZhi
// 			4/9/2019
//

import (
	"log"
	"os"
	"sync"
	"sync/atomic"
)

//无效游标
const (
	INVALIDCURSORPOS int32 = iota - 1
)

//TimeWheel 时间轮盘，处理超时会话
type TimeWheel interface {
	//SetTimer 指定所在定时器
	SetTimer(t ScopedTimer)
	//GetTimer 返回所在定时器
	GetTimer() ScopedTimer

	//UpdateWheel 定时器tick调用
	UpdateWheel() (v []int64)

	//PushBucket 登陆成功压入桶
	PushBucket(val int64, timeout int32) int32

	//UpdateBucket oldcuror 元素当前所在桶位置
	//UpdateBucket timeout 心跳超时清理时间
	//UpdateBucket 心跳间隔时间(interval)收到心跳包调用更新桶元素
	UpdateBucket(oldcuror int32, val int64, timeout int32) int32
}

//Bucket 容纳会话的桶
type Bucket struct {
	v map[int64]bool
	l *sync.Mutex
}

//
func newBucket() *Bucket {
	return &Bucket{v: map[int64]bool{}, l: &sync.Mutex{}}
}

//Add 添加元素到桶
func (s *Bucket) Add(val int64) {
	s.l.Lock()
	s.v[val] = true
	s.l.Unlock()
}

//Remove 从桶里面移除
func (s *Bucket) Remove(val int64) bool {
	s.l.Lock()
	if _, ok := s.v[val]; ok {
		delete(s.v, val)
		s.l.Unlock()
		return true
	}
	s.l.Unlock()
	return false
}

//Pop 弹出桶里面所有元素
func (s *Bucket) Pop() (v []int64) {
	s.l.Lock()
	//取出桶内所有id
	for id := range s.v {
		v = append(v, id)
	}
	//清空桶
	if len(s.v) > 0 {
		s.v = map[int64]bool{}
	}
	s.l.Unlock()
	return
}

//DefTimeWheel 时间轮实现
type DefTimeWheel struct {
	pid    uint32      //协程
	cursor int32       //秒针
	size   int32       //轮盘大小 [0 1 2 3 4 5 6 7 8 9] size = 10
	ring   []*Bucket   //环形数组
	t      ScopedTimer //所在定时器
}

//NewTimeWheel 轮盘大小(size) >=
//NewTimeWheel 心跳超时清理时间(timeout) >
//NewTimeWheel 心跳间隔时间(interval)
//----------------------------------------------------------
func NewTimeWheel(pid uint32, size int32) TimeWheel {
	s := &DefTimeWheel{pid: pid, size: size, ring: make([]*Bucket, size)}
	for i := int32(0); i < s.size; i++ {
		s.ring[i] = newBucket()
	}
	return s
}

//SetTimer 指定所在定时器
func (s *DefTimeWheel) SetTimer(t ScopedTimer) {
	s.t = t
}

//GetTimer 返回所在定时器
func (s *DefTimeWheel) GetTimer() ScopedTimer {
	return s.t
}

//UpdateWheel 定时器tick调用
func (s *DefTimeWheel) UpdateWheel() (v []int64) {
	// defer func() {
	// 	if r := recover(); r != nil {
	// 		log.Fatalln(debug.Stack())
	// 	}
	// }()
	s.cursor = atomic.AddInt32(&s.cursor, 1) % s.size
	//log.Printf("--- *** PID[%07d] [%05d] UpdateWheel[检测] size:%d cursor:%d\n", os.Getpid(), s.pid, s.size, s.cursor)
	v = s.ring[s.cursor].Pop()
	return
}

//PushBucket 登陆成功压入桶
//PushBucket return newcursor 初始游标位置
func (s *DefTimeWheel) PushBucket(val int64, timeout int32) int32 {
	// defer func() {
	// 	if r := recover(); r != nil {
	// 		log.Fatalln(debug.Stack())
	// 	}
	// }()
	newcursor := (s.cursor + timeout) % s.size
	//log.Printf("--- *** PID[%07d] [%05d] PushBucket[{{{压入}}}] size:%d csursor:%d newcursor:%d\n", os.Getpid(), s.pid, s.size, s.cursor, newcursor)
	bucket := s.ring[newcursor]
	bucket.Add(val)
	return newcursor
}

//UpdateBucket oldcuror 元素当前所在桶位置
//UpdateBucket timeout 心跳超时清理时间
//UpdateBucket 心跳间隔时间(interval)收到心跳包调用更新桶元素
//UpdateBucket return newcursor 新的游标位置
//----------------------------------------------------------
func (s *DefTimeWheel) UpdateBucket(oldcuror int32, val int64, timeout int32) int32 {
	// defer func() {
	// 	if r := recover(); r != nil {
	// 		log.Fatalln(debug.Stack())
	// 	}
	// }()
	if oldcuror == INVALIDCURSORPOS {
		//log.Printf("--- *** PID[%07d] [%05d] UpdateBucket oldcuror == INVALIDCURSORPOS\n", os.Getpid(), s.pid)
		return INVALIDCURSORPOS
	}
	newcursor := INVALIDCURSORPOS
	bucket := s.ring[oldcuror]
	//先从原来的桶中尝试移除
	if true == bucket.Remove(val) {
		//移除成功，说明没有被超时清理
		newcursor = (s.cursor + timeout) % s.size
		bucket = s.ring[newcursor]
		bucket.Add(val)
		log.Printf("--- *** PID[%07d] [%05d] UpdateBucket[{{{更新}}}][SUCC] size:%d cursor:%d oldcursor:%d newcursor:%d\n", os.Getpid(), s.pid, s.size, s.cursor, oldcuror, newcursor)
	} else {
		//移除失败，已经被超时清理
		//log.Printf("--- *** PID[%07d] [%05d] UpdateBucket[{{{更新}}}][FAILED] size:%d cursor:%d oldcursor:%d newcursor:%d\n", os.Getpid(), s.pid, s.size, s.cursor, oldcuror, newcursor)
	}

	return newcursor
}
