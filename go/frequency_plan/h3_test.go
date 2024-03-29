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

package frequency_plan_test

import (
	"testing"

	"github.com/ThingsIXFoundation/frequency-plan/go/frequency_plan"
	h3light "github.com/ThingsIXFoundation/h3-light"
	"github.com/stretchr/testify/assert"
)

func TestIsValidBandForHex(t *testing.T) {

	actual := frequency_plan.IsValidBandForHex(frequency_plan.EU868, h3light.MustCellFromString("8b1fa5db57b6fff"))
	expected := true

	assert.Equal(t, expected, actual, "Eindhoven has EU868")

	actual = frequency_plan.IsValidBandForHex(frequency_plan.EU868, h3light.MustCellFromString("8b2a10728bb1fff"))
	expected = false

	assert.Equal(t, expected, actual, "New York does not have EU868")

}
