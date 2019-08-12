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

import "net/http"

// Client a Tencent cloud SMS client
type Client struct {
	client      *http.Client // http client
	appID       string
	appKey      string
	IDGenerator func() string
}

const (
	smsNotification = iota
	smsMarket
)

const (
	qSmsBase                = `https://yun.tim.qq.com/v5/`
	singleSmsRoutine        = `tlssmssvr/sendsms`
	multiSmsRoutine         = `tlssmssvr/sendmultisms2`
	statusPullRoutine       = `tlssmssvr/pullstatus`
	mobileStatusPullRoutine = `tlssmssvr/pullstatus4mobile`
	voicePromptRoutine      = `tlsvoicesvr/sendvoiceprompt`
	voiceRoutine            = `tlsvoicesvr/sendcvoice`
	uploadVoiceFileRoutine  = `tlsvoicesvr/uploadvoicefile`
	sendVoiceFileRoutine    = `tlsvoicesvr/sendfvoice`
)

type number struct {
	CountryCode    string `json:"nationcode"`
	NationalNUmber string `json:"mobile"`
}

// smsSingleMessage wrapper of SMS message with only one receiver
type smsSingleMessage struct {
	MessageType    int      `json:"type"`
	Receiver       number   `json:"tel"`
	Message        string   `json:"msg,omitempty"`
	Extend         string   `json:"extend"`
	Time           int64    `json:"time"`
	LocalMessageID string   `json:"ext"`
	Digest         string   `json:"sig"`
	Sign           string   `json:"sign,omitempty"`
	TemplateID     string   `json:"tpl_id,omitempty"`
	TemplateParams []string `json:"params,omitempty"`
	random         int
}

// smsMultiMessage wrapper of SMS message with multiple receivers
type smsMultiMessage struct {
	MessageType    int      `json:"type"`
	Receivers      []number `json:"tel"`
	Message        string   `json:"msg,omitempty"`
	Extend         string   `json:"extend"`
	Time           int64    `json:"time"`
	LocalMessageID string   `json:"ext"`
	Digest         string   `json:"sig"`
	Sign           string   `json:"sign,omitempty"`
	TemplateID     string   `json:"tpl_id,omitempty"`
	TemplateParams []string `json:"params,omitempty"`
	random         int
}

func (message smsMultiMessage) GetRandom() int {
	return message.random
}

// Result send result of single SMS
// {
//   "result": 0,
//   "errmsg": "OK",
//   "ext": "6b6a09e7-be31-48cd-8171-0057ccbe550d",
//   "sid": "2018:-9173686749108840660",
//   "fee": 1
// }
type Result struct {
	Status         int    `json:"result"`
	ErrMsg         string `json:"errmsg"`
	LocalMessageID string `json:"ext"`
	MessageID      string `json:"sid,omitempty"`
	Fee            int    `json:"fee,omitempty"`
}

// MultiResult result for multiple receivers
type MultiResult struct {
	Status         int          `json:"result"`
	ErrMsg         string       `json:"errMsg"`
	LocalMessageID string       `json:"ext,omitempty"`
	CallID         string       `json:"callid,omitempty"`
	ReportCount    int          `json:"count,omitempty"`
	MTResults      []SendResult `json:"detail,omitempty"`
	Reports        []Report     `json:"data,omitempty"`
	FileID         string       `json:"fid,omitempty"`
}

// SendResult result for single sms sending
type SendResult struct {
	Status      int    `json:"result"`
	ErrMsg      string `json:"errmsg"`
	CountryCode string `json:"nationcode,omitempty"`
	Mobile      string `json:"mobile,omitempty"`
	MessageID   string `json:"sid,omitempty"`
	Fee         int    `json:"fee,omitempty"`
}

// Report represents the SMS Report message
type Report struct {
	Status        string `json:"report_status,omitempty"`
	DeliverTime   string `json:"user_receive_time,omitempty"`
	CountryCode   string `json:"nationcode,omitempty"`
	Mobile        string `json:"mobile,omitempty"`
	MessageID     string `json:"sid,omitempty"`
	DeliverStatus string `json:"errmsg,omitempty"`
	Description   string `json:"description,omitempty"`
	PullType      int    `json:"pull_type,omitempty"`
}

// smsStatus request for status check
type smsStatus struct {
	Digest      string `json:"sig"`
	Time        int64  `json:"time"`
	MessageType int    `json:"type"`
	MaxNumber   int64  `json:"max"`
	random      int
}

type mobileStatus struct {
	Digest         string `json:"sig"`
	MessageType    int    `json:"type"`
	Time           int64  `json:"time"`
	MaxNumber      int64  `json:"max"`
	StartTime      int64  `json:"begin_time"`
	EndTime        int64  `json:"end_time"`
	CountryCode    string `json:"nationcode"`
	NationalNumber string `json:"mobile"`
	random         int
}

type voicePrompt struct {
	Receiver       number `json:"tel"`
	PromptType     int    `json:"prompttype"`
	Message        string `json:"promptfile"`
	Duration       int64  `json:"playtimes"`
	Digest         string `json:"sig"`
	LocalMessageID string `json:"ext"`
	Time           int64  `json:"time"`
	random         int
}

type voice struct {
	Receiver       number   `json:"tel"`
	TemplateID     string   `json:"tpl_id,omitempty"`
	TemplateParams []string `json:"params,omitempty"`
	FileID         string   `json:"fid,omitempty"`
	Duration       int64    `json:"playtimes"`
	Digest         string   `json:"sig"`
	LocalMessageID string   `json:"ext"`
	Time           int64    `json:"time"`
	random         int
}
