package httpsnusk

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"os"
)

// Out is the writer where requests and responses are written
var Out io.Writer = os.Stdout

// RoundTripper is a http.RoundTripper
type RoundTripper struct {
	Transport    http.RoundTripper
	DumpRequest  func(*http.Request)
	DumpResponse func(*http.Response)
}

// DefaultRoundTripper is a RoundTripper that prints the request/response to Out
var DefaultRoundTripper = &RoundTripper{
	Transport: http.DefaultTransport,

	DumpRequest: func(req *http.Request) {
		doPrint(httputil.DumpRequestOut(req, true))
	},

	DumpResponse: func(resp *http.Response) {
		doPrint(httputil.DumpResponse(resp, true))
	},
}

// RoundTrip implements the http.RoundTripper interface
func (rt *RoundTripper) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	rt.DumpRequest(req)

	defer func() {
		if err == nil {
			rt.DumpResponse(resp)
		}
	}()

	return rt.Transport.RoundTrip(req)
}

// Handler func is a http.HandlerFunc
type HandlerFunc struct {
	Handler      http.HandlerFunc
	DumpRequest  func(*http.Request)
	DumpResponse func(*http.Response)
}

// HandlerFunc404 is a HandlerFunc that that prints the request/respones to Out
var HandlerFunc404 = HandlerFunc{
	Handler: http.NotFound,

	DumpRequest: func(req *http.Request) {
		doPrint(httputil.DumpRequest(req, true))
	},

	DumpResponse: func(resp *http.Response) {
		doPrint(httputil.DumpResponse(resp, true))
	},
}

// ServeHTTP implements the HandlerFunc interface
func (hf HandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	hf.DumpRequest(r)

	recorder := httptest.NewRecorder()
	hf.Handler(recorder, r)

	hf.DumpResponse(recorder.Result())

	copyResponseWriter(w, recorder)
}

func copyResponseWriter(dst http.ResponseWriter, src *httptest.ResponseRecorder) error {
	dst.WriteHeader(src.Code)

	for k, v := range src.Header() {
		dst.Header()[k] = v
	}

	_, err := io.Copy(dst, src.Body)
	return err
}

func doPrint(b []byte, err error) {
	if err != nil {
		fmt.Fprintln(Out, err)
	}
	fmt.Fprintln(Out, string(b))
}
