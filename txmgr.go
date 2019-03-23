package goethx

import (
	"context"
	"fmt"
	"math/big"
	"sync"
	"time"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/sirupsen/logrus"
)

// ethCli represents the interface we need to assess a tx
type ethCli interface {
	HeaderByNumber(ctx context.Context, number *big.Int) (*types.Header, error)
	ethereum.TransactionReader
}

// TxMgr will allow listening to transactions
type TxMgr struct {
	Logger          logrus.StdLogger
	Cli             ethCli
	BlockCountValid int64
	PollingInterval time.Duration
	PollingTimeOut  time.Duration
	TxProcessing    map[common.Hash]bool
	Mutex           *sync.Mutex
}

// NewTxMgr will return an TxListener entity
func NewTxMgr(
	lg logrus.StdLogger,
	cli ethCli,
	bcv int64,
	pi time.Duration,
	pto time.Duration,
) *TxMgr {
	return &TxMgr{
		Logger:          lg,
		Cli:             cli,
		BlockCountValid: bcv,
		PollingInterval: pi,
		PollingTimeOut:  pto,
		TxProcessing:    make(map[common.Hash]bool),
		Mutex:           &sync.Mutex{},
	}
}

// TxMsg is the message sent through the channel
type TxMsg struct {
	Hash   common.Hash
	Status TxStatus
	Err    error
}

// MonitorTx will monitor a tx until success, error or timeout
func (txm *TxMgr) MonitorTx(
	ctx context.Context,
	txH common.Hash,
	chTx chan TxMsg,
) {
	timeout := time.After(txm.PollingTimeOut)

	var succTxBlock *big.Int
	t := time.NewTicker(txm.PollingInterval)

	for c := t.C; ; {
		var errT error
		var txS TxStatus
		txS, succTxBlock, errT = txm.checkTx(ctx, txH, succTxBlock)
		if errT != nil {
			chTx <- TxMsg{
				Hash:   txH,
				Status: TxNil,
				Err: fmt.Errorf(
					"MonitorTx(%s): %v - tx status %s",
					txH.String(), errT, txS.String(),
				),
			}
			return
		}
		if txS == TxSuccess {
			chTx <- TxMsg{
				Hash:   txH,
				Status: TxSuccess,
				Err:    nil,
			}
			return
		}

		select {
		case <-c:
			continue
		case <-timeout:
			chTx <- TxMsg{
				Hash:   txH,
				Status: TxTimeOut,
				Err: fmt.Errorf(
					"MonitorTx(%s): time out after %s with status %s",
					txH.String(),
					txm.PollingTimeOut,
					txS.String(),
				),
			}
			return
		}
	}
}

func (txm *TxMgr) checkTx(
	ctx context.Context,
	txH common.Hash,
	succTxBlock *big.Int,
) (TxStatus, *big.Int, error) {
	if txm.checkTxStatus(ctx, txH) == TxSuccess {
		if succTxBlock == nil {
			hdrSuccessBlock, errH := txm.Cli.HeaderByNumber(ctx, nil)
			if errH != nil {
				return TxNil, nil,
					fmt.Errorf("checkTx(%s): %v", txH.String(), errH)
			}
			succTxBlock = hdrSuccessBlock.Number
		}
		isEnough, errB := txm.enoughBlocksSince(ctx, succTxBlock)
		if errB != nil {
			return TxNil, succTxBlock,
				fmt.Errorf("checkTx(%s): %v", txH.String(), errB)
		}
		if isEnough {
			return TxSuccess, succTxBlock, nil
		} else {
			return TxSuccessNotEnoughBlocks, succTxBlock, nil
		}
	}
	return TxPending, nil, nil
}

func (txm *TxMgr) enoughBlocksSince(
	ctx context.Context,
	bn *big.Int,
) (bool, error) {
	// Get current block number
	hdr, errH := txm.Cli.HeaderByNumber(ctx, nil)
	if errH != nil {
		return false, fmt.Errorf("enoughBlocksSince(%s): %v", bn, errH)
	}

	var expBN big.Int
	if expBN.Add(bn, big.NewInt(txm.BlockCountValid)).Cmp(hdr.Number) <= 0 {
		return true, nil
	}
	return false, nil
}

// checkTxStatus returns the status if the tx
func (txm *TxMgr) checkTxStatus(
	ctx context.Context,
	txH common.Hash,
) TxStatus {
	_, isPending, errTbH := txm.Cli.TransactionByHash(ctx, txH)
	if errTbH == ethereum.NotFound {
		return TxNotFound
	}
	if errTbH != nil {
		return TxError
	}
	if isPending {
		return TxPending
	}

	rec, errTxR := txm.Cli.TransactionReceipt(ctx, txH)
	if errTxR == ethereum.NotFound {
		return TxNotFound
	}
	if errTxR != nil {
		return TxError
	}
	// if incomplete receipt
	if rec == nil {
		return TxPending
	}
	// if rec.Status == types.ReceiptStatusFailed {
	// 	return TxError
	// }
	if rec.Status == types.ReceiptStatusSuccessful {
		return TxSuccess
	}

	return TxPending
}
