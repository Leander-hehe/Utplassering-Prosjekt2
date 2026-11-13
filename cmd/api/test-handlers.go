package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
)

func (app *application) TestAllHandlers(w http.ResponseWriter, r *http.Request) { // Convert the byte slice to `io.Reader`

	Testhandler(app.GetHandler, httptest.NewRequest("GET", "/", nil), w)

	TestHandlerWithBody(app.PostHandler, "POST", "/", map[string]string{"id": "test", "name": "test", "description": "test"}, w)

	Testhandler(app.DeleteHandler, httptest.NewRequest("DELETE", "/1", nil), w)

	TestHandlerWithBody(app.PutHandler, "PUT", "/1", map[string]string{"id": "test", "name": "test", "description": "test"}, w)

	// !!! Item with id=1 might not exist
	TestHandlerWithBody(app.PatchHandler, "PATCH", "/1", [][]string{{"name", "test"}}, w)
}

func Testhandler(handler func(w http.ResponseWriter, r *http.Request), req *http.Request, testResp http.ResponseWriter) {

	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Test", "true")

	res := httptest.NewRecorder()

	handler(res, req)

	resp := res.Result()

	fmt.Println(req.Method + ": " + strconv.Itoa(resp.StatusCode))
	fmt.Fprintln(testResp, req.Method+": "+strconv.Itoa(resp.StatusCode))

}

func TestHandlerWithBody(handler func(w http.ResponseWriter, r *http.Request), method string, target string, input any, testResp http.ResponseWriter) {
	body, err := createBody(input)
	if err == nil {
		Testhandler(handler, httptest.NewRequest(method, target, body), testResp)
	} else {
		fmt.Fprintln(testResp, method+": 400")
	}
}

func createBody(input any) (*bytes.Reader, error) {
	bodyBytes, err := json.Marshal(input) // Create a byte slice of the JSON data
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(bodyBytes), nil // Convert the byte slice to `io.Reader`
}
