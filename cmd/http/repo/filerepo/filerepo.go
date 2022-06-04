package filerepo

import (
	"database/sql"
	"fmt"
	"goframework/cmd/http/dao"
	"goframework/cmd/http/model"
	"time"

	"github.com/jianbo-zh/go-log"
)

const (
	TypeSmallFile  = 1 // <=50MB
	TypeNormalFile = 2 // >50MB <=500MB
	TypeLargeFile  = 3 // >500MB
)

var repoLogger = log.Logger("repo")

func GetPendingHashs(lastId, limit int64, ftype int) (hashs []PendingHash, err error) {

	var rows *sql.Rows

	if ftype > 0 {
		rows, err = dao.DB().Raw("select id, hash from ipfsqueue where id > ? and status = ? and ftype = ? order by id asc limit ?",
			lastId, model.StatusHashPending, ftype, limit,
		).Rows()
	} else {
		rows, err = dao.DB().Raw("select id, hash from ipfsqueue where id > ? and status = ? order by id asc limit ?",
			lastId, model.StatusHashPending, limit,
		).Rows()
	}

	if err != nil {
		err = fmt.Errorf("get pending queue error: %w", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var phash PendingHash
		err = rows.Scan(&phash.Id, &phash.Hash)
		if err != nil {
			err = fmt.Errorf("rows scan error: %w", err)
			return
		}

		hashs = append(hashs, phash)
	}
	err = rows.Err()
	if err != nil {
		err = fmt.Errorf("rows error: %w", err)
		return
	}

	return hashs, nil
}

func UpdateQueueHashPending() error {
	err := dao.DB().Exec(
		"update ipfsqueue set status = ?, utime = ? where status = ? or status = ?",
		model.StatusHashPending, time.Now().Unix(), model.StatusHashInQueue, model.StatusHashFailed).Error

	if err != nil {
		return fmt.Errorf("db update queue hash to pending error: %w", err)
	}

	return nil
}

func UpdatePendingHashInQueue(id int64) error {
	err := dao.DB().Exec("update ipfsqueue set status = ?, utime = ? where id = ?", model.StatusHashInQueue, time.Now().Unix(), id).Error
	if err != nil {
		return fmt.Errorf("db update pending hash in queue error: %w", err)
	}

	return nil
}

func UpdatePendingHashFailed(hash string) error {
	res := dao.DB().Exec("update ipfsqueue set status = ?, utime = ? where hash = ?", model.StatusHashFailed, time.Now().Unix(), hash)
	if res.Error != nil {
		return fmt.Errorf("db update failed hash error: %w", res.Error)
	}

	if res.RowsAffected == 0 {
		repoLogger.Errorf("rows affected zero1, hash: %s", hash)
	}

	return nil
}

func UpdatePendingHashSuccess(hash, cid string) error {
	res := dao.DB().Exec("update ipfsqueue set status = ?, cid = ?, utime = ? where hash = ?", model.StatusHashSuccess, cid, time.Now().Unix(), hash)
	if res.Error != nil {
		return fmt.Errorf("db update pending hash error: %w", res.Error)
	}

	if res.RowsAffected == 0 {
		repoLogger.Errorf("rows affected zero2, hash: %s", hash)
	}

	return nil
}

func DeletePendingHash(id int64) error {
	err := dao.DB().Exec("delete from ipfsqueue where id = ?", id).Error
	if err != nil {
		return fmt.Errorf("delete pending hash error: %w", err)
	}

	return nil
}

func AddPendHash(phash string, fsize int64) error {
	if phash == "" {
		return fmt.Errorf("phash can't is empty string")
	}

	// check if exists
	var id int64
	err := dao.DB().Raw("select id from ipfsqueue where hash = ? and status != ?", phash, model.StatusHashFailed).Scan(&id).Error
	if err != nil {
		return fmt.Errorf("get ipfsqueue error: %w", err)
	}

	if id > 0 {
		return nil
	}

	var ftype int
	switch {
	case fsize > 524288000:
		// 大于500MB
		ftype = TypeLargeFile
	case fsize <= 52428800:
		// 小于等于50MB
		ftype = TypeSmallFile
	default:
		ftype = TypeNormalFile
	}

	ptime := time.Now().Unix()
	err = dao.DB().Exec("insert into ipfsqueue(hash, cid, ftype, fsize, status, ptime, utime) values(?, ?, ?, ?, ?, ?, ?)",
		phash, "", ftype, fsize, model.StatusHashPending, ptime, ptime).Error
	if err != nil {
		return fmt.Errorf("add ipfsqueue error: %w, hash: %s", err, phash)
	}

	return nil
}
