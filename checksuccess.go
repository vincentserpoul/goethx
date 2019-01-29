package transaction

import (
	"context"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// checkTxSuccess returns isTxSuccessful, isTxValid, error
func checkTxSuccess(
	ctx context.Context,
	cli ethereum.TransactionReader,
	txH common.Hash,
) (bool, error) {
	rec, errTxR := cli.TransactionReceipt(ctx, txH)
	if errTxR == ethereum.NotFound {
		return false, nil
	}
	// if incomplete receipt
	if rec == nil {
		return false, nil
	}
	// if rec.Status == types.ReceiptStatusFailed {
	// 	return false, nil
	// }
	if rec.Status == types.ReceiptStatusSuccessful {
		return true, nil
	}

	return false, nil
}
