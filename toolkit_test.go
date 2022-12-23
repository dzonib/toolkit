package toolkit

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

var jsonTests = []struct {
	name          string
	json          string
	errorExpected bool
	maxSize       int
	allowUnknown  bool
}{
	{name: "good json", json: `{"foo": "bar"}`, errorExpected: false, maxSize: 1024, allowUnknown: false},
	{name: "badly formatted json", json: `{"foo":"}`, errorExpected: true, maxSize: 1024, allowUnknown: false},
	{name: "incorrect type", json: `{"foo": 1}`, errorExpected: true, maxSize: 1024, allowUnknown: false},
	{name: "incorrect type", json: `{1: 1}`, errorExpected: true, maxSize: 1024, allowUnknown: false},
	{name: "two json files", json: `{"foo": "bar"}{"alpha": "beta"}`, errorExpected: true, maxSize: 1024, allowUnknown: false},
	{name: "empty body", json: ``, errorExpected: true, maxSize: 1024, allowUnknown: false},
	{name: "syntax error in json", json: `{"foo": 1"}`, errorExpected: true, maxSize: 1024, allowUnknown: false},
	{name: "unknown field in json", json: `{"fooo": "bar"}`, errorExpected: true, maxSize: 1024, allowUnknown: false},
	{name: "allow unknown field in json", json: `{"fooo": "bar"}`, errorExpected: false, maxSize: 1024, allowUnknown: true},
	{name: "missing field name", json: `{jack: "bar"}`, errorExpected: true, maxSize: 1024, allowUnknown: false},
	{name: "file too large", json: `{"foo": "bar"}`, errorExpected: true, maxSize: 5, allowUnknown: false},
	{name: "not json", json: `Hello, world`, errorExpected: true, maxSize: 1024, allowUnknown: false},
}

func Test_ReadJSON(t *testing.T) {
	var testTools Tools

	for _, e := range jsonTests {
		// set max file size
		testTools.MaxJSONSize = e.maxSize

		// allow/disallow unknown fields
		testTools.AllowUnknownFields = e.allowUnknown

		// declare a variable to read the decoded json into
		var decodedJSON struct {
			Foo string `json:"foo"`
		}

		// create a request with the body
		req, err := http.NewRequest("POST", "/", bytes.NewReader([]byte(e.json)))
		if err != nil {
			t.Log("Error", err)
		}

		// create a test response recorder, which satisfies the requirements
		// for a ResponseWriter
		rr := httptest.NewRecorder()


		// run a sub-test to call readJSON and check for an error
		t.Run(e.name, func(t *testing.T) {
			err = testTools.ReadJSON(rr, req, &decodedJSON)

			// if we expect an error, but do not get one, something went wrong
			if e.errorExpected && err == nil {
				t.Errorf("error expected, but none received")
			}

			// if we do not expect an error, but get one, something went wrong
			if !e.errorExpected && err != nil {
				t.Errorf("error not expected, but one received: %s", err.Error())
			}
		})
		req.Body.Close()
	}
}


func TestTools_WriteJSON(t *testing.T) {
	// create a variable of type toolbox.Tools, and just use the defaults.
	var testTools Tools

	rr := httptest.NewRecorder()
	payload := JSONResponse{
		Error:   false,
		Message: "foo",
	}

	headers := make(http.Header)

	headers.Add("FOO", "BAR")

	err := testTools.WriteJSON(rr, http.StatusOK, payload, headers)
	
	if err != nil {
		t.Errorf("failed to write JSON: %v", err)
	}
}