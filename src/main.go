package main

import (
	"fmt"

	"github.com/poteto0/jagaimo/net/linux/http"
)

func main() {
	client := http.NewHttpClient()

	res, err := client.GET(
		"localhost", 8000, "/test.html",
	)

	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	fmt.Println(res)
}
