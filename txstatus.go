package goethx

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
	// TxSuccess is used when tx is successfully included into a block
	TxSuccess TxStatus = 1000
)
