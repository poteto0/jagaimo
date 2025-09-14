package main

import (
	"fmt"

	"github.com/poteto0/jagaimo/core"
	"github.com/poteto0/jagaimo/core/browser"
)

var testResponse = `HTTP/1.1 200 OK
Data: xx xx xx


<html>
<head></head>
<body>
  <h1 id="title">H1 title</h1>
  <h2 class="class">H2 title</h2>
  <p>Test text.</p>
  <p>
    <a href="example.com">Link1</a>
    <a href="example.com">Link2</a>
  </p>
</body>
</html>
`

func main() {
	saba := browser.NewBrowser()
	page := saba.CurrentPage()

	response, _ := core.NewHttpResponse(testResponse)
	domString := page.ReceiveResponse(
		response,
	)

	fmt.Println(domString)
}
