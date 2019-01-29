package transaction

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

// PublicBlockCountValid represents after how many blocks
// a tx is considered valid
const PublicBlockCountValid = 3

// PublicPollingInterval represents the number of seconds between each tx check
const PublicPollingInterval = 2000

// PublicPollingTimeOut represents the number of seconds
// after we stop looking at tx
const PublicPollingTimeOut = 7200000

// WaitForSuccessfulTx will wait for a tx to be successful
func WaitForSuccessfulTx(
	ctx context.Context,
	cli *ethclient.Client,
	txH common.Hash,
	blockCountValid int64,
	pollingInterval time.Duration,
) error {
	var succTxBlock *big.Int
	t := time.NewTicker(pollingInterval)
	for c := t.C; ; {
		if succTxBlock == nil {
			isSuccess, err := checkTxSuccess(ctx, cli, txH)
			if err != nil {
				return fmt.Errorf("WaitForSuccessfulTx(%s): %v",
					txH.String(), err)
			}
			if isSuccess {
				// Get current block number
				hdrSuccess, errH := cli.HeaderByNumber(ctx, nil)
				if errH != nil {
					return fmt.Errorf(
						"WaitForSuccessfulTx(%s): %v",
						txH.String(), errH,
					)
				}
				succTxBlock = hdrSuccess.Number
			}
		} else {
			isEnough, err := enoughBlocksSince(
				ctx,
				cli,
				blockCountValid,
				succTxBlock,
			)
			if err != nil {
				return fmt.Errorf("WaitForSuccessfulTx(%s): %v", txH, err)
			}
			if isEnough {
				return nil
			}
		}
		select {
		case <-c:
			continue
		case <-ctx.Done():
			return fmt.Errorf("WaitForSuccessfulTx(%s): timeout", txH)
		}
	}
}

func enoughBlocksSince(
	ctx context.Context,
	cli *ethclient.Client,
	blockCountValid int64,
	bn *big.Int,
) (bool, error) {
	// Get current block number
	hdr, errH := cli.HeaderByNumber(ctx, nil)
	if errH != nil {
		return false, fmt.Errorf("enoughBlocksSince(%d): %v", bn, errH)
	}

	var expBN big.Int
	if expBN.Add(bn, big.NewInt(blockCountValid)).Cmp(hdr.Number) == 1 {
		return true, nil
	}

	return false, nil
}
