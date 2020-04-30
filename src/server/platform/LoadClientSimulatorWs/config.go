package main

import (
	"fmt"
	"server/platform/util"
)

//
type iniConfig struct {
	flag         int
	children     int
	httpaddr     string
	httpaddr1    string
	httptimeout  int
	wsaddr       string
	md5code      string
	descode      string
	isdecrypt    int
	agentid      int
	dynamic      int
	wssl         int
	numMailbox   int
	totalClients int
	numClients   int
	numClients2  int
	numClients3  int
	baseAccount  int64
	deltaClients int
	deltaTime    int
	heartbeat    int
	timeout      int
	subgameID    int
	subroomID    int
	tokenprefix  string
	tokenstart   int
	tokenend     int
}

//
func readini(filename string) (c *iniConfig) {
	ini := util.Ini{}
	if err := ini.Load("conf.ini"); err != nil {
		fmt.Printf("load %s err: [%s]\n", filename, err.Error())
		return
	}
	c = &iniConfig{}
	c.flag = ini.GetInt("flag", "flag")
	c.children = ini.GetInt("children", "num")
	c.httpaddr = ini.GetString("httpaddr", "httpaddr")
	c.httpaddr1 = ini.GetString("httpaddr", "httpaddr1")
	c.isdecrypt = ini.GetInt("httpaddr", "isdecrypt")
	c.httptimeout = ini.GetInt("httpaddr", "timeout")
	c.wsaddr = ini.GetString("wsaddr", "wsaddr")
	c.dynamic = ini.GetInt("wsaddr", "dynamic")
	c.wssl = ini.GetInt("wsaddr", "wssl")
	c.numMailbox = ini.GetInt("mailboxs", "mailboxs")
	c.totalClients = ini.GetInt("clients", "totalClients")
	c.numClients = ini.GetInt("clients", "numClients")
	c.numClients2 = ini.GetInt("clients", "numClients2")
	c.numClients3 = ini.GetInt("clients", "numClients3")
	c.baseAccount = ini.GetInt64("baseTest", "baseTest")
	c.deltaClients = ini.GetInt("deltaClients", "deltaClients")
	c.deltaTime = ini.GetInt("deltaTime", "deltaTime")
	c.heartbeat = ini.GetInt("heartbeat", "interval")
	c.timeout = ini.GetInt("heartbeat", "timeout")
	c.subgameID = ini.GetInt("subgame", "gameID")
	c.subroomID = ini.GetInt("subgame", "roomID")
	c.tokenprefix = ini.GetString("baseTest", "prefix")
	c.tokenstart = ini.GetInt("baseTest", "tokenstart")
	c.tokenend = ini.GetInt("baseTest", "tokenend")
	c.agentid = ini.GetInt("baseTest", "agentid")
	c.md5code = ini.GetString("baseTest", "md5code")
	c.descode = ini.GetString("baseTest", "descode")
	// log.Println("children: ", c.children)
	// log.Println("httpaddr: ", c.httpaddr)
	// log.Println("wsaddr: ", c.wsaddr)
	// log.Println("numMailbox: ", c.numMailbox)
	// log.Println("numClient: ", c.numClient)
	// log.Println("baseAccount: ", c.baseAccount)
	// log.Println("deltaClients: ", c.deltaClients)
	// log.Println("deltaTime: ", c.deltaTime)
	return
}
