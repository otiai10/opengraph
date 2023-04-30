package opengraph

import (
	"context"
	"github.com/otiai10/opengraph/v2/http_fetchers"
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
	"golang.org/x/net/html"
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
	When(t, "strict flag is on", func(t *testing.T) {
		og, err := Fetch(testserver.URL+"/case/03_image", Intent{Strict: true})
		Expect(t, err).ToBe(nil)
		Expect(t, og.Title).ToBe("")
		og, err = Fetch(testserver.URL+"/case/03_image", Intent{Strict: false})
		Expect(t, err).ToBe(nil)
		Expect(t, og.Title).ToBe("TEST: 03_image (title tag)")
	})
	When(t, "context is specified", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
		intent := Intent{
			Context: ctx,
		}
		ogp, err := Fetch(testserver.URL+"/slow/100", intent)
		cancel()
		Expect(t, err).ToBe(nil)
		Expect(t, ogp.Title).ToBe("Hello! Open Graph!!")
		ctx, cancel = context.WithTimeout(context.Background(), 200*time.Millisecond)
		_, err = Fetch(testserver.URL+"/slow/300", intent)
		cancel()
		Expect(t, err).Not().ToBe(nil)
	})
}

func TestOpenGraph_Fetch(t *testing.T) {
	ogp := &OpenGraph{}
	err := ogp.Fetch()
	Expect(t, err.Error()).ToBe("no URL given yet")

	ogp.Intent.URL = testserver.URL + "/case/01_hello"
	err = ogp.Fetch()
	Expect(t, err).ToBe(nil)

	When(t, "ogp already has an error", func(t *testing.T) {
		ogp := New(":INVALID_URL")
		err := ogp.Fetch()
		Expect(t, err).Not().ToBe(nil)
	})

	When(t, "custom http client is given", func(t *testing.T) {
		ogp.Intent.HTTPFetcher = http_fetchers.DefaultSimpleHTTPFetcher
		err := ogp.Fetch()
		Expect(t, err).ToBe(nil)
	})
}

func TestFetchVideo(t *testing.T) {
	og, err := Fetch(testserver.URL + "/case/05_video")
	Expect(t, err).ToBe(nil)
	Because(t, "og:video:url - Identical to og:video", func(t *testing.T) {
		Expect(t, len(og.Video)).ToBe(1)
		Expect(t, og.Video).ToBe([]Video{
			{
				URL:       "https://www.youtube.com/embed/1MxA0i2rxQo",
				SecureURL: "https://www.youtube.com/embed/1MxA0i2rxQo",
				Type:      "text/html",
				Width:     1280,
				Height:    720,
			},
		})
	})
}

func TestOpenGraph_Parse(t *testing.T) {
	ogp := New("")
	r := strings.NewReader(`<html><meta property="og:title" content="test_test"></html>`)
	err := ogp.Parse(r)
	Expect(t, err).ToBe(nil)
	Expect(t, ogp.Title).ToBe("test_test")
}

func TestOpenGraph_Walk(t *testing.T) {
	res, err := http.Get(testserver.URL + "/case/01_hello")
	Expect(t, err).ToBe(nil)

	node, err := html.Parse(res.Body)
	Expect(t, err).ToBe(nil)

	og := New(testserver.URL + "/case/01_hello")
	err = og.Walk(node)
	Expect(t, err).ToBe(nil)
	Expect(t, og.Title).ToBe("Hello! Open Graph!!")
	res.Body.Close()

	again := New(testserver.URL + "/case/01_hello")
	err = again.Walk(node)
	Expect(t, err).ToBe(nil)
	Expect(t, again.Title).ToBe("Hello! Open Graph!!")
}

func TestOpenGraph_ToAbs(t *testing.T) {
	ogp := New(testserver.URL + "/case/01_hello")
	err := ogp.Fetch()
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

	ogp = New(testserver.URL + "/case/03_image")
	err = ogp.Fetch()
	Expect(t, err).ToBe(nil)
	err = ogp.ToAbs()
	Expect(t, err).ToBe(nil)
	u, err = url.Parse(ogp.Image[0].URL)
	Expect(t, err).ToBe(nil)
	Expect(t, u.IsAbs()).ToBe(true)
	Expect(t, u.Host).ToBe("www-cdn.jtvnw.net")
	Expect(t, u.Path).ToBe("/images/twitch_logo3.jpg")
	Expect(t, u.String()).ToBe("http://www-cdn.jtvnw.net/images/twitch_logo3.jpg")
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
		mmst.Render(w).HTML("01_hello", nil)
	})
	return httptest.NewServer(r)
}
