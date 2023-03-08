package def

//require param
/**
  //b is require
  GET(func(a int, b def.StringReq) {
		fmt.Println(a, b)
   }, "/require")

*/
type (
	Int[T any] struct {
		V int
		_ T
	}
	Int8[T any] struct {
		V int8
		_ T
	}
	Int16[T any] struct {
		V int16
		_ T
	}
	Int32[T any] struct {
		V int32
		_ T
	}
	Int64[T any] struct {
		V int64
		_ T
	}
	String[T any] struct {
		V string
		_ T
	}
)

func (i Int[T]) Int() int {
	return i.V
}

func (i Int8[T]) Int8() int8 {
	return i.V
}

func (i Int16[T]) Int16() int16 {
	return i.V
}

func (i Int8[T]) Int32() int32 {
	return int32(i.V)
}

func (i Int64[T]) Int64() int64 {
	return i.V
}

func (i String[T]) String() string {
	return i.V
}
