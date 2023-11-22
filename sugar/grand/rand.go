/**
User: cr-mao
Date: 2023/10/31 13:36
Email: crmao@qq.com
Desc: rand.go
*/
package grand

import (
	"crypto/rand"
	"math/big"
)

// 真随机
func RealRandNum(num int64) int64 {
	n, _ := rand.Int(rand.Reader, big.NewInt(num))
	return n.Int64()
}
