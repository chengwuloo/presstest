package main

//
// Created by andy_ro@qq.com
// 			4/9/2019
//

import (
	"reflect"
	"server/pb/BenCiBaoMa"
	"server/pb/Brnn"
	"server/pb/ErBaGang"
	"server/pb/GameServer"
	"server/pb/Game_Common"
	"server/pb/HallServer"
	"server/pb/HongHei"
	"server/pb/Longhu"

	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/codec"
	_ "github.com/davyxu/cellnet/codec/gogopb"
)

//
func ENWORD(mainID, subID int) int {
	return ((0xFF & mainID) << 8) | (0xFF & subID)
}

//
func DEWORD(cmd int) (mainID, subID int) {
	mainID = (0xFF & (cmd >> 8))
	subID = (0xFF & cmd)
	return
}

//
const (
	SubCmdID = iota + ((0xFF & 20) << 8) | (0xFF & 19)
)

//init 由系统在main之前调用，改由main解析命令行之后手动调用register
func init() {
	//register()
}

//register 注册消息，main解析命令行之后手动调用
func register() {
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: codec.MustGetCodec("gogopb"),
		Type:  reflect.TypeOf((*HallServer.GetPlayingGameInfoMessage)(nil)).Elem(),
		ID:    ENWORD(1, 9),
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: codec.MustGetCodec("gogopb"),
		Type:  reflect.TypeOf((*HallServer.GetPlayingGameInfoMessageResponse)(nil)).Elem(),
		ID:    ENWORD(1, 10),
	})
	//心跳 - 请求
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: codec.MustGetCodec("gogopb"),
		Type:  reflect.TypeOf((*Game_Common.KeepAliveMessage)(nil)).Elem(),
		ID:    ENWORD(2, 1),
	})
	//应答
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: codec.MustGetCodec("gogopb"),
		Type:  reflect.TypeOf((*Game_Common.KeepAliveMessageResponse)(nil)).Elem(),
		ID:    ENWORD(2, 2),
	})
	//登陆 - 请求
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: codec.MustGetCodec("gogopb"),
		Type:  reflect.TypeOf((*HallServer.LoginMessage)(nil)).Elem(), 
		ID:    ENWORD(2, 3),
	})
	//应答
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: codec.MustGetCodec("gogopb"),
		Type:  reflect.TypeOf((*HallServer.LoginMessageResponse)(nil)).Elem(),
		ID:    ENWORD(2, 4),
	})
	//游戏信息 - 请求
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: codec.MustGetCodec("gogopb"),
		Type:  reflect.TypeOf((*HallServer.GetGameMessage)(nil)).Elem(),
		ID:    ENWORD(2, 5),
	})
	//应答
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: codec.MustGetCodec("gogopb"),
		Type:  reflect.TypeOf((*HallServer.GetGameMessageResponse)(nil)).Elem(),
		ID:    ENWORD(2, 6),
	})
	//游戏IP - 请求
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: codec.MustGetCodec("gogopb"),
		Type:  reflect.TypeOf((*HallServer.GetGameServerMessage)(nil)).Elem(),
		ID:    ENWORD(2, 7),
	})
	//应答
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: codec.MustGetCodec("gogopb"),
		Type:  reflect.TypeOf((*HallServer.GetGameServerMessageResponse)(nil)).Elem(),
		ID:    ENWORD(2, 8),
	})
	//上下分通知
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: codec.MustGetCodec("gogopb"),
		Type:  reflect.TypeOf((*HallServer.OrderNotifyMessage)(nil)).Elem(),
		ID:    ENWORD(2, 66),
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: codec.MustGetCodec("gogopb"),
		Type:  reflect.TypeOf((*Game_Common.ProxyNotifyOrderScoreMessage)(nil)).Elem(),
		ID:    ENWORD(1, 11),
	})
	//进入房间 - 请求
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: codec.MustGetCodec("gogopb"),
		Type:  reflect.TypeOf((*GameServer.MSG_C2S_UserEnterMessage)(nil)).Elem(),
		ID:    ENWORD(3, 3),
	})
	//应答
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: codec.MustGetCodec("gogopb"),
		Type:  reflect.TypeOf((*GameServer.MSG_S2C_UserEnterMessageResponse)(nil)).Elem(),
		ID:    ENWORD(3, 4),
	})
	//玩家进入返回
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: codec.MustGetCodec("gogopb"),
		Type:  reflect.TypeOf((*GameServer.MSG_S2C_UserBaseInfo)(nil)).Elem(),
		ID:    ENWORD(3, 5),
	})
	//玩家积分信息
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: codec.MustGetCodec("gogopb"),
		Type:  reflect.TypeOf((*GameServer.MSG_S2C_UserScoreInfo)(nil)).Elem(),
		ID:    ENWORD(3, 6),
	})
	//玩家状态
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: codec.MustGetCodec("gogopb"),
		Type:  reflect.TypeOf((*GameServer.MSG_S2C_GameUserStatus)(nil)).Elem(),
		ID:    ENWORD(3, 7),
	})
	//玩家就绪 - 请求
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: codec.MustGetCodec("gogopb"),
		Type:  reflect.TypeOf((*GameServer.MSG_C2S_UserReadyMessage)(nil)).Elem(),
		ID:    ENWORD(3, 8),
	})
	//应答
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: codec.MustGetCodec("gogopb"),
		Type:  reflect.TypeOf((*GameServer.MSG_S2C_UserReadyMessageResponse)(nil)).Elem(),
		ID:    ENWORD(3, 29),
	})
	//玩家离开 - 请求
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: codec.MustGetCodec("gogopb"),
		Type:  reflect.TypeOf((*GameServer.MSG_C2S_UserLeftMessage)(nil)).Elem(),
		ID:    ENWORD(3, 9),
	})
	//应答
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: codec.MustGetCodec("gogopb"),
		Type:  reflect.TypeOf((*GameServer.MSG_C2S_UserLeftMessageResponse)(nil)).Elem(),
		ID:    ENWORD(3, 10),
	})
	//游戏服通道消息
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: codec.MustGetCodec("gogopb"),
		Type:  reflect.TypeOf((*GameServer.MSG_CSC_Passageway)(nil)).Elem(),
		ID:    SubCmdID,
	})
	switch int32(*subGameID) {
	case GGames.ByName["抢庄牛牛"].ID:
		{

		}
	case GGames.ByName["扎金花"].ID:
		{

		}
	case GGames.ByName["21点"].ID:
		{

		}
	case GGames.ByName["三公"].ID:
		{

		}
	case GGames.ByName["二八杠"].ID:
		{
			//开始游戏
			cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
				Codec: codec.MustGetCodec("gogopb"),
				Type:  reflect.TypeOf((*ErBaGang.CMD_S_GameStart)(nil)).Elem(),
				ID:    ENWORD(4, 100),
			})
			//结束游戏
			cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
				Codec: codec.MustGetCodec("gogopb"),
				Type:  reflect.TypeOf((*ErBaGang.CMD_S_GameEnd)(nil)).Elem(),
				ID:    ENWORD(4, 101),
			})
			//开始游戏场景
			cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
				Codec: codec.MustGetCodec("gogopb"),
				Type:  reflect.TypeOf((*ErBaGang.CMD_S_Scene_GameStart)(nil)).Elem(),
				ID:    ENWORD(4, 102),
			})
			//结束游戏场景
			cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
				Codec: codec.MustGetCodec("gogopb"),
				Type:  reflect.TypeOf((*ErBaGang.CMD_S_Scene_GameEnd)(nil)).Elem(),
				ID:    ENWORD(4, 103),
			})
			//玩家列表
			// cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
			// 	Codec: codec.MustGetCodec("gogopb"),
			// 	Type:  reflect.TypeOf((*ErBaGang.CMD_S_PlayerList)(nil)).Elem(),
			// 	ID:    ENWORD(4, 104),
			// })
			//下注成功
			cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
				Codec: codec.MustGetCodec("gogopb"),
				Type:  reflect.TypeOf((*ErBaGang.CMD_S_PlaceJetSuccess)(nil)).Elem(),
				ID:    ENWORD(4, 105),
			})
			//下注失败
			cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
				Codec: codec.MustGetCodec("gogopb"),
				Type:  reflect.TypeOf((*ErBaGang.CMD_S_PlaceJettonFail)(nil)).Elem(),
				ID:    ENWORD(4, 106),
			})
			//开始下注
			cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
				Codec: codec.MustGetCodec("gogopb"),
				Type:  reflect.TypeOf((*ErBaGang.CMD_S_GameJetton)(nil)).Elem(),
				ID:    ENWORD(4, 112),
			})
			//开始游戏场景
			cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
				Codec: codec.MustGetCodec("gogopb"),
				Type:  reflect.TypeOf((*ErBaGang.CMD_S_Scene_GameJetton)(nil)).Elem(),
				ID:    ENWORD(4, 113),
			})
			//玩家下注
			cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
				Codec: codec.MustGetCodec("gogopb"),
				Type:  reflect.TypeOf((*ErBaGang.CMD_C_PlaceJet)(nil)).Elem(),
				ID:    ENWORD(4, 107),
			})
			//玩家申请列表
			cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
				Codec: codec.MustGetCodec("gogopb"),
				Type:  reflect.TypeOf((*ErBaGang.CMD_C_AskList)(nil)).Elem(),
				ID:    ENWORD(4, 108),
			})
			//
			cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
				Codec: codec.MustGetCodec("gogopb"),
				Type:  reflect.TypeOf((*ErBaGang.CMD_C_ReJetton)(nil)).Elem(),
				ID:    ENWORD(4, 109),
			})
			//
			cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
				Codec: codec.MustGetCodec("gogopb"),
				Type:  reflect.TypeOf((*ErBaGang.CMD_C_PlayerList)(nil)).Elem(),
				ID:    ENWORD(4, 110),
			})
			//
			// cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
			// 	Codec: codec.MustGetCodec("gogopb"),
			// 	Type:  reflect.TypeOf((*ErBaGang.CMD_S_PlayerList)(nil)).Elem(),
			// 	ID:    ENWORD(4, 111),
			// })
			//
			cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
				Codec: codec.MustGetCodec("gogopb"),
				Type:  reflect.TypeOf((*ErBaGang.CMD_S_Jetton_Broadcast)(nil)).Elem(),
				ID:    ENWORD(4, 114),
			})
		}
	case GGames.ByName["龙虎斗"].ID:
		{
			//玩家下注
			cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
				Codec: codec.MustGetCodec("gogopb"),
				Type:  reflect.TypeOf((*Longhu.CMD_C_PlaceJet)(nil)).Elem(),
				ID:    ENWORD(4, 100),
			})
			//获取玩家在线列表
			cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
				Codec: codec.MustGetCodec("gogopb"),
				Type:  reflect.TypeOf((*Longhu.CMD_C_PlayerList)(nil)).Elem(),
				ID:    ENWORD(4, 103),
			})
			//同步TIME
			cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
				Codec: codec.MustGetCodec("gogopb"),
				Type:  reflect.TypeOf((*Longhu.CMD_C_SyncTime_Req)(nil)).Elem(),
				ID:    ENWORD(4, 104),
			})
			//同步TIME
			cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
				Codec: codec.MustGetCodec("gogopb"),
				Type:  reflect.TypeOf((*Longhu.CMD_S_SyncTime_Res)(nil)).Elem(),
				ID:    ENWORD(4, 105),
			})
			//游戏空闲
			cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
				Codec: codec.MustGetCodec("gogopb"),
				Type:  reflect.TypeOf((*Longhu.CMD_Scene_StatusFree)(nil)).Elem(),
				ID:    ENWORD(4, 120),
			})
			//游戏开始
			cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
				Codec: codec.MustGetCodec("gogopb"),
				Type:  reflect.TypeOf((*Longhu.CMD_S_GameStart)(nil)).Elem(),
				ID:    ENWORD(4, 121),
			})
			//用户下注
			cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
				Codec: codec.MustGetCodec("gogopb"),
				Type:  reflect.TypeOf((*Longhu.CMD_S_PlaceJetSuccess)(nil)).Elem(),
				ID:    ENWORD(4, 122),
			})
			//当局游戏结束
			cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
				Codec: codec.MustGetCodec("gogopb"),
				Type:  reflect.TypeOf((*Longhu.CMD_S_GameEnd)(nil)).Elem(),
				ID:    ENWORD(4, 123),
			})
			//游戏记录
			cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
				Codec: codec.MustGetCodec("gogopb"),
				Type:  reflect.TypeOf((*Longhu.CMD_S_GameRecord)(nil)).Elem(),
				ID:    ENWORD(4, 127),
			})
			//下注失败
			cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
				Codec: codec.MustGetCodec("gogopb"),
				Type:  reflect.TypeOf((*Longhu.CMD_S_PlaceJettonFail)(nil)).Elem(),
				ID:    ENWORD(4, 128),
			})
			//玩家在线列表返回
			cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
				Codec: codec.MustGetCodec("gogopb"),
				Type:  reflect.TypeOf((*Longhu.CMD_S_PlayerList)(nil)).Elem(),
				ID:    ENWORD(4, 130),
			})
			//开始下注
			cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
				Codec: codec.MustGetCodec("gogopb"),
				Type:  reflect.TypeOf((*Longhu.CMD_S_StartPlaceJetton)(nil)).Elem(),
				ID:    ENWORD(4, 139),
			})
			//
			cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
				Codec: codec.MustGetCodec("gogopb"),
				Type:  reflect.TypeOf((*Longhu.CMD_S_Jetton_Broadcast)(nil)).Elem(),
				ID:    ENWORD(4, 114),
			})
		}
	case GGames.ByName["百人牛牛"].ID:
		{
			//玩家下注
			cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
				Codec: codec.MustGetCodec("gogopb"),
				Type:  reflect.TypeOf((*Brnn.CMD_C_PlaceJet)(nil)).Elem(),
				ID:    ENWORD(4, 100),
			})
			//获取玩家在线列表
			cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
				Codec: codec.MustGetCodec("gogopb"),
				Type:  reflect.TypeOf((*Brnn.CMD_C_PlayerList)(nil)).Elem(),
				ID:    ENWORD(4, 103),
			})
			//同步TIME
			cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
				Codec: codec.MustGetCodec("gogopb"),
				Type:  reflect.TypeOf((*Brnn.CMD_C_SyncTime_Req)(nil)).Elem(),
				ID:    ENWORD(4, 104),
			})
			//同步TIME
			cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
				Codec: codec.MustGetCodec("gogopb"),
				Type:  reflect.TypeOf((*Brnn.CMD_S_SyncTime_Res)(nil)).Elem(),
				ID:    ENWORD(4, 105),
			})
			//游戏空闲
			cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
				Codec: codec.MustGetCodec("gogopb"),
				Type:  reflect.TypeOf((*Brnn.CMD_Scene_StatusFree)(nil)).Elem(),
				ID:    ENWORD(4, 120),
			})
			//游戏开始
			cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
				Codec: codec.MustGetCodec("gogopb"),
				Type:  reflect.TypeOf((*Brnn.CMD_S_GameStart)(nil)).Elem(),
				ID:    ENWORD(4, 121),
			})
			//用户下注
			cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
				Codec: codec.MustGetCodec("gogopb"),
				Type:  reflect.TypeOf((*Brnn.CMD_S_PlaceJetSuccess)(nil)).Elem(),
				ID:    ENWORD(4, 122),
			})
			//当局游戏结束
			cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
				Codec: codec.MustGetCodec("gogopb"),
				Type:  reflect.TypeOf((*Brnn.CMD_S_GameEnd)(nil)).Elem(),
				ID:    ENWORD(4, 123),
			})
			//游戏记录
			cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
				Codec: codec.MustGetCodec("gogopb"),
				Type:  reflect.TypeOf((*Brnn.CMD_S_GameRecord)(nil)).Elem(),
				ID:    ENWORD(4, 127),
			})
			//下注失败
			cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
				Codec: codec.MustGetCodec("gogopb"),
				Type:  reflect.TypeOf((*Brnn.CMD_S_PlaceJettonFail)(nil)).Elem(),
				ID:    ENWORD(4, 128),
			})
			//玩家在线列表返回
			cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
				Codec: codec.MustGetCodec("gogopb"),
				Type:  reflect.TypeOf((*Brnn.CMD_S_PlayerList)(nil)).Elem(),
				ID:    ENWORD(4, 130),
			})
			//开始下注
			cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
				Codec: codec.MustGetCodec("gogopb"),
				Type:  reflect.TypeOf((*Brnn.CMD_S_StartPlaceJetton)(nil)).Elem(),
				ID:    ENWORD(4, 139),
			})
			//
			cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
				Codec: codec.MustGetCodec("gogopb"),
				Type:  reflect.TypeOf((*Brnn.CMD_S_Jetton_Broadcast)(nil)).Elem(),
				ID:    ENWORD(4, 140),
			})
		}
	case GGames.ByName["红黑大战"].ID:
		{
			//玩家下注
			cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
				Codec: codec.MustGetCodec("gogopb"),
				Type:  reflect.TypeOf((*HongHei.CMD_C_PlaceJet)(nil)).Elem(),
				ID:    ENWORD(4, 100),
			})
			//获取玩家在线列表
			cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
				Codec: codec.MustGetCodec("gogopb"),
				Type:  reflect.TypeOf((*HongHei.CMD_C_PlayerList)(nil)).Elem(),
				ID:    ENWORD(4, 103),
			})
			//下注玩家
			cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
				Codec: codec.MustGetCodec("gogopb"),
				Type:  reflect.TypeOf((*HongHei.CMD_C_ReJetton)(nil)).Elem(),
				ID:    ENWORD(4, 104),
			})
			//游戏空闲
			cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
				Codec: codec.MustGetCodec("gogopb"),
				Type:  reflect.TypeOf((*HongHei.CMD_Scene_StatusFree)(nil)).Elem(),
				ID:    ENWORD(4, 120),
			})
			//游戏开始
			cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
				Codec: codec.MustGetCodec("gogopb"),
				Type:  reflect.TypeOf((*HongHei.CMD_S_GameStart)(nil)).Elem(),
				ID:    ENWORD(4, 121),
			})
			//用户下注
			cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
				Codec: codec.MustGetCodec("gogopb"),
				Type:  reflect.TypeOf((*HongHei.CMD_S_PlaceJetSuccess)(nil)).Elem(),
				ID:    ENWORD(4, 122),
			})
			//当局游戏结束
			cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
				Codec: codec.MustGetCodec("gogopb"),
				Type:  reflect.TypeOf((*HongHei.CMD_S_GameEnd)(nil)).Elem(),
				ID:    ENWORD(4, 123),
			})
			//游戏记录
			cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
				Codec: codec.MustGetCodec("gogopb"),
				Type:  reflect.TypeOf((*HongHei.CMD_S_GameRecord)(nil)).Elem(),
				ID:    ENWORD(4, 127),
			})
			//下注失败
			cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
				Codec: codec.MustGetCodec("gogopb"),
				Type:  reflect.TypeOf((*HongHei.CMD_S_PlaceJettonFail)(nil)).Elem(),
				ID:    ENWORD(4, 128),
			})
			//玩家在线列表返回
			cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
				Codec: codec.MustGetCodec("gogopb"),
				Type:  reflect.TypeOf((*HongHei.CMD_S_PlayerList)(nil)).Elem(),
				ID:    ENWORD(4, 130),
			})
			//开始下注
			cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
				Codec: codec.MustGetCodec("gogopb"),
				Type:  reflect.TypeOf((*HongHei.CMD_S_StartPlaceJetton)(nil)).Elem(),
				ID:    ENWORD(4, 139),
			})
			//
			cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
				Codec: codec.MustGetCodec("gogopb"),
				Type:  reflect.TypeOf((*HongHei.CMD_S_Jetton_Broadcast)(nil)).Elem(),
				ID:    ENWORD(4, 140),
			})
		}
	case GGames.ByName["抢庄牌九"].ID:
		{

		}
	case GGames.ByName["奔驰宝马"].ID:
		{
			//開始游戲(SUB_S_GameStart)
			cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
				Codec: codec.MustGetCodec("gogopb"),
				Type:  reflect.TypeOf((*BenCiBaoMa.CMD_S_GameStart)(nil)).Elem(),
				ID:    ENWORD(4, 100),
			})
			//游戲結束(SUB_S_GameEnd)
			cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
				Codec: codec.MustGetCodec("gogopb"),
				Type:  reflect.TypeOf((*BenCiBaoMa.CMD_S_GameEnd)(nil)).Elem(),
				ID:    ENWORD(4, 101),
			})
			//開始游戲場景(SUB_S_SCENE_START)
			cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
				Codec: codec.MustGetCodec("gogopb"),
				Type:  reflect.TypeOf((*BenCiBaoMa.CMD_S_Scene_GameStart)(nil)).Elem(),
				ID:    ENWORD(4, 102),
			})
			//結束游戲場景(SUB_S_SCENE_END)
			cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
				Codec: codec.MustGetCodec("gogopb"),
				Type:  reflect.TypeOf((*BenCiBaoMa.CMD_S_Scene_GameEnd)(nil)).Elem(),
				ID:    ENWORD(4, 103),
			})
			//游戏结束时返回玩家结算信息(SUB_S_PLAYERLIST)
			cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
				Codec: codec.MustGetCodec("gogopb"),
				Type:  reflect.TypeOf((*BenCiBaoMa.CMD_S_UserWinList)(nil)).Elem(),
				ID:    ENWORD(4, 104),
			})
			//下注成功返回(SUB_S_JETTON_SUCCESS)
			cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
				Codec: codec.MustGetCodec("gogopb"),
				Type:  reflect.TypeOf((*BenCiBaoMa.CMD_S_PlaceJetSuccess)(nil)).Elem(),
				ID:    ENWORD(4, 105),
			})
			//下注失败返回(SUB_S_JETTON_FAIL)
			cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
				Codec: codec.MustGetCodec("gogopb"),
				Type:  reflect.TypeOf((*BenCiBaoMa.CMD_S_PlaceJettonFail)(nil)).Elem(),
				ID:    ENWORD(4, 106),
			})
			//玩家下注(SUB_C_USER_JETTON)
			cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
				Codec: codec.MustGetCodec("gogopb"),
				Type:  reflect.TypeOf((*BenCiBaoMa.CMD_C_PlaceJet)(nil)).Elem(),
				ID:    ENWORD(4, 107),
			})
			//玩家申請排行榜信息(SUB_C_USER_ASKLIST)应该暂时用不上了
			cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
				Codec: codec.MustGetCodec("gogopb"),
				Type:  reflect.TypeOf((*BenCiBaoMa.CMD_C_AskList)(nil)).Elem(),
				ID:    ENWORD(4, 108),
			})
			//续押命令(SUB_C_USER_REPEAT_JETTON)
			cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
				Codec: codec.MustGetCodec("gogopb"),
				Type:  reflect.TypeOf((*BenCiBaoMa.CMD_C_ReJetton)(nil)).Elem(),
				ID:    ENWORD(4, 109),
			})
			//客户端获取玩家列表命令(SUB_C_QUERY_PLAYERLIST)
			cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
				Codec: codec.MustGetCodec("gogopb"),
				Type:  reflect.TypeOf((*BenCiBaoMa.CMD_C_PlayerList)(nil)).Elem(),
				ID:    ENWORD(4, 110),
			})
			//获取玩家在线列表(SUB_S_QUERY_PLAYLIST)
			cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
				Codec: codec.MustGetCodec("gogopb"),
				Type:  reflect.TypeOf((*BenCiBaoMa.CMD_S_PlayerList)(nil)).Elem(),
				ID:    ENWORD(4, 111),
			})
			//获取当前状态(SUB_C_QUERY_CUR_STATE)
			cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
				Codec: codec.MustGetCodec("gogopb"),
				Type:  reflect.TypeOf((*BenCiBaoMa.CMD_C_CurState)(nil)).Elem(),
				ID:    ENWORD(4, 112),
			})
			//服务器返回当前状态(SUB_S_QUERY_CUR_STATE)
			cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
				Codec: codec.MustGetCodec("gogopb"),
				Type:  reflect.TypeOf((*BenCiBaoMa.CMD_S_CurState)(nil)).Elem(),
				ID:    ENWORD(4, 113),
			})
			//续押下注成功返回(SUB_S_REPEAT_JETTON_OK)
			cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
				Codec: codec.MustGetCodec("gogopb"),
				Type:  reflect.TypeOf((*BenCiBaoMa.CMD_S_RepeatJetSuccess)(nil)).Elem(),
				ID:    ENWORD(4, 114),
			})
			//广播一段时间内其他所有玩家的下注(SUB_S_JETTON_SUCCESS_OTHER)
			cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
				Codec: codec.MustGetCodec("gogopb"),
				Type:  reflect.TypeOf((*BenCiBaoMa.CMD_S_OtherPlaceJetSuccess)(nil)).Elem(),
				ID:    ENWORD(4, 115),
			})
		}
	}
}

//Msg 消息结构
type Msg struct {
	ver     uint16
	sign    uint16
	encType uint8
	mainID  uint8
	subID   uint8
	msg     interface{}
}

//
func newMsg(mainID uint8, subID uint8, msg interface{}) *Msg {
	return &Msg{
		ver:     0x0001,
		sign:    0x5F5F,
		encType: 0x02,
		mainID:  mainID,
		subID:   subID,
		msg:     msg,
	}
}
