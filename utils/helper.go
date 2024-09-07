package utils

import (
	"bytes"
	"encoding/binary"
)

func IntToHex(num int64) ([]byte, error) {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, num)
	if err != nil {
		return nil, CatchErr(err)
	}

	return buff.Bytes(), nil
}
