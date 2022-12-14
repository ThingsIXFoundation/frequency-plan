#!/bin/sh
# Copyright 2022 Stichting ThingsIX Foundation
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#
# SPDX-License-Identifier: Apache-2.0


# Adds a license to all code files, requires the https://github.com/google/addlicense tool to be installed, add ignores for code that doesn't need a license.
# Example for Go: addlicense -c "Stichting ThingsIX Foundation" -s -ignore "external/**" .
addlicense -c "Stichting ThingsIX Foundation" -s .