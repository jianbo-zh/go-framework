package gencidtask

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"time"

	"github.com/gabriel-vasile/mimetype"
	"github.com/jianbo-zh/go-errors"
	goimg "github.com/penglonghua/go-image"
)

const (
	// MaxLimitSize = int64(20971520) // 20MB
	MaxLimitSize = 500 * 1024 * 1024 // 500MB
)

// getFileThumbnailAndCover 获取文件的缩略图和封面图cid
func getThumbnailAndCover(filePath, writeDir, hash string) (
	fileMimeType, fileExt, thumbnailFile, coverFile string, err error) {

	defer func() {
		if err != nil {
			imgLog.Errorf("get thumbnail and cover error: %s", err.Error())
		}
	}()

	st := time.Now()

	imgLog.Infof("handle thumbnail start, hash: %s, file: %s", hash, filePath)

	var pBytes []byte
	pBytes, err = getFileHeader(filePath)
	if err != nil {
		err = fmt.Errorf("get file header, hash: %s, error: %w", hash, err)
		return
	}

	mimeType := mimetype.Detect(pBytes)

	fileMimeType = mimeType.String()
	fileExt = mimeType.Extension()
	imgLog.Infof("get mimetype, hash: %s, mime: %s", hash, fileMimeType)

	var isMatchImage bool
	isMatchImage, err = regexp.Match("^image", []byte(fileMimeType))
	if err != nil {
		err = fmt.Errorf("mimetype detectFile, hash: %s, error: %w", hash, err)
		return
	}

	var isMatchVideo bool
	isMatchVideo, err = regexp.Match("^video", []byte(fileMimeType))
	if err != nil {
		err = fmt.Errorf("mimetype detectFile1, hash: %s, error: %w", hash, err)
		return
	}

	// not handle other
	if !isMatchImage && !isMatchVideo {
		imgLog.Infof("not match image and video, hash: %s", hash)
		return
	}

	if isMatchVideo {
		imgLog.Infof("match video, hash: %s", hash)
		var fInfo fs.FileInfo
		fInfo, err = os.Stat(filePath)
		if err != nil {
			err = fmt.Errorf("os stat file, hash: %s, error: %w", hash, err)
			return
		}

		if fInfo.Size() > MaxLimitSize {
			var tmpFilePath string
			tmpFilePath, err = getFileChunk(filePath, writeDir)
			if err != nil {
				err = fmt.Errorf("get file chunk, hash: %s, error: %w", hash, err)
				return
			}
			defer os.Remove(tmpFilePath)

			filePath = tmpFilePath
			imgLog.Infof("chunk video, hash: %s, file: %s", hash, filePath)
		}
	}

	imgLog.Infof("generate thumbnail start, hash: %s, file: %s, dir: %s", hash, filePath, writeDir)
	ir, err := goimg.ImageAndSave(filePath, writeDir)
	if err != nil {
		err = fmt.Errorf("goimg.ImageAndSave hash: %s, error: %w", hash, err)
		return
	}

	imgLog.Infof("handle thumbnail start, hash: %s, file: %s, sec: %f", hash, filePath, time.Since(st).Seconds())

	return fileMimeType, fileExt, ir.ThumbnailImgPath, ir.CoverImgPath, nil
}

// getFileHeader 获取文件头部切片
func getFileHeader(filePath string) ([]byte, error) {
	oFile, err := os.Open(filePath)
	if err != nil {
		return nil, errors.New("open file error").With(errors.Inner(err))
	}
	defer oFile.Close()

	// 读 256KB，来判断 mimeType
	pBytes := make([]byte, 262144)
	pn, err := oFile.ReadAt(pBytes, 0)
	if err != nil && !errors.Is(err, io.EOF) {
		return nil, errors.New("read file header error").With(errors.Inner(err))
	}

	return pBytes[0:pn], nil
}

// getFileChunk 获取文件头部来
func getFileChunk(filePath, writeDir string) (string, error) {
	oFile, err := os.Open(filePath)
	if err != nil {
		return "", errors.New("open file error").With(errors.Inner(err))
	}
	defer oFile.Close()

	tmpFilePath := filepath.Join(writeDir, "tmp")
	tFile, err := os.OpenFile(tmpFilePath, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0755)
	if err != nil {
		return "", errors.New("open tmp file error").With(errors.Inner(err))
	}
	defer tFile.Close()

	_, err = io.CopyN(tFile, oFile, MaxLimitSize)
	if err != nil {
		return "", errors.New("copy max limit file error").With(errors.Inner(err))
	}

	return tmpFilePath, nil
}
