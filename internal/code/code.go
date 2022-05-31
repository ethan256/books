package code

import (
	_ "embed"
)

//go:embed code.go
var ByteCodeFile []byte

const (
	ServerError        = 10101
	TooManyRequests    = 10102
	ParamBindError     = 10103
	AuthorizationError = 10104
	UrlSignError       = 10105
	CacheSetError      = 10106
	CacheGetError      = 10107
	CacheDelError      = 10108
	CacheNotExist      = 10109
	ResubmitError      = 10110
	RBACError          = 10113
	RedisConnectError  = 10114
	MySQLConnectError  = 10115
	WriteConfigError   = 10116

	BookNotFoundError     = 20101
	BookKindNotExistError = 20101
	BookSaveError         = 20101
	BookUpdateError       = 20101

	UserNotFoundError = 20201
)

func Text(code int) string {
	return zhCNText[code]
}
