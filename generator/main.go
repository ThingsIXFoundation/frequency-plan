// Copyright 2022 Stichting ThingsIX Foundation
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// SPDX-License-Identifier: Apache-2.0

//go:build cgo

package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"runtime"
	"sync"

	"github.com/ThingsIXFoundation/frequency-plan/go/frequency_plan"
	"github.com/ThingsIXFoundation/frequency-plan/go/h3tools"
	"github.com/biter777/countries"
	"github.com/gammazero/workerpool"
	"github.com/sirupsen/logrus"
	"github.com/tidwall/geojson"
	"github.com/uber/h3-go/v4"
)

func loadLandBoundaries() map[countries.CountryCode][]*geojson.Feature {
	logrus.Infof("loading land boundaries")
	cfm := map[countries.CountryCode][]*geojson.Feature{}
	for countryCode := range frequency_plan.CountryToPlan {
		resp, err := http.Get(fmt.Sprintf("https://nominatim.openstreetmap.org/search.php?country=%s&countrycodes=%s&polygon_geojson=1&format=geojson", url.QueryEscape(countryCode.Alpha2()), url.QueryEscape(countryCode.Alpha2())))
		if err != nil {
			logrus.WithError(err).Errorf("could not fetch boundaries of %s (%s)", countryCode, countryCode.Alpha2())
		}
		defer resp.Body.Close()
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			logrus.WithError(err).Errorf("could not fetch boundaries of %s (%s)", countryCode, countryCode.Alpha2())
		}

		obj, err := geojson.Parse(string(bodyBytes), &geojson.ParseOptions{})
		if err != nil {
			logrus.WithError(err).Errorf("could not fetch boundaries of %s (%s)", countryCode, countryCode.Alpha2())
		}

		fc := obj.(*geojson.FeatureCollection)

		if len(fc.Base()) <= 0 {
			logrus.Errorf("could not find boundaries of coyntry of %s (%s)", countryCode, countryCode.Alpha2())
		} else {
			logrus.Infof("downloaded boundaries of country of %s (%s)", countryCode, countryCode.Alpha2())
		}

		for _, countryObj := range fc.Base() {
			cf := countryObj.(*geojson.Feature)

			cfm[countryCode] = append(cfm[countryCode], cf)
		}

	}

	logrus.Infof("completed loading land boundaries")

	return cfm
}

func loadEEZBoundaries() map[countries.CountryCode][]*geojson.Feature {
	logrus.Infof("loading EEZ boundaries")

	type eezBoundariesProperties struct {
		Properties struct {
			MRGID json.Number `json:"MRGID"`
		}
	}
	f, err := os.ReadFile("eez_v11.geojson")
	if err != nil {
		logrus.WithError(err).Fatal("could not open eez file")
	}

	obj, err := geojson.Parse(string(f), &geojson.ParseOptions{})
	if err != nil {
		logrus.WithError(err).Fatal("could not parse eez file")
	}

	fc := obj.(*geojson.FeatureCollection)

	cfm := map[countries.CountryCode][]*geojson.Feature{}
	eezs := fc.Base()
	for _, eezObj := range eezs {
		cf := eezObj.(*geojson.Feature)
		cp := &eezBoundariesProperties{}
		err := json.Unmarshal([]byte(cf.Members()), cp)
		if err != nil {
			logrus.WithError(err).Error("error while reading properties")
		}

		countryCode, ok := frequency_plan.EEZToCountry[cp.Properties.MRGID.String()]
		if !ok {
			continue
		}

		logrus.Infof("Adding eez for %s", countryCode)

		if _, ok := cfm[countryCode]; ok {
			logrus.Errorf("country has multiple eez definitions: %s", countryCode)
		}

		cfm[countryCode] = append(cfm[countryCode], cf)
	}

	logrus.Infof("completed loading EEZ boundaries")

	return cfm
}

func writeH3IndexFiles(frequencyHex map[string][]h3.Cell) {
	for frequencyPlan, index := range frequencyHex {
		filename := fmt.Sprintf("../go/frequency_plan/%s.h3", frequencyPlan)
		f, err := os.Create(filename)
		if err != nil {
			logrus.WithError(err).Fatalf("could not open file: %s", filename)
		}

		for _, cell := range index {
			err = binary.Write(f, binary.LittleEndian, int64(cell))
			if err != nil {
				logrus.WithError(err).Fatalf("could not write to file: %s", filename)
			}
		}
	}
}

func writeH3GeojsonFiles(frequencyHex map[string][]h3.Cell) {
	for frequencyPlan, index := range frequencyHex {
		filename := fmt.Sprintf("../go/frequency_plan/%s.geojson", frequencyPlan)
		f, err := os.Create(filename)
		if err != nil {
			logrus.WithError(err).Fatalf("could not open file: %s", filename)
		}

		geoj := h3tools.H3CellsToMultiPolygon(index)

		f.Write([]byte(geoj.JSON()))
	}
}

func generateH3Index() {
	wp := workerpool.New(runtime.NumCPU())
	logrus.Infof("running h3 index generator using %d workers", runtime.NumCPU())

	landBoundaries := loadLandBoundaries()

	eezBoundaries := loadEEZBoundaries()

	// Resolution to use for the frequency-plan map
	resolution := 6

	// Plan containing the frequency per country
	frequencyPlan := frequency_plan.CountryToPlan

	// Index of h3-hexes that for a certain frequency-plan
	frequencyHex := map[string][]h3.Cell{}

	// Mutex to protect the frequency hex
	frequencyHexMutex := sync.Mutex{}

	// Iterate over all countries with it's plan and build a H3 map
	for country, plan := range frequencyPlan {
		country := country
		plan := plan
		wp.Submit(func() {
			country := country
			plan := plan
			logrus.Infof("adding land for %s with frequency-plan; %s", country, plan)
			cells := []h3.Cell{}

			// Get the land boundaries of the country
			landBoundary := landBoundaries[country]
			if landBoundary == nil {
				logrus.Errorf("could not get land boundaries for %s", country)
			} else {
				for _, feature := range landBoundary {
					landObj := feature.Base()
					cells = append(cells, h3tools.ObjectToH3(landObj, resolution)...)
				}
			}

			logrus.Infof("adding EEZ for %s with frequency-plan; %s", country, plan)
			eezBoundary := eezBoundaries[country]
			if eezBoundary == nil {
				logrus.Errorf("could not get EEZ boundaries for %s", country)
			} else {
				for _, feature := range eezBoundary {
					eezObj := feature.Base()
					cells = append(cells, h3tools.ObjectToH3(eezObj, resolution)...)
				}
			}

			logrus.Infof("adding for %s", country)
			frequencyHexMutex.Lock()
			frequencyHex[plan] = append(frequencyHex[plan], cells...)
			frequencyHexMutex.Unlock()
			logrus.Infof("completed for %s", country)
		})
	}

	// Wait for all workers to complete
	wp.StopWait()

	for plan, _ := range frequencyHex {
		frequencyHex[plan] = h3tools.DebupAndCompactCells(frequencyHex[plan])
	}

	writeH3IndexFiles(frequencyHex)

}

func main() {
	generateH3Index()
}
