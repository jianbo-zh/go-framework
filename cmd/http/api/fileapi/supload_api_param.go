package fileapi

type SgUploadParams struct {
	Hash     string `json:"hash" form:"hash"`         // 文件哈希值(md5)，小写
	FileName string `json:"fileName" form:"fileName"` // 文件名称，例如： hell.mp4
	FileSize int64  `json:"fileSize" form:"fileSize"` // 文件大小（字节）
}
