// Copyright 2020 The Reed Developers
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

package vmcommon

const (
	OpPushData32 = 0x40
	OpPushData64 = 0x41

	OpDup         = 0x76
	OpHash256     = 0xaa
	OpEqualVerify = 0x88
	OpCheckSig    = 0xac
)
