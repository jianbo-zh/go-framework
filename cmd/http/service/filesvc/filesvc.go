package filesvc

import (
	"encoding/json"
	"errors"
	"fmt"
	"goframework/pkg/mmux"
	"io"
	"os"
	"path/filepath"

	"github.com/jianbo-zh/go-config"
)

const (
	StatusHashPending = 1
	StatusHashSuccess = 2
	StatusHashInQueue = 3
	StatusHashFailed  = 4
)

const (
	TypeSmallFile  = 1 // <=50MB
	TypeNormalFile = 2 // >50MB <=500MB
	TypeLargeFile  = 3 // >500MB
)

func GetHashMfi(phash string) (mfi MpUploadFileInfo, err error) {

	fileUploadDir, err := GetHashUploadDir(phash)
	if err != nil {
		err = fmt.Errorf("get upload dir error: %w", err)
		return
	}

	// TODO: 互斥锁
	mmux.Lock(phash)
	defer mmux.Unlock(phash)

	fi, err := os.OpenFile(filepath.Join(fileUploadDir, "stat.info"), os.O_CREATE|os.O_RDWR, 0755)
	if err != nil {
		err = fmt.Errorf("open file save dir error: %w", err)
		return
	}
	defer fi.Close()

	bytes, err := io.ReadAll(fi)
	if err != nil {
		err = fmt.Errorf("read file info error: %w", err)
		return
	}

	err = json.Unmarshal(bytes, &mfi)
	if err != nil {
		err = fmt.Errorf("json unmarshal file info error: %w", err)
		return
	}

	return mfi, nil
}

func GetHashUploadDir(hash string) (string, error) {

	storageDir := config.GetString("storage.dir", "")
	if storageDir == "" {
		return "", errors.New("upload dir not set")
	}

	return filepath.Join(storageDir, hash[0:2], hash[2:4], hash[:]), nil
}
