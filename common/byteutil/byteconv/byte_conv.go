// Copyright 2020 The Reed Developers
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

package byteconv

import (
	"bytes"
	"encoding/binary"
)

func ByteToUint64(b []byte) uint64 {
	buf := bytes.NewBuffer(b)
	var data uint64
	binary.Read(buf, binary.BigEndian, &data)
	return data
}

func Uint64ToBytes(n uint64) []byte {
	buf := bytes.NewBuffer([]byte{})
	binary.Write(buf, binary.BigEndian, n)
	return buf.Bytes()
}
