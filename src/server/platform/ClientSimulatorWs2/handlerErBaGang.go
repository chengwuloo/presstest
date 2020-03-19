package main

import (
	"server/pb/ErBaGang"
	"server/pb/GameServer"
	"server/pb/Game_Common"
	"strconv"

	"github.com/davyxu/cellnet/codec"
)

//sendPlayerPlaceJetErBaGang 玩家下注
//-------------------------------------------------------------
func sendPlayerPlaceJetErBaGang(peer Session, area int32, score float64) {
	reqdata := &GameServer.MSG_CSC_Passageway{}
	val, _ := strconv.ParseUint("F5F5F5F5", 16, 32)
	reqdata.Header = &Game_Common.Header{}
	reqdata.Header.Sign = int32(val)
	msg := newMsg(
		uint8(Game_Common.MAINID_MAIN_MESSAGE_CLIENT_TO_GAME_LOGIC),
		uint8(ErBaGang.SUBID_SUB_C_USER_JETTON),
		reqdata)
	p := &ErBaGang.CMD_C_PlaceJet{}
	p.CbJettonArea = area  //筹码区域
	p.LJettonScore = score //加注数目
	data, _, _ := codec.EncodeMessage(p, nil)
	reqdata.PassData = data[:]
	peer.Write(msg)
}
