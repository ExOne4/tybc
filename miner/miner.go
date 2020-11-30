// Copyright 2020 The Reed Developers
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

package miner

import (
	"bytes"
	"github.com/reed/consensus/pow"
	"github.com/reed/errors"
	"github.com/reed/log"
	"github.com/reed/types"
	"strconv"
	"sync"
)

var (
	startErr = errors.New("miner start error")
)

const (
	maxTries = ^uint64(0)
)

type Miner struct {
	sync.Mutex
	working          bool
	blockReceptionCh <-chan *types.Block
	blockSendCh      chan<- *types.Block
	stopWorkCh       <-chan struct{}
}

func NewMiner(submitCh <-chan *types.Block) *Miner {
	return &Miner{
		working:          false,
		blockReceptionCh: submitCh,
	}
}

func (m *Miner) Start() error {
	m.Lock()
	defer m.Unlock()

	if m.working {
		return errors.Wrap(startErr, "miner has started.")
	}

	go m.work()

	return nil
}

func (m *Miner) work() {
	var block types.Block
	//calc difficulty
	block.BigNumber = pow.GetNextDifficulty(&block)

	for {
		born, stop := m.generateBlock(&block)
		if born {
			//broadcast new blockmanager
		}
		if stop {
			break
		}
	}
}

func (m *Miner) generateBlock(block *types.Block) (born bool, stop bool) {
	extraNonce := uint64(0)
loop:
	for {
		select {
		case rblock := <-m.blockReceptionCh:
			log.Logger.Infof("Received a blockmanager from blockReception channel.id=%x", rblock.GetHash())
			// receive a new blockmanager from remote node
			// blockmanager = fetch laest blockmanager
			block = rblock
			born = true
			break loop
		case <-m.stopWorkCh:
			log.Logger.Info("Received a stop single,stop miner...")
			stop = true
			break loop
		default:
			//just for no blockmanager,do nothing
		}

		if pow.CheckProofOfWork(block.BigNumber, block.GetHash()) {
			born = true
			break loop
		} else {
			if block.Nonce == maxTries {
				//reset nonce
				block.Nonce = 0

				//change coinbase tx's scriptSig and continue
				extraNonce++
				m.incrementExtraNonce(extraNonce, block)
			} else {
				block.Nonce++
			}
		}
	}
	return
}

func (m *Miner) incrementExtraNonce(extraNonce uint64, cblock *types.Block) {
	txs := *cblock.Transactions
	txs[0].TxInput[0].ScriptSig = bytes.Join([][]byte{txs[0].TxInput[0].ScriptSig, []byte(strconv.FormatUint(extraNonce, 10))}, []byte{})
}
