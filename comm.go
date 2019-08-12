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
	"github.com/google/uuid"
	"github.com/nyaruka/phonenumbers"
	"strconv"
)

func (c *Client) generateMessageID() string {
	if c.IDGenerator != nil {
		return c.IDGenerator()
	}
	return uuid.New().String()
}

func (c *Client) routine(routine string, rand int) string {
	return fmt.Sprintf("%s%s?sdkappid=%s&random=%d", qSmsBase, routine, c.appID, rand)
}

func (c *Client) routineWithTime(routine string, rand int, time int64) string {
	return fmt.Sprintf("%s%s?sdkappid=%s&random=%d&time=%d", qSmsBase, routine, c.appID, rand, time)
}

func parseNumber(mobile string) (string, string, error) {
	phoneNumber, err := phonenumbers.Parse(mobile, "")

	if err != nil {
		return "", "", fmt.Errorf("invalid mobile number : [mobile = %s, parse error = %v]", mobile, err)
	}

	return strconv.FormatInt(int64(phoneNumber.GetCountryCode()), 10),
		strconv.FormatUint(phoneNumber.GetNationalNumber(), 10), nil
}

func parseNumbers(mobiles []string) ([]number, []string) {
	var numbers []number
	var receivers []string

	for _, mobile := range mobiles {

		phoneNumber, err := phonenumbers.Parse(mobile, "")

		if err != nil {
			continue
		}

		nationalNumber := strconv.FormatUint(phoneNumber.GetNationalNumber(), 10)
		numbers = append(numbers, number{
			CountryCode:    strconv.FormatInt(int64(phoneNumber.GetCountryCode()), 10),
			NationalNUmber: nationalNumber,
		})
		receivers = append(receivers, nationalNumber)
	}
	return numbers, receivers
}
