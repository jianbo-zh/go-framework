package gencidtask

import (
	"goframework/cmd/http/service/filesvc"
	"time"
)

func largeFileListen() {
	lastId := int64(0)
	for {
		// 防止请求过快
		time.Sleep(15 * time.Second)
		checkLog.Infof("[al] lastId: %d", lastId)

		pendingHashs, err := filesvc.GetPendingHashs(lastId, 1000, filesvc.TypeLargeFile)
		if err != nil {
			checkLog.Errorf("[al] get pending hashs error: %s", err.Error())
			continue
		}

		if len(pendingHashs) == 0 {
			lastId = 0

		} else {
			for _, phash := range pendingHashs {
				// lastId.
				lastId = phash.Id

				mfi, err := filesvc.GetHashMfi(phash.Hash)
				if err != nil {
					checkLog.Errorf("[al] get hashs mfi error: %s", err.Error())
					filesvc.DeletePendingHash(phash.Id)
					continue
				}

				if mfi.FileCid != "" {
					filesvc.UpdatePendingHashSuccess(phash.Hash, mfi.FileCid)
					continue
				}

				err = filesvc.UpdatePendingHashInQueue(phash.Id)
				if err != nil {
					checkLog.Errorf("[al] push task to queue but updata db error: %s, id: %d, hash: %s", err.Error(), phash.Id, phash.Hash)

				} else {
					checkLog.Infof("[al] push task to queue, id: %d, hash: %s", phash.Id, phash.Hash)
				}

				largeFileChannel <- mfi
			}
		}
	}
}

func normalFileListen() {

	lastId := int64(0)
	for {
		// 防止请求过快
		time.Sleep(15 * time.Second)
		checkLog.Infof("[am] lastId: %d", lastId)

		pendingHashs, err := filesvc.GetPendingHashs(lastId, 100, filesvc.TypeNormalFile)
		if err != nil {
			checkLog.Errorf("[am] get pending hashs error: %s", err.Error())
			continue
		}

		if len(pendingHashs) == 0 {
			lastId = 0

		} else {
			for _, phash := range pendingHashs {
				// lastId.
				lastId = phash.Id

				mfi, err := filesvc.GetHashMfi(phash.Hash)
				if err != nil {
					checkLog.Errorf("[am] get hashs mfi error: %s", err.Error())
					filesvc.DeletePendingHash(phash.Id)
					continue
				}

				err = filesvc.UpdatePendingHashInQueue(phash.Id)
				if err != nil {
					checkLog.Errorf("[am] push task to queue but updata db error: %s, id: %d, hash: %s", err.Error(), phash.Id, phash.Hash)

				} else {
					checkLog.Infof("[am] push task to queue, id: %d, hash: %s", phash.Id, phash.Hash)
				}

				normalFileChannel <- mfi
			}
		}
	}
}

func smallFileListen() {

	lastId := int64(0)
	for {
		// 防止请求过快
		time.Sleep(15 * time.Second)
		checkLog.Infof("[as] lastId: %d", lastId)

		pendingHashs, err := filesvc.GetPendingHashs(lastId, 100, filesvc.TypeSmallFile)
		if err != nil {
			checkLog.Errorf("[as] get pending hashs error: %s", err.Error())
			continue
		}

		if len(pendingHashs) == 0 {
			lastId = 0

		} else {
			for _, phash := range pendingHashs {
				// lastId.
				lastId = phash.Id

				mfi, err := filesvc.GetHashMfi(phash.Hash)
				if err != nil {
					checkLog.Errorf("[as] get hashs mfi error: %s", err.Error())
					filesvc.DeletePendingHash(phash.Id)
					continue
				}

				err = filesvc.UpdatePendingHashInQueue(phash.Id)
				if err != nil {
					checkLog.Errorf("[as] push task to queue but updata db error: %s, id: %d, hash: %s", err.Error(), phash.Id, phash.Hash)

				} else {
					checkLog.Infof("[as] push task to queue, id: %d, hash: %s", phash.Id, phash.Hash)
				}

				smallFileChannel <- mfi
			}
		}
	}
}
