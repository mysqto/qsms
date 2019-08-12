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
	"github.com/nyaruka/phonenumbers"
	"strconv"
	"time"
)

func (c *Client) newMultiSMS(smsType int, mobiles []string, message string, sign string, templateID string,
	parameters []string) (*smsMultiMessage, error) {

	rand := random()
	now := time.Now().Unix()

	numbers, receivers := parseNumbers(mobiles)

	return &smsMultiMessage{
		MessageType:    smsType,
		Receivers:      numbers,
		Message:        message,
		TemplateID:     templateID,
		TemplateParams: parameters,
		Extend:         "",
		LocalMessageID: c.generateMessageID(),
		Time:           now,
		Sign:           sign,
		Digest:         calSign(c.appKey, rand, now, receivers),
		random:         rand,
	}, nil
}

func (c *Client) newSMS(smsType int, mobile string, message string, sign string, templateID string,
	parameters []string) (*smsSingleMessage, error) {

	rand := random()
	now := time.Now().Unix()

	phoneNumber, err := phonenumbers.Parse(mobile, "")

	if err != nil {
		return nil, err
	}

	nationalNumber := strconv.FormatInt(int64(phoneNumber.GetNationalNumber()), 10)

	return &smsSingleMessage{
		MessageType: smsType,
		Receiver: number{
			CountryCode:    strconv.FormatInt(int64(phoneNumber.GetCountryCode()), 10),
			NationalNUmber: nationalNumber,
		},
		Message:        message,
		TemplateID:     templateID,
		TemplateParams: parameters,
		Extend:         "",
		LocalMessageID: c.generateMessageID(),
		Time:           now,
		Sign:           sign,
		Digest:         calSign(c.appKey, rand, now, []string{nationalNumber}),
		random:         rand,
	}, nil
}

// SendSMS send single notification SMS without template and sign, which is available for enterprise account
func (c *Client) SendSMS(mobile, content string) (*Result, error) {
	smsMessage, _ := c.newSMS(smsNotification, mobile, content, "", "", nil)
	var result Result
	err := c.sendMessage(smsMessage, c.routine(singleSmsRoutine, smsMessage.random), &result)
	return &result, err
}

// SendSMSWithTemplate send single notification SMS with template and sign, which is default for personal account
func (c *Client) SendSMSWithTemplate(mobile string, sign string, templateID string,
	parameters []string) (*Result, error) {
	smsMessage, _ := c.newSMS(smsNotification, mobile, "", sign, templateID, parameters)
	var result Result
	err := c.sendMessage(smsMessage, c.routine(singleSmsRoutine, smsMessage.random), &result)
	return &result, err

}

// SendMultipleSMS multiple receiver version of SendSMS
func (c *Client) SendMultipleSMS(mobile []string, content string) (*MultiResult, error) {
	smsMessage, _ := c.newMultiSMS(smsNotification, mobile, content, "", "", nil)
	var result MultiResult
	err := c.sendMessage(smsMessage, c.routine(multiSmsRoutine, smsMessage.random), &result)
	return &result, err
}

// SendMultipleSMSWithTemplate multiple receiver version of SendSMSWithTemplate
func (c *Client) SendMultipleSMSWithTemplate(mobile []string, sign, templateID string,
	parameters []string) (*MultiResult, error) {
	smsMessage, _ := c.newMultiSMS(smsNotification, mobile, "", sign, templateID, parameters)
	var result MultiResult
	err := c.sendMessage(smsMessage, c.routine(multiSmsRoutine, smsMessage.random), &result)
	return &result, err
}

// SendMarketSMS send single market SMS without template and sign, which is available for enterprise account
func (c *Client) SendMarketSMS(mobile, content string) (*Result, error) {
	smsMessage, _ := c.newSMS(smsMarket, mobile, content, "", "", nil)
	var result Result
	err := c.sendMessage(smsMessage, c.routine(singleSmsRoutine, smsMessage.random), &result)
	return &result, err
}

// SendMarketSMSWithTemplate send single market SMS with template and sign, which is default for personal account
func (c *Client) SendMarketSMSWithTemplate(mobile string, sign string, templateID string,
	parameters []string) (*Result, error) {
	smsMessage, _ := c.newSMS(smsMarket, mobile, "", sign, templateID, parameters)
	var result Result
	err := c.sendMessage(smsMessage, c.routine(singleSmsRoutine, smsMessage.random), &result)
	return &result, err
}

// SendMultipleMarketSMS multiple receiver version of SendMarketSMS
func (c *Client) SendMultipleMarketSMS(mobile []string,
	content string) (*MultiResult, error) {
	smsMessage, _ := c.newMultiSMS(smsMarket, mobile, content, "", "", nil)
	var result MultiResult
	err := c.sendMessage(smsMessage, c.routine(multiSmsRoutine, smsMessage.random), &result)
	return &result, err
}

// SendMultipleMarketSMSWithTemplate multiple receiver version of SendMarketSMSWithTemplate
func (c *Client) SendMultipleMarketSMSWithTemplate(mobile []string,
	sign, templateID string,
	parameters []string) (*MultiResult, error) {
	smsMessage, _ := c.newMultiSMS(smsMarket, mobile, "", sign, templateID, parameters)
	var result MultiResult
	err := c.sendMessage(smsMessage, c.routine(multiSmsRoutine, smsMessage.random), &result)
	return &result, err
}

func (c *Client) newSmsStatus(messageType int, maxNumber int64) *smsStatus {
	rand := random()
	now := time.Now().Unix()

	return &smsStatus{
		Digest:      calSign(c.appKey, rand, now, nil),
		Time:        now,
		MessageType: messageType,
		MaxNumber:   maxNumber,
		random:      rand,
	}
}

func (c *Client) newMobileStatus(mobile string, messageType int,
	startTime int64, endTime int64,
	maxNumber int64) (*mobileStatus, error) {
	rand := random()
	now := time.Now().Unix()

	countryCode, nationalNumber, err := parseNumber(mobile)

	if err != nil {
		return nil, err
	}

	return &mobileStatus{
		Digest:         calSign(c.appKey, rand, now, nil),
		Time:           now,
		MessageType:    messageType,
		MaxNumber:      maxNumber,
		StartTime:      startTime,
		EndTime:        endTime,
		CountryCode:    countryCode,
		NationalNumber: nationalNumber,
		random:         rand,
	}, nil
}

// SMSStatusPull Pull SMS message status, this routine is enterprise only, please contact Tencent to open this feature
// see
func (c *Client) SMSStatusPull() (*MultiResult, error) {
	smsStatus := c.newSmsStatus(smsNotification, 16)
	var result MultiResult
	err := c.sendMessage(smsStatus, c.routine(statusPullRoutine, smsStatus.random), &result)
	return &result, err
}

// MobileStatusPull pull SMS messages status for single mobile.
func (c *Client) MobileStatusPull(mobile string, startTime int64, endTime int64) (*MultiResult, error) {
	mobileStatus, err := c.newMobileStatus(mobile, smsNotification, startTime, endTime, 16)

	if err != nil {
		return nil, err
	}
	var result MultiResult
	err = c.sendMessage(mobileStatus, c.routine(mobileStatusPullRoutine, mobileStatus.random), &result)
	return &result, err
}
