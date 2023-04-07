package pagecache_test

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"strings"
	"testing"

	"git.sr.ht/~jamesponddotco/pagecache-go"
)

func TestSaveResponse(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		handler     http.HandlerFunc
		modifyResp  func(*http.Response)
		expectedErr error
	}{
		{
			name: "successful_save",
			handler: func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprintln(w, "Hello, World!")
			},
			modifyResp:  nil,
			expectedErr: nil,
		},
		{
			name: "invalid_request",
			handler: func(w http.ResponseWriter, r *http.Request) {
				// Intentionally left empty.
			},
			modifyResp: func(resp *http.Response) {
				resp.Request.URL = nil
			},
			expectedErr: pagecache.ErrInvalidResponse,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ts := httptest.NewServer(tt.handler)
			defer ts.Close()

			resp, err := http.Get(ts.URL)
			if err != nil {
				t.Fatalf("failed to get response: %v", err)
			}

			if tt.modifyResp != nil {
				tt.modifyResp(resp)
			}

			reqBytes, respBytes, err := pagecache.SaveResponse(resp)

			if tt.expectedErr != nil {
				if err == nil || err.Error() != tt.expectedErr.Error() {
					t.Fatalf("expected error '%v', got '%v'", tt.expectedErr, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			dumpedReq, err := http.NewRequest(resp.Request.Method, ts.URL, http.NoBody)
			if err != nil {
				t.Fatalf("failed to create new request for comparison: %v", err)
			}
			dumpedReqBytes, err := httputil.DumpRequestOut(dumpedReq, true)
			if err != nil {
				t.Fatalf("failed to dump request for comparison: %v", err)
			}

			if !bytes.Equal(reqBytes, dumpedReqBytes) {
				t.Errorf("request dump mismatch: got %q, want %q", reqBytes, dumpedReqBytes)
			}

			respDump, err := httputil.DumpResponse(resp, true)
			if err != nil {
				t.Fatalf("failed to dump response for comparison: %v", err)
			}

			if !bytes.Equal(respBytes, respDump) {
				t.Errorf("response dump mismatch: got %q, want %q", respBytes, respDump)
			}
		})
	}
}

func TestLoadResponse(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		handler        http.HandlerFunc
		modifyRequest  func([]byte) []byte
		modifyResponse func([]byte) []byte
		expectedErr    error
	}{
		{
			name: "successful_load",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte("Hello, World!"))
			},
			modifyRequest:  nil,
			modifyResponse: nil,
			expectedErr:    nil,
		},
		{
			name: "invalid_request",
			handler: func(w http.ResponseWriter, r *http.Request) {
				// Intentionally left empty.
			},
			modifyRequest: func(req []byte) []byte {
				return req[:len(req)-1] // Remove the last byte to make it invalid.
			},
			modifyResponse: nil,
			expectedErr:    errors.New("malformed MIME header"),
		},
		{
			name: "invalid_response",
			handler: func(w http.ResponseWriter, r *http.Request) {
				// Intentionally left empty.
			},
			modifyRequest: nil,
			modifyResponse: func(resp []byte) []byte {
				return resp[:len(resp)-1] // Remove the last byte to make it invalid.
			},
			expectedErr: errors.New("malformed MIME header"),
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			server := httptest.NewServer(tt.handler)
			defer server.Close()

			req, err := http.NewRequest(http.MethodGet, server.URL, http.NoBody)
			if err != nil {
				t.Fatalf("unable to create request: %v", err)
			}

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Fatalf("unable to make request: %v", err)
			}

			requestDump, responseDump, err := pagecache.SaveResponse(resp)
			if err != nil {
				t.Fatalf("unable to save response: %v", err)
			}

			if tt.modifyRequest != nil {
				requestDump = tt.modifyRequest(requestDump)
			}

			if tt.modifyResponse != nil {
				responseDump = tt.modifyResponse(responseDump)
			}

			loadedResp, err := pagecache.LoadResponse(requestDump, responseDump)

			if tt.expectedErr != nil {
				if err == nil || !strings.Contains(err.Error(), tt.expectedErr.Error()) {
					t.Fatalf("expected error containing '%v', got '%v'", tt.expectedErr, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if loadedResp.StatusCode != resp.StatusCode {
				t.Errorf("status code mismatch: got %d, want %d", loadedResp.StatusCode, resp.StatusCode)
			}

			loadedBody, _ := io.ReadAll(loadedResp.Body)
			originalBody, _ := io.ReadAll(resp.Body)

			if !bytes.Equal(loadedBody, originalBody) {
				t.Errorf("body mismatch: got %q, want %q", loadedBody, originalBody)
			}
		})
	}
}
