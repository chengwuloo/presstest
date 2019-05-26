package main

//
// Created by andy_ro@qq.com
// 			4/9/2019
//

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"runtime/debug"
	"server/pb/GameServer"
	"server/pb/Game_Common"
	"server/pb/HallServer"
	"server/platform/util"
	"strconv"
	"strings"
)

//ParallLoginRequest 发起并发连接/登陆请求
//-------------------------------------------------------------
func ParallLoginRequest() {
	go func() {
		//起始时间戳
		timestart = TimeNowMilliSec()
		for i := 0; i < *totalClients; i++ {
			//进入访问资源
			gSemLogin.Enter()
			//HTTP请求token
			token, ipaddr, err := HTTPGetToken(*httpaddr, *baseAccount+int64(i))
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
			//token := *tokenprefix + fmt.Sprintf("%d", *tokenstart+i)
			client.(*DefWSClient).Token = token
			client.(*DefWSClient).Account = *baseAccount + int64(i)
			//连接游戏大厅
			if *dynamic == 0 {
				client.ConnectTCP(*wsaddr)
			} else {
				client.ConnectTCP(ipaddr)
			}
		}
	}()
}

//ParallEnterRoomRequest 发起并发进房间请求
//-------------------------------------------------------------
func ParallEnterRoomRequest() {
	go func() {
		//起始时间戳
		timestart = TimeNowMilliSec()
		//游戏类型和房间都有效，则进入房间
		p, ok := GGames.Exist(int32(*subGameID))
		if 0 != *subGameID && 0 != *subRoomID && ok && p.Exist(int32(*subRoomID)) {
			for i := 0; i < *totalClients; i++ {
				//进入访问资源
				gSemEnter.Enter()
				sesID := PopPeer()
				if sesID > 0 {
					peer := gSessMgr.Get(sesID)
					if peer != nil {
						//登陆成功，进入房间
						// client := peer.GetCtx(TagUserInfo).(*DefWSClient)
						// client.GameID = int32(*subGameID) //保存
						// client.RoomID = int32(*subRoomID)
						// reqEnterRoom(peer, client.GameID, client.RoomID, client.Pwd[:])
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
	Data     *HTTPAuthData `json:"data,omitempty"`
}

//
type HTTPAuthData struct {
	Code    float64 `json:"code,omitempty"`
	Account string  `json:"account,omitempty"`
	URL     string  `json:"url,omitempty"`
	IP      string  `json:"domain,omitempty"`
	Port    string  `json:"port,omitempty"`
	Token   string  `json:"token,omitempty"`
}

//HTTPGetToken 客户端 - 查询token
//HTTPGetToken ipaddr 网关ipaddr
//-------------------------------------------------------------
func HTTPGetToken(httpaddr string, account int64) (token, ipaddr string, e error) {
	defer func() {
		if r := recover(); r != nil {
			log.Fatalln(debug.Stack())
		}
	}()
	requrl := fmt.Sprintf("http://%s/GameHandle?testAccount=%d", httpaddr, account)
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
	//log.Println(str)
	str = strings.Replace(str, "\\", "", -1)
	//str = str[1 : len(str)-1]
	body = util.Str2Byte(str)
	var authResult HTTPAuthResult
	if err := util.Byte2JSON(body, &authResult); err != nil {
		//log.Println("----->>>> ", str)
		log.Printf("--- *** PID[%07d] HTTPGetToken Byte2JSON %v", os.Getpid(), err)
		e = err
		return
	}
	values, _ := url.ParseQuery(authResult.Data.URL)
	//log.Println(values)
	for url := range values {
		sub := url[strings.Index(url, "?")+1:]
		dic := map[string]string{}
		for {
			s := strings.Index(sub, "=")
			if s == -1 {
				break
			}
			p := strings.Index(sub, "&")
			if p == -1 {
				dic[sub[0:s]] = sub[s+1:]
				break
			} else {
				dic[sub[0:s]] = sub[s+1 : p]
			}
			sub = sub[p+1:]
		}
		authResult.Data.IP = dic["domain"]
		authResult.Data.Port = dic["port"]
		authResult.Data.Token = dic["token"]
		break
	}
	//token
	token = authResult.Data.Token
	//网关ipaddr
	ipaddr = authResult.Data.IP + ":" + authResult.Data.Port
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
