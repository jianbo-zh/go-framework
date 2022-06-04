package gencidtask

type MpUploadChunkInfo struct {
	ChunkIndex int64 `json:"chunkIndex"` // 分片序号
	ChunkSize  int64 `json:"chunkSize"`  // 分片大小
}

type MpUploadFileInfo struct {
	FileName      string              `json:"fileName"` // 文件名称，例如： hell.mp4
	FileSize      int64               `json:"fileSize"` // 文件大小（字节）
	ChunkSize     int64               `json:"chunkSize"`
	UploadChunks  []MpUploadChunkInfo `json:"uploadChunks"`
	FileCid       string              `json:"fileCid"`       // 图片、视频或文件CID
	ThumbnailCid  string              `json:"thumbnailCid"`  // 图片、视频缩略图CID
	CoverCid      string              `json:"coverCid"`      // 视频封面CID
	UploadState   int64               `json:"uploadState"`   // 上传状态 0-未上传完 1-已上传完
	FileMimeType  string              `json:"fileMimeType"`  // MIME-Type
	FileExtension string              `json:"fileExtension"` // 后缀
	FileHash      string              `json:"fileHash"`      // 文件Hash
	QueueExpireTs int64               `json:"queueExpireTs"` // 队列处理过期时间
}
