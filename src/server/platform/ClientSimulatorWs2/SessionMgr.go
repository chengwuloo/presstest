package main

//
// Created by YangZhi
// 			4/9/2019
//

import (
	"runtime"
	"sync"
	"sync/atomic"
)

//
type SessionMgr interface {
	Add(conn interface{}) Session
	Remove(peer Session)
	Get(sesID int64) Session
	Count() int64
	CloseAll()
	Wait()
}

//
var gSessMgr = newSessionMgr()

//
type defaultSessionMgr struct {
	peers map[int64]Session
	l     *sync.Mutex
	c     int64
}

//
func newSessionMgr() SessionMgr {
	return &defaultSessionMgr{l: &sync.Mutex{}, peers: map[int64]Session{}}
}

func (s *defaultSessionMgr) Add(conn interface{}) Session {
	s.l.Lock()
	defer s.l.Unlock()
	peer := newSession(conn)
	s.peers[peer.ID()] = peer
	atomic.AddInt64(&s.c, 1)
	return peer
}

//
func (s *defaultSessionMgr) Remove(peer Session) {
	s.l.Lock()
	defer s.l.Unlock()
	if _, ok := s.peers[peer.ID()]; ok {
		delete(s.peers, peer.ID())
		atomic.AddInt64(&s.c, -1)
	}
}

//
func (s *defaultSessionMgr) Get(sesID int64) Session {
	s.l.Lock()
	defer s.l.Unlock()
	if peer, ok := s.peers[sesID]; ok {
		return peer
	}
	return nil
}

//
func (s *defaultSessionMgr) Count() int64 {
	return atomic.LoadInt64(&s.c)
}

//
func (s *defaultSessionMgr) CloseAll() {
	s.l.Lock()
	defer s.l.Unlock()
	for _, peer := range s.peers {
		peer.Close()
	}
}

func (s *defaultSessionMgr) Wait() {
	for atomic.LoadInt64(&s.c) != 0 {
		runtime.Gosched()
	}
}
