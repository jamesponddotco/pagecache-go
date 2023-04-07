package pagecache

import (
	"bufio"
	"bytes"
	"fmt"
	"net/http"
	"net/http/httputil"

	"git.sr.ht/~jamesponddotco/xstd-go/xerrors"
)

// ErrInvalidResponse is returned when a response is invalid.
const ErrInvalidResponse xerrors.Error = "invalid response"

// SaveResponse saves an HTTP response and the request that generated it.
func SaveResponse(resp *http.Response) (request, response []byte, err error) {
	if resp == nil || resp.Request == nil || resp.Request.URL == nil {
		return nil, nil, fmt.Errorf("%w", ErrInvalidResponse)
	}

	request, err = httputil.DumpRequestOut(resp.Request, true)
	if err != nil {
		return nil, nil, fmt.Errorf("%w", err)
	}

	response, err = httputil.DumpResponse(resp, true)
	if err != nil {
		return nil, nil, fmt.Errorf("%w", err)
	}

	return request, response, nil
}

// LoadResponse loads an HTTP response from a saved request and response.
func LoadResponse(request, response []byte) (*http.Response, error) {
	req, err := http.ReadRequest(bufio.NewReader(bytes.NewReader(request)))
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	resp, err := http.ReadResponse(bufio.NewReader(bytes.NewReader(response)), req)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return resp, nil
}
