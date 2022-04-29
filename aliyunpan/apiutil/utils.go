// Copyright (c) 2020 tickstep.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package apiutil

import (
	"fmt"
	jsoniter "github.com/json-iterator/go"
	uuid "github.com/satori/go.uuid"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

const (
	FileNameSpecialChars = "\\/:*?\"<>|"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func Timestamp() int {
	// millisecond
	return int(time.Now().UTC().UnixNano() / 1e6)
}

func Rand() string {
	randStr := &strings.Builder{}
	fmt.Fprintf(randStr, "%d_%d", rand.Int63n(1e5), rand.Int63n(1e10))
	return randStr.String()
}

func DateOfGmtStr() string {
	return time.Now().UTC().Format(http.TimeFormat)
}

func XRequestId() string {
	u4 := uuid.NewV4()
	return strings.ToUpper(u4.String())
}

func Uuid() string {
	u4 := uuid.NewV4()
	return u4.String()
}

// CheckFileNameValid 检测文件名是否有效，包含特殊字符则无效
func CheckFileNameValid(name string) bool {
	if name == "" {
		return true
	}
	return !strings.ContainsAny(name, FileNameSpecialChars)
}

// UtcTime2LocalFormat UTC时间转换为本地时间
func UtcTime2LocalFormat(timeStr string) string {
	if timeStr == "" {
		return ""
	}
	t, _ := time.Parse(time.RFC3339, timeStr)
	timeUint := t.In(time.Local).Unix()
	return time.Unix(timeUint, 0).Format("2006-01-02 15:04:05")
}

// LocalTime2UtcFormat 本地时间转换为UTC时间
func LocalTime2UtcFormat(utcTimeStr string) string {
	if utcTimeStr == "" {
		return ""
	}
	t, _ := time.ParseInLocation("2006-01-02 15:04:05", utcTimeStr, time.Local)
	return t.UTC().Format("2006-01-02T15:04:05.000Z07:00")
}

// UnixTime2LocalFormat 时间戳转换为本地时间字符串
func UnixTime2LocalFormat(unixTime int64) string {
	return time.Unix(unixTime/1000, 0).Format("2006-01-02 15:04:05")
}

// AddCommonHeader 增加公共header
func AddCommonHeader(headers map[string]string) map[string]string {
	commonHeaders := map[string]string{
		"accept":       "application/json, text/plain, */*",
		"referer":      "https://www.aliyundrive.com/",
		"origin":       "https://www.aliyundrive.com",
		"content-type": "application/json;charset=UTF-8",
		"user-agent":   "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
	}
	if headers == nil {
		return commonHeaders
	}

	// merge
	for k, v := range headers {
		commonHeaders[strings.ToLower(k)] = v
	}
	return commonHeaders
}

func GetMapSet(param interface{}) map[string]interface{} {
	if param == nil {
		return nil
	}

	r, _ := jsoniter.MarshalToString(param)
	m := map[string]interface{}{}
	jsoniter.Unmarshal([]byte(r), &m)
	return m
}
