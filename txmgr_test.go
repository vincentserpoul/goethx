package goethx

import (
	"context"
	"errors"
	"log"
	"math/big"
	"testing"
	"time"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/sirupsen/logrus"
)

func TestNewTxMgr(t *testing.T) {
	tests := []struct {
		name   string
		logger logrus.StdLogger
		cli    ethCli
		bcv    int64
		pi     time.Duration
		pto    time.Duration
		want   *TxMgr
	}{
		{
			name:   "working new Tx manager",
			logger: logrus.New(),
			cli:    &ethCliMock{},
			bcv:    1,
			pi:     1 * time.Second,
			pto:    1 * time.Second,
			want: &TxMgr{
				BlockCountValid: 1,
				PollingInterval: 1 * time.Second,
				PollingTimeOut:  1 * time.Second,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewTxMgr(tt.logger, tt.cli, tt.bcv, tt.pi, tt.pto)
			if got.Logger == nil {
				t.Error("NewTxMgr(): empty logger")
			}
			if got.Cli == nil {
				t.Error("NewTxMgr(): empty ethcli")
			}
			if tt.want.BlockCountValid != got.BlockCountValid {
				t.Errorf("NewTxMgr() = %v, want %v",
					got.BlockCountValid, tt.want.BlockCountValid,
				)
			}
			if tt.want.PollingInterval != got.PollingInterval {
				t.Errorf("NewTxMgr() = %v, want %v",
					got.PollingInterval, tt.want.PollingInterval,
				)
			}
			if tt.want.PollingTimeOut != got.PollingTimeOut {
				t.Errorf("NewTxMgr() = %v, want %v",
					got.PollingTimeOut, tt.want.PollingTimeOut,
				)
			}
		})
	}
}

func TestTxMgr_MonitorTx(t *testing.T) {
	tests := []struct {
		name            string
		txLocked        bool
		txTimeOut       bool
		isTBHPending    bool
		isTRFailure     bool
		withHBNErr      bool
		tbhErr          error
		trErr           error
		trNoReceipt     bool
		notEnoughBlocks bool
		expectErr       bool
	}{
		{
			name: "usual case, existing successful tx",
		},
		{
			name:       "header by number error",
			withHBNErr: true,
			expectErr:  true,
		},
		{
			name:      "tx locked error",
			txLocked:  true,
			expectErr: true,
		},
		{
			name:      "tx timeout",
			txTimeOut: true,
			expectErr: true,
		},
		{
			name:         "tx always pending",
			isTBHPending: true,
			expectErr:    true,
		},
		{
			name:        "receipt failure",
			isTRFailure: true,
			expectErr:   true,
		},
		{
			name:      "tbh hash not found",
			tbhErr:    ethereum.NotFound,
			expectErr: true,
		},
		{
			name:      "tbh error",
			tbhErr:    errors.New("error tbh"),
			expectErr: true,
		},
		{
			name:      "tr receipt error",
			trErr:     errors.New("error tbh"),
			expectErr: true,
		},
		{
			name:      "tr receipt not found",
			trErr:     ethereum.NotFound,
			expectErr: true,
		},
		{
			name:        "tr receipt not found",
			trNoReceipt: true,
			expectErr:   true,
		},
		{
			name:            "not enough blocks to validate",
			notEnoughBlocks: true,
			expectErr:       true,
		},
	}

	txH := common.HexToHash("0x1234566890")
	ctx := context.Background()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			timeout := 1 * time.Millisecond
			if tt.txTimeOut {
				timeout = 0
				tt.isTBHPending = true
			}
			bcv := int64(1)
			if tt.notEnoughBlocks {
				bcv = 10000000000
			}
			txMon := NewTxMgr(
				&log.Logger{},
				&ethCliMock{
					txH:          txH,
					isTBHPending: tt.isTBHPending,
					isTRFailure:  tt.isTRFailure,
					withHBNErr:   tt.withHBNErr,
					tbhErr:       tt.tbhErr,
					trErr:        tt.trErr,
					trNoReceipt:  tt.trNoReceipt,
				},
				bcv,
				100*time.Microsecond,
				timeout,
			)
			chTx := make(chan TxMsg)
			if tt.txLocked {
				_ = txMon.lock(ctx, txH)
			}
			go txMon.MonitorTx(ctx, txH, chTx)
			msg := <-chTx
			ctx.Done()
			if (msg.Err != nil) != tt.expectErr {
				t.Errorf("Expected error %t, got %v", tt.expectErr, msg.Err)
			}
		})
	}
}

type ethCliMock struct {
	txH          common.Hash
	isTBHPending bool
	isTRFailure  bool
	blockNum     int64
	withHBNErr   bool
	tbhErr       error
	trErr        error
	trNoReceipt  bool
}

// HeaderByNumber is returning a Header
func (ecm *ethCliMock) HeaderByNumber(
	ctx context.Context,
	number *big.Int,
) (*types.Header, error) {
	if ecm.withHBNErr {
		return nil, errors.New("hbn fail")
	}
	ecm.blockNum++
	return &types.Header{
		Number: big.NewInt(ecm.blockNum),
	}, nil
}

// TransactionByHash returns a tx
func (ecm *ethCliMock) TransactionByHash(
	ctx context.Context,
	txHash common.Hash,
) (tx *types.Transaction, isPending bool, err error) {
	return nil, ecm.isTBHPending, ecm.tbhErr
}

// TransactionReceipt returns the receipt of a transaction
func (ecm *ethCliMock) TransactionReceipt(
	ctx context.Context,
	txHash common.Hash,
) (*types.Receipt, error) {
	if ecm.trNoReceipt {
		return nil, nil
	}
	trStatus := types.ReceiptStatusSuccessful
	if ecm.isTRFailure {
		trStatus = types.ReceiptStatusFailed
	}
	return &types.Receipt{
		Status: trStatus,
	}, ecm.trErr
}
