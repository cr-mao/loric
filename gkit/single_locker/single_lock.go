/**
User: cr-mao
Date: 2023/10/11 09:29
Email: crmao@qq.com
Desc: 单进程全局锁， 不适用分布式锁。 都在同一进程，适合用。
*/
package user_lock

import "sync"

var lockMap = &sync.Map{}

func TryLock(key string) bool {
	if len(key) <= 0 {
		return false
	}
	_, loaded := lockMap.LoadOrStore(key, 1)
	return !loaded
}

func Unlock(key string) {
	if len(key) <= 0 {
		return
	}
	lockMap.Delete(key)
}

func HasLock(key string) bool {
	_, ok := lockMap.Load(key)
	return ok
}
