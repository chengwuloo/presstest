package main

//
// Created by andy_ro@qq.com
// 			4/9/2019
//

import (
	"log"
	"math/rand"
	"os"
	"reflect"
	"server/pb/BenCiBaoMa"
	"server/pb/Brnn"
	"server/pb/ErBaGang"
	"server/pb/GameServer"
	"server/pb/Game_Common"
	"server/pb/HallServer"
	"server/pb/HongHei"
	"server/pb/Longhu"
	"server/platform/util"
	"sync/atomic"
)

//
type Player struct {
	entry *sentry
	i     int64
}

//
func newPlayer(s *sentry) *Player {
	return &Player{entry: s}
}

//randPlaceJet 随机下注
func (s *Player) randPlaceJet(peer Session) {
	switch int32(*subGameID) {
	case GGames.ByName["二八杠"].ID:
		{
			//用户主动下注 [1,3]
			x := rand.Intn(3) + 1
			sendPlayerPlaceJetErBaGang(peer, int32(x), 1000)
		}
	case GGames.ByName["龙虎斗"].ID:
		{
			//用户主动下注 [1,3]
			x := rand.Intn(3) + 1
			sendPlayerPlaceJetLonghu(peer, int32(x), 2000)
		}
	case GGames.ByName["百人牛牛"].ID:
		{
			//用户主动下注 [1,4]
			x := rand.Intn(4) + 1
			sendPlayerPlaceJetBrnn(peer, int32(x), 10*100)
		}
	case GGames.ByName["红黑大战"].ID:
		{
			//用户主动下注 [0,2]
			x := rand.Intn(3)
			sendPlayerPlaceJetHongHei(peer, int32(x), 100)
		}
	case GGames.ByName["奔驰宝马"].ID:
		{
			//用户主动下注 [0,7]
			x := rand.Intn(7)
			sendPlayerPlaceJetBcbm(peer, int32(x), 200)
			client := peer.GetCtx(TagUserInfo).(*DefWSClient)
			s.entry.RemoveTimer(client.TimerID1)
			//进入访问资源
			//if gSemJetton.Enter() {
			client.TimerID1 = s.entry.RunAfter(2000, client)
			//} else {
			//	s.entry.GetCell().Append(func() {
			//		if gSemJetton.Enter() {
			//			client.TimerID1 = s.entry.RunAfter(500, client)
			//		}
			//	})
			//}
		}
	}
}

//OnTimer tick定时器及心跳定时器
//-------------------------------------------------------------
func (s *Player) OnTimer(timerID uint32, dt int32, args interface{}) bool {
	if s.entry.tickID == timerID {
		//tick检查，轮询当前Cell上所有会话
		peerIDs := s.entry.GetTimeWheel().UpdateWheel()
		for _, id := range peerIDs {
			peer := gSessMgr.Get(id)
			if peer == nil {
				return false
			}
			client := peer.GetCtx(TagUserInfo).(*DefWSClient)
			log.Printf("--- *** PID[%07d] player[%d:%d:%s] :: OnTimer 心跳超时 !!!!!!!!!!!!!!!",
				os.Getpid(), client.UserID, client.Account, client.Token)
			peer.Close() //超时关闭连接
		}
	} else if args != nil {
		if client, ok := args.(*DefWSClient); ok {
			peer := gSessMgr.Get(client.ID())
			if peer == nil {
				return false
			}
			//发送心跳包
			if client.HeartID == timerID {
				sendKeepAlive(peer, client.Token)
				//不需要了
				return false
			}
			if client.TimerID1 == timerID {
				//离开释放资源
				gSemJetton.Leave()
				s.randPlaceJet(peer)
				return false
			}
		}
	}
	return true
}

// OnTimerOrder - 来自服务端 - 上下分请求
//-------------------------------------------------------------
func (s *Player) OnTimerOrder(timerID uint32, dt int32, args interface{}) bool {

	return true
}

// resultPlayerLogin - 来自服务端 - 玩家登陆结果
//-------------------------------------------------------------
func (s *Player) resultPlayerLogin(msg interface{}, peer Session) {
	rspdata, ok := msg.(*HallServer.LoginMessageResponse)
	if !ok {
		log.Fatalln("resultPlayerLogin ", ok)
	}
	util.Logy("UserClient", "Player", "resultPlayerLogin", rspdata)
	if rspdata.RetCode == 0 {
		atomic.AddInt64(&gClientsSucc, 1)
		AddSuccPeer(peer.ID())
		client := peer.GetCtx(TagUserInfo).(*DefWSClient)
		//登陆成功，保存用户数据
		client.UserID = rspdata.UserId
		client.HeadID = rspdata.HeadId
		client.Nick = rspdata.NickName
		client.Score = rspdata.Score
		client.Pwd = rspdata.GamePass
		client.AgentID = rspdata.AgentId
		//登陆成功，压入桶元素
		//client.Cursor = s.entry.GetTimeWheel().PushBucket(peer.ID(), int32(*timeout)/1000)
		//登陆成功，间隔发送心跳包
		//client.HeartID = s.entry.RunAfter(int32(*heartbeat), client)
		//游戏类型和房间都有效，则进入房间
		// p, ok := GGames.Exist(int32(*subGameID))
		// if 0 != *subGameID && 0 != *subRoomID && ok && p.Exist(int32(*subRoomID)) {
		// 	//登陆成功，获取游戏列表
		// 	reqGameListInfo(peer)
		// }
		GLoginedUsers.AddUserInfo(client.UserID, client.Account, client.AgentID)
		//发送上下分请求
		//sendHTTPOrderRequest(peer)
	} else {
		atomic.AddInt64(&gClientsFailed, 1)
		//失败关闭
		peer.Close()
	}
	c := atomic.AddInt64(&gClients, 1)
	if c%int64(*numClients) == 0 {
		t := TimeNowMilliSec()
		d := TimeDiff(t, timestart)
		timestart = t
		elapsed += d
		if c >= int64(*totalClients) {
			log.Printf("--- *** PID[%07d] resultPlayerLogin[%03d][%03d][%03d] elapsed:%dms All:%dms\n", os.Getpid(), gClients, gClientsSucc, gClientsFailed, d, elapsed)
		} else {
			log.Printf("--- *** PID[%07d] resultPlayerLogin[%03d][%03d][%03d] elapsed:%dms\n", os.Getpid(), gClients, gClientsSucc, gClientsFailed, d)
		}
	}
	//不需要了关闭连接
	peer.Close()
	//离开释放资源
	gSemLogin.Leave()
}

//
// resultKeepAlive - 来自服务端 - 应答心跳包
//-------------------------------------------------------------
func (s *Player) resultKeepAlive(msg interface{}, peer Session) {
	s.i++
	//if s.i%10 == 0 {
	//rspdata := msg.(*Game_Common.KeepAliveMessageResponse)
	//util.Logy("UserClient", "Player", "resultKeepAlive", rspdata)
	//}
	client := peer.GetCtx(TagUserInfo).(*DefWSClient)
	if client != nil {
		//收到心跳包，更新桶元素
		client.Cursor = s.entry.GetTimeWheel().UpdateBucket(client.Cursor, peer.ID(), int32(*timeout)/1000)
		//继续发送心跳包
		client.HeartID = s.entry.RunAfter(int32(*heartbeat), client)
	}
}

