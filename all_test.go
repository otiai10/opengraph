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

	s := dummyServer()
	og, err := Fetch(s.URL)

	Expect(t, err).ToBe(nil)
	Expect(t, og.Title).ToBe("Hello! Open Graph!!")
	Expect(t, og.Type).ToBe("website")
	Expect(t, og.URL.Source).ToBe(s.URL)
	Expect(t, len(og.Image)).ToBe(1)

	Expect(t, og.Image[0].URL).ToBe("/images/01.png")
	Expect(t, og.Favicon).ToBe("/images/01.favicon.png")
	og.ToAbsURL()
	Expect(t, og.Image[0].URL).ToBe(s.URL + "/images/01.png")
	Expect(t, og.Favicon).ToBe(s.URL + "/images/01.favicon.png")
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
