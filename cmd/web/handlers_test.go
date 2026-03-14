package main

import (
	"bytes"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"snippetbox.project.sope/internal/assert"
)

// func TestPing(t *testing.T){
// 	rr := httptest.NewRecorder()

// 	r, err := http.NewRequest(http.MethodGet, "/", nil)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	ping(rr, r)

// 	rs := rr.Result()

// 	assert.Equal(t, rs.StatusCode, http.StatusOK)

// 	defer rs.Body.Close()

// 	body, err := io.ReadAll(rs.Body)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	body = bytes.TrimSpace(body)

// 	assert.Equal(t, string(body), "OK")

// }

func TestPing(t *testing.T){

	app := &application{
		logger: slog.New(slog.DiscardHandler),
	}

	ts := httptest.NewTLSServer(app.routes())
	defer ts.Close()

	

	rs, err := ts.Client().Get(ts.URL + "/ping")
	if err != nil{
		t.Fatal(err)
	}

	assert.Equal(t, rs.StatusCode, http.StatusOK)

	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}

	body = bytes.TrimSpace(body)
	assert.Equal(t, string(body), "OK")

}