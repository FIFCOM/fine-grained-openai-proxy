package main

import (
	"fine-grained-openai-proxy/conf"
	"fine-grained-openai-proxy/ctrl"
	"flag"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	allowAdmin := flag.Bool("admin", false, "set to true if allow admin or server operation")
	flag.Parse()

	r := gin.Default()
	ctrl.InitRouter(r, *allowAdmin)
	if err := r.Run(conf.Port); err != nil {
		fmt.Printf("Start fine-grained openai proxy error : %s", err.Error())
		os.Exit(1)
	}
}