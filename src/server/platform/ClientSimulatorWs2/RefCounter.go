package main

//
// Created by andy_ro@qq.com
// 			4/9/2019
//

import (
	"sync/atomic"
)

//RefCounter 引用计数
type RefCounter struct {
	refc int64
}

//AddRef 引用计数加1
func (s *RefCounter) AddRef(name string) {
	atomic.AddInt64(&s.refc, 1)
	//log.Printf("%v::AddRef refc == %v", name, newc)
}

//Release 引用计数减1
func (s *RefCounter) Release(name string) int64 {
	//if atomic.LoadInt64(&s.refc) <= 0 {
	//	log.Panicf("%v::Release refc == %v", name, s.refc)
	//}
	newc := atomic.AddInt64(&s.refc, -1)
	//log.Printf("%v::Release refc == %v", name, newc)
	return newc
}

//RefCount 引用计数
func (s *RefCounter) RefCount(name string) int64 {
	curc := atomic.LoadInt64(&s.refc)
	//log.Printf("%v::RefCount refc == %v", name, curc)
	return curc
}

//NewRefCounter 创建引用计数
func NewRefCounter() *RefCounter {
	//初始化引用计数0
	return &RefCounter{refc: int64(0)}
}