//resultGameInfo 服务端返回 - 取游戏信息
//-------------------------------------------------------------
func (s *Player) resultGameListInfo(msg interface{}, peer Session) {
	rspdata, ok := msg.(*HallServer.GetGameMessageResponse)
	if !ok {
		log.Fatalln("resultGameListInfo ", ok)
	}
	util.Log("UserClient", "Player", "resultGameListInfo", rspdata)
	// reqGameserverInfo(peer,
	// 	GGames.ByName["龙虎斗"].ID,
	// 	GGames.ByName["龙虎斗"].ByName["初级房"])
	//*subGameID = GGames.ByName["红黑大战"].ID
	//*subRoomID = GGames.ByName["红黑大战"].ByName["体验房"]
	//reqGameserverInfo(peer, int32(*subGameID), int32(*subRoomID))
}

//resultGameserverInfo 服务端返回 - 获取游戏IP
//-------------------------------------------------------------
func (s *Player) resultGameserverInfo(msg interface{}, peer Session) {
	rspdata, ok := msg.(*HallServer.GetGameServerMessageResponse)
	if !ok {
		log.Fatalln("resultGameserverInfo ", ok)
	}
	util.Log("UserClient", "Player", "resultGameserverInfo", rspdata)
	if rspdata.RetCode == 0 {
		client := peer.GetCtx(TagUserInfo).(*DefWSClient)
		//进入房间
		reqEnterRoom(peer, client.GameID, client.RoomID, client.Pwd[:])
	} else {

	}
}

//resultPlayerEnterRoom 服务端返回 - 进入房间
//-------------------------------------------------------------
func (s *Player) resultPlayerEnterRoom(msg interface{}, peer Session) {
	rspdata, ok := msg.(*GameServer.MSG_S2C_UserEnterMessageResponse)
	if !ok {
		log.Fatalln("resultPlayerEnterRoom ", ok)
	}
	util.Log("UserClient", "Player", "resultPlayerEnterRoom", rspdata)
	if rspdata.RetCode == 0 {
		atomic.AddInt64(&gClientsSucc, 1)
		//玩家就绪
		reqPlayerReady(peer)
	} else {
		atomic.AddInt64(&gClientsFailed, 1)
		client := peer.GetCtx(TagUserInfo).(*DefWSClient)
		//进入房间失败///////////////////////////////////
		log.Printf("--- *** PID[%07d] player[%d:%d:%s] :: resultPlayerEnterRoom 进入房间失败 !!!!!!!!!!!!!! \n%v\n %v\n",
			os.Getpid(),
			client.UserID,
			client.Account,
			client.Token,
			reflect.TypeOf(rspdata).Elem(), util.JSON2Str(rspdata))
	}
	c := atomic.AddInt64(&gClients, 1)
	if c%int64(*numClients2) == 0 {
		t := TimeNowMilliSec()
		d := TimeDiff(t, timestart)
		timestart = t
		elapsed += d
		if c >= int64(*totalClients) {
			log.Printf("--- *** PID[%07d] resultPlayerEnterRoom[%03d][%03d][%03d] elapsed:%dms All:%dms\n", os.Getpid(), gClients, gClientsSucc, gClientsFailed, d, elapsed)
			gStep = StepAll
		} else {
			log.Printf("--- *** PID[%07d] resultPlayerEnterRoom[%03d][%03d][%03d] elapsed:%dms\n", os.Getpid(), gClients, gClientsSucc, gClientsFailed, d)
		}
	}
	//离开释放资源
	gSemEnter.Leave()
}

//onPlayerEnterNotify 服务端返回 - 玩家进入返回
//-------------------------------------------------------------
func (s *Player) onPlayerEnterNotify(msg interface{}, peer Session) {
	rspdata, ok := msg.(*GameServer.MSG_S2C_UserBaseInfo)
	if !ok {
		log.Fatalln("onPlayerEnterNotify ", ok)
	}
	util.Log("UserClient", "Player", "onPlayerEnterNotify", rspdata)
}

//onPlayerScoreNotify 服务端返回 - 玩家积分信息
//-------------------------------------------------------------
func (s *Player) onPlayerScoreNotify(msg interface{}, peer Session) {
	rspdata, ok := msg.(*GameServer.MSG_S2C_UserScoreInfo)
	if !ok {
		log.Fatalln("onPlayerScoreNotify ", ok)
	}
	util.Log("UserClient", "Player", "onPlayerScoreNotify", rspdata)
}

//onPlayerStatusNotify 服务端返回 - 玩家状态
//-------------------------------------------------------------
func (s *Player) onPlayerStatusNotify(msg interface{}, peer Session) {
	rspdata, ok := msg.(*GameServer.MSG_S2C_GameUserStatus)
	if !ok {
		log.Fatalln("onPlayerStatusNotify ", ok)
	}
	util.Log("UserClient", "Player", "onPlayerStatusNotify", rspdata)
}

//resultPlayerReady 服务端返回 - 玩家就绪
//-------------------------------------------------------------
func (s *Player) resultPlayerReady(msg interface{}, peer Session) {
	rspdata, ok := msg.(*GameServer.MSG_S2C_UserReadyMessageResponse)
	if !ok {
		log.Fatalln("resultPlayerReady ", ok)
	}
	util.Log("UserClient", "Player", "resultPlayerReady", rspdata)
}

//resultPlayerLeave 服务端返回 - 离开响应
//-------------------------------------------------------------
func (s *Player) resultPlayerLeave(msg interface{}, peer Session) {
	rspdata, ok := msg.(*GameServer.MSG_C2S_UserLeftMessageResponse)
	if !ok {
		log.Fatalln("resultPlayerLeave ", ok)
	}
	util.Log("UserClient", "Player", "resultPlayerLeave", rspdata)
}

//resultOrderNotify 服务端返回 - 上下分通知
//-------------------------------------------------------------
// func (s *Player) resultOrderNotify(msg interface{}, peer Session) {
// 	rspdata, ok := msg.(*HallServer.OrderNotifyMessage)
// 	if !ok {
// 		log.Fatalln("resultOrderNotify ", ok)
// 	}
// 	util.Logy("UserClient", "Player", "resultOrderNotify", rspdata)
// 	//发送上下分请求
// 	//sendHTTPOrderRequest(peer)
// }

//resultProxyOrderNotify 服务端返回 - 上下分通知
//-------------------------------------------------------------
func (s *Player) resultProxyOrderNotify(msg interface{}, peer Session) {
	rspdata, ok := msg.(*Game_Common.ProxyNotifyOrderScoreMessage)
	if !ok {
		log.Fatalln("resultProxyOrderNotify ", ok)
	}
	util.Logy("UserClient", "Player", "resultProxyOrderNotify", rspdata)
	//发送上下分请求
	//sendHTTPOrderRequest(peer)
}

