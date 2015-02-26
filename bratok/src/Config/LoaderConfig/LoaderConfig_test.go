package LoaderConfig

import (
	. "gopkg.in/check.v1"
	"testing"
)

func TestLoaderConfig(t *testing.T) {
	TestingT(t)
}

type LoaderConfigTestsSuite struct{}

var _ = Suite(&LoaderConfigTestsSuite{})

func (s *LoaderConfigTestsSuite) Test_LoaderConfig_Load_RemoutConfig(c *C) {

	f := NewLoaderConfig()

	err := f.Load("http://127.0.0.1/test.conf")

	// LoaderConfig tries Load "http://127.0.0.1/test.conf"
	c.Assert(err, IsNil)
}

func (s *LoaderConfigTestsSuite) Test_LoaderConfig_Not_Load_EmptyURL(c *C) {

	f := NewLoaderConfig()

	err := f.Load("")

	// LoaderConfig tries Load EmptyURL
	c.Assert(err, NotNil)
}
