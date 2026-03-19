package main

import (
	"bytes"
	"html"
	"io"
	"log/slog"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	"snippetbox.project.sope/internal/models/mocks"

	"testing"
)

var csrfTokenRX = regexp.MustCompile(`<input type='hidden' name='csrf_token' value='(.+)'>`)

func extractCSRFToken(t *testing.T, body string) string {
	matches := csrfTokenRX.FindStringSubmatch(body)
	
	if len(matches) < 2 {
		t.Fatal("no csrf token found in body")
	}

	return html.UnescapeString(string(matches[1]))
}

func newTestApplication(t *testing.T) *application {

	newTemplateCache, err := newTemplateCache()
	if err != nil{
		t.Fatal(err)
	}

	formDecoder := form.NewDecoder()

	sessionManger := scs.New()
	sessionManger.Lifetime = 12 * time.Hour
	sessionManger.Cookie.Secure = true
	


	return &application{
		logger: slog.New(slog.DiscardHandler),
		snippets: &mocks.SnippetModel{},
		users: &mocks.UserModel{},
		templateCache: newTemplateCache,
		formDecoder: formDecoder,
		sessionManager: sessionManger,
	}
}

type testServer struct{
	*httptest.Server
}

func newTestServer(t *testing.T, h http.Handler) *testServer{
	ts := httptest.NewTLSServer(h)

	jar , err:= cookiejar.New(nil)
	if err != nil{
		t.Fatal(err)
	}

	ts.Client().Jar = jar

	ts.Client().CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	return &testServer{ts}
}


func (ts *testServer) get(t *testing.T, urlPath string)(int, http.Header,string){
	rs, err := ts.Client().Get(ts.URL + urlPath)
	if err != nil {
		t.Fatal(err)
	}

	defer rs.Body.Close()

	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}

	body = bytes.TrimSpace(body)
	return rs.StatusCode, rs.Header, string(body)
}

func (ts *testServer) postForm(t *testing.T, urlPath string, form url.Values) (int, http.Header, string){
	r, err := http.NewRequest("POST", ts.URL + urlPath, strings.NewReader(form.Encode()))
	if err != nil {
		t.Fatal(err)
	}

	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Referer", ts.URL)

	rs, err := ts.Client().Do(r)
	if err != nil {
		t.Fatal(err)
	}

	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}

	body = bytes.TrimSpace(body)

	return rs.StatusCode, rs.Header, string(body)

}