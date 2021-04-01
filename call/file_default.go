package call

import (
	"gitee.com/fast_api/api/public"
	"github.com/sirupsen/logrus"
	"mime/multipart"
	"reflect"
)

type FileType struct{}

func (f FileType) Mapper(param public.ParamWarp) reflect.Value {
	newT := reflect.New(param.PTyp)
	reader, err := param.MultipartReader()
	if err != nil {
		logrus.Error(err)
	}
	newT.Elem().Set(reflect.ValueOf(*reader))
	return newT.Elem()
}

func (f FileType) Register() []reflect.Type {
	return []reflect.Type{reflect.TypeOf((*multipart.Reader)(nil)).Elem()}
}
