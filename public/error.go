package public

import (
	"crypto/md5"
	"io"
	"math/big"
)

type Error struct {
	ErrorMessage string `json:"error"`
	Code         int    `json:"code"`
}

func NewError(string2 string) Error {
	h := md5.New()
	io.WriteString(h, string2)
	b := h.Sum(nil)
	return NewErrorCode(string2, int(big.NewInt(0).SetBytes(b[:4]).Int64()))
}

func NewErrorCode(string2 string, code int) Error {
	return Error{ErrorMessage: string2, Code: code}
}

func (e Error) Error() string {
	return e.ErrorMessage
}
