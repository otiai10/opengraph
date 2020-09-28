package opengraph

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	mmst "github.com/otiai10/marmoset"
	. "github.com/otiai10/mint"
)

var testserver *httptest.Server

func TestMain(m *testing.M) {
	testserver = createTestServer()
	code := m.Run()
	os.Exit(code)
}

func TestNew(t *testing.T) {
	og := New("https://github.com")
	Expect(t, og).TypeOf("*opengraph.OpenGraph")
}

func TestFetch(t *testing.T) {
	og, err := Fetch(testserver.URL + "/case/01_hello")
	Expect(t, err).ToBe(nil)
	Expect(t, og).TypeOf("*opengraph.OpenGraph")
	Expect(t, og.Description).ToBe("This description should be preferred")

	When(t, "invalid URL provided", func(t *testing.T) {
		_, err := Fetch(":SOME_INVALID_URL")
		Expect(t, err).Not().ToBe(nil)
	})

	When(t, "og:image provided", func(t *testing.T) {
		og, err := Fetch(testserver.URL + "/case/03_image")
		Expect(t, err).ToBe(nil)
		Expect(t, og).TypeOf("*opengraph.OpenGraph")
	})
	When(t, "structured image properties provided", func(t *testing.T) {
		og, err := Fetch(testserver.URL + "/case/04_image_props")
		Expect(t, err).ToBe(nil)
		Expect(t, og).TypeOf("*opengraph.OpenGraph")
	})
	When(t, "content-type is NOT text/html", func(t *testing.T) {
		_, err := Fetch(testserver.URL + "/case/00_data.json")
		Expect(t, err).Not().ToBe(nil)
		Expect(t, err.Error()).ToBe("Content type must be text/html")
	})
}

func TestOpenGraph_Fetch(t *testing.T) {
	ogp := &OpenGraph{}
	err := ogp.Fetch(nil)
	Expect(t, err.Error()).ToBe("no URL given yet")

	ogp.Intent.URL = testserver.URL + "/case/01_hello"
	err = ogp.Fetch(nil)
	Expect(t, err).ToBe(nil)

	When(t, "ogp already has an error", func(t *testing.T) {
		ogp := New(":INVALID_URL")
		err := ogp.Fetch(nil)
		Expect(t, err).Not().ToBe(nil)
	})
}

func TestOpenGraph_Parse(t *testing.T) {
	ogp := New("")
	r := strings.NewReader(`<html><meta property="og:title" content="test_test"></html>`)
	err := ogp.Parse(r)
	Expect(t, err).ToBe(nil)
	Expect(t, ogp.Title).ToBe("test_test")
}

func TestOpenGraph_ToAbs(t *testing.T) {
	ogp := New(testserver.URL + "/case/01_hello")
	err := ogp.Fetch(nil)
	Expect(t, err).ToBe(nil)
	u, err := url.Parse(ogp.Image[0].URL)
	Expect(t, err).ToBe(nil)
	Expect(t, u.IsAbs()).ToBe(false)
	err = ogp.ToAbs()
	Expect(t, err).ToBe(nil)
	u, err = url.Parse(ogp.Image[0].URL)
	Expect(t, err).ToBe(nil)
	Expect(t, u.IsAbs()).ToBe(true)
	u, err = url.Parse(ogp.Audio[0].URL)
	Expect(t, err).ToBe(nil)
	Expect(t, u.IsAbs()).ToBe(true)
}

// This server is ONLY for testing.
func createTestServer() *httptest.Server {
	mmst.LoadViews("./test/html")
	r := mmst.NewRouter()
	r.GET("/case/(?P<name>[_a-zA-Z0-9]+)", func(w http.ResponseWriter, r *http.Request) {
		name := r.FormValue("name")
		switch {
		case strings.HasSuffix(name, ".json"):
			mmst.Render(w).JSON(200, mmst.P{})
		case name == "notfound":
			w.WriteHeader(http.StatusNotFound)
		default:
			mmst.Render(w).HTML(name, nil)
		}
	})
	r.GET("/slow/(?P<msec>[0-9]+)", func(w http.ResponseWriter, r *http.Request) {
		msec, err := strconv.Atoi(r.FormValue("msec"))
		if err != nil {
			w.WriteHeader(404)
			return
		}
		time.Sleep(time.Duration(msec) * time.Millisecond)
		mmst.Render(w).HTML("02", nil)
	})
	return httptest.NewServer(r)
}
