package LoaderConfig

import (
	"Config/Config"
	"Config/ReadFlags"
	"Net/HTTPLoader"
	"fmt"
	. "gopkg.in/check.v1"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLoaderConfig(t *testing.T) {
	TestingT(t)
}

type LoaderConfigTestsSuite struct{}

var _ = Suite(&LoaderConfigTestsSuite{})

const (
	testConfig string = `{"host":"test_host"}`
)

func (s *LoaderConfigTestsSuite) TestLoaderConfigLoadRemoutConfig(c *C) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, testConfig)
	}))
	defer ts.Close()

	config := Config.New(ReadFlags.New())
	loader := HTTPLoader.NewHTTPLoader()

	f := NewLoaderConfig(config, loader)

	err := f.LoadURL(ts.URL)
	c.Assert(err, IsNil)
}

func (s *LoaderConfigTestsSuite) TestLoaderConfigNotLoadEmptyFile(c *C) {

	config := Config.New(ReadFlags.New())
	loader := HTTPLoader.NewHTTPLoader()

	f := NewLoaderConfig(config, loader)

	err := f.LoadFile("")

	// LoaderConfig tries Load EmptyURL
	c.Assert(err, NotNil)
}

func (s *LoaderConfigTestsSuite) TestLoaderConfigNotLoadFile(c *C) {

	config := Config.New(ReadFlags.New())
	loader := HTTPLoader.NewHTTPLoader()

	f := NewLoaderConfig(config, loader)

	err := f.LoadFile("")

	// LoaderConfig tries Load EmptyURL
	c.Assert(err, NotNil)
}

func (s *LoaderConfigTestsSuite) TestLoaderConfigNotLoadEmptyURL(c *C) {

	config := Config.New(ReadFlags.New())
	loader := HTTPLoader.NewHTTPLoader()

	f := NewLoaderConfig(config, loader)

	err := f.LoadURL("")

	// LoaderConfig tries Load EmptyURL
	c.Assert(err, NotNil)
}
