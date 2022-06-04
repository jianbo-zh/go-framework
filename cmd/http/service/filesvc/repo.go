package filesvc

import (
	"goframework/cmd/http/repo/filerepo"
)

func GetPendingHashs(lastId, limit int64, ftype int) (hashs []filerepo.PendingHash, err error) {
	return filerepo.GetPendingHashs(lastId, limit, ftype)
}

func UpdateQueueHashPending() error {
	return filerepo.UpdateQueueHashPending()
}

func UpdatePendingHashInQueue(id int64) error {
	return filerepo.UpdatePendingHashInQueue(id)
}

func UpdatePendingHashFailed(hash string) error {
	return filerepo.UpdatePendingHashFailed(hash)
}

func UpdatePendingHashSuccess(hash, cid string) error {
	return filerepo.UpdatePendingHashSuccess(hash, cid)
}

func DeletePendingHash(id int64) error {
	return filerepo.DeletePendingHash(id)
}

func AddPendHash(phash string, fsize int64) error {
	return filerepo.AddPendHash(phash, fsize)
}
