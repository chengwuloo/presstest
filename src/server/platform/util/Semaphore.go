package util

//
// Created by andy_ro@qq.com
// 			4/9/2019
//

import "sync"

//Semaphore 信号量互斥访问控制
type Semaphore struct {
	w        *sync.Mutex
	l        *sync.Mutex
	c        *sync.Cond
	avail    int64
	initsize int64
}

//NewSemaphore 初始化initsize个并发访问资源
func NewSemaphore(initsize int64) *Semaphore {
	s := &Semaphore{initsize: initsize, avail: initsize, l: &sync.Mutex{}, w: &sync.Mutex{}}
	s.c = sync.NewCond(s.l)
	return s
}

//Enter 进入访问资源
func (s *Semaphore) Enter() {
wait:
	s.wait()
	s.w.Lock()
	if s.avail > 0 {
		s.avail--
		s.w.Unlock()
	} else {
		s.w.Unlock()
		goto wait
	}
}

//Leave 离开释放资源
func (s *Semaphore) Leave() {
	s.w.Lock()
	if s.avail < s.initsize {
		s.avail++
		if s.avail == 1 {
			s.c.Signal()
		}
	}
	s.w.Unlock()
}

//wait 等待资源
func (s *Semaphore) wait() {
	s.l.Lock()
	if s.avail == 0 {
		s.c.Wait()
	}
	s.l.Unlock()
}

//FreeSemaphore 信号量互斥访问控制
type FreeSemaphore struct {
	l        *sync.Mutex
	avail    int64
	initsize int64
}

//NewFreeSemaphore 初始化initsize个并发访问资源
func NewFreeSemaphore(initsize int64) *FreeSemaphore {
	s := &FreeSemaphore{initsize: initsize, avail: initsize, l: &sync.Mutex{}}
	return s
}

//Enter 进入访问资源
func (s *FreeSemaphore) Enter() (bv bool) {
	s.l.Lock()
	if s.avail > 0 {
		s.avail--
		bv = true
	}
	s.l.Unlock()
	return
}

//Leave 离开释放资源
func (s *FreeSemaphore) Leave() {
	s.l.Lock()
	if s.avail < s.initsize {
		s.avail++
	}
	s.l.Unlock()
}

//
var gSem = NewSemaphore(10)
var ix = 10

//
func TestSemaphore() {

	for i := 0; i < 100; i++ {
		go func() {
			for {
				gSem.Enter()
				ix--
				println("1======= ", ix)
				ix++
				println("2======= ", ix)
				gSem.Leave()
			}
		}()
	}
}

//
func OnInputTestSemaphore(str string) int {
	switch str {
	case "w":
		{
			for i := 0; i < 30; i++ {
				gSem.Leave()
			}
			return 0
		}
	}
	return 0
}
