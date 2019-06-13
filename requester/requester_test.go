package requester_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	. "prometheus-metrics-exporter/pmeerrors"
	. "prometheus-metrics-exporter/requester"
	"testing"
	"time"
)

const (
	mimeType      = "json"
	timeoutInSecs = 10
)

// Timeout test
func Test_GetContent_Timeout(t *testing.T) {
	// for the fake http server
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		time.Sleep(time.Second * 2)
	}

	// creating fake test server
	ts := httptest.NewServer(http.HandlerFunc(handler))

	defer ts.Close()

	_, _, err := GetContent(ts.URL, mimeType, 1)

	if err != nil && err == err.(ErrorRequestTimeOut) {
		t.Log("OK. Timeout as expected")
	} else if err != nil && err != err.(*ErrorRequestTimeOut) {
		t.Fatalf("Test failed with an unexpexted error: %s", err)
	} else {
		t.Fatal("Test failed unexpectedly.")
	}

}

// Response status code 404
func Test_GetContent_ResponseCode404(t *testing.T) {

	// for the fake http server
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}

	// creating fake test server
	ts := httptest.NewServer(http.HandlerFunc(handler))
	defer ts.Close()

	_, _, err := GetContent(ts.URL, mimeType, timeoutInSecs)

	if err != nil && err == err.(ErrorRequestResponseStatus404) {
		t.Log("OK. As expected the response code was 404")
	} else if err != nil && err != err.(ErrorRequestResponseStatus404) {
		t.Fatalf("Test failed with an unexpexted error: %s", err)
	} else {
		t.Fatal("Test failed unexpectedly.")
	}

}

// Response status code 500
func Test_GetContent_ResponseCode500(t *testing.T) {

	// for the fake http server
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}

	// creating fake test server
	ts := httptest.NewServer(http.HandlerFunc(handler))
	defer ts.Close()

	_, _, err := GetContent(ts.URL, mimeType, timeoutInSecs)

	if err != nil && err == err.(ErrorRequestResponseStatus500) {
		t.Log("OK. As expected the response code was 500")
	} else if err != nil && err != err.(ErrorRequestResponseStatus500) {
		t.Fatalf("Test failed with an unexpexted error: %s", err)
	} else {
		t.Fatal("Test failed unexpectedly.")
	}

}

// Accepted response status code
func Test_GetContent_ResponseCodeNot200(t *testing.T) {
	// for the fake http server
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadGateway)
	}

	// creating fake test server
	ts := httptest.NewServer(http.HandlerFunc(handler))
	defer ts.Close()

	_, _, err := GetContent(ts.URL, mimeType, timeoutInSecs)

	if err != nil && err == err.(ErrorRequestResponseStatusNot200) {
		t.Log("OK. As expected the response code was not 200")
	} else if err != nil && err != err.(ErrorRequestResponseStatusNot200) {
		t.Fatalf("Method failed unexpectedly: %s", err)
	} else {
		t.Fatal("Test failed unexpectedly.")
	}

}

// status 200 with no content type
func Test_GetContent_ResponseCode200NoContentType(t *testing.T) {
	// for the fake http server
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}

	// creating fake test server
	ts := httptest.NewServer(http.HandlerFunc(handler))
	defer ts.Close()

	_, _, err := GetContent(ts.URL, mimeType, timeoutInSecs)

	if err != nil && err == err.(ErrorRequestContentTypeParse) {
		t.Log("OK. As expected no content type found")
	} else if err != nil && err != err.(ErrorRequestContentTypeParse) {
		t.Fatalf("Method failed unexpectedly: %s", err)
	} else {
		t.Fatal("Test failed unexpectedly.")
	}

}

// status 200 with a not accepted content type
func Test_GetContent_NotAcceptedContentType(t *testing.T) {
	// for the fake http server
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
	}

	// creating fake test server
	ts := httptest.NewServer(http.HandlerFunc(handler))
	defer ts.Close()

	_, _, err := GetContent(ts.URL, mimeType, timeoutInSecs)

	if err != nil && err == err.(ErrorRequestInvalidContentTypeFound) {
		t.Log("OK. As expected invalid content type found")
	} else if err != nil {
		t.Fatalf("Method failed unexpectedly: %s", err)
	} else {
		t.Fatal("Test failed unexpectedly.")
	}

}

// Request body read error
func Test_GetContent_RequestBodyReadError(t *testing.T) {

	// for the fake http server
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Content-Length", "1")
		w.WriteHeader(http.StatusOK)
	}

	// creating fake test server
	ts := httptest.NewServer(http.HandlerFunc(handler))
	defer ts.Close()

	_, _, err := GetContent(ts.URL, mimeType, timeoutInSecs)

	if err != nil && err == err.(ErrorRequestUnableToReadBody) {
		t.Log("OK. Unable to read body as expected.")
	} else if err != nil {
		t.Fatalf("Method failed unexpectedly: %s", err)
	} else {
		t.Fatal("Test failed unexpectedly.")
	}

}

// Happy path, all ok, response 200, ok content type and config etc.
func Test_GetContent_OK(t *testing.T) {

	responseOKBytes := []byte(`{"response": "ok"}`)

	// for the fake http server
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, err := w.Write(responseOKBytes)

		if err != nil {
			t.Fatalf("Test failed unexpectedly: %s", err.Error())
		}
	}

	// creating fake test server
	ts := httptest.NewServer(http.HandlerFunc(handler))
	defer ts.Close()

	receivedContent, receivedMimeType, err := GetContent(ts.URL, mimeType, timeoutInSecs)

	if err == nil && bytes.Compare(receivedContent, responseOKBytes) == 0 && receivedMimeType == mimeType {
		t.Log("OK. Ideal line of events")
	} else if err != nil {
		t.Fatalf("Method failed unexpectedly: %s", err)
	} else if receivedMimeType != mimeType {
		t.Fatalf("Unexpected mime type. Expected \"%s\" but got \"%s\"", mimeType, receivedMimeType)
	} else if bytes.Compare(receivedContent, responseOKBytes) == 0 {
		t.Fatalf("Mismatching content. Expected \"%s\" but got \"%s\"", responseOKBytes, receivedContent)
	} else {
		t.Fatal("Test failed unexpectedly.")
	}

}
