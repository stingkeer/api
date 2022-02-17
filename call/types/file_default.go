package types

import (
	"gitee.com/fast_api/api/def"
	"gitee.com/fast_api/api/log"
	"mime/multipart"
	"reflect"
)

type FileType struct{}

func (f FileType) Mapper(param def.ParamWarp) reflect.Value {
	newT := reflect.New(param.PTyp)
	reader, err := param.MultipartReader()
	if err != nil {
		log.Error(err)
	}
	newT.Elem().Set(reflect.ValueOf(*reader))
	return newT.Elem()
}

func (f FileType) Register() []reflect.Type {
	return []reflect.Type{reflect.TypeOf((*multipart.Reader)(nil)).Elem()}
}
