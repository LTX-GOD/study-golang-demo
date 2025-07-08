package main

import (
	"fmt"

	"github.com/zsm/go-bookstore/pkg/routes"
)

func main() {
	r := routes.Router()
	fmt.Println("Server is running on localhost:9010")
	r.Run(":9010")
}
