package function

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

// TestHandle ensures that Handle executes without error and returns the
// HTTP 200 status code indicating no errors.
func TestHandle(t *testing.T) {
	data := url.Values{}
	data.Set("name", "foo")
	data.Set("surname", "bar")
	var (
		w   = httptest.NewRecorder()

		req = httptest.NewRequest("POST", "http://example.com/test", strings.NewReader(data.Encode()))

		res *http.Response
		err error
	)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// Invoke the Handler via a standard Go http.Handler
	func(w http.ResponseWriter, req *http.Request) {
		Handle(context.Background(), w, req)
	}(w, req)

	res = w.Result()
	defer res.Body.Close()

	// Assert postconditions
	if err != nil {
		t.Fatalf("unepected error in Handle: %v", err)
	}
	if res.StatusCode != 200 {
		t.Fatalf("unexpected response code: %v", res.StatusCode)
	}
}
