package middlewares

import (
	"fmt"
	"log"
	"net/http"
	"runtime"
	"strings"

	"github.com/lz-nsc/cupcake"
)

func trace(message string) string {
	var pcs [32]uintptr
	//Skip first three callers
	n := runtime.Callers(3, pcs[:])

	var str strings.Builder
	str.WriteString(message + "\nTraceback:")

	for _, pc := range pcs[:n] {
		fn := runtime.FuncForPC(pc)
		file, line := fn.FileLine(pc)
		str.WriteString(fmt.Sprintf("\n\t%s:%d", file, line))
	}
	return str.String()
}

func Recovery(handler cupcake.HandlerFunc) cupcake.HandlerFunc {
	return cupcake.HandlerFunc(func(resp *cupcake.Response, req *cupcake.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("%s\n\n", trace(fmt.Sprintf("%s", err)))
				resp.Error(http.StatusInternalServerError, "Internal Server Error")
			}
		}()
		handler(resp, req)
	})
}
