package main

//
// Created by andy_ro@qq.com
// 			4/9/2019
//
import (
	//sync.Pool sync.Mutex sync.RWMutex sync.WaitGroup sync.Cond sync.Once
	"sync"
)

//UserInfo 用户信息
type UserInfo struct {
	UserID  int64
	Account int64
	AgentID uint32
}

// UserList 用户列表
type UserList struct {
	Players []*UserInfo
	l       *sync.Mutex
}

// NewUserList 创建用户列表
func NewUserList() *UserList {
	return &UserList{Players: []*UserInfo{}, l: &sync.Mutex{}}
}

//AddUserInfo 添加用户信息
func (s *UserList) AddUserInfo(userID int64, account int64, agentID uint32) {
	s.l.Lock()
	s.Players = append(s.Players, &UserInfo{UserID: userID, Account: account, AgentID: agentID})
	s.l.Unlock()
}

//GLoginedUsers 已登陆用户信息
var GLoginedUsers = NewUserList()
