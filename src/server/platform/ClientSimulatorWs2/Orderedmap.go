package main

//
// Created by andy_ro@qq.com
// 			4/9/2019
//

import (
	"container/list"
)

// Pair ...
type Pair struct {
	key interface{}
	val interface{}
}

// Orderedmap ...
type Orderedmap struct {
	list *list.List
}

// newOrderedmap ...
func newOrderedmap() *Orderedmap {
	return &Orderedmap{list: list.New()}
}

// insert 关键字排序插入排序表
func (s *Orderedmap) insert(key interface{}, value interface{}, compare func(a, b interface{}) bool) {
	pos := s.list.Front()
	for ; pos != nil; pos = pos.Next() {
		if !compare(key, pos.Value.(*Pair).key) {
			data := &Pair{key: key, val: value}
			s.list.InsertBefore(data, pos)
			break
		}
	}
	if pos == nil {
		data := &Pair{key: key, val: value}
		s.list.PushBack(data)
	}
}

// top 栈顶节点(key,val)
func (s *Orderedmap) top() (interface{}, interface{}) {
	if elem := s.list.Front(); elem != nil {
		data := elem.Value.(*Pair)
		return data.key, data.val
	}
	return nil, nil
}

// front 栈顶节点
func (s *Orderedmap) front() *list.Element {
	return s.list.Front()
}

// pop 移除栈顶节点
func (s *Orderedmap) pop() {
	if elem := s.list.Front(); elem != nil {
		s.list.Remove(elem)
	}
}

// empty 判空
func (s *Orderedmap) empty() bool {
	return s.list.Len() == 0
}
