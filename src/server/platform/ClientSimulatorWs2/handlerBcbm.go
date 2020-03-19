package main

import (
	"server/pb/BenCiBaoMa"
	"server/pb/GameServer"
	"server/pb/Game_Common"
	"strconv"

	"github.com/davyxu/cellnet/codec"
)

//sendPlayerPlaceJetBcbm 玩家下注
//-------------------------------------------------------------
func sendPlayerPlaceJetBcbm(peer Session, area int32, score int32) {
	reqdata := &GameServer.MSG_CSC_Passageway{}
	val, _ := strconv.ParseUint("F5F5F5F5", 16, 32)
	reqdata.Header = &Game_Common.Header{}
	reqdata.Header.Sign = int32(val)
	msg := newMsg(
		uint8(Game_Common.MAINID_MAIN_MESSAGE_CLIENT_TO_GAME_LOGIC),
		uint8(BenCiBaoMa.SUBID_SUB_C_USER_JETTON),
		reqdata)
	p := &BenCiBaoMa.CMD_C_PlaceJet{}
	p.CbJettonArea = area  //筹码区域
	p.LJettonScore = score //加注数目
	data, _, _ := codec.EncodeMessage(p, nil)
	reqdata.PassData = data[:]
	peer.Write(msg)

	//client := peer.GetCtx(TagUserInfo).(*DefWSClient)
	//util.Logz("UserClient", "Player", client.UserID, client.Account, "sendPlayerPlaceJetBcbm", p)
}
