package main

//
// Created by andy_ro@qq.com
// 			4/9/2019
//

import (
	"fmt"
	"log"
	"net"
	"time"
)

//
type DefTCPClient struct {
	sesID int64
}

//
func NewDefTCPClient() TCPClient {
	return &DefTCPClient{}
}

//会话ID
func (s *DefTCPClient) ID() int64 {
	return s.sesID
}

//会话
func (s *DefTCPClient) Session() Session {
	return gSessMgr.Get(s.sesID)
}

//关闭
func (s *DefTCPClient) Close() {
	peer := gSessMgr.Get(s.sesID)
	if peer != nil {
		peer.Close()
	}
}

//写
func (s *DefTCPClient) Write(msg interface{}) {
	peer := gSessMgr.Get(s.sesID)
	if peer != nil {
		peer.Write(msg)
	}
}

//连接
func (s *DefTCPClient) ConnectTCP(address string) {
	conn, err := net.DialTimeout("tcp", address, 3*time.Second)
	if err != nil {
		fmt.Println(err)
		return
	}
	peer := gSessMgr.Add(conn)
	peer.SetOnConnected(s.onConnected)
	peer.SetOnClosed(s.onClosed)
	peer.SetOnMessage(s.onMessage)
	peer.SetOnError(s.onError)
	peer.SetCloseCallback(s.remove)
	peer.OnEstablished()
}

//
func (s *DefTCPClient) remove(peer Session) {
	peer.OnDestroyed()
}

//
func (s *DefTCPClient) onConnected(peer Session) {
	s.sesID = peer.ID()
	//laddr := peer.Conn().(*websocket.Conn).LocalAddr().String()
	//raddr := peer.Conn().(*websocket.Conn).RemoteAddr().String()
	//peerID := peer.ID()
	//log.Printf("--- *** TCPClient :: onConnected [%05d:%s]%s %d\n", peerID, laddr, raddr, gSessMgr.Count())

	//s.Close()
}

//
func (s *DefTCPClient) onMessage(msg interface{}, peer Session) {
	log.Println("--- *** TCPClient :: onMessage ", msg)

}

//
func (s *DefTCPClient) onClosed(peer Session) {
	s.sesID = 0
	//log.Println("--- *** TCPClient :: onClosed")
}

//
func (s *DefTCPClient) onError(peer Session, err error) {
	log.Println("--- *** TCPClient :: onError ", err)
}
