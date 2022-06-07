//  This file is part of the eliona project.
//  Copyright © 2022 LEICOM iTEC AG. All Rights Reserved.
//  ______ _ _
// |  ____| (_)
// | |__  | |_  ___  _ __   __ _
// |  __| | | |/ _ \| '_ \ / _` |
// | |____| | | (_) | | | | (_| |
// |______|_|_|\___/|_| |_|\__,_|
//
//  THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING
//  BUT NOT LIMITED  TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
//  NON INFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM,
//  DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
//  OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package http

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/eliona-smart-building-assistant/go-eliona/log"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

// NewRequestWithBarrier creates a new request for the given url. The url have is authenticated with a barrier token.
func NewRequestWithBarrier(url string, token string) (*http.Request, error) {
	return newRequestWithBarrier(url, nil, "GET", token)
}

// NewPostRequestWithBarrier creates a new request for the given url. The url have is authenticated with a barrier token.
func NewPostRequestWithBarrier(url string, body any, token string) (*http.Request, error) {
	return newRequestWithBarrier(url, body, "POST", token)
}

func newRequestWithBarrier(url string, body any, method string, token string) (*http.Request, error) {

	// Create a new request
	request, err := newRequest(url, body, method)
	if err != nil {
		log.Error("Http", "Error creating request %s: %v", url, err)
		return nil, err
	}

	request.Header.Set("Authorization", "Bearer "+token)
	return request, nil
}

// NewRequest creates a new request for the given url. The url have to provide free access without any
// authentication. For authentication use other functions like NewGetRequestWithBarrier.
func NewRequest(url string) (*http.Request, error) {
	return newRequest(url, nil, "GET")
}

// NewPostRequest creates a new request for the given url and the body as payload. The url have to provide free access without any
// authentication. For authentication use other functions like NewPostRequestWithBarrier.
func NewPostRequest(url string, body any) (*http.Request, error) {
	return newRequest(url, body, "POST")
}

func newRequest(url string, body any, method string) (*http.Request, error) {

	// Create payload if used
	var payload []byte = nil
	if body != nil {
		var err error
		payload, err = json.Marshal(body)
		if err != nil {
			log.Error("Kafka", "Failed to marshal body: %s", err.Error())
			return nil, err
		}
	}

	// Create a new request with payload if used
	var buffer *bytes.Buffer = nil
	if body != nil {
		buffer = bytes.NewBuffer(payload)
	}

	request, err := http.NewRequest(method, url, buffer)
	if err != nil {
		log.Error("Http", "Error creating request %s: %v", url, err)
		return nil, err
	}

	return request, nil
}

// Read returns the payload returned from the request
func Read(request *http.Request, timeout int, checkCertificate bool) ([]byte, error) {

	// creates a http client with timeout and tsl security configuration
	httpClient := http.Client{
		Timeout: time.Duration(timeout) * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: !checkCertificate},
		},
	}

	// start the request
	response, err := httpClient.Do(request)
	if err != nil {
		log.Error("Http", "Error request to %s: %v", request.URL, err)
		return nil, err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Error("Http", "Error closing request for %s: %v", request.URL, err)
		}
	}(response.Body)

	// read the complete payload
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Error("Http", "Error body from %s: %v", request.URL, err)
		return nil, err
	}

	// returns the payload as string, if the status code is OK
	if response.StatusCode == http.StatusOK {
		return body, nil
	} else {
		log.Error("Http", "Error request code %d for request to %s.", response.StatusCode, request.URL)
		return nil, fmt.Errorf("error request code %d for request to %s", response.StatusCode, request.URL)
	}
}
