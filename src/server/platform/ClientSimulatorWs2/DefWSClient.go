package main

//
// Created by andy_ro@qq.com
// 			4/9/2019
//

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	_ "github.com/davyxu/cellnet/codec/gogopb"

	"github.com/gorilla/websocket"
)

//
const (
	TagUserInfo = iota + 1001
)

//DefWSClient 保存客户端用户信息
type DefWSClient struct {
	HeartID  uint32 //心跳ID
	TimerID1 uint32 //延迟下注
	Cursor   int32  //时间轮游标
	SesID    int64  //会话ID
	Token    string //访问令牌
	UserID   int64  //用户ID
	Account  int64  //账号
	Nick     string //昵称
	HeadURL  string //头像
	HeadID   uint32 //
	Score    int64  //分数
	Pwd      []byte //密码
	AgentID  uint32 //代理服务器ID
	GameID   int32  //游戏类型
	RoomID   int32  //游戏房间
}

//
func NewDefWSClient() WSClient {
	return &DefWSClient{}
}

//
func (s *DefWSClient) ID() int64 {
	return s.SesID
}

//
func (s *DefWSClient) Session() Session {
	return gSessMgr.Get(s.SesID)
}

//
func (s *DefWSClient) Close() {
	peer := gSessMgr.Get(s.SesID)
	if peer != nil {
		peer.Close()
	}
}

//
func (s *DefWSClient) Write(msg interface{}) {
	peer := gSessMgr.Get(s.SesID)
	if peer != nil {
		peer.Write(msg)
	}
}

//
func (s *DefWSClient) ConnectTCP(address string) {
	dialer := websocket.Dialer{}
	dialer.Proxy = http.ProxyFromEnvironment
	dialer.HandshakeTimeout = 3 * time.Second
	u := url.URL{Scheme: "ws", Host: address, Path: "/ws"}
	conn, _, err := dialer.Dial(u.String(), nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	peer := gSessMgr.Add(conn)
	if peer != nil {
		peer.SetOnConnected(s.onConnected)
		peer.SetOnClosed(s.onClosed)
		peer.SetOnMessage(s.onMessage)
		peer.SetOnError(s.onError)
		peer.SetOnWritten(s.onWritten)
		peer.SetCloseCallback(s.remove)
		peer.OnEstablished()
	} else {
		conn.Close()
	}
}

//
func (s *DefWSClient) remove(peer Session) {
	peer.OnDestroyed()
}

//
func (s *DefWSClient) onConnected(peer Session) {
	s.SesID = peer.ID()
	peer.SetCtx(TagUserInfo, s)
	//if 0 == gSessMgr.Count()%int64(*deltaClients) {
	//timediff := TimeDiff(timenow, timestart2)
	//timestart2 = timenow
	//log.Printf("--- *** PID[%07d] WSClient [%05d:%d]%s %05d delta = %dms\n", os.Getpid(), peer.ID(), s.Account, s.Token, gSessMgr.Count(), timediff)
	//}
	sendPlayerLogin(peer, s.Token)
}

//
func (s *DefWSClient) onMessage(msg interface{}, peer Session) {
	//log.Printf("--- *** PID[%07d] WSClient :: onMessage %v\n", os.Getpid(), msg)
	DecodecAddReadTask(msg, peer)
}

//
func (s *DefWSClient) onClosed(peer Session) {
	s.SesID = 0
	//peer.SetCtx(TagUserInfo, nil)
	//if 0 == gSessMgr.Count()%500 {
	//log.Printf("--- *** PID[%07d] WSClient :: onClosed[%v]", os.Getpid(), peer.ID())
	//}
}

//
func (s *DefWSClient) onError(peer Session, err error) {
	log.Printf("--- *** PID[%07d] WSClient :: onError %d", os.Getpid(), gSessMgr.Count())
}

//
func (s *DefWSClient) onWritten(msg interface{}, peer Session) {
	h, ok := msg.(*Msg)
	if !ok || h == nil {
		return
	}
	//log.Println("--- *** PID[%07d] WSClient :: onWritten ", h.msg)
}
