package common

type TxRequest struct {
	TxId   int64
	Flag   string
	Action string
	Key    string
	Val    []byte
}
