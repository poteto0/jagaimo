package main

import (
	"fmt"

	"github.com/poteto0/jagaimo/net/linux/http"
)

func main() {
	client := http.NewHttpClient()

	res, err := client.GET(
		"example.com", 80, "",
	)

	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	fmt.Println(res)
}
