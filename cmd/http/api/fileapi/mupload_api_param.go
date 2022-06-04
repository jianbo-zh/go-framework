package fileapi

type MpUploadInitParams struct {
	Hash     string `json:"hash"`     // 文件哈希值(sha1)，小写
	FileName string `json:"fileName"` // 文件名称，例如： hell.mp4
	FileSize int64  `json:"fileSize"` // 文件大小（字节）
}

type MpUploadChunkParams struct {
	Hash       string `json:"hash" form:"hash"`             // 文件哈希值(sha1)，小写
	ChunkIndex int64  `json:"chunkIndex" form:"chunkIndex"` // 分片序号
	ChunkSize  int64  `json:"chunkSize" form:"chunkSize"`   // 分片大小
}

type MpUploadFinishParams struct {
	Hash string `json:"hash"` // 文件哈希值(sha1)，小写
}

type MpUploadCheckParams struct {
	Hashs []string `json:"hashs"`
}

type MpUploadCheckInfo struct {
	FileName      string `json:"fileName"`      // 文件名称，例如： hell.mp4
	FileSize      int64  `json:"fileSize"`      // 文件大小（字节）
	FileCid       string `json:"fileCid"`       // 图片、视频或文件CID
	ThumbnailCid  string `json:"thumbnailCid"`  // 图片、视频缩略图CID
	CoverCid      string `json:"coverCid"`      // 视频封面CID
	FileMimeType  string `json:"fileMimeType"`  // MIME-Type
	FileExtension string `json:"fileExtension"` // 后缀
	FileHash      string `json:"fileHash"`      // 文件Hash
	Error         string `json:"error"`         // 错误消息
}
