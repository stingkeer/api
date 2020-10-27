package api

import (
	"fmt"
	"reflect"
	"runtime"
	"sync"
	"testing"
)

func TestDef(t *testing.T) {
	PackApi()
}

func aa() {

}

func show(f interface{}) {
	v := reflect.ValueOf(f)
	fmt.Println(runtime.FuncForPC(v.Pointer()).Name())
}

func TestFunc(t *testing.T) {
	go func() {
		show(aa)
	}()
}

func TestAvg(t *testing.T) {
	averageDo(4, 8, func(start, end int, g *sync.WaitGroup) {
		for i := start; i < end; i++ {
			fmt.Print(i, " ")
		}
		fmt.Println()
		g.Done()
	})
}
