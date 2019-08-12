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

package main

import (
	"github.com/mysqto/log"
	"github.com/mysqto/qsms"
	"os"
	"regexp"
	"strings"
)

func extractCode(message string) []string {
	// normally verification code is 6-digit or 4-digit
	// telegram is 5-digit
	regex := `(\d{6}|\d{5}|\d{4})`
	re := regexp.MustCompile(regex)
	return re.FindAllString(message, -1)
}

func main() {


	log.New(os.Stderr, "", log.Lfull)
	appID, appKey := os.Getenv("APP_ID"), os.Getenv("APP_KEY")
	sign, templateID := os.Getenv("SIGN"), os.Getenv("TEMPLATE_ID")
	receiver, message := os.Getenv("PHONE_NUMBER"), os.Getenv("MESSAGE")

	if len(appID) == 0 || len(appKey) == 0 ||
		len(sign) == 0 || len(templateID) == 0 ||
		len(receiver) == 0 || len(message) == 0 {
		log.Fatalf("fatal: invalid parameters, current[appID = %v, appKey = %v, sign = %v,"+
			" templateID = %v, receiver = %v, message = %v]",
			appID, appKey, sign, templateID, receiver, message)
	}
	message = strings.ToLower(message)
	log.Debugf("message = %v", message)

	if !(strings.Contains(message, "验证码") ||
	    strings.Contains(message, "密碼") ||
	    strings.Contains(message, "驗證碼") ||
		strings.Contains(message, "code") ||
		strings.Contains(message, "password") ||
		strings.Contains(message, "verification") ||
		strings.Contains(message, "verify")) {
		log.Fatalf("message [%s] does not contain any verification code", message)
	}

	receivers := strings.Split(receiver, ",")

	for index, receiver := range receivers {
		if !strings.HasPrefix(receiver, "+86") {
			receivers[index] = "+86" + receiver
		}
	}
	if len(receivers) == 0 {
		log.Fatalf("invalid receivers : %v", receiver)
	}

	client := qsms.NewClient(appID, appKey)
	codes := extractCode(message)

	log.Infof("codes = %v", codes)

	if len(codes) == 0 {
		log.Fatalf("message [%s] contains no code", message)
	}

	/*
	 * =================================================================================================================
	 * Send SMS with single receiver
	 * =================================================================================================================
	 */
	result, err := client.SendSMSWithTemplate(receivers[0], sign, templateID, codes)
	if err != nil || result.Status != 0 {
		log.Warnf("error sending sms : [err = %v, result = %v]", err, result)
	} else {
		// result.MessageID is server-side messageId, can be paired with result.LocalMessageID to identify single message
		log.Infof("sms send with : status = %s, messageID = %s", result.ErrMsg, result.MessageID)
	}
}