//---------------------------------------------------------------------------------------------
//二八杠
//---------------------------------------------------------------------------------------------

//onGameStartErBaGang 开始游戏
//-------------------------------------------------------------
func (s *Player) onGameStartErBaGang(msg interface{}, peer Session) {
	rspdata, ok := msg.(*ErBaGang.CMD_S_GameStart)
	if !ok {
		log.Fatalln("onGameStartErBaGang ", ok)
	}
	client := peer.GetCtx(TagUserInfo).(*DefWSClient)
	util.Logx("UserClient", "Player", client.UserID, client.Account, "onGameStartErBaGang", rspdata)
	// log.Printf("--- *** PID[%07d] player[%d:%d:%s] :: onGameStartErBaGang \n%v\n %v\n",
	// 	os.Getpid(),
	// 	client.UserID,
	// 	client.Account,
	// 	client.Token,
	// 	reflect.TypeOf(rspdata).Elem(), util.JSON2Str(rspdata))
}

//onGameEndErBaGang 结束游戏
//-------------------------------------------------------------
func (s *Player) onGameEndErBaGang(msg interface{}, peer Session) {
	rspdata, ok := msg.(*ErBaGang.CMD_S_GameEnd)
	if !ok {
		log.Fatalln("onGameEndErBaGang ", ok)
	}
	client := peer.GetCtx(TagUserInfo).(*DefWSClient)
	util.Logx("UserClient", "Player", client.UserID, client.Account, "onGameEndErBaGang", rspdata)
}

//onSceneGameStartErBaGang 开始游戏场景
//-------------------------------------------------------------
func (s *Player) onSceneGameStartErBaGang(msg interface{}, peer Session) {
	rspdata, ok := msg.(*ErBaGang.CMD_S_Scene_GameStart)
	if !ok {
		log.Fatalln("onSceneGameStartErBaGang ", ok)
	}
	client := peer.GetCtx(TagUserInfo).(*DefWSClient)
	util.Logx("UserClient", "Player", client.UserID, client.Account, "onSceneGameStartErBaGang", rspdata)
}

//onSceneGameEndErBaGang 结束游戏场景
//-------------------------------------------------------------
func (s *Player) onSceneGameEndErBaGang(msg interface{}, peer Session) {
	rspdata, ok := msg.(*ErBaGang.CMD_S_Scene_GameEnd)
	if !ok {
		log.Fatalln("onSceneGameEndErBaGang ", ok)
	}
	client := peer.GetCtx(TagUserInfo).(*DefWSClient)
	util.Logx("UserClient", "Player", client.UserID, client.Account, "onSceneGameEndErBaGang", rspdata)
}

//onPlayerListErBaGang 玩家列表
//-------------------------------------------------------------
func (s *Player) onPlayerListErBaGang(msg interface{}, peer Session) {
	// rspdata, ok := msg.(*ErBaGang.CMD_S_PlayerList)
	// if !ok {
	// 	log.Fatalln(ok)
	// }
	// client := peer.GetCtx(TagUserInfo).(*DefWSClient)
	// util.Logx("UserClient", "Player", client.UserID, client.Account, "onPlayerListErBaGang", rspdata)
}

//onPlaceJetSuccessErBaGang 下注成功
//-------------------------------------------------------------
func (s *Player) onPlaceJetSuccessErBaGang(msg interface{}, peer Session) {
	rspdata, ok := msg.(*ErBaGang.CMD_S_PlaceJetSuccess)
	if !ok {
		log.Fatalln("onPlaceJetSuccessErBaGang ", ok)
	}
	client := peer.GetCtx(TagUserInfo).(*DefWSClient)
	util.Logx("UserClient", "Player", client.UserID, client.Account, "onPlaceJetSuccessErBaGang", rspdata)
	//下注成功///////////////////////////////////
	// log.Printf("--- *** PID[%07d] player[%d:%d:%s] :: onPlaceJetSuccessErBaGang \n%v\n %v\n",
	// 	os.Getpid(),
	// 	client.UserID,
	// 	client.Account,
	// 	client.Token,
	// 	reflect.TypeOf(rspdata).Elem(), util.JSON2Str(rspdata))
	//离开释放资源
	gSemJetton.Leave()
}

//onPlaceJettonFailErBaGang 下注失败
//-------------------------------------------------------------
func (s *Player) onPlaceJettonFailErBaGang(msg interface{}, peer Session) {
	rspdata, ok := msg.(*ErBaGang.CMD_S_PlaceJettonFail)
	if !ok {
		log.Fatalln("onPlaceJettonFailErBaGang ", ok)
	}
	client := peer.GetCtx(TagUserInfo).(*DefWSClient)
	//util.Logx("UserClient", "Player", client.UserID, client.Account, "onPlaceJettonFailErBaGang", rspdata)
	//下注失败///////////////////////////////////
	log.Printf("--- *** PID[%07d] player[%d:%d:%s] :: onPlaceJettonFailErBaGang 下注失败 !!!!!!!!!!!!!! \n%v\n %v\n",
		os.Getpid(),
		client.UserID,
		client.Account,
		client.Token,
		reflect.TypeOf(rspdata).Elem(), util.JSON2Str(rspdata))
	//离开释放资源
	gSemJetton.Leave()
}

//onGameJettonErBaGang 开始下注
//-------------------------------------------------------------
func (s *Player) onGameJettonErBaGang(msg interface{}, peer Session) {
	rspdata, ok := msg.(*ErBaGang.CMD_S_GameJetton)
	if !ok {
		log.Fatalln("onGameJettonErBaGang ", ok)
	}
	client := peer.GetCtx(TagUserInfo).(*DefWSClient)
	util.Logx("UserClient", "Player", client.UserID, client.Account, "onGameJettonErBaGang", rspdata)
	//进入访问资源
	if gSemJetton.Enter() {
		//用户主动下注 [1,3]
		//x := rand.Intn(3) + 1
		//sendPlayerPlaceJetErBaGang(peer, int32(x), 100)
		client.TimerID1 = s.entry.RunAfter(2000, client)
	} else {
		s.entry.GetCell().Append(func() {
			if gSemJetton.Enter() {
				client.TimerID1 = s.entry.RunAfter(2000, client)
			}
		})
	}
}

//
//onSceneGameJettonErBaGang 开始游戏场景
//-------------------------------------------------------------
func (s *Player) onSceneGameJettonErBaGang(msg interface{}, peer Session) {
	rspdata, ok := msg.(*ErBaGang.CMD_S_Scene_GameJetton)
	if !ok {
		log.Fatalln("onSceneGameJettonErBaGang ", ok)
	}
	client := peer.GetCtx(TagUserInfo).(*DefWSClient)
	util.Logx("UserClient", "Player", client.UserID, client.Account, "onSceneGameJettonErBaGang", rspdata)
}

