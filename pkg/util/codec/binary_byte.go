package codec

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"errors"
	"fmt"

	"github.com/bytedance/sonic"

	"github.com/hanzug/goS/types"
)

// BinaryWrite 将所有的类型 转成byte buffer类型，易于存储// TODO change
// func BinaryWrite(v any) (buf *bytes.Buffer, err error) {
func BinaryWrite(buf *bytes.Buffer, v any) (err error) {
	size := binary.Size(v)
	// log.Debug("docid size:", size)
	fmt.Println("size", size)
	if size <= 0 {
		zap.S().Errorf("encodePostings binary.Size err,size: %v", size)
		return
	}

	err = binary.Write(buf, binary.LittleEndian, v)
	if err != nil {
		fmt.Println(err)
	}

	return
}

// GobWrite 将所有的类型 转成 bytes.Buffer 类型，易于存储// TODO change
func GobWrite(v any) (buf *bytes.Buffer, err error) {
	if v == nil {
		err = errors.New("BinaryWrite the value is nil")
		return
	}
	buf = new(bytes.Buffer)
	enc := gob.NewEncoder(buf)
	if err = enc.Encode(v); err != nil {
		return
	}

	return
}

// DecodePostings 解码 return *PostingsList postingslen err
func DecodePostings(buf []byte) (p *types.InvertedIndexValue, err error) {
	p = new(types.InvertedIndexValue)
	err = sonic.Unmarshal(buf, &p)

	return
}

// EncodePostings 编码
func EncodePostings(postings *types.InvertedIndexValue) (buf []byte, err error) {
	buf, err = sonic.Marshal(postings)
	if err != nil {
		zap.S().Errorf("sonic.Marshal err:%v,postings:%+v", err, postings)
		return
	}

	return
}

// BinaryEncoding 二进制编码
func BinaryEncoding(buf *bytes.Buffer, v any) (err error) {
	err = gob.NewEncoder(buf).Encode(v)
	return
}

// BinaryDecoding 二进制解码
func BinaryDecoding(buf *bytes.Buffer, v any) (err error) {
	err = gob.NewDecoder(buf).Decode(v)

	return
}
