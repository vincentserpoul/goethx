# Transaction listener for ethereum

[![Documentation](https://godoc.org/github.com/vincentserpoul/goethx?status.svg)](http://godoc.org/github.com/vincentserpoul/goethx) [![Go Report Card](https://goreportcard.com/badge/github.com/vincentserpoul/goethx)](https://goreportcard.com/report/github.com/vincentserpoul/goethx) [![Coverage Status](https://coveralls.io/repos/github/vincentserpoul/goethx/badge.svg?branch=master)](https://coveralls.io/github/vincentserpoul/goethx?branch=master) [![CircleCI](https://circleci.com/gh/vincentserpoul/goethx.svg?style=svg)](https://circleci.com/gh/vincentserpoul/goethx) [![Maintainability](https://api.codeclimate.com/v1/badges/937d15e44061eeb32877/maintainability)](https://codeclimate.com/github/vincentserpoul/goethx/maintainability)

## Usage

```golang
txm, err := goethx.NewTxMgr(&logger.Log, ethClient, 6, 15* time.Minute, 1 * time.Second)
if err != nil {
    log.Fatalf("goethx.NewTxMgr: %v", err)
}
chTx := make(chan goethx.TxMsg)
go txm.MonitorTx(common.HexToHash("0x123"), chTx)
msg := <- chTx
if msg.Err != nil {
    log.Fatalf("goethx.MonitorTx(%s): %v", common.HexToHash("0x123").String(), msg.Err)
}
```
