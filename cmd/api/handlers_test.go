package main

import (
	"io"
	"net/http"
	"strings"
	"testing"
)

// func TestGetTaskById(t *testing.T) {
// 	app := newTestApplication(t)

// 	ts := newTestServer(t, app.routes())
// 	defer ts.Close()

// 	testToken := "testtoken"

// 	req := httptest.NewRequest(http.MethodGet, ts.URL+"/api/v1/tasks/idone", nil)
// 	req.Header.Set("Authorization", "Bearer testtoken")

// 	resp, err := ts.Client().Do(req)
// 	if err != nil {
// 		t.Error(err)
// 	}

// 	// tests := []struct {
// 	// 	name     string
// 	// 	urlPath  string
// 	// 	wantCode int
// 	// 	wantBody string
// 	// }{
// 	// 	{
// 	// 		name:     "Valid ID",
// 	// 		urlPath:  "/api/v1/tasks/idone",
// 	// 		wantCode: http.StatusOK,
// 	// 		wantBody: "Green Mocks & Ham",
// 	// 	},
// 	// 	// {
// 	// 	// 	name:     "Non-existent ID",
// 	// 	// 	urlPath:  "/api/v1/tasks/non-exist",
// 	// 	// 	wantCode: http.StatusNotFound,
// 	// 	// },
// 	// 	// {
// 	// 	// 	name:     "Negative ID",
// 	// 	// 	urlPath:  "/api/v1/tasks/-1",
// 	// 	// 	wantCode: http.StatusNotFound,
// 	// 	// },
// 	// 	// {
// 	// 	// 	name:     "Decimal ID",
// 	// 	// 	urlPath:  "/api/v1/tasks/1.23",
// 	// 	// 	wantCode: http.StatusNotFound,
// 	// 	// },
// 	// 	// {
// 	// 	// 	name:     "Empty ID",
// 	// 	// 	urlPath:  "/api/v1/tasks/",
// 	// 	// 	wantCode: http.StatusNotFound,
// 	// 	// },
// 	// }

// 	// for _, tt := range tests {
// 	// 	t.Run(tt.name, func(t *testing.T) {
// 	// 		code, _, body := ts.makeGetRequestWithToken(t, tt.urlPath, testToken)

// 	// 		fmt.Println("inside of range loop in testtaskbyid")

// 	// 		assert.Equal(t, code, tt.wantCode)

// 	// 		fmt.Println("after assert.Equal of range loop in testtaskbyid")

// 	// 		if tt.wantBody != "" {
// 	// 			assert.StringContains(t, body, tt.wantBody)
// 	// 		}

// 	// 		fmt.Println("after assert.StringContains of range loop in testtaskbyid")

// 	// 	})
// 	// }
// }

func TestGetTaskById(t *testing.T) {
	app := newTestApplication(t)

	ts := newTestServer(t, app.routes())
	defer ts.Close()

	testToken := "testtoken"

	req, err := http.NewRequest(http.MethodGet, ts.URL+"/api/v1/tasks/idone", nil)
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Authorization", "Bearer "+testToken)

	resp, err := ts.Client().Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		t.Errorf("want status %d; got %d", http.StatusOK, resp.StatusCode)
	}

	// Optionally read and check the body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	// Example: check that the body contains some expected string
	if !strings.Contains(string(body), "expected content") {
		t.Errorf("response body does not contain expected content; got: %s", string(body))
	}
}
