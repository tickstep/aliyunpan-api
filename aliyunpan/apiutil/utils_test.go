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
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRand(t *testing.T) {
	r := Rand()
	fmt.Println(r)
	assert.Equal(t, 16, len(r))
}

func TestDateOfGmtStr(t *testing.T) {
	r := DateOfGmtStr()
	fmt.Println(r)
}

func TestUtcTime2LocalFormat(t *testing.T) {
	r := UtcTime2LocalFormat("2021-07-29T23:18:07.000Z")
	fmt.Println(r) // 2021-07-30 07:18:07
}

func TestLocalTime2UtcFormat(t *testing.T) {
	r := LocalTime2UtcFormat("2021-07-30 07:18:07")
	fmt.Println(r) // 2021-07-29T23:18:07.000Z
}

func TestUnixTime2LocalFormat(t *testing.T) {
	r := UnixTime2LocalFormat(1650793433058)
	fmt.Println(r) // 2022-04-24 17:43:53
}
