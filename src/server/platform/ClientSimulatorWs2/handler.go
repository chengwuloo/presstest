package main

//
// Created by YangZhi
// 			4/9/2019
//

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"runtime/debug"
	"server/pb/GameServer"
	"server/pb/Game_Common"
	"server/pb/HallServer"
	"server/platform/util"
	"strconv"
	"strings"
	"sync/atomic"
)

//
var x = int64(0)

//
var k int

//ParallLoginRequest 发起并发连接/登陆请求
//-------------------------------------------------------------
func ParallLoginRequest() {
	log.Printf("ParallLoginRequest %d ...", atomic.AddInt64(&x, 1))
	go func() {
		//起始时间戳
		timestart = TimeNowMilliSec()
		//for i := 0; i < *numClient; i++ {
		for i := 0; i < *totalClient; i++ {
			gSemLogin.Enter()
			//HTTP请求token
			token, err := HTTPGetToken(*httpaddr, *baseAccount+int64(k))
			if token == "" || err != nil {
				continue
			}
			//当前时间戳
			//timenow = TimeNowMilliSec()
			// timdiff := TimeDiff(timenow, timestart)
			// if timdiff >= int32(*deltaTime) {
			// 	timestart = timenow
			// 	c := gSessMgr.Count()
			// 	delteConn := c - curConn
			// 	curConn = c
			// 	log.Printf("--- *** detla = %dms deltaClients = %03d", timdiff, delteConn)
			// }
			//websocket客户端
			client := NewDefWSClient()
			//token := *tokenprefix + fmt.Sprintf("%d", *tokenstart+k)
			client.(*DefWSClient).Token = token
			client.(*DefWSClient).Account = *baseAccount + int64(k)
			k++
			//连接游戏大厅
			client.ConnectTCP(*wsaddr)
		}
	}()
}

//ParallEnterRoomRequest 发起并发进房间请求
//-------------------------------------------------------------
func ParallEnterRoomRequest() {
	log.Printf("ParallEnterRoomRequest %d ...", atomic.AddInt64(&x, 1))
	go func() {
		//起始时间戳
		timestart = TimeNowMilliSec()
		for i := 0; i < *numClients2; i++ {
			//游戏类型和房间都有效，则进入房间
			p, ok := GGames.Exist(int32(*subGameID))
			if 0 != *subGameID && 0 != *subRoomID && ok && p.Exist(int32(*subRoomID)) {
				sesID := PopPeer()
				if sesID > 0 {
					peer := gSessMgr.Get(sesID)
					if peer != nil {
						//登陆成功，获取游戏列表
						reqGameListInfo(peer)
					}
				}
			}
		}
	}()
}

//
type HTTPAuthResult struct {
	Type     float64       `json:"type,omitempty"`
	Maintype string        `json:"maintype,omitempty"`
	Data     *HTTPAuthData `json:"d,omitempty"`
}

//
type HTTPAuthData struct {
	Code float64 `json:"code,omitempty"`
	URL  string  `json:"url,omitempty"`
}

//HTTPGetToken 客户端 - 查询token
//-------------------------------------------------------------
func HTTPGetToken(httpaddr string, account int64) (token string, e error) {
	defer func() {
		if r := recover(); r != nil {
			log.Fatalln(debug.Stack())
		}
	}()
	//requrl := fmt.Sprintf("http://%s/GameHandle?testAccount=%d", httpaddr, account)
	requrl := fmt.Sprintf("http://%s?testAccount=%d", httpaddr, account)
	//log.Printf("--- *** PID[%07d] HTTPGetToken >>> %v", os.Getpid(), requrl)
	rsp, err := http.Get(requrl)
	if err != nil {
		log.Printf("--- *** PID[%07d] HTTPGetToken httpGet %v\n", os.Getpid(), err)
		e = err
		return
	}
	defer rsp.Body.Close()
	body, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		e = err
		log.Printf("--- *** PID[%07d] HTTPGetToken ReadAll %v\n", os.Getgid(), err)
		return
	}
	str := util.Byte2Str(body)
	str = strings.Replace(str, "\\", "", -1)
	str = str[1 : len(str)-1]
	body = util.Str2Byte(str)

	var authResult HTTPAuthResult
	if err := util.Byte2JSON(body, &authResult); err != nil {
		//log.Println("----->>>> ", str)
		log.Printf("--- *** PID[%07d] HTTPGetToken Byte2JSON %v", os.Getpid(), err)
		e = err
		return
	}
	url := authResult.Data.URL
	pos := strings.Index(url, "=")
	token = url[pos+1:]
	//log.Printf("--- *** PID[%07d] token >>> %v", os.Getpid(), token)
	return
}

//sendPlayerLogin 登陆大厅
//-------------------------------------------------------------
func sendPlayerLogin(peer Session, token string) {
	reqdata := &HallServer.LoginMessage{}
	val, _ := strconv.ParseUint("F5F5F5F5", 16, 32)
	reqdata.Header = &Game_Common.Header{}
	reqdata.Header.Sign = int32(val)
	reqdata.Session = token
	msg := newMsg(
		uint8(Game_Common.MAINID_MAIN_MESSAGE_CLIENT_TO_HALL),
		uint8(Game_Common.MESSAGE_CLIENT_TO_HALL_SUBID_CLIENT_TO_HALL_LOGIN_MESSAGE_REQ),
		reqdata)
	util.Log("UserClient", "DefWSClient", "sendPlayerLogin", reqdata)
	peer.Write(msg)
}

