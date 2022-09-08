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

package h3tools

import (
	"reflect"

	"github.com/sirupsen/logrus"
	"github.com/tidwall/geojson"
	"github.com/tidwall/geojson/geometry"
	"github.com/uber/h3-go/v4"
)

func ObjectToH3(obj geojson.Object, maxResolution int) []h3.Cell {
	if poly, ok := obj.(*geojson.Polygon); ok {
		return PolyToH3(poly, maxResolution)
	} else if multiPoly, ok := obj.(*geojson.MultiPolygon); ok {
		return MultiPolyToH3(multiPoly, maxResolution)
	} else {
		logrus.Fatalf("invalid type: %s", reflect.TypeOf(poly))
		return nil
	}
}

func PolyToH3(poly *geojson.Polygon, maxResolution int) []h3.Cell {
	ret := []h3.Cell{}
	for _, res0 := range h3.Res0Cells() {
		ret = append(ret, polyToH3(poly.Base(), res0, maxResolution)...)
	}

	ret = append(ret, h3.PolygonToCells(PolyToH3Poly(poly.Base()), maxResolution)...)
	return ret
}

func polyToH3(poly *geometry.Poly, cell h3.Cell, maxResolution int) []h3.Cell {
	ret := []h3.Cell{}
	cellPoly := H3CellToPoly(cell)
	if poly.IntersectsPoly(cellPoly) {
		if maxResolution == cell.Resolution() {
			ret = append(ret, cell)
		} else {
			for _, child := range cell.ImmediateChildren() {
				ret = append(ret, polyToH3(poly, child, maxResolution)...)
			}
		}
	} else if poly.ContainsPoly(cellPoly) {
		ret = append(ret, cell.Children(maxResolution)...)
	}

	return ret
}

func PolyContainsPointOfH3(poly *geometry.Poly, cell h3.Cell) bool {
	for _, point := range H3CellToPoints(cell) {
		if poly.ContainsPoint(point) {
			return true
		}
	}

	return false
}

func H3CellToPoints(cell h3.Cell) []geometry.Point {
	ret := []geometry.Point{}
	for _, point := range cell.Boundary() {
		ret = append(ret, geometry.Point{X: point.Lng, Y: point.Lat})
	}

	return ret
}

func H3CellToPoly(cell h3.Cell) *geometry.Poly {
	var points []geometry.Point
	for _, h3p := range cell.Boundary() {
		points = append(points, geometry.Point{X: h3p.Lng, Y: h3p.Lat})
	}

	points = fixTransmeridianLoop(points)

	return geometry.NewPoly(points, nil, nil)
}

func H3CellToPoint(cell h3.Cell) geometry.Point {
	return geometry.Point{X: cell.LatLng().Lng, Y: cell.LatLng().Lat}
}

func fixTransmeridianLoop(points []geometry.Point) []geometry.Point {
	isTransmedian := false
	for i := range points {
		if points[i].X-points[(i+1)%len(points)].X > 180.0 {
			isTransmedian = true
			break
		}
	}

	if !isTransmedian {
		return points
	}

	for i := range points {
		if points[i].X < 0 {
			points[i].X += 360.0
		}
	}

	return points
}

func PolyToH3Poly(poly *geometry.Poly) h3.GeoPolygon {
	ret := h3.GeoPolygon{}
	ret.GeoLoop = RingToH3GeoLoop(poly.Exterior)
	for _, hole := range poly.Holes {
		ret.Holes = append(ret.Holes, RingToH3GeoLoop(hole))
	}

	return ret
}

func RingToH3GeoLoop(ring geometry.Ring) h3.GeoLoop {
	var ret h3.GeoLoop
	for i := 0; i < ring.NumPoints(); i++ {
		p := ring.PointAt(i)
		ret = append(ret, h3.LatLng{Lat: p.Y, Lng: p.X})
	}

	return ret
}

func H3CellsToMultiPolygon(cells []h3.Cell) *geojson.MultiPolygon {
	ps := []*geometry.Poly{}
	for _, cell := range cells {
		ps = append(ps, H3CellToPoly(cell))
	}
	mp := geojson.NewMultiPolygon(ps)

	return mp
}

func MultiPolyToH3(multipoly *geojson.MultiPolygon, maxResolution int) []h3.Cell {
	ret := []h3.Cell{}

	for _, poly := range multipoly.Children() {
		ret = append(ret, PolyToH3(poly.(*geojson.Polygon), maxResolution)...)
	}

	ret = DedupCells(ret)

	return ret
}

func DedupCells(cells []h3.Cell) []h3.Cell {
	dd := map[h3.Cell]bool{}
	ret := []h3.Cell{}

	for _, c := range cells {
		if _, ok := dd[c]; !ok {
			dd[c] = true
			ret = append(ret, c)
		}
	}

	return ret

}

func DebupAndCompactCells(cells []h3.Cell) []h3.Cell {
	if len(cells) <= 1 {
		return cells
	}
	ret := DedupCells(cells)

	ret = h3.CompactCells(ret)

	return ret
}
