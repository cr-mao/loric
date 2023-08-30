/**
User: cr-mao
Date: 2023/8/30 11:43
Email: crmao@qq.com
Desc: 全局变量都在这个包里面
*/
package global

import "sync"

type UserManager struct {
	UserMap sync.Map // user_token => &User
}

type User struct {
	UserId   int64  // 用户id
	UserName string // 用户名
	LiMengId int64  // 所加入的联盟id
}

var UserId = int64(0)
var UserData *UserManager

func init() {
	UserData = &UserManager{}
}
