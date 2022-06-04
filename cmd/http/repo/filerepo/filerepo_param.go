package filerepo

type PendingHash struct {
	Id   int64  `json:"id"`   // id
	Hash string `json:"hash"` // hash
}
