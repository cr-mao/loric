/*
*
User: cr-mao
Date: 2023/8/30 15:40
Email: crmao@qq.com
Desc: slice.go
*/
package gslice

import (
	"strconv"
	"strings"
)

func InSliceInt64(src int64, bocket []int64) bool {
	length := len(bocket)
	for i := 0; i < length; i++ {
		if src == bocket[i] {
			return true
		}
	}
	return false
}

func InSliceInt32(src int32, bocket []int32) bool {
	length := len(bocket)
	for i := 0; i < length; i++ {
		if src == bocket[i] {
			return true
		}
	}
	return false
}

// int32切片转 字符串
func SliceInt32ToString(src []int32, split string) string {
	var res string
	for _, v := range src {
		s1 := strconv.Itoa(int(v))
		res += s1 + split
	}

	return strings.TrimRight(res, split)
}
