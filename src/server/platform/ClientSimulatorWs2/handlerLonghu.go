package main

import (
	"server/pb/GameServer"
	"server/pb/Game_Common"
	"server/pb/Longhu"
	"strconv"

	"github.com/davyxu/cellnet/codec"
)

//sendPlayerPlaceJetLonghu 玩家下注
//-------------------------------------------------------------
func sendPlayerPlaceJetLonghu(peer Session, area int32, score int64) {
	reqdata := &GameServer.MSG_CSC_Passageway{}
	val, _ := strconv.ParseUint("F5F5F5F5", 16, 32)
	reqdata.Header = &Game_Common.Header{}
	reqdata.Header.Sign = int32(val)
	msg := newMsg(
		uint8(Game_Common.MAINID_MAIN_MESSAGE_CLIENT_TO_GAME_LOGIC),
		uint8(Longhu.SUBID_SUB_C_PLACE_JETTON),
		reqdata)
	p := &Longhu.CMD_C_PlaceJet{}
	p.CbJettonArea = area  //筹码区域
	p.LJettonScore = score //加注数目
	data, _, _ := codec.EncodeMessage(p, nil)
	reqdata.PassData = data[:]
	peer.Write(msg)
}