//onQueryPlayerListErBaGang
//-------------------------------------------------------------
func (s *Player) onQueryPlayerListErBaGang(msg interface{}, peer Session) {
	// rspdata, ok := msg.(*ErBaGang.CMD_S_PlayerList)
	// if !ok {
	// 	log.Fatalln(ok)
	// }
	// client := peer.GetCtx(TagUserInfo).(*DefWSClient)
	// util.Logx("UserClient", "Player", client.UserID, client.Account, "onQueryPlayerListErBaGang", rspdata)
}

//onJettonBroadcastErBaGang
//-------------------------------------------------------------
func (s *Player) onJettonBroadcastErBaGang(msg interface{}, peer Session) {
	rspdata, ok := msg.(*ErBaGang.CMD_S_Jetton_Broadcast)
	if !ok {
		log.Fatalln("onJettonBroadcastErBaGang ", ok)
	}
	client := peer.GetCtx(TagUserInfo).(*DefWSClient)
	util.Logx("UserClient", "Player", client.UserID, client.Account, "onJettonBroadcastErBaGang", rspdata)
}

//---------------------------------------------------------------------------------------------
//龙虎斗
//---------------------------------------------------------------------------------------------

//onSyncTimeLonghu 同步TIME
//-------------------------------------------------------------
func (s *Player) onSyncTimeLonghu(msg interface{}, peer Session) {
	rspdata, ok := msg.(*Longhu.CMD_S_SyncTime_Res)
	if !ok {
		log.Fatalln("onSyncTimeLonghu ", ok)
	}
	client := peer.GetCtx(TagUserInfo).(*DefWSClient)
	util.Logx("UserClient", "Player", client.UserID, client.Account, "onSyncTimeLonghu", rspdata)
}

//onSceneStatusFreeLonghu 游戏空闲
//-------------------------------------------------------------
func (s *Player) onSceneStatusFreeLonghu(msg interface{}, peer Session) {
	rspdata, ok := msg.(*Longhu.CMD_Scene_StatusFree)
	if !ok {
		log.Fatalln("onSceneStatusFreeLonghu ", ok)
	}
	client := peer.GetCtx(TagUserInfo).(*DefWSClient)
	util.Logx("UserClient", "Player", client.UserID, client.Account, "onSceneStatusFreeLonghu", rspdata)
}

//onGameStartLonghu 游戏开始
//-------------------------------------------------------------
func (s *Player) onGameStartLonghu(msg interface{}, peer Session) {
	rspdata, ok := msg.(*Longhu.CMD_S_GameStart)
	if !ok {
		log.Fatalln("onGameStartLonghu ", ok)
	}
	client := peer.GetCtx(TagUserInfo).(*DefWSClient)
	util.Logx("UserClient", "Player", client.UserID, client.Account, "onGameStartLonghu", rspdata)
}

//onPlaceJetSuccessLonghu 用户下注
//-------------------------------------------------------------
func (s *Player) onPlaceJetSuccessLonghu(msg interface{}, peer Session) {
	rspdata, ok := msg.(*Longhu.CMD_S_PlaceJetSuccess)
	if !ok {
		log.Fatalln("onPlaceJetSuccessLonghu ", ok)
	}
	client := peer.GetCtx(TagUserInfo).(*DefWSClient)
	util.Logx("UserClient", "Player", client.UserID, client.Account, "onPlaceJetSuccessLonghu", rspdata)
	//下注成功///////////////////////////////////
	// log.Printf("--- *** PID[%07d] player[%d:%d:%s] :: onPlaceJetSuccessLonghu \n%v\n %v\n",
	// 	os.Getpid(),
	// 	client.UserID,
	// 	client.Account,
	// 	client.Token,
	// 	reflect.TypeOf(rspdata).Elem(), util.JSON2Str(rspdata))
	//离开释放资源
	gSemJetton.Leave()
}

//onGameEndLonghu 当局游戏结束
//-------------------------------------------------------------
func (s *Player) onGameEndLonghu(msg interface{}, peer Session) {
	rspdata, ok := msg.(*Longhu.CMD_S_GameEnd)
	if !ok {
		log.Fatalln("onGameEndLonghu ", ok)
	}
	client := peer.GetCtx(TagUserInfo).(*DefWSClient)
	util.Logx("UserClient", "Player", client.UserID, client.Account, "onGameEndLonghu", rspdata)
}

//onGameRecordLonghu 游戏记录
//-------------------------------------------------------------
func (s *Player) onGameRecordLonghu(msg interface{}, peer Session) {
	rspdata, ok := msg.(*Longhu.CMD_S_GameRecord)
	if !ok {
		log.Fatalln("onGameRecordLonghu ", ok)
	}
	client := peer.GetCtx(TagUserInfo).(*DefWSClient)
	util.Logx("UserClient", "Player", client.UserID, client.Account, "onGameRecordLonghu", rspdata)
}

//onPlaceJettonFailLonghu 下注失败
//-------------------------------------------------------------
func (s *Player) onPlaceJettonFailLonghu(msg interface{}, peer Session) {
	rspdata, ok := msg.(*Longhu.CMD_S_PlaceJettonFail)
	if !ok {
		log.Fatalln("onPlaceJettonFailLonghu ", ok)
	}
	client := peer.GetCtx(TagUserInfo).(*DefWSClient)
	//util.Logx("UserClient", "Player", client.UserID, client.Account, "onPlaceJettonFailLonghu", rspdata)
	//下注失败///////////////////////////////////
	log.Printf("--- *** PID[%07d] player[%d:%d:%s] :: onPlaceJettonFailLonghu 下注失败 !!!!!!!!!!!!!! \n%v\n %v\n",
		os.Getpid(),
		client.UserID,
		client.Account,
		client.Token,
		reflect.TypeOf(rspdata).Elem(), util.JSON2Str(rspdata))
	//离开释放资源
	gSemJetton.Leave()
}

//onQueryPlayerListLonghu 玩家在线列表返回
//-------------------------------------------------------------
func (s *Player) onQueryPlayerListLonghu(msg interface{}, peer Session) {
	rspdata, ok := msg.(*Longhu.CMD_S_PlayerList)
	if !ok {
		log.Fatalln("onQueryPlayerListLonghu ", ok)
	}
	client := peer.GetCtx(TagUserInfo).(*DefWSClient)
	util.Logx("UserClient", "Player", client.UserID, client.Account, "onQueryPlayerListLonghu", rspdata)
}

//onStartPlaceJettonLonghu 开始下注
//-------------------------------------------------------------
func (s *Player) onStartPlaceJettonLonghu(msg interface{}, peer Session) {
	rspdata, ok := msg.(*Longhu.CMD_S_StartPlaceJetton)
	if !ok {
		log.Fatalln(ok)
	}
	client := peer.GetCtx(TagUserInfo).(*DefWSClient)
	util.Logx("UserClient", "Player", client.UserID, client.Account, "onStartPlaceJettonLonghu", rspdata)
	//进入访问资源
	if gSemJetton.Enter() {
		//用户主动下注 [1,3]
		//x := rand.Intn(3) + 1
		//sendPlayerPlaceJetLonghu(peer, int32(x), 100)
		client.TimerID1 = s.entry.RunAfter(2000, client)
	} else {
		s.entry.GetCell().Append(func() {
			if gSemJetton.Enter() {
				client.TimerID1 = s.entry.RunAfter(2000, client)
			}
		})
	}
}

