/**
User: cr-mao
Date: 2023/11/1 19:22
Email: crmao@qq.com
Desc: gstring.go
*/
package gstring

import (
	"strconv"
	"strings"
)

func String2Int32Slice(s string, split string) []int32 {
	if s == "" {
		return []int32{}
	}
	strList := strings.Split(s, split)
	var res = make([]int32, 0)
	for _, v := range strList {
		vInt, _ := strconv.Atoi(v)
		res = append(res, int32(vInt))
	}
	return res
}
