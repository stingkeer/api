package dwarf

import (
	"debug/gosym"
	"fmt"
	"os"
	"regexp"
	"testing"
)

type D struct {
	_ string
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
	Init()
	dd := &D{}
	kk, _ := exeHelper.LookFun(dd.Show)
	fmt.Println(kk)
}

func TestFunBase(t *testing.T) {
	Init()
	pp, _ := exeHelper.LookFun(Add)
	fmt.Println(pp)
}

func TestRegex(t *testing.T) {
	regexp.MustCompile("$gitee.com/fast_api/.+$")
}

func TestGoSym(t *testing.T) {
	b := gosym.Sym{Name: "main.init.0.func1"}
	fmt.Println(b.PackageName())
}
