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
	"strings"

	"github.com/brocaar/lorawan"
	"github.com/brocaar/lorawan/band"
)

type BandName string

type BlockchainFrequencyPlan uint

const (
	Invalid BandName = "INVALID"
	EU868   BandName = "EU868"
	US915   BandName = "US915"
	CN779   BandName = "CN779"
	EU433   BandName = "EU433"
	AU915   BandName = "AU915"
	CN470   BandName = "CN470"
	AS923   BandName = "AS923"
	AS923_2 BandName = "AS923-2"
	AS923_3 BandName = "AS923-3"
	KR920   BandName = "KR920"
	IN865   BandName = "IN865"
	RU864   BandName = "RU864"
	AS923_4 BandName = "AS923-4"
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

func FromBlockchain(in BlockchainFrequencyPlan) BandName {
	switch in {
	case 0:
		return Invalid
	case 1:
		return EU868
	case 2:
		return US915
	case 3:
		return CN779
	case 4:
		return EU433
	case 5:
		return AU915
	case 6:
		return CN470
	case 7:
		return AS923
	case 8:
		return AS923_2
	case 9:
		return AS923_3
	case 10:
		return KR920
	case 11:
		return IN865
	case 12:
		return RU864
	case 13:
		return AS923_4
	default:
		return Invalid
	}
}

func (b BandName) ToBlockchain() BlockchainFrequencyPlan {
	switch b {
	case EU868:
		return 1
	case US915:
		return 2
	case CN779:
		return 3
	case EU433:
		return 4
	case AU915:
		return 5
	case CN470:
		return 6
	case AS923:
		return 7
	case AS923_2:
		return 8
	case AS923_3:
		return 9
	case KR920:
		return 10
	case IN865:
		return 11
	case RU864:
		return 12
	case AS923_4:
		return 13
	default:
		return 0
	}
}

func (b *BandName) UnmarshalText(text []byte) error {
	switch strings.ToUpper(string(text)) {
	case string(EU868):
		*b = EU868
		return nil
	case string(US915):
		*b = US915
		return nil
	case string(CN779):
		*b = CN779
		return nil
	case string(EU433):
		*b = EU433
		return nil
	case string(AU915):
		*b = AU915
		return nil
	case string(CN470):
		*b = CN470
		return nil
	case string(AS923):
		*b = AS923
		return nil
	case string(AS923_2):
		*b = AS923_2
		return nil
	case string(AS923_3):
		*b = AS923_3
		return nil
	case string(KR920):
		*b = KR920
		return nil
	case string(IN865):
		*b = IN865
		return nil
	case string(RU864):
		*b = RU864
		return nil
	case string(AS923_4):
		*b = AS923_4
		return nil
	default:
		return fmt.Errorf(`unknown frequency plan "%s"`, text)
	}
}

func GetBand(commonName string) (band.Band, error) {
	switch commonName {
	case string(EU868):
		b, err := band.GetConfig(band.EU868, false, lorawan.DwellTimeNoLimit)
		if err != nil {
			return nil, err
		}
		err = b.AddChannel(867100000, 0, 5)
		if err != nil {
			return nil, err
		}

		err = b.AddChannel(867300000, 0, 5)
		if err != nil {
			return nil, err
		}

		err = b.AddChannel(867500000, 0, 5)
		if err != nil {
			return nil, err
		}

		err = b.AddChannel(867700000, 0, 5)
		if err != nil {
			return nil, err
		}

		err = b.AddChannel(867900000, 0, 5)
		if err != nil {
			return nil, err
		}

		return b, nil
	case string(AU915):
		b, err := band.GetConfig(band.AU915, false, lorawan.DwellTime400ms)
		if err != nil {
			return nil, err
		}
		err = b.AddChannel(916800000, 0, 5)
		if err != nil {
			return nil, err
		}

		err = b.AddChannel(917000000, 0, 5)
		if err != nil {
			return nil, err
		}

		err = b.AddChannel(917200000, 0, 5)
		if err != nil {
			return nil, err
		}

		err = b.AddChannel(917400000, 0, 5)
		if err != nil {
			return nil, err
		}

		err = b.AddChannel(917600000, 0, 5)
		if err != nil {
			return nil, err
		}

		err = b.AddChannel(917800000, 0, 5)
		if err != nil {
			return nil, err
		}

		err = b.AddChannel(918000000, 0, 5)
		if err != nil {
			return nil, err
		}

		err = b.AddChannel(918200000, 0, 5)
		if err != nil {
			return nil, err
		}

		err = b.AddChannel(917500000, 6, 6)
		if err != nil {
			return nil, err
		}

		return b, nil	
	default:
		return nil, fmt.Errorf("%s is not yet supported", commonName)
	}
}
