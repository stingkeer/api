package def

//require param
/**
  //b is require
  GET(func(a int, b def.StringReq) {
		fmt.Println(a, b)
   }, "/require")

*/
//
type (
	IntReq    Int[int]
	Int8Req   Int8[int8]
	Int16Req  Int16[int16]
	Int32Req  Int32[int32]
	Int64Req  Int64[int64]
	StringReq String[string]
)

func (i IntReq) Int() int {
	return i.V
}

func (i Int8Req) Int8() int8 {
	return i.V
}

func (i Int16Req) Int16() int16 {
	return i.V
}

// Deprecated: an Int32() accessor on Int8Req was a copy-paste mistake. Use
// Int32Req.Int32() or int32(i.V).
func (i Int8Req) Int32() int32 {
	return int32(i.V)
}

func (i Int32Req) Int32() int32 {
	return i.V
}

func (i Int64Req) Int64() int64 {
	return i.V
}

func (i StringReq) String() string {
	return i.V
}
