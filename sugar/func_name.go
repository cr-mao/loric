/**
User: cr-mao
Date: 2023/8/23 17:16
Email: crmao@qq.com
Desc: func.go
*/
package sugar

import (
	"reflect"
	"runtime"
)

func NameOfFunction(f interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
}
