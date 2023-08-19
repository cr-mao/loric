/**
User: cr-mao
Date: 2023/8/17 16:55
Email: crmao@qq.com
Desc: uuid.go
*/
package sugar

import (
	"github.com/google/uuid"
)

func UUID() (string, error) {
	id, err := uuid.NewUUID()
	if err != nil {
		return "", err
	}
	return id.String(), err
}
