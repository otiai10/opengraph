package opengraph

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/otiai10/marmoset"
	. "github.com/otiai10/mint"
)

func TestNew(t *testing.T) {
	og := New(dummyServer().URL)
	Expect(t, og).TypeOf("*opengraph.OpenGraph")
	Expect(t, og.Error).ToBe(nil)

	When(t, "invalid url is given", func(t *testing.T) {
		og := New(":invalid_url")
		Expect(t, og.Error).Not().ToBe(nil)
	})
}

func TestFetch(t *testing.T) {
	When(t, "invalid scheme is given", func(t *testing.T) {
		_, err := Fetch(":invalid_url")
		Expect(t, err).Not().ToBe(nil)
	})
	When(t, "invalid url is given", func(t *testing.T) {
		_, err := Fetch("htt://xxx/yyy")
		Expect(t, err).Not().ToBe(nil)
	})
	og, err := Fetch(dummyServer().URL)
	Expect(t, err).ToBe(nil)
	Expect(t, og.Title).ToBe("Hello! Open Graph!!")
}

func dummyServer() *httptest.Server {
	marmoset.LoadViews("./testdata/html")
	r := marmoset.NewRouter()
	r.GET("/", func(w http.ResponseWriter, r *http.Request) {
		marmoset.Render(w).HTML("01", nil)
	})
	r.GET("/case/01", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})
	return httptest.NewServer(r)
}
