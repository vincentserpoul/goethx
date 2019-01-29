package goethx

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
)

// TxMgrI represent the interface to listen to eth tx
type TxMgrI interface {
	MonitorTx(
		ctx context.Context,
		txH common.Hash,
		chTx chan<- TxMsg,
	) (bool, error)
}
