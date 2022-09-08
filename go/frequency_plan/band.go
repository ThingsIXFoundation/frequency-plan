// Copyright 2022 Stichting ThingsIX Foundation
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// SPDX-License-Identifier: Apache-2.0

package frequency_plan

import (
	"fmt"

	"github.com/brocaar/lorawan"
	"github.com/brocaar/lorawan/band"
)

type BandName string

const (
	EU868   BandName = "EU868"
	US915   BandName = "US915"
	CN779   BandName = "CN779"
	AU915   BandName = "AU915"
	AS923   BandName = "AS923"
	AS923_2 BandName = "AS923-2"
	AS923_3 BandName = "AS923-3"
	AS923_4 BandName = "AS923-4"
	KR920   BandName = "KR920"
	IN865   BandName = "IN865"
	RU864   BandName = "RU864"
)

var AllBands []BandName = []BandName{
	EU868,
	US915,
	CN779,
	AU915,
	AS923,
	AS923_2,
	AS923_3,
	AS923_4,
	KR920,
	IN865,
	RU864,
}

func GetBand(commonName string) (band.Band, error) {
	switch commonName {
	case string(EU868):
		b, err := band.GetConfig(band.EU868, false, lorawan.DwellTimeNoLimit)
		if err != nil {
			return nil, err
		}
		b.AddChannel(867100000, 0, 5)
		b.AddChannel(867300000, 0, 5)
		b.AddChannel(867500000, 0, 5)
		b.AddChannel(867700000, 0, 5)
		b.AddChannel(867900000, 0, 5)
	default:
		return nil, fmt.Errorf("%s is not yet supported", commonName)
	}
}
