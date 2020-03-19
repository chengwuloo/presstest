package main

import (
	"server/pb/GameServer"
	"server/pb/Game_Common"
	"server/pb/HongHei"
	"strconv"

	"github.com/davyxu/cellnet/codec"
)

//sendPlayerPlaceJetHongHei 玩家下注
//-------------------------------------------------------------
func sendPlayerPlaceJetHongHei(peer Session, area int32, score float64) {
	reqdata := &GameServer.MSG_CSC_Passageway{}
	val, _ := strconv.ParseUint("F5F5F5F5", 16, 32)
	reqdata.Header = &Game_Common.Header{}
	reqdata.Header.Sign = int32(val)
	msg := newMsg(
		uint8(Game_Common.MAINID_MAIN_MESSAGE_CLIENT_TO_GAME_LOGIC),
		uint8(HongHei.SUBID_SUB_C_PLACE_JETTON),
		reqdata)
	p := &HongHei.CMD_C_PlaceJet{}
	p.CbJettonArea = area  //筹码区域
	p.LJettonScore = score //加注数目
	data, _, _ := codec.EncodeMessage(p, nil)
	reqdata.PassData = data[:]
	peer.Write(msg)
}