//onJettonBroadcastLonghu
//-------------------------------------------------------------
func (s *Player) onJettonBroadcastLonghu(msg interface{}, peer Session) {
	rspdata, ok := msg.(*Longhu.CMD_S_Jetton_Broadcast)
	if !ok {
		log.Fatalln("onJettonBroadcastLonghu ", ok)
	}
	client := peer.GetCtx(TagUserInfo).(*DefWSClient)
	util.Logx("UserClient", "Player", client.UserID, client.Account, "onJettonBroadcastLonghu", rspdata)
}

//---------------------------------------------------------------------------------------------
//百人牛牛
//---------------------------------------------------------------------------------------------

//onSyncTimeBrnn 服务端返回 - 同步TIME
//-------------------------------------------------------------
func (s *Player) onSyncTimeBrnn(msg interface{}, peer Session) {
	rspdata, ok := msg.(*Brnn.CMD_S_SyncTime_Res)
	if !ok {
		log.Fatalln("onSyncTimeBrnn ", ok)
	}
	client := peer.GetCtx(TagUserInfo).(*DefWSClient)
	util.Logx("UserClient", "Player", client.UserID, client.Account, "onSyncTimeBrnn", rspdata)
}

//onSceneStatusFreeBrnn 服务端返回 - 游戏空闲
//-------------------------------------------------------------
func (s *Player) onSceneStatusFreeBrnn(msg interface{}, peer Session) {
	rspdata, ok := msg.(*Brnn.CMD_Scene_StatusFree)
	if !ok {
		log.Fatalln("onSceneStatusFreeBrnn ", ok)
	}
	client := peer.GetCtx(TagUserInfo).(*DefWSClient)
	util.Logx("UserClient", "Player", client.UserID, client.Account, "onSceneStatusFreeBrnn", rspdata)
}

//onGameStartBrnn 服务端返回 - 游戏开始
//-------------------------------------------------------------
func (s *Player) onGameStartBrnn(msg interface{}, peer Session) {
	rspdata, ok := msg.(*Brnn.CMD_S_GameStart)
	if !ok {
		log.Fatalln("onGameStartBrnn ", ok)
	}
	client := peer.GetCtx(TagUserInfo).(*DefWSClient)
	util.Logx("UserClient", "Player", client.UserID, client.Account, "onGameStartBrnn", rspdata)
}

//onPlaceJetSuccessBrnn 服务端返回 - 用户下注
//-------------------------------------------------------------
func (s *Player) onPlaceJetSuccessBrnn(msg interface{}, peer Session) {
	rspdata, ok := msg.(*Brnn.CMD_S_PlaceJetSuccess)
	if !ok {
		log.Fatalln("onPlaceJetSuccessBrnn ", ok)
	}
	client := peer.GetCtx(TagUserInfo).(*DefWSClient)
	util.Logx("UserClient", "Player", client.UserID, client.Account, "onPlaceJetSuccessBrnn", rspdata)
	//下注成功///////////////////////////////////
	// log.Printf("--- *** PID[%07d] player[%d:%d:%s] onPlaceJetSuccessBrnn 下注成功 !!!!!!!!!!!!!! \n",
	// 	os.Getpid(),
	// 	client.UserID,
	// 	client.Account,
	// 	client.Token)
	//离开释放资源
	gSemJetton.Leave()
}

//onGameEndBrnn 服务端返回 - 当局游戏结束
//-------------------------------------------------------------
func (s *Player) onGameEndBrnn(msg interface{}, peer Session) {
	rspdata, ok := msg.(*Brnn.CMD_S_GameEnd)
	if !ok {
		log.Fatalln("onGameEndBrnn ", ok)
	}
	client := peer.GetCtx(TagUserInfo).(*DefWSClient)
	util.Logx("UserClient", "Player", client.UserID, client.Account, "onGameEndBrnn", rspdata)
}

//onGameRecordBrnn 服务端返回 - 游戏记录
//-------------------------------------------------------------
func (s *Player) onGameRecordBrnn(msg interface{}, peer Session) {
	rspdata, ok := msg.(*Brnn.CMD_S_GameRecord)
	if !ok {
		log.Fatalln("onGameRecordBrnn ", ok)
	}
	client := peer.GetCtx(TagUserInfo).(*DefWSClient)
	util.Logx("UserClient", "Player", client.UserID, client.Account, "onGameRecordBrnn", rspdata)
}

//onPlaceJettonFailBrnn 服务端返回 - 下注失败
//-------------------------------------------------------------
func (s *Player) onPlaceJettonFailBrnn(msg interface{}, peer Session) {
	rspdata, ok := msg.(*Brnn.CMD_S_PlaceJettonFail)
	if !ok {
		log.Fatalln("onPlaceJettonFailBrnn ", ok)
	}
	client := peer.GetCtx(TagUserInfo).(*DefWSClient)
	//util.Logx("UserClient", "Player", client.UserID, client.Account, "onPlaceJettonFailBrnn", rspdata)
	//下注失败///////////////////////////////////
	log.Printf("--- *** PID[%07d] player[%d:%d:%s] :: onPlaceJettonFailBrnn 下注失败 !!!!!!!!!!!!!! \n%v\n %v\n",
		os.Getpid(),
		client.UserID,
		client.Account,
		client.Token,
		reflect.TypeOf(rspdata).Elem(), util.JSON2Str(rspdata))
	//离开释放资源
	gSemJetton.Leave()
}

//onPlayerListBrnn 服务端返回 - 玩家在线列表返回
//-------------------------------------------------------------
func (s *Player) onPlayerListBrnn(msg interface{}, peer Session) {
	rspdata, ok := msg.(*Brnn.CMD_S_PlayerList)
	if !ok {
		log.Fatalln("onPlayerListBrnn ", ok)
	}
	client := peer.GetCtx(TagUserInfo).(*DefWSClient)
	util.Logx("UserClient", "Player", client.UserID, client.Account, "onPlayerListBrnn", rspdata)
}

//onStartJettonBrnn 服务端返回 - 开始下注
//-------------------------------------------------------------
func (s *Player) onStartJettonBrnn(msg interface{}, peer Session) {
	rspdata, ok := msg.(*Brnn.CMD_S_StartPlaceJetton)
	if !ok {
		log.Fatalln("onStartJettonBrnn ", ok)
	}
	client := peer.GetCtx(TagUserInfo).(*DefWSClient)
	util.Logx("UserClient", "Player", client.UserID, client.Account, "onStartJettonBrnn", rspdata)
	//进入访问资源
	if gSemJetton.Enter() {
		//用户主动下注 [0,5]
		//x := rand.Intn(6)
		//sendPlayerPlaceJetBrnn(peer, int32(x), 100)
		client.TimerID1 = s.entry.RunAfter(2000, client)
	} else {
		s.entry.GetCell().Append(func() {
			if gSemJetton.Enter() {
				client.TimerID1 = s.entry.RunAfter(2000, client)
			}
		})
	}
}

