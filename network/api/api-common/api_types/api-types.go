package api_types

import (
	"pandora-pay/blockchain/blocks/block"
	"pandora-pay/blockchain/info"
	"pandora-pay/blockchain/transactions/transaction"
	"pandora-pay/helpers"
)

type APIBlockWithTxs struct {
	Block *block.Block       `json:"block"`
	Txs   []helpers.HexBytes `json:"txs"`
}

type APIBlockchain struct {
	Height            uint64 `json:"height"`
	Hash              string `json:"hash"`
	PrevHash          string `json:"prevHash"`
	KernelHash        string `json:"kernelHash"`
	PrevKernelHash    string `json:"prevKernelHash"`
	Timestamp         uint64 `json:"timestamp"`
	TransactionsCount uint64 `json:"transactions"`
	Target            string `json:"target"`
	TotalDifficulty   string `json:"totalDifficulty"`
}

type APIBlockchainSync struct {
	SyncTime uint64 `json:"syncTime"`
}

type APITransaction struct {
	Tx           *transaction.Transaction `json:"tx,omitempty"`
	TxSerialized helpers.HexBytes         `json:"serialized,omitempty"`
	Mempool      bool                     `json:"mempool,omitempty"`
	Info         *info.TxInfo             `json:"info,omitempty"`
}

type APIMempoolAnswer struct {
	Count  int                `json:"count"`
	Hashes []helpers.HexBytes `json:"hashes"`
}

type APISubscriptionNotification struct {
	Key  helpers.HexBytes `json:"key"`
	Data helpers.HexBytes `json:"tx"`
}

type APIAccountTxs struct {
	Next uint64             `json:"next"`
	Txs  []helpers.HexBytes `json:"txs"`
}