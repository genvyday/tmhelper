package main

import (
	"fmt"
	"os"
	"runtime"
	"encoding/base64"

	"golang.org/x/term"
	"github.com/dop251/goja"
	"tmhelper/tmhelper"
)

type JS struct {
	vm *goja.Runtime
	cptkey  []byte
	timeout int
}
type TermProc struct {
	vm *goja.Runtime
	xe *tmhelper.TMHelper
}

func NewJS() *JS {
	sf:= &JS{
		vm: goja.New(),
		cptkey: nil,
		timeout:10,
	}
    sf.vm.Set("tmh", sf)
    return sf
}

func (sf *JS) Run(jsCode string) goja.Value {
	v, err := sf.vm.RunString(jsCode)
	if err != nil {
		panic(err)
	}
    return v;
}
func (sf *JS) CptKey(value goja.FunctionCall) goja.Value {
    ekey:=value.Argument(0).String()
    chk:=value.Argument(1).String()
    enc:=tmhelper.EncText("123",ekey)
    valid:=enc==chk;
    if(valid){
        sf.cptkey=tmhelper.GenKey([]byte(ekey),32)
    }
	return sf.vm.ToValue(valid)
}

func (sf *JS) Dec(value goja.FunctionCall) goja.Value {
    str := value.Argument(0)
    if sf.cptkey==nil||len(sf.cptkey)==0{
        return str
    }
    encData,_:=base64.RawURLEncoding.DecodeString(str.String())
	plain:=string(tmhelper.AesDec(encData,sf.cptkey))
	return sf.vm.ToValue(plain)
}

func (sf *JS) Goos(value goja.FunctionCall) goja.Value {
    return sf.vm.ToValue(runtime.GOOS)
}
func (sf *JS) Enc(value goja.FunctionCall) goja.Value {
    str := value.Argument(0)
    if sf.cptkey==nil||len(sf.cptkey)==0{
        return str
    }
	encData:=tmhelper.AesEnc([]byte(str.String()),sf.cptkey)
	ret:=base64.RawURLEncoding.EncodeToString(encData)
	return sf.vm.ToValue(ret)
}
func (sf *JS) SetTimeout(value goja.FunctionCall) goja.Value {
	sec := value.Argument(0).ToInteger()
	sf.timeout=int(sec)
	return sf.vm.ToValue(sf)
}
func (sf *JS) Pwd(call goja.FunctionCall) goja.Value {
	str := call.Argument(0)
	fmt.Print(str.String())
	bpwd,_ := term.ReadPassword(int(os.Stdin.Fd()))
	return sf.vm.ToValue(string(bpwd))
}
func (sf *JS) Input(call goja.FunctionCall) goja.Value {
    prompt:=call.Argument(0).String()
    fmt.Print(prompt)
	return sf.vm.ToValue(tmhelper.ReadStr(os.Stdin))
}
func (sf *JS) Println(call goja.FunctionCall) goja.Value {
	str := call.Argument(0)
	fmt.Println(str.String())
	return str
}
func (sf *JS) Print(call goja.FunctionCall) goja.Value {
	str := call.Argument(0)
	fmt.Print(str.String())
	return str
}
func (sf *JS) NewTerm(call goja.FunctionCall) goja.Value {
	tp:=&TermProc{
	    vm:sf.vm,
	    xe:tmhelper.NewTMHelper(),
	}
    timeout:=sf.timeout
	arg:=call.Argument(0)
	if arg!=goja.Undefined(){
	    timeout=int(arg.ToInteger())
	}
    tp.xe.SetTimeout(timeout)
	return sf.vm.ToValue(tp)
}
func formatArgs(value goja.Value) []string {
	arrayInterface, ok := value.Export().([]interface{})
	if !ok {
		errorf("run args error: must string array")
		return nil
	}

	var args []string
	for _, item := range arrayInterface {
		str, ok := item.(string)
		if !ok {
			errorf("run args error: must string array")
			return nil
		}
		args = append(args, str)
	}
	return args
}

func (tp *TermProc) Exec(value goja.FunctionCall) goja.Value {
	args := formatArgs(value.Argument(0))
	fmt.Println(args)
	tp.xe.Run(args)
	return tp.vm.ToValue(tp)
}
// func (rule [][]string) map[string]any{"idx": idx, "str": str}
func (tp *TermProc) Matchs(value goja.FunctionCall) goja.Value {
	rule := formatRule(value.Argument(0))
	idx, str := tp.xe.Matchs(rule)
	return tp.vm.ToValue(map[string]any{"idx": idx, "str": str})
}

func (tp *TermProc) Term(_ goja.FunctionCall) goja.Value {
	tp.xe.Term()
	return tp.vm.ToValue(tp)
}
func (tp *TermProc) Expect(call goja.FunctionCall) goja.Value {
	str := call.Argument(0)
	tp.xe.Expect(str.String())
	return str
}
func (tp *TermProc) ValRaw(call goja.FunctionCall) goja.Value {
    ret :=tp.xe.ValRaw()
    return tp.vm.ToValue(ret)
}
func (tp *TermProc) ValHex(call goja.FunctionCall) goja.Value {
    ret :=tp.xe.ValHex()
    return tp.vm.ToValue(ret)
}
func (tp *TermProc) ReadStr(call goja.FunctionCall) goja.Value {
	str := call.Argument(0)
	ret :=tp.xe.ReadPty(str.String())
	return tp.vm.ToValue(ret)
}
func (tp *TermProc) WaitDone(call goja.FunctionCall) goja.Value {
	arg:=call.Argument(0)
	str:=""
	if arg!=goja.Undefined(){
	    str=arg.String()
	}
    tp.xe.WaitRelayExit(str)
    return arg;
}
func (tp *TermProc) Exit(call goja.FunctionCall) goja.Value {
	tp.xe.Exit()
	return tp.vm.ToValue(tp)
}
func (tp *TermProc) Ok(value goja.FunctionCall) goja.Value {
	return tp.vm.ToValue(tp.xe.Ok())
}
func (tp *TermProc) Input(call goja.FunctionCall) goja.Value {
    prompt:=call.Argument(0).String()
	return tp.vm.ToValue(tp.xe.ReadInput(prompt))
}
func formatRule(value goja.Value) [][]string {
	var rule [][]string

	// 断言Value是一个数组
	arr1, ok := value.Export().([]interface{})
	if !ok {
		errorf("matchs args error: must two-dimensional string array")
	}

	// 遍历第一层数组
	for i := 0; i < len(arr1); i++ {
		arr2, ok := arr1[i].([]interface{})
		if !ok {
			errorf("matchs args error: must two-dimensional string array")
		}

		var innerArray []string

		// 遍历内层数组
		for _, elem := range arr2 {
			strVal, ok := elem.(string)
			if !ok {
				errorf("matchs args error: must two-dimensional string array")

			}
			innerArray = append(innerArray, strVal)
		}

		rule = append(rule, innerArray)
	}

	return rule
}