//onJettonBroadcastBrnn 服务端返回
//-------------------------------------------------------------
func (s *Player) onJettonBroadcastBrnn(msg interface{}, peer Session) {
	rspdata, ok := msg.(*Brnn.CMD_S_Jetton_Broadcast)
	if !ok {
		log.Fatalln("onJettonBroadcastBrnn ", ok)
	}
	client := peer.GetCtx(TagUserInfo).(*DefWSClient)
	util.Logx("UserClient", "Player", client.UserID, client.Account, "onJettonBroadcastBrnn", rspdata)
}

//---------------------------------------------------------------------------------------------
//红黑大战
//---------------------------------------------------------------------------------------------

//onSceneStatusFreeHongHei 服务端返回 - 游戏空闲
//-------------------------------------------------------------
func (s *Player) onSceneStatusFreeHongHei(msg interface{}, peer Session) {
	rspdata, ok := msg.(*HongHei.CMD_Scene_StatusFree)
	if !ok {
		log.Fatalln("onSceneStatusFreeHongHei ", ok)
	}
	client := peer.GetCtx(TagUserInfo).(*DefWSClient)
	util.Logx("UserClient", "Player", client.UserID, client.Account, "onSceneStatusFreeHongHei", rspdata)
}

//onGameStartHongHei 服务端返回 - 游戏开始
//-------------------------------------------------------------
func (s *Player) onGameStartHongHei(msg interface{}, peer Session) {
	rspdata, ok := msg.(*HongHei.CMD_S_GameStart)
	if !ok {
		log.Fatalln("onGameStartHongHei ", ok)
	}
	client := peer.GetCtx(TagUserInfo).(*DefWSClient)
	util.Logx("UserClient", "Player", client.UserID, client.Account, "onGameStartHongHei", rspdata)
}

//onPlaceJetSuccessHongHei 服务端返回 - 用户下注
//-------------------------------------------------------------
func (s *Player) onPlaceJetSuccessHongHei(msg interface{}, peer Session) {
	rspdata, ok := msg.(*HongHei.CMD_S_PlaceJetSuccess)
	if !ok {
		log.Fatalln("onPlaceJetSuccessHongHei ", ok)
	}
	client := peer.GetCtx(TagUserInfo).(*DefWSClient)
	util.Logx("UserClient", "Player", client.UserID, client.Account, "onPlaceJetSuccessHongHei", rspdata)
	//下注成功///////////////////////////////////
	// log.Printf("--- *** PID[%07d] player[%d:%d:%s] onPlaceJetSuccessHongHei 下注成功 !!!!!!!!!!!!!! \n",
	// 	os.Getpid(),
	// 	client.UserID,
	// 	client.Account,
	// 	client.Token)
	//离开释放资源
	gSemJetton.Leave()
}

//onGameEndHongHei 服务端返回 - 当局游戏结束
//-------------------------------------------------------------
func (s *Player) onGameEndHongHei(msg interface{}, peer Session) {
	rspdata, ok := msg.(*HongHei.CMD_S_GameEnd)
	if !ok {
		log.Fatalln("onGameEndHongHei ", ok)
	}
	client := peer.GetCtx(TagUserInfo).(*DefWSClient)
	util.Logx("UserClient", "Player", client.UserID, client.Account, "onGameEndHongHei", rspdata)
}

//onGameRecordHongHei 服务端返回 - 游戏记录
//-------------------------------------------------------------
func (s *Player) onGameRecordHongHei(msg interface{}, peer Session) {
	rspdata, ok := msg.(*HongHei.CMD_S_GameRecord)
	if !ok {
		log.Fatalln("onGameRecordHongHei ", ok)
	}
	client := peer.GetCtx(TagUserInfo).(*DefWSClient)
	util.Logx("UserClient", "Player", client.UserID, client.Account, "onGameRecordHongHei", rspdata)
}

//onPlaceJettonFailHongHei 服务端返回 - 下注失败
//-------------------------------------------------------------
func (s *Player) onPlaceJettonFailHongHei(msg interface{}, peer Session) {
	rspdata, ok := msg.(*HongHei.CMD_S_PlaceJettonFail)
	if !ok {
		log.Fatalln("onPlaceJettonFailHongHei ", ok)
	}
	client := peer.GetCtx(TagUserInfo).(*DefWSClient)
	//util.Logx("UserClient", "Player", client.UserID, client.Account, "onPlaceJettonFailHongHei", rspdata)
	//下注失败///////////////////////////////////
	log.Printf("--- *** PID[%07d] player[%d:%d:%s] :: onPlaceJettonFailHongHei 下注失败 !!!!!!!!!!!!!! \n%v\n %v\n",
		os.Getpid(),
		client.UserID,
		client.Account,
		client.Token,
		reflect.TypeOf(rspdata).Elem(), util.JSON2Str(rspdata))
	//离开释放资源
	gSemJetton.Leave()
}

//onPlayerListHongHei 服务端返回 - 玩家在线列表返回
//-------------------------------------------------------------
func (s *Player) onPlayerListHongHei(msg interface{}, peer Session) {
	rspdata, ok := msg.(*HongHei.CMD_S_PlayerList)
	if !ok {
		log.Fatalln("onPlayerListHongHei ", ok)
	}
	client := peer.GetCtx(TagUserInfo).(*DefWSClient)
	util.Logx("UserClient", "Player", client.UserID, client.Account, "onPlayerListHongHei", rspdata)
}

//onStartJettonHongHei 服务端返回 - 开始下注
//-------------------------------------------------------------
func (s *Player) onStartJettonHongHei(msg interface{}, peer Session) {
	rspdata, ok := msg.(*HongHei.CMD_S_StartPlaceJetton)
	if !ok {
		log.Fatalln("onStartJettonHongHei ", ok)
	}
	client := peer.GetCtx(TagUserInfo).(*DefWSClient)
	util.Logx("UserClient", "Player", client.UserID, client.Account, "onStartJettonHongHei", rspdata)
	//进入访问资源
	if gSemJetton.Enter() {
		//用户主动下注 [0,2]
		//x := rand.Intn(3)
		//sendPlayerPlaceJetHongHei(peer, int32(x), 100)
		client.TimerID1 = s.entry.RunAfter(2000, client)
	} else {
		s.entry.GetCell().Append(func() {
			if gSemJetton.Enter() {
				client.TimerID1 = s.entry.RunAfter(2000, client)
			}
		})
	}
}

