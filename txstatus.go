package goethx

import "fmt"

// TxStatus represent the status of the tx
type TxStatus int

const (
	// TxNil is used when status is empty
	TxNil TxStatus = 0
	// TxPending is used when the tx is still waiting to be mined
	TxPending TxStatus = 10
	// TxError is used when the tx is not successful
	TxError TxStatus = 100
	// TxNotFound is used when no tx has been found
	TxNotFound TxStatus = 101
	// TxTimeOut is used when tx has not been included in a block
	// within the polling timeout
	TxTimeOut TxStatus = 102
	// TxSuccessNotEnoughBlocks is when tx has been included
	// but not eough blocks have been appended after
	TxSuccessNotEnoughBlocks = 103
	// TxSuccess is used when tx is successfully included into a block
	TxSuccess TxStatus = 1000
)

// String fits stringer interface
func (txs TxStatus) String() string {
	switch txs {
	case TxNil:
		return "no tx status"
	case TxPending:
		return "tx pending"
	case TxError:
		return "tx error"
	case TxTimeOut:
		return "tx didn't not succeed within time"
	case TxSuccessNotEnoughBlocks:
		return "tx was successful but not enough blocks added before timeout"
	case TxSuccess:
		return "tx successful"
	default:
		return fmt.Sprintf("status %d doesn't exist", txs)
	}
}