//sendKeepAlive 发送心跳包
//-------------------------------------------------------------
func sendKeepAlive(peer Session, token string) {
	reqdata := &Game_Common.KeepAliveMessage{}
	val, _ := strconv.ParseUint("F5F5F5F5", 16, 32)
	reqdata.Header = &Game_Common.Header{}
	reqdata.Header.Sign = int32(val)
	reqdata.Session = token //
	msg := newMsg(
		uint8(Game_Common.MAINID_MAIN_MESSAGE_CLIENT_TO_HALL),
		uint8(Game_Common.MESSAGE_CLIENT_TO_SERVER_SUBID_KEEP_ALIVE_REQ),
		reqdata)
	//util.Log("UserClient", "Player", "sendKeepAlive", reqdata)
	peer.Write(msg)
}

//reqGameListInfo 取游戏信息
//-------------------------------------------------------------
func reqGameListInfo(peer Session) {
	reqdata := &HallServer.GetGameMessage{}
	val, _ := strconv.ParseUint("F5F5F5F5", 16, 32)
	reqdata.Header = &Game_Common.Header{}
	reqdata.Header.Sign = int32(val)
	msg := newMsg(
		uint8(Game_Common.MAINID_MAIN_MESSAGE_CLIENT_TO_HALL),
		uint8(Game_Common.MESSAGE_CLIENT_TO_HALL_SUBID_CLIENT_TO_HALL_GET_GAME_ROOM_INFO_REQ),
		reqdata)
	util.Log("UserClient", "Player", "reqGameListInfo", reqdata)
	peer.Write(msg)
}

//reqGameserverInfo 获取游戏IP
//-------------------------------------------------------------
func reqGameserverInfo(peer Session, gameID, roomID int32) {
	reqdata := &HallServer.GetGameServerMessage{}
	val, _ := strconv.ParseUint("F5F5F5F5", 16, 32)
	reqdata.Header = &Game_Common.Header{}
	reqdata.Header.Sign = int32(val)
	reqdata.GameId = uint32(gameID) //游戏类型
	reqdata.RoomId = uint32(roomID) //游戏房间
	client := peer.GetCtx(TagUserInfo).(*DefWSClient)
	client.GameID = gameID //保存
	client.RoomID = roomID
	msg := newMsg(
		uint8(Game_Common.MAINID_MAIN_MESSAGE_CLIENT_TO_HALL),
		uint8(Game_Common.MESSAGE_CLIENT_TO_HALL_SUBID_CLIENT_TO_HALL_GET_GAME_SERVER_MESSAGE_REQ),
		reqdata)
	util.Log("UserClient", "Player", "reqGameserverInfo", reqdata)
	log.Printf("--- *** PID[%07d] player[%d:%d:%s] %d:%s %d:%s \n",
		os.Getpid(),
		client.UserID,
		client.Account,
		client.Token,
		//
		gameID,
		GGames.ByID[gameID].Name,
		//
		roomID,
		GGames.ByID[gameID].ByID[roomID])
	peer.Write(msg)
}

//reqEnterRoom 进入房间
//-------------------------------------------------------------
func reqEnterRoom(peer Session, gameID, roomID int32, pwd []byte) {
	reqdata := &GameServer.MSG_C2S_UserEnterMessage{}
	val, _ := strconv.ParseUint("F5F5F5F5", 16, 32)
	reqdata.Header = &Game_Common.Header{}
	reqdata.Header.Sign = int32(val)
	reqdata.GameId = gameID          //游戏类型
	reqdata.RoomId = roomID          //游戏房间
	reqdata.DynamicPassword = pwd[:] //动态密码
	msg := newMsg(
		uint8(Game_Common.MAINID_MAIN_MESSAGE_CLIENT_TO_GAME_SERVER),
		uint8(GameServer.SUBID_SUB_C2S_ENTER_ROOM_REQ),
		reqdata)
	util.Log("UserClient", "Player", "reqEnterRoom", reqdata)
	peer.Write(msg)
}

//reqPlayerReady 玩家就绪
//-------------------------------------------------------------
func reqPlayerReady(peer Session) {
	reqdata := &GameServer.MSG_C2S_UserReadyMessage{}
	val, _ := strconv.ParseUint("F5F5F5F5", 16, 32)
	reqdata.Header = &Game_Common.Header{}
	reqdata.Header.Sign = int32(val)
	msg := newMsg(
		uint8(Game_Common.MAINID_MAIN_MESSAGE_CLIENT_TO_GAME_SERVER),
		uint8(GameServer.SUBID_SUB_C2S_USER_READY_REQ),
		reqdata)
	util.Log("UserClient", "Player", "reqPlayerReady", reqdata)
	peer.Write(msg)
}

//reqPlayerLeave 玩家离开
//-------------------------------------------------------------
func reqPlayerLeave(peer Session, userID int32, gameID, roomID, Type int32) {
	reqdata := &GameServer.MSG_C2S_UserLeftMessage{}
	val, _ := strconv.ParseUint("F5F5F5F5", 16, 32)
	reqdata.Header = &Game_Common.Header{}
	reqdata.Header.Sign = int32(val)
	reqdata.UserId = uint32(userID)
	reqdata.GameId = gameID
	reqdata.RoomId = roomID
	reqdata.Type = Type
	msg := newMsg(
		uint8(Game_Common.MAINID_MAIN_MESSAGE_CLIENT_TO_GAME_SERVER),
		uint8(GameServer.SUBID_SUB_C2S_USER_LEFT_REQ),
		reqdata)
	//util.Log("UserClient", "Player", "reqPlayerLeave", reqdata)
	peer.Write(msg)
}
