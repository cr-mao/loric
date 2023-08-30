/**
User: cr-mao
Date: 2023/8/30 15:40
Email: crmao@qq.com
Desc: slice.go
*/
package slice

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