//onJettonBroadcastHongHei 服务端返回
//-------------------------------------------------------------
func (s *Player) onJettonBroadcastHongHei(msg interface{}, peer Session) {
	rspdata, ok := msg.(*HongHei.CMD_S_Jetton_Broadcast)
	if !ok {
		log.Fatalln("onJettonBroadcastHongHei ", ok)
	}
	client := peer.GetCtx(TagUserInfo).(*DefWSClient)
	util.Logx("UserClient", "Player", client.UserID, client.Account, "onJettonBroadcastHongHei", rspdata)
}

//---------------------------------------------------------------------------------------------
//奔驰宝马
//---------------------------------------------------------------------------------------------

//開始游戲
func (s *Player) onGameStartBcbm(msg interface{}, peer Session) {
	rspdata, ok := msg.(*BenCiBaoMa.CMD_S_GameStart)
	if !ok {
		log.Fatalln("onGameStartBcbm ", ok)
	}
	client := peer.GetCtx(TagUserInfo).(*DefWSClient)
	util.Logx("UserClient", "Player", client.UserID, client.Account, "onGameStartBcbm", rspdata)
	s.entry.RemoveTimer(client.TimerID1)
	//进入访问资源
	//if gSemJetton.Enter() {
	client.TimerID1 = s.entry.RunAfter(2000, client)
	//} else {
	//	s.entry.GetCell().Append(func() {
	//		if gSemJetton.Enter() {
	//			client.TimerID1 = s.entry.RunAfter(500, client)
	//		}
	//	})
	//}
}

//游戲結束
func (s *Player) onGameEndBcbm(msg interface{}, peer Session) {
	rspdata, ok := msg.(*BenCiBaoMa.CMD_S_GameEnd)
	if !ok {
		log.Fatalln("onGameEndBcbm ", ok)
	}
	client := peer.GetCtx(TagUserInfo).(*DefWSClient)
	util.Logx("UserClient", "Player", client.UserID, client.Account, "onGameEndBcbm", rspdata)
	//取消定时器
	s.entry.RemoveTimer(client.TimerID1)
}

//開始游戲場景
func (s *Player) onSceneGameStartBcbm(msg interface{}, peer Session) {
	rspdata, ok := msg.(*BenCiBaoMa.CMD_S_Scene_GameStart)
	if !ok {
		log.Fatalln("onSceneGameStartBcbm ", ok)
	}
	client := peer.GetCtx(TagUserInfo).(*DefWSClient)
	util.Logx("UserClient", "Player", client.UserID, client.Account, "onSceneGameStartBcbm", rspdata)
}

//結束游戲場景
func (s *Player) onSceneGameEndBcbm(msg interface{}, peer Session) {
	rspdata, ok := msg.(*BenCiBaoMa.CMD_S_Scene_GameEnd)
	if !ok {
		log.Fatalln("onSceneGameEndBcbm ", ok)
	}
	client := peer.GetCtx(TagUserInfo).(*DefWSClient)
	util.Logx("UserClient", "Player", client.UserID, client.Account, "onSceneGameEndBcbm", rspdata)
}

//游戏结束时返回玩家结算信息
func (s *Player) onUserWinListBcbm(msg interface{}, peer Session) {
	rspdata, ok := msg.(*BenCiBaoMa.CMD_S_UserWinList)
	if !ok {
		log.Fatalln("onUserWinListBcbm ", ok)
	}
	client := peer.GetCtx(TagUserInfo).(*DefWSClient)
	util.Logx("UserClient", "Player", client.UserID, client.Account, "onUserWinListBcbm", rspdata)
}

//下注成功返回
func (s *Player) onPlaceJetSuccessBcbm(msg interface{}, peer Session) {
	rspdata, ok := msg.(*BenCiBaoMa.CMD_S_PlaceJetSuccess)
	if !ok {
		log.Fatalln("onPlaceJetSuccessBcbm ", ok)
	}
	client := peer.GetCtx(TagUserInfo).(*DefWSClient)
	util.Logx("UserClient", "Player", client.UserID, client.Account, "onPlaceJetSuccessBcbm", rspdata)
}

//下注失败返回
func (s *Player) onPlaceJettonFailBcbm(msg interface{}, peer Session) {
	rspdata, ok := msg.(*BenCiBaoMa.CMD_S_PlaceJettonFail)
	if !ok {
		log.Fatalln("onPlaceJettonFailBcbm ", ok)
	}
	client := peer.GetCtx(TagUserInfo).(*DefWSClient)
	util.Logx("UserClient", "Player", client.UserID, client.Account, "onPlaceJettonFailBcbm", rspdata)
}

//获取玩家在线列表
func (s *Player) onPlayerListBcbm(msg interface{}, peer Session) {
	rspdata, ok := msg.(*BenCiBaoMa.CMD_S_PlayerList)
	if !ok {
		log.Fatalln("onPlayerListBcbm ", ok)
	}
	client := peer.GetCtx(TagUserInfo).(*DefWSClient)
	util.Logx("UserClient", "Player", client.UserID, client.Account, "onPlayerListBcbm", rspdata)
}

//服务器返回当前状态
func (s *Player) onUserCurStateBcbm(msg interface{}, peer Session) {
	rspdata, ok := msg.(*BenCiBaoMa.CMD_S_CurState)
	if !ok {
		log.Fatalln("onUserCurStateBcbm ", ok)
	}
	client := peer.GetCtx(TagUserInfo).(*DefWSClient)
	util.Logx("UserClient", "Player", client.UserID, client.Account, "onUserCurStateBcbm", rspdata)
}

//续押下注成功返回
func (s *Player) onPlaceRepeatJetSuccessBcbm(msg interface{}, peer Session) {
	rspdata, ok := msg.(*BenCiBaoMa.CMD_S_RepeatJetSuccess)
	if !ok {
		log.Fatalln("onPlaceRepeatJetSuccessBcbm ", ok)
	}
	client := peer.GetCtx(TagUserInfo).(*DefWSClient)
	util.Logx("UserClient", "Player", client.UserID, client.Account, "onPlaceRepeatJetSuccessBcbm", rspdata)
}

//广播一段时间内其他所有玩家的下注
func (s *Player) onOtherPlaceJetSuccessBcbm(msg interface{}, peer Session) {
	rspdata, ok := msg.(*BenCiBaoMa.CMD_S_OtherPlaceJetSuccess)
	if !ok {
		log.Fatalln("onOtherPlaceJetSuccessBcbm ", ok)
	}
	client := peer.GetCtx(TagUserInfo).(*DefWSClient)
	util.Logx("UserClient", "Player", client.UserID, client.Account, "onOtherPlaceJetSuccessBcbm", rspdata)
}

