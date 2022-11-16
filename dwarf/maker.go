package dwarf

import (
	"debug/dwarf"
	"debug/elf"
	"debug/macho"
	"debug/pe"
	"fmt"
	"log"
	"os"
	"reflect"
	"runtime"
	"strings"
	"time"
)

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

type Params []string

type DwarfMaker struct {
	openData func() *dwarf.Reader
	r        *dwarf.Reader
	debug    map[string]Params
}

func NewDwarfMaker() *DwarfMaker {
	return &DwarfMaker{debug: make(map[string]Params, 1000)}
}

func (h *DwarfMaker) load(exe *string) {
	path := ""
	if exe != nil {
		path = *exe
	} else {
		if dll, err := os.Executable(); err == nil {
			path = dll
		}
	}
	if path == "" {
		panic("DwarfMaker load path == nil")
	}
	fmt.Printf("system os %s\n", runtime.GOOS)
	switch runtime.GOOS {
	case "windows":
		h.openData = func() *dwarf.Reader {
			f, e := pe.Open(path)
			checkError(e)
			defer f.Close()
			data, e := f.DWARF()
			checkError(e)
			return data.Reader()
		}
	case "linux":
		h.openData = func() *dwarf.Reader {
			f, e := elf.Open(path)
			checkError(e)
			defer f.Close()
			data, e := f.DWARF()
			checkError(e)
			return data.Reader()
		}

	case "darwin":
		h.openData = func() *dwarf.Reader {
			f, e := macho.Open(path)
			checkError(e)
			defer f.Close()
			data, e := f.DWARF()
			checkError(e)
			return data.Reader()
		}

	}
}

func (h *DwarfMaker) Init(exe *string) {
	now := time.Now()
	h.load(exe)
	h.r = h.openData()
	tempName := ""
	for r, _ := h.r.Next(); r != nil; r, _ = h.r.Next() {
		if rName := r.Val(dwarf.AttrName); r.Tag == dwarf.TagSubprogram && rName != nil {
			//sym := gosym.Sym{
			//	Name: rName.(string),
			//}
			tempName = rName.(string)
			h.debug[tempName] = Params{}
			if strings.HasPrefix(tempName, "gitee.com/fast_api/api") {
				fmt.Println(tempName)
			}

		}
		if r.Tag != dwarf.TagFormalParameter {
			continue
		}
		if v, b := r.Val(dwarf.AttrVarParam).(bool); b && v {
			continue
		}
		n := r.Val(dwarf.AttrName)
		if n == nil {
			continue
		}
		if strings.HasPrefix(tempName, "gitee.com/fast_api/api") {
			fmt.Println(tempName, r.Val(dwarf.AttrVarParam), r.Val(dwarf.AttrDeclLine), n)
		}
		if _, b := h.debug[tempName]; b {
			h.debug[tempName] = append(h.debug[tempName], n.(string))
		}

	}
	h.r = nil
	log.Printf("DwarfMaker init use %s", time.Since(now))
}

func (h *DwarfMaker) LookFun(inf interface{}) *MethodMeta {
	v := reflect.ValueOf(inf)
	fName := runtime.FuncForPC(v.Pointer()).Name()
	if v.Kind() != reflect.Func {
		return nil
	}
	tp := v.Type()
	//argsNum := tp.NumIn()
	if params, b := h.debug[fName]; b {
		var args []ArgsMeta
		for i, vName := range params {
			args = append(args, ArgsMeta{
				Order: i,
				Name:  vName,
				Typ:   tp.In(i),
			})
		}
		return &MethodMeta{
			MethodName: fName,
			Args:       args,
			Ret:        nil,
		}
	}
	return nil

}
