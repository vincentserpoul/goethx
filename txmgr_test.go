package goethx

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/sirupsen/logrus"
)

type ethCliMock struct{}

func (e *ethCliMock) HeaderByNumber(
	ctx context.Context,
	number *big.Int,
) (*types.Header, error) {
	return &types.Header{}, nil
}

func (e *ethCliMock) TransactionByHash(
	ctx context.Context,
	txHash common.Hash,
) (tx *types.Transaction, isPending bool, err error) {
	return &types.Transaction{}, true, nil
}

func (e *ethCliMock) TransactionReceipt(
	ctx context.Context,
	txHash common.Hash,
) (*types.Receipt, error) {
	return &types.Receipt{}, nil
}

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
