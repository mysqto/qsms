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
	"fmt"
	"time"
)

func (c *Client) newVoicePrompt(mobile string, content string, duration int64) (*voicePrompt, error) {
	rand := random()
	now := time.Now().Unix()

	countryCode, nationalNumber, err := parseNumber(mobile)

	if err != nil {
		return nil, err
	}

	return &voicePrompt{
		Receiver: number{
			CountryCode:    countryCode,
			NationalNUmber: nationalNumber,
		},
		PromptType:     2, // voice prompt type, currently value is 2
		Message:        content,
		Duration:       duration,
		Digest:         calSign(c.appKey, rand, now, []string{nationalNumber}),
		Time:           now,
		LocalMessageID: c.generateMessageID(),
		random:         rand,
	}, nil
}

// SendVoicePrompt send voice message
func (c *Client) SendVoicePrompt(mobile string, content string, duration int64) (*Result, error) {

	voicePrompt, err := c.newVoicePrompt(mobile, content, duration)

	if err != nil {
		return nil, err
	}

	var result Result

	err = c.sendMessage(voicePrompt, c.routine(voicePromptRoutine, voicePrompt.random), &result)

	return &result, err

}

func (c *Client) newVoice(mobile string, fileID string, templateID string, templateParams []string, duration int64) (*voice, error) {
	rand := random()
	now := time.Now().Unix()

	countryCode, nationalNumber, err := parseNumber(mobile)

	if err != nil {
		return nil, err
	}

	return &voice{
		Receiver: number{
			CountryCode:    countryCode,
			NationalNUmber: nationalNumber,
		},
		FileID:         fileID,
		TemplateID:     templateID,
		TemplateParams: templateParams,
		Duration:       duration,
		Digest:         calSign(c.appKey, rand, now, []string{nationalNumber}),
		Time:           now,
		LocalMessageID: c.generateMessageID(),
		random:         rand,
	}, nil
}

// SendVoiceWithTemplate send voice with template and template parameters
func (c *Client) SendVoiceWithTemplate(mobile string, templateID string,
	templateParams []string, duration int64) (*Result, error) {
	voice, err := c.newVoice(mobile, "", templateID, templateParams, duration)

	if err != nil {
		return nil, err
	}
	var result Result
	err = c.sendMessage(voice, c.routine(voiceRoutine, voice.random), &result)

	return &result, err
}

// SendVoiceWithFile send voice with file
func (c *Client) SendVoiceWithFile(mobile string, file string) (*Result, error) {

	rand := random()
	now := time.Now().Unix()

	contentType := getMIMEType(file)
	sha256Sum, err := shaSum(file)

	if err != nil {
		return nil, err
	}

	headers := map[string]interface{}{
		`Content-Type`:   contentType,
		`x-content-sha1`: sha256Sum,
		`Authorization`:  calAuth(c.appKey, rand, now, sha256Sum),
	}

	result, err := c.upload(c.routineWithTime(uploadVoiceFileRoutine, rand, now), file, headers)

	if err != nil {
		return nil, err
	}

	if result.Status != 0 {
		return nil, fmt.Errorf("error uploading %s : %v", file, result.ErrMsg)
	}

	voice, err := c.newVoice(mobile, result.FileID, "", nil, 3)

	if err != nil {
		return nil, err
	}
	var response Result

	err = c.sendMessage(voice, c.routine(sendVoiceFileRoutine, voice.random), &response)

	return &response, err
}
