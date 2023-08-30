/**
User: cr-mao
Date: 2023/8/30 11:43
Email: crmao@qq.com
Desc: 全局变量都在这个包里面
*/
package manager

import (
	"errors"
	"sync"
)

type LianmengManager struct {
	LimengInfos map[int64]map[int64]*User
	locker      sync.RWMutex
}

type User struct {
	UserId     int64
	UserName   string
	LianmengId int64
}

var Manager *LianmengManager

func init() {
	Manager = &LianmengManager{
		LimengInfos: make(map[int64]map[int64]*User),
	}
}

func GetManager() *LianmengManager {
	return Manager
}

func (m *LianmengManager) GetUserByUserIdAndLianMengId(userId int64, lianmengId int64) (*User, error) {
	m.locker.RLock()
	defer m.locker.RUnlock()
	limengInfo, ok := m.LimengInfos[lianmengId]
	if !ok {
		return nil, errors.New("lianmengId not exists")
	}
	user, ok := limengInfo[userId]
	if !ok {
		return nil, errors.New("lianmeng userId not exists")
	}
	return user, nil
}

func (m *LianmengManager) AddUser(user *User) {
	if user == nil {
		return
	}
	m.locker.Lock()
	defer m.locker.Unlock()
	limengInfo, ok := m.LimengInfos[user.LianmengId]
	if !ok {
		m.LimengInfos[user.LianmengId] = make(map[int64]*User)
		limengInfo = m.LimengInfos[user.LianmengId]
	}
	limengInfo[user.UserId] = user
}

func (m *LianmengManager) DelUser(user *User) {
	if user == nil {
		return
	}
	m.locker.Lock()
	defer m.locker.Unlock()
	limengInfo, ok := m.LimengInfos[user.LianmengId]
	if !ok {
		return
	}
	if limengInfo != nil {
		delete(limengInfo, user.UserId)
	}
}

func (m *LianmengManager) GetUserIdListByLianmengId(lianmengId int64, userId int64) []int64 {
	m.locker.RLock()
	defer m.locker.RUnlock()
	limengInfo, ok := m.LimengInfos[lianmengId]
	if !ok {
		return nil
	}
	var res = make([]int64, 0, len(limengInfo))
	for id := range limengInfo {
		if userId == id {
			continue
		}
		res = append(res, id)
	}
	return res
}
