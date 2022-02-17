package def

//require param
/**
  //b is require
  GET(func(a int, b def.StringReq) {
		fmt.Println(a, b)
   }, "/require")

*/

type (
	IntReq    int
	Int8Req   int8
	Int16Req  int16
	Int32Req  int32
	Int64Req  int64
	StringReq string
)

func (i IntReq) Int() int {
	return int(i)
}

func (i Int8Req) Int8() int8 {
	return int8(i)
}

func (i Int16Req) Int16() int16 {
	return int16(i)
}

func (i Int8Req) Int32() int32 {
	return int32(i)
}

func (i Int8Req) Int64() int64 {
	return int64(i)
}

func (i StringReq) String() string {
	return string(i)
}
