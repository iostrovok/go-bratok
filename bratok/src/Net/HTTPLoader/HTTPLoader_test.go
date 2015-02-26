package HTTPLoader

import (
	"fmt"
	. "gopkg.in/check.v1"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHTTPLoader(t *testing.T) {
	TestingT(t)
}

type HTTPLoaderTestsSuite struct{}

var _ = Suite(&HTTPLoaderTestsSuite{})

// TestHTTPLoaderCanLoadPOST test POST/GET requests
func (s *HTTPLoaderTestsSuite) TestHTTPLoaderCanLoadPOST(c *C) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `{"name":"Mike"}`)
	}))
	defer ts.Close()

	f := NewHTTPLoader()

	body, err := f.Load("", "POST", ts.URL, map[string]interface{}{})

	c.Assert(err, IsNil)
	c.Assert(string(body), Equals, `{"name":"Mike"}`+"\n")
}

func (s *HTTPLoaderTestsSuite) TestHTTPLoaderCanLoadGET(c *C) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `{"name":"Mike"}`)
	}))
	defer ts.Close()

	f := NewHTTPLoader()

	body, err := f.Load("", "GET", ts.URL, map[string]interface{}{})

	c.Assert(err, IsNil)
	c.Assert(string(body), Equals, `{"name":"Mike"}`+"\n")
}

func (s *HTTPLoaderTestsSuite) TestHTTPLoaderCanLoadBadURL(c *C) {
	f := NewHTTPLoader()

	_, err := f.Load("", "POST", "", map[string]interface{}{})

	c.Assert(err, NotNil)
}

func (s *HTTPLoaderTestsSuite) TestHTTPLoaderCanLoadWrongURL(c *C) {
	f := NewHTTPLoader()

	_, err := f.Load("", "POST", "httpqweqwewqeqw://asdsadasdas", map[string]interface{}{})

	c.Assert(err, NotNil)
}

func (s *HTTPLoaderTestsSuite) TestHTTPLoaderCanLoadWrongType(c *C) {
	f := NewHTTPLoader()

	_, err := f.Load("", "DDD", "", map[string]interface{}{})

	c.Assert(err, NotNil)
}

// TestHTTPLoaderCanLoadPOST test GET requests
func (s *HTTPLoaderTestsSuite) TestHTTPLoaderCanLoad404(c *C) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `{"name":"Mike"}`)
	}))
	defer ts.Close()

	f := NewHTTPLoader()

	_, err := f.Load("", "POST", ts.URL+"u", map[string]interface{}{})

	c.Assert(err, NotNil)
}

// TestHTTPLoaderCanLoadPOST test GET requests
func (s *HTTPLoaderTestsSuite) TestHTTPLoaderCanLoadJSON404(c *C) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `{"name":"Mike"}`)
	}))
	defer ts.Close()

	f := NewHTTPLoader()

	_, err := f.LoadJson("POST", ts.URL+"u", map[string]interface{}{})

	c.Assert(err, NotNil)
}

func (s *HTTPLoaderTestsSuite) TestHTTPLoaderCanLoadJson(c *C) {

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `{"name":"Mike"}`)
	}))

	defer ts.Close()

	f := NewHTTPLoader()

	body, err := f.LoadJson("POST", ts.URL, map[string]interface{}{"e": 1})

	c.Logf("err: %s\n", err)
	c.Logf("body: %+v\n", body)

	c.Assert(err, IsNil)
}
