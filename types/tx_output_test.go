// Copyright 2020 The Reed Developers
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

package types

import (
	"bytes"
	"encoding/hex"
	"github.com/reed/common/byteutil/byteconv"
	"github.com/reed/crypto"
	"github.com/reed/vm/vmcommon"
	"testing"
)

func TestTxOutput_GenerateID(t *testing.T) {
	var addr, _ = hex.DecodeString("c27f26c8bf818e5509abacfc20206d43fc0db6a415f20d48726eb8cd2888f68e")
	var scriptPK, _ = hex.DecodeString("bf8776efb3367228d115c325a623b3fe6b359a87e45d25c98a506e203b0ec5b1fc0db6a415f20d48726eb8cd2888f68")
	amt := uint64(199)
	icb := false
	output := &TxOutput{
		IsCoinBase: icb,
		Amount:     amt,
		Address:    addr,
		ScriptPk:   scriptPK,
	}

	id := output.GenerateID()

	var datas [][]byte
	split := []byte(":")
	datas = append(datas, []byte{0}, split, addr, split, byteconv.Uint64ToBytes(amt), split, scriptPK)

	h := crypto.Sha256(datas...)

	if !bytes.Equal(id.Bytes(), h) {
		t.Fatalf("GenerateID error")
	}

}

func TestTxOutput_SetLockingScript(t *testing.T) {
	var pub, _ = hex.DecodeString("b12049d709358dc427433050625aa2135163181ccc320f22859d7c065ecc9dcb")
	o := TxOutput{
		Address: pub,
	}

	script := o.GenerateLockingScript()

	//pk + 5op
	if len(script) != (32 + 5) {
		t.Errorf("script len error;actual len=%d", len(script))
	}
	if script[0] != vmcommon.OpDup {
		t.Error("script first part is not OpDup")
	}
	if script[1] != vmcommon.OpHash256 {
		t.Error("script second part is not OpHash256")
	}

	if script[2] != vmcommon.OpPushData32 {
		t.Error("script third part is not OpPushData32")
	}

	pubHash := crypto.Sha256(pub)

	if !bytes.Equal(script[3:35], pubHash) {
		t.Error("script fourth part is not Hash data")
	}

	if script[35] != vmcommon.OpEqualVerify {
		t.Error("script fifth part is not OpEqualVerify")
	}
	if script[36] != vmcommon.OpCheckSig {
		t.Error("script sixth part is not OpCheckSig")
	}

}
