package httpsnusk

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"regexp"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello"))
}

func replacePortAndDate(s string) string {
	s = regexp.MustCompile("Date: .*").ReplaceAllString(s, "Date: <DATE>")
	s = regexp.MustCompile("Host: .*").ReplaceAllString(s, "Host: <HOST>")

	// Whitespace fixes. https://github.com/golang/go/issues/26460
	s = regexp.MustCompile("\r\n").ReplaceAllString(s, "\n")
	s = regexp.MustCompile("\n\n+").ReplaceAllString(s, "\n\n")
	return s
}

func ExampleHandlerFunc() {
	// Capture output in buffer instead of os.Stdout
	var b bytes.Buffer
	Out = &b

	// Copy the 404-handler and inject our handler
	f := HandlerFunc404
	f.Handler = helloHandler

	server := httptest.NewServer(f)
	defer server.Close()

	// Execute the request
	http.Get(server.URL)

	// Print the captured output from the request
	fmt.Print(replacePortAndDate(b.String()))
	// Output:
	// GET / HTTP/1.1
	// Host: <HOST>
	// Accept-Encoding: gzip
	// User-Agent: Go-http-client/1.1
	//
	//
	// HTTP/1.1 200 OK
	// Connection: close
	// Content-Type: text/plain; charset=utf-8
	//
	// hello
}

func ExampleRoundTripper() {
	// Capture output in buffer instead of os.Stdout
	var b bytes.Buffer
	Out = &b

	server := httptest.NewServer(http.HandlerFunc(helloHandler))

	cli := http.Client{
		Transport: DefaultRoundTripper,
	}

	// Execute the request
	cli.Get(server.URL)

	// Print the captured output from the request
	fmt.Println(replacePortAndDate(b.String()))
	// Output:
	// GET / HTTP/1.1
	// Host: <HOST>
	// User-Agent: Go-http-client/1.1
	// Accept-Encoding: gzip
	//
	//
	// HTTP/1.1 200 OK
	// Content-Length: 5
	// Content-Type: text/plain; charset=utf-8
	// Date: <DATE>
	//
	// hello
}
