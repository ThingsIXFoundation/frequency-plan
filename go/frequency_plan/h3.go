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
	"embed"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log"

	h3light "github.com/ThingsIXFoundation/h3-light"
)

//go:embed *.h3
var h3index embed.FS
var h3cache map[BandName][]h3light.Cell = make(map[BandName][]h3light.Cell)

func init() {
	err := makeH3Cache()
	if err != nil {
		log.Fatal("could not load frequency_plan h3 cache")
	}
}

func IsValidBandForHex(band BandName, hex h3light.Cell) bool {
	// To check if the band is valid for a certain hex we have to check
	// the hex is contained in the band h3 index. However since the
	// h3index is compacted we also have to downscale the hex to the same
	// resolution
	for _, bandCell := range h3cache[band] {
		if hex.Parent(bandCell.Resolution()) == bandCell {
			return true
		}
	}

	return false
}

func makeH3Cache() error {
	for _, band := range AllBands {
		f, err := h3index.Open(fmt.Sprintf("%s.h3", string(band)))
		if err != nil {
			if errors.Is(err, fs.ErrNotExist) {
				continue
			}

			return err
		}

		for {
			var cell int64
			err = binary.Read(f, binary.LittleEndian, &cell)
			if err != nil {
				if err == io.EOF {
					break
				}

				return err
			}

			h3cache[band] = append(h3cache[band], h3light.Cell(cell))
		}
	}

	return nil
}
