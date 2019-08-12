// Copyright (c) 2019, Chen Lei <my@mysq.to>
// All rights reserved.

// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are met:

// 1. Redistributions of source code must retain the above copyright notice, this
//    list of conditions and the following disclaimer.
// 2. Redistributions in binary form must reproduce the above copyright notice,
//    this list of conditions and the following disclaimer in the documentation
//    and/or other materials provided with the distribution.

// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
// ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
// WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
// DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT OWNER OR CONTRIBUTORS BE LIABLE FOR
// ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
// (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
// LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
// ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
// SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

package qsms

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/mysqto/log"
	"github.com/mysqto/to"
)

func (c *Client) post(url string, payload interface{}) ([]byte, error) {

	data, err := json.Marshal(payload)

	if err != nil {
		log.Warnf("error marshal payload %v : %v", payload, err)
		return nil, err
	}

	body := bytes.NewBuffer(data)

	req, err := http.NewRequest("POST", url, body)

	if err != nil {
		log.Warnf("error creating POST request with %v : %v", url, err)
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json;charset=utf-8")

	resp, err := c.client.Do(req)
	if err != nil {
		log.Warnf("error doing POST request with %v : %v", url, err)
		return nil, err
	}

	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}

func getMIMEType(path string) string {

	mimeType := "application/octet-stream"

	file, err := os.Open(path)

	if err != nil {
		return mimeType
	}

	defer file.Close()

	// Only the first 512 bytes are used to sniff the content type.
	buffer := make([]byte, 512)

	_, err = file.Read(buffer)

	if err != nil {
		log.Warnf("error rading file :%v", err)
		return mimeType
	}

	mimeType = strings.ToLower(http.DetectContentType(buffer))

	// Always returns a valid content-type and "application/octet-stream" if no others seemed to match.
	return mimeType
}

func (c *Client) upload(url string, path string, headers map[string]interface{}) (*MultiResult, error) {
	file, err := os.Open(path)

	if err != nil {
		log.Warnf("error opening file %s : %v", path, err)
		return nil, err
	}

	defer file.Close()

	req, err := http.NewRequest("POST", url, file)

	for key, value := range headers {
		req.Header.Add(key, to.String(value))
	}

	if err != nil {
		log.Warnf("error creating POST request with %v : %v", url, err)
		return nil, err
	}

	resp, err := c.client.Do(req)

	if err != nil {
		log.Warnf("error doing POST request with %v : %v", url, err)
		return nil, err
	}

	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Warnf("error doing POST request with %v : %v", url, err)
		return nil, err
	}

	var result *MultiResult

	err = json.Unmarshal(data, &result)

	if err != nil {
		return nil, err
	}

	return result, nil
}

// sendMessage the real send routine
func (c *Client) sendMessage(message interface{}, url string, response interface{}) error {

	data, err := c.post(url, message)

	if err != nil {
		return err
	}

	err = json.Unmarshal(data, response)

	if err != nil {
		return err
	}

	return nil
}

// NewClient create a Tencent cloud SMS client
func NewClient(appID, appKey string) *Client {
	return &Client{
		appID:  appID,
		appKey: appKey,
		client: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}
