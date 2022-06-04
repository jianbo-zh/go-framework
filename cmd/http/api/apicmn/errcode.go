package apicmn

const (
	ErrUnknown  = 5000 // 未知错误
	ErrAuthFail = 6000 // token 认证失败
	ErrParamErr = 6300 // 参数错误

)

const (
	ErrUploadInit   = 6410 + iota // 上传初始化错误
	ErrUploadChunk                // 上传失败
	ErrUploadFinish               // 上传完成失败
)
