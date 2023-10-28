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
	port := flag.Int64("port", 8080, "set the port to listen")
	adminToken := flag.String("admin-token", "123456789", "set the admin token")
	flag.Parse()

	if *port != 8080 && *port > 0 && *port < 65535 {
		conf.Port = fmt.Sprintf(":%d", *port)
	}

	conf.AdminToken = *adminToken

	r := gin.Default()
	ctrl.InitRouter(r, *allowAdmin)
	if err := r.Run(conf.Port); err != nil {
		fmt.Printf("Start fine-grained openai proxy error : %s", err.Error())
		os.Exit(1)
	}
}