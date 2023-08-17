/**
User: cr-mao
Date: 2023/8/13 16:43
Email: crmao@qq.com
Desc: go_call.go
*/
package sugar

import (
	"runtime"

	"github.com/cr-mao/loric/log"
)

func SafeGo(fn func()) {
	defer func() {
		if err := recover(); err != nil {
			switch err.(type) {
			case runtime.Error:
				log.Error(err)
			default:
				log.Errorf("panic error: %v", err)
			}
		}
	}()

	fn()
}
