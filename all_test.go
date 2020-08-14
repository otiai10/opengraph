package opengraph

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

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
	Expect(t, og.URL.Value).ToBe("https://haisai.party/")

	b := bytes.NewBuffer(nil)
	json.NewEncoder(b).Encode(og)
	Expect(t, strings.Trim(b.String(), "\n")).ToBe(fmt.Sprintf(
		`{"Policy":{"TrustedTags":["meta","link","title"]},"Title":"はいさいナイト","Type":"website","URL":{"Source":"%s","Scheme":"http","Opaque":"","User":null,"Host":"%s","Path":"","RawPath":"","ForceQuery":false,"RawQuery":"","Fragment":"","Value":"%s"},"SiteName":"","Image":[],"Video":[],"Audio":[],"Description":"All Genre Music Party","Determiner":"","Locale":"","LocaleAlt":[],"Favicon":"/favicon.ico"}`,
		s.URL,
		strings.Replace(s.URL, "http://", "", -1),
		og.URL.Value,
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

func TestFetch_04(t *testing.T) {
	s := dummyServer(4)
	og, err := Fetch(s.URL)
	Expect(t, err).ToBe(nil)
	Expect(t, len(og.Image)).ToBe(1)
	Expect(t, og.Image[0].URL).ToBe("/images/01.png")
}

func TestFetch_05_TrustedTags(t *testing.T) {
	s := dummyServer(2)
	res, err := http.Get(s.URL)
	Expect(t, err).ToBe(nil)
	defer res.Body.Close()
	og := New(s.URL)
	og.Policy.TrustedTags = []string{HTMLMetaTag}
	err = og.Parse(res.Body)
	Expect(t, err).ToBe(nil)
	Expect(t, og.Title).Not().ToBe("はいさいナイト")
	Expect(t, og.Title).ToBe("")
}

func TestFetchWithContext(t *testing.T) {
	s := dummySlowServer(time.Millisecond * 300)
	defer s.Close()

	ctx500ms, cancel500ms := context.WithTimeout(context.Background(), time.Millisecond*500)
	_, err := FetchWithContext(ctx500ms, s.URL)
	Expect(t, err).ToBe(nil)
	cancel500ms()

	ctx100ms, cancel100ms := context.WithTimeout(context.Background(), time.Millisecond*100)
	_, err = FetchWithContext(ctx100ms, s.URL)
	Expect(t, err).Match(context.DeadlineExceeded.Error())
	cancel100ms()
}

func dummyServer(id int) *httptest.Server {
	marmoset.LoadViews("./test/html")
	r := marmoset.NewRouter()
	r.GET("/", func(w http.ResponseWriter, r *http.Request) {
		marmoset.Render(w).HTML(fmt.Sprintf("%02d", id), nil)
	})
	r.GET("/case/01", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})
	return httptest.NewServer(r)
}

func dummySlowServer(d time.Duration) *httptest.Server {
	var h = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(d)
		t, _ := template.ParseFiles("./test/html/02.html")
		t.Execute(w, nil)
	})
	return httptest.NewServer(h)
}
