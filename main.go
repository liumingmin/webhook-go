package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
)

var (
	help = flag.Bool("h", false, "The Help")
	port = flag.String("p", "7442", "HTTP Server Port, Default `7442`")
)

func init() {
	flag.Usage = usage
}

func usage() {
	fmt.Fprintf(os.Stderr, `Webhook version: webhook/1.0.0
	Usage: hook [-h] [-s signal] [-p port]

	Options: 
	`)
	flag.PrintDefaults()
}

func main() {
	// 命令行的参数获取
	flag.Parse()
	if *help {
		flag.Usage()
	} else {
		deamonHTTP()
	}

}

func deamonHTTP() {
	serveHTTP()
}

func serveHTTP() {
	mux := http.NewServeMux()
	mux.HandleFunc("/gitee", gitee)
	mux.HandleFunc("/coding", coding)
	mux.HandleFunc("/gogs", gogs)
	mux.HandleFunc("/", index)
	log.Fatalln(http.ListenAndServe(":"+*port, mux))
}
