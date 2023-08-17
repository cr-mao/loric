package packet

import (
	"bytes"
	"encoding/binary"

	"github.com/cr-mao/loric/errors"
	"github.com/cr-mao/loric/log"
)

var (
	ErrSeqOverflow    = errors.New("seq overflow")
	ErrRouteOverflow  = errors.New("route overflow")
	ErrInvalidMessage = errors.New("invalid message")
	ErrBufferTooLarge = errors.New("buffer too large")
)

type Packer interface {
	// Pack 打包
	Pack(message *Message) ([]byte, error)
	// Unpack 解包
	Unpack(data []byte) (*Message, error)
}

type defaultPacker struct {
	opts *options
}

func NewPacker(opts ...Option) Packer {
	o := defaultOptions()
	for _, opt := range opts {
		opt(o)
	}
	if o.bufferBytes < 0 {
		log.Fatalf("the number of buffer bytes must be greater than or equal to 0, and give %d", o.bufferBytes)
	}

	return &defaultPacker{opts: o}
}

// Pack 打包消息
func (p *defaultPacker) Pack(message *Message) ([]byte, error) {
	if message.Route > 32767 || message.Route < int32(-1<<(8*p.opts.routeBytes-1)) {
		return nil, ErrRouteOverflow
	}

	if p.opts.seqBytes > 0 {
		if message.Seq > 2147483647 || message.Seq < -2147483648 {
			return nil, ErrSeqOverflow
		}
	}
	if len(message.Buffer) > p.opts.bufferBytes {
		return nil, ErrBufferTooLarge
	}
	var (
		err error
		//ln  = p.opts.routeBytes + p.opts.seqBytes + len(message.Buffer)
		ln  = 6 + len(message.Buffer)
		buf = &bytes.Buffer{}
	)
	buf.Grow(ln)
	err = binary.Write(buf, binary.BigEndian, int16(message.Route))
	if err != nil {
		return nil, err
	}
	err = binary.Write(buf, binary.BigEndian, message.Seq)
	if err != nil {
		return nil, err
	}
	err = binary.Write(buf, binary.BigEndian, message.Buffer)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// Unpack 解包消息
func (p *defaultPacker) Unpack(data []byte) (*Message, error) {
	var (
		err error
		ln  = len(data) - p.opts.routeBytes - p.opts.seqBytes
		//ln     = len(data) - 6
		reader = bytes.NewReader(data)
	)
	if ln < 0 {
		return nil, ErrInvalidMessage
	}
	message := &Message{Buffer: make([]byte, ln)}
	var route int16
	if err = binary.Read(reader, binary.BigEndian, &route); err != nil {
		return nil, err
	}
	message.Route = int32(route)

	var seq int32
	if err = binary.Read(reader, binary.BigEndian, &seq); err != nil {
		return nil, err
	}
	message.Seq = seq
	err = binary.Read(reader, binary.BigEndian, &message.Buffer)
	if err != nil {
		return nil, err
	}
	return message, nil
}
