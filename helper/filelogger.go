package helper

import (
	"os"
	"io"
	"fmt"
	"github.com/kataras/iris"
)

type fileloggerMiddleware struct {
	path string
}

func (l *fileloggerMiddleware) Serve(ctx *iris.Context) {
	go saveToFile(l.path, getInfo(ctx))
}

func  getInfo(ctx *iris.Context) string{
	var path, method, status, ip, request string
	path = ctx.Path()
	method = ctx.Method()
	ctx.Next()
	result := fmt.Sprintf("%s %s %s %s \n %s \n\n", path, method, status, ip, request)
	return result
}

func saveToFile(filename string, body string){
	f, err := os.OpenFile( filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666 )
	if err == nil{
		defer f.Close()
		io.WriteString(f, body)
	}
}

func New(path string) iris.HandlerFunc {
	l := &fileloggerMiddleware{path: path}
	return l.Serve
}