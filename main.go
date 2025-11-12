package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"tmhelper/tmhelper"
)
const CPTENV ="TMHCPTKEY"
func main() {
	jsPath := flag.String("f", "", "file, 运行tmhelper js文件")
	jxPath := flag.String("x", "", "file, 运行tmhelper js文件")
	jsCode := flag.String("c", "", "code, 运行tmhelper js代码")
	encText := flag.String("e", "", "encrypt,加密文本")
	decText := flag.String("d", "", "decrypt,解密文本")
	flag.Parse()
    pwd:=os.Getenv(CPTENV)
	code := []byte{}
	if *jsPath != "" {
		code = readFromFile(*jsPath)
	} else if *jxPath != "" {
     	code = readFromFile(*jxPath)
     	fmt.Print("Please enter password decrypt key: ")
     	pwd=tmhelper.ReadStr(os.Stdin)
    } else if *jsCode != "" {
		code = []byte(*jsCode)
	} else if *encText != "" {
		encryptText(*encText)
		return
	}  else if *decText != "" {
      	decryptText(*decText)
      	return
    }else {
		code = readFromStdin()
	}
	js := NewJS()
	js.Run(string(code),pwd)
}

func readFileCode(jsPath string) (b []byte) {
    b = readFromFile(jsPath)
	return b
}

func readFromFile(jsPath string) []byte {
	b, err := os.ReadFile(jsPath)
	if err != nil {
		errorf("read js file error: %v", err)
	}

	return b
}

func encryptText(plain string) {
    fmt.Println(tmhelper.EncText(plain,os.Getenv(CPTENV)))
}
func decryptText(encstr string) {
    fmt.Println(tmhelper.DecText(encstr,os.Getenv(CPTENV)))
}

func readFromStdin() []byte {
	b, err := io.ReadAll(os.Stdin)
	if err != nil {
		errorf("read stdin error: %v", err)
	}
	return b
}

func errorf(format string, vals ...any) {
	fmt.Fprintf(os.Stderr, format+"\n", vals...)
	os.Exit(1)
}
