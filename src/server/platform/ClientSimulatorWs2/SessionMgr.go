package main

//
// Created by andy_ro@qq.com
// 			4/9/2019
//

import (
	"sync"
)

//
type SessionMgr interface {
	Add(conn interface{}) Session
	Remove(peer Session)
	Get(sesID int64) Session
	Count() int64
	Stop()
	Wait()
}

//
var gSessMgr = newSessionMgr()

//
type defaultSessionMgr struct {
	peers map[int64]Session
	l     *sync.Mutex
	c     *sync.Cond
	exit  bool
}

//
func newSessionMgr() SessionMgr {
	s := &defaultSessionMgr{l: &sync.Mutex{}, peers: map[int64]Session{}}
	s.c = sync.NewCond(s.l)
	s.exit = false
	return s

}

func (s *defaultSessionMgr) Add(conn interface{}) Session {
	s.l.Lock()
	if !s.exit {
		peer := newSession(conn)
		s.peers[peer.ID()] = peer
		//引用计数加1
		peer.AddRef("SessionMgr")
		s.l.Unlock()
		return peer
	}
	s.l.Unlock()
	return nil
}

//
func (s *defaultSessionMgr) Remove(peer Session) {
	s.l.Lock()
	if _, ok := s.peers[peer.ID()]; ok {
		//if 0 != peer.RefCount("SessionMgr") {
		//	log.Panicf("SessionMgr::Remove error")
		//}
		//引用计数为0，移除peer
		delete(s.peers, peer.ID())
		//log.Printf("--- *** SessionMgr::Remove peers = %v", len(s.peers))
		if s.exit && len(s.peers) == 0 {
			s.c.Signal()
		}
	}
	s.l.Unlock()
}

//
func (s *defaultSessionMgr) Get(sesID int64) Session {
	s.l.Lock()
	if peer, ok := s.peers[sesID]; ok {
		s.l.Unlock()
		return peer
	}
	s.l.Unlock()
	return nil
}

//Count 有效连接数
func (s *defaultSessionMgr) Count() int64 {
	return 0
}

//
func (s *defaultSessionMgr) Stop() {
	s.l.Lock()
	s.exit = true
	for _, peer := range s.peers {
		peer.Close()
	}
	s.l.Unlock()
}

func (s *defaultSessionMgr) Wait() {
	s.l.Lock()
	s.c.Wait()
	//log.Printf("SessionMgr::Wait exit...")
	s.l.Unlock()
}