//
func (s *Player) registerModuleHandler(handlers CmdCallbacks) {
	//登陆返回
	handlers[uint32(ENWORD(2, 4))] = s.resultPlayerLogin
	//心跳返回
	handlers[uint32(ENWORD(2, 2))] = s.resultKeepAlive
	//游戏信息返回
	handlers[uint32(ENWORD(2, 6))] = s.resultGameListInfo
	//获取游戏IP
	handlers[uint32(ENWORD(2, 8))] = s.resultGameserverInfo
	//进入房间
	handlers[uint32(ENWORD(3, 4))] = s.resultPlayerEnterRoom
	//玩家进入返回
	handlers[uint32(ENWORD(3, 5))] = s.onPlayerEnterNotify
	//玩家积分信息
	handlers[uint32(ENWORD(3, 6))] = s.onPlayerScoreNotify
	//玩家状态
	handlers[uint32(ENWORD(3, 7))] = s.onPlayerStatusNotify
	//玩家就绪
	handlers[uint32(ENWORD(3, 29))] = s.resultPlayerReady
	//离开响应
	handlers[uint32(ENWORD(3, 10))] = s.resultPlayerLeave
	//上下分通知
	handlers[uint32(ENWORD(2, 66))] = s.resultOrderNotify
	handlers[uint32(ENWORD(1, 11))] = s.resultProxyOrderNotify
	switch int32(*subGameID) {
	case GGames.ByName["龙虎斗"].ID:
		{
			//同步TIME
			handlers[uint32(ENWORD(4, 105))] = s.onSyncTimeLonghu
			//游戏空闲
			handlers[uint32(ENWORD(4, 120))] = s.onSceneStatusFreeLonghu
			//游戏开始
			handlers[uint32(ENWORD(4, 121))] = s.onGameStartLonghu
			//用户下注
			handlers[uint32(ENWORD(4, 122))] = s.onPlaceJetSuccessLonghu
			//当局游戏结束
			handlers[uint32(ENWORD(4, 123))] = s.onGameEndLonghu
			//游戏记录
			handlers[uint32(ENWORD(4, 127))] = s.onGameRecordLonghu
			//下注失败
			handlers[uint32(ENWORD(4, 128))] = s.onPlaceJettonFailLonghu
			//玩家在线列表返回
			handlers[uint32(ENWORD(4, 130))] = s.onQueryPlayerListLonghu
			//开始下注
			handlers[uint32(ENWORD(4, 139))] = s.onStartPlaceJettonLonghu
			//
			handlers[uint32(ENWORD(4, 114))] = s.onJettonBroadcastLonghu
		}
	case GGames.ByName["二八杠"].ID:
		{
			//开始游戏
			handlers[uint32(ENWORD(4, 100))] = s.onGameStartErBaGang
			//结束游戏
			handlers[uint32(ENWORD(4, 101))] = s.onGameEndErBaGang
			//开始游戏场景
			handlers[uint32(ENWORD(4, 102))] = s.onSceneGameStartErBaGang
			//结束游戏场景
			handlers[uint32(ENWORD(4, 103))] = s.onSceneGameEndErBaGang
			//玩家列表
			handlers[uint32(ENWORD(4, 104))] = s.onPlayerListErBaGang
			//下注成功
			handlers[uint32(ENWORD(4, 105))] = s.onPlaceJetSuccessErBaGang
			//下注失败
			handlers[uint32(ENWORD(4, 106))] = s.onPlaceJettonFailErBaGang
			//开始下注
			handlers[uint32(ENWORD(4, 112))] = s.onGameJettonErBaGang
			//开始游戏场景
			handlers[uint32(ENWORD(4, 113))] = s.onSceneGameJettonErBaGang
			//
			handlers[uint32(ENWORD(4, 111))] = s.onQueryPlayerListErBaGang
			//
			handlers[uint32(ENWORD(4, 114))] = s.onJettonBroadcastErBaGang
		}
	case GGames.ByName["百人牛牛"].ID:
		{
			//同步TIME
			handlers[uint32(ENWORD(4, 105))] = s.onSyncTimeBrnn
			//游戏空闲
			handlers[uint32(ENWORD(4, 120))] = s.onSceneStatusFreeBrnn
			//游戏开始
			handlers[uint32(ENWORD(4, 121))] = s.onGameStartBrnn
			//用户下注
			handlers[uint32(ENWORD(4, 122))] = s.onPlaceJetSuccessBrnn
			//当局游戏结束
			handlers[uint32(ENWORD(4, 123))] = s.onGameEndBrnn
			//游戏记录
			handlers[uint32(ENWORD(4, 127))] = s.onGameRecordBrnn
			//下注失败
			handlers[uint32(ENWORD(4, 128))] = s.onPlaceJettonFailBrnn
			//玩家在线列表返回
			handlers[uint32(ENWORD(4, 130))] = s.onPlayerListBrnn
			//开始下注
			handlers[uint32(ENWORD(4, 139))] = s.onStartJettonBrnn
			//
			handlers[uint32(ENWORD(4, 140))] = s.onJettonBroadcastBrnn
		}
	case GGames.ByName["红黑大战"].ID:
		{
			//游戏空闲
			handlers[uint32(ENWORD(4, 120))] = s.onSceneStatusFreeHongHei
			//游戏开始
			handlers[uint32(ENWORD(4, 121))] = s.onGameStartHongHei
			//用户下注
			handlers[uint32(ENWORD(4, 122))] = s.onPlaceJetSuccessHongHei
			//当局游戏结束
			handlers[uint32(ENWORD(4, 123))] = s.onGameEndHongHei
			//游戏记录
			handlers[uint32(ENWORD(4, 127))] = s.onGameRecordHongHei
			//下注失败
			handlers[uint32(ENWORD(4, 128))] = s.onPlaceJettonFailHongHei
			//玩家在线列表返回
			handlers[uint32(ENWORD(4, 130))] = s.onPlayerListHongHei
			//开始下注
			handlers[uint32(ENWORD(4, 139))] = s.onStartJettonHongHei
			//
			handlers[uint32(ENWORD(4, 140))] = s.onJettonBroadcastHongHei
		}
	case GGames.ByName["奔驰宝马"].ID:
		{
			//開始游戲
			handlers[uint32(ENWORD(4, 100))] = s.onGameStartBcbm
			//游戲結束
			handlers[uint32(ENWORD(4, 101))] = s.onGameEndBcbm
			//開始游戲場景
			handlers[uint32(ENWORD(4, 102))] = s.onSceneGameStartBcbm
			//結束游戲場景
			handlers[uint32(ENWORD(4, 103))] = s.onSceneGameEndBcbm
			//游戏结束时返回玩家结算信息
			handlers[uint32(ENWORD(4, 104))] = s.onUserWinListBcbm
			//下注成功返回
			handlers[uint32(ENWORD(4, 105))] = s.onPlaceJetSuccessBcbm
			//下注失败返回
			handlers[uint32(ENWORD(4, 106))] = s.onPlaceJettonFailBcbm
			//获取玩家在线列表
			handlers[uint32(ENWORD(4, 111))] = s.onPlayerListBcbm
			//服务器返回当前状态
			handlers[uint32(ENWORD(4, 113))] = s.onUserCurStateBcbm
			//续押下注成功返回
			handlers[uint32(ENWORD(4, 114))] = s.onPlaceRepeatJetSuccessBcbm
			//广播一段时间内其他所有玩家的下注
			handlers[uint32(ENWORD(4, 115))] = s.onOtherPlaceJetSuccessBcbm
		}
	}
}
