package opengraph

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/otiai10/marmoset"
	. "github.com/otiai10/mint"
)

func TestNew(t *testing.T) {
	og := New(dummyServer(1).URL)
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

	s := dummyServer(1)
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

func TestFetch_02(t *testing.T) {
	s := dummyServer(2)
	og, err := Fetch(s.URL)

	Expect(t, err).ToBe(nil)
	Expect(t, og.Title).ToBe("はいさいナイト")
	Expect(t, og.Description).ToBe("All Genre Music Party")

	b := bytes.NewBuffer(nil)
	json.NewEncoder(b).Encode(og)
	Expect(t, strings.Trim(b.String(), "\n")).ToBe(fmt.Sprintf(
		`{"Title":"はいさいナイト","Type":"website","URL":{"Source":"%s","Scheme":"http","Opaque":"","User":null,"Host":"%s","Path":"","RawPath":"","ForceQuery":false,"RawQuery":"","Fragment":""},"SiteName":"","Image":[],"Video":[],"Audio":[],"Description":"All Genre Music Party","Determiner":"","Locale":"","LocaleAlt":[],"Favicon":"/favicon.ico"}`,
		s.URL,
		strings.Replace(s.URL, "http://", "", -1),
	))
}

func TestFetch_03(t *testing.T) {
	s := dummyServer(3)
	og, err := Fetch(s.URL)
	Expect(t, err).ToBe(nil)
	err = og.ToAbsURL().Fulfill()
	Expect(t, err).ToBe(nil)
	Expect(t, og.Image[0].URL).ToBe("http://www-cdn.jtvnw.net/images/twitch_logo3.jpg")
}

func dummyServer(id int) *httptest.Server {
	marmoset.LoadViews("./testdata/html")
	r := marmoset.NewRouter()
	r.GET("/", func(w http.ResponseWriter, r *http.Request) {
		marmoset.Render(w).HTML(fmt.Sprintf("%02d", id), nil)
	})
	r.GET("/case/01", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})
	return httptest.NewServer(r)
}
