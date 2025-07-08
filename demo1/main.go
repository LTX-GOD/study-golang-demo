package main

import (
	"fmt"
	"log"
	"net/http"
)

func formHandler(w http.ResponseWriter, r *http.Request) {
	// w 为 http.ResponseWriter 的实例，用于向客户端返回响应。
	// r 为 http.Request 的实例，包含了客户端的请求信息。
	err := r.ParseForm()
	if err != nil {
		fmt.Println(w, "ParseForm() err:%v", err)
		return
	}
	fmt.Fprintln(w, "POST request successful")

	name := r.FormValue("name")
	address := r.FormValue("address")
	fmt.Fprintf(w, "Name=%s\n", name)
	fmt.Fprintf(w, "address=%s\n", address)

}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	//请求路径
	if r.URL.Path != "/hello" {
		http.Error(w, "404 not found", http.StatusNotFound)
		return
	}
	//请求方法
	if r.Method != "GET" {
		http.Error(w, "method is not supported", http.StatusNotFound)
		return
	}
	fmt.Fprintf(w, "Hello")
}

func main() {
	fileServer := http.FileServer(http.Dir("./static")) //挂载静态文件
	http.Handle("/", fileServer)                        //挂载到根目录

	//路径绑定函数
	http.HandleFunc("/form", formHandler)
	http.HandleFunc("/hello", helloHandler)

	fmt.Println("server is running on 8080")

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}

}
