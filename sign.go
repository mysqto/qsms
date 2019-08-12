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
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/mysqto/log"
	"io"
	"math/rand"
	"os"
	"strings"
)

func random() int {
	min := 100000
	max := 999999
	return rand.Intn(max-min) + min
}

// sha256Sum returns sha256 check sum of string
func sha256Sum(v string) string {
	hash := sha256.New()
	hash.Write([]byte(v))
	return hex.EncodeToString(hash.Sum(nil))
}

// shaSum returns sha256 check sum of a file
func shaSum(path string) (string, error) {
	var sha256Digest string
	file, err := os.Open(path)

	if err != nil {
		return sha256Digest, err
	}

	defer file.Close()

	hash := sha256.New()

	if _, err := io.Copy(hash, file); err != nil {
		log.Warnf("error calculating md5sun :%v", err)
		return sha256Digest, nil
	}

	hashInBytes := hash.Sum(nil)
	sha256Digest = hex.EncodeToString(hashInBytes)

	return sha256Digest, nil
}

func calSign(appKey string, rand int, time int64, numbers []string) string {
	query := fmt.Sprintf("appkey=%s&random=%d&time=%d", appKey, rand, time)

	if len(numbers) > 0 {
		query = fmt.Sprintf("%s&mobile=%s", query, strings.Join(numbers, ","))
	}
	return sha256Sum(query)
}

func calAuth(appKey string, rand int, time int64, sha1sum string) string {
	query := fmt.Sprintf("appkey=%s&random=%d&time=%d&content-sha1=%s", appKey, rand, time, sha1sum)
	return sha256Sum(query)
}
