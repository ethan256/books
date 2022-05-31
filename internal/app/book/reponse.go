package book

import (
	"fmt"

	"github.com/ethan256/books/internal/code"
)

type Response struct {
	RetCode int         `json:"retcode"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type RequestUUID string

type Failure struct {
	RetCode int    `json:"retcode"`
	Reason  string `json:"reason"`
}

func BuildFailure(retcode int, msg string) *Failure {
	if msg != "" {
		msg = fmt.Sprintf("%s: %s", code.Text(retcode), msg)
	} else {
		msg = code.Text(retcode)
	}

	return &Failure{
		RetCode: retcode,
		Reason:  msg,
	}
}
