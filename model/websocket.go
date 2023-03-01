package model

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

func BytesCombine(pBytes ...[]byte) []byte {
	return bytes.Join(pBytes, []byte(""))
}

func WebSocketUnpack(data []byte) []byte {
	en_bytes := []byte("")
	cn_bytes := make([]int, 0)

	v := data[1] & 0x7f
	p := 0
	switch v {
	case 0x7e:
		p = 4
	case 0x7f:
		p = 10
	default:
		p = 2
	}
	mask := data[p : p+4]
	data_tmp := data[p+4:]
	nv := ""
	nv_bytes := []byte("")
	nv_len := 0

	for k, v := range data_tmp {

		nv = string(int(v ^ mask[k%4]))
		// nv = fmt.Sprintf("%d", int(v^mask[k%4]))
		nv_bytes = []byte(nv)
		nv_len = len(nv_bytes)
		if nv_len == 1 {
			en_bytes = BytesCombine(en_bytes, nv_bytes)
		} else {
			en_bytes = BytesCombine(en_bytes, []byte("%s"))
			cn_bytes = append(cn_bytes, int(v^mask[k%4]))
		}
	}

	//处理中文
	cn_str := make([]interface{}, 0)
	if len(cn_bytes) > 2 {
		clen := len(cn_bytes)
		count := int(clen / 3)

		for i := 0; i < count; i++ {
			mm := i * 3

			hh := make([]byte, 3)
			h1, _ := IntToBytes(cn_bytes[mm], 1)
			h2, _ := IntToBytes(cn_bytes[mm+1], 1)
			h3, _ := IntToBytes(cn_bytes[mm+2], 1)
			hh[0] = h1[0]
			hh[1] = h2[0]
			hh[2] = h3[0]

			cn_str = append(cn_str, string(hh))
		}
		// TODO string to []byte
		new := string(bytes.Replace(en_bytes, []byte("%s%s%s"), []byte("%s"), -1))
		return []byte(fmt.Sprintf(new, cn_str...))

	}
	return en_bytes
}

func WebSocketPack(data []byte) []byte {
	lenth := len(data)
	token := string(0x81)
	if lenth < 126 {
		token += string(lenth)
	}
	bb, _ := IntToBytes(0x81, 1)
	b0 := bb[0]
	b1 := byte(0)
	framePos := 0
	// fmt.Println("长度", lenth)
	switch {
	case lenth >= 65536:
		writeBuf := make([]byte, 10)
		writeBuf[framePos] = b0
		writeBuf[framePos+1] = b1 | 127
		binary.BigEndian.PutUint64(writeBuf[framePos+2:], uint64(lenth))

		return BytesCombine(writeBuf, data)
	case lenth > 125:
		fmt.Println("》125")
		writeBuf := make([]byte, 4)
		writeBuf[framePos] = b0
		writeBuf[framePos+1] = b1 | 126
		binary.BigEndian.PutUint16(writeBuf[framePos+2:], uint16(lenth))
		fmt.Println(writeBuf)
		return BytesCombine(writeBuf, data)
	default:
		writeBuf := make([]byte, 2)
		writeBuf[framePos] = b0
		writeBuf[framePos+1] = b1 | byte(lenth)

		return BytesCombine(writeBuf, data)
	}
}

// 整形转换成字节
func IntToBytes(n int, b byte) ([]byte, error) {
	switch b {
	case 1:
		tmp := int8(n)
		bytesBuffer := bytes.NewBuffer([]byte{})
		binary.Write(bytesBuffer, binary.BigEndian, &tmp)
		return bytesBuffer.Bytes(), nil
	case 2:
		tmp := int16(n)
		bytesBuffer := bytes.NewBuffer([]byte{})
		binary.Write(bytesBuffer, binary.BigEndian, &tmp)
		return bytesBuffer.Bytes(), nil
	case 3, 4:
		tmp := int32(n)
		bytesBuffer := bytes.NewBuffer([]byte{})
		binary.Write(bytesBuffer, binary.BigEndian, &tmp)
		return bytesBuffer.Bytes(), nil
	}
	return nil, fmt.Errorf("IntToBytes b param is invaild")
}
