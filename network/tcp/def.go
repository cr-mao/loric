package tcp

import (
	"encoding/binary"
	"io"
	"net"

	"github.com/cr-mao/loric/errors"
)

var ErrMsgSizeToBig = errors.New("msg size to big")

const sizeBytes = 4

const (
	closeSig   int = iota // 关闭信号
	dataPacket            // 数据包
)

type chWrite struct {
	typ int
	msg []byte
}

// 执行写入操作
func write(conn net.Conn, msg []byte) error {
	buf := make([]byte, sizeBytes+len(msg))

	binary.BigEndian.PutUint32(buf, uint32(len(msg)))
	copy(buf[sizeBytes:], msg)

	_, err := conn.Write(buf)
	return err
}

// 执行读取操作
func read(conn net.Conn, maxMsgSize uint32) ([]byte, error) {
	buf := make([]byte, sizeBytes)
	if _, err := io.ReadFull(conn, buf); err != nil {
		return nil, err
	}

	size := binary.BigEndian.Uint32(buf)
	if size == 0 {
		return nil, nil
	}
	if maxMsgSize > 0 && size > maxMsgSize {
		return nil, ErrMsgSizeToBig
	}

	buf = make([]byte, size)
	if _, err := io.ReadFull(conn, buf); err != nil {
		return nil, err
	}

	return buf, nil
}
