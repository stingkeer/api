package dwarf

import (
	"fmt"
	"os"
	"testing"
)

type D struct {
}

var exeHelper *DwarfMaker

func Init() {
	path, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exeHelper = NewDwarfMaker()
	exeHelper.Init(&path)
}

func (*D) Show(shangzebei string) {
	fmt.Println("aaa")
}

func Add(int2 int, s string) {

}

func TestError(t *testing.T) {
	dd := &D{}
	kk := exeHelper.LookFun(dd.Show)
	fmt.Println(kk)
}

func TestFunBase(t *testing.T) {
	Init()
	pp := exeHelper.LookFun(Add)
	fmt.Println(pp)
}
