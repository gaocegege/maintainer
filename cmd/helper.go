// Copyright © 2017 Maintainer Authors
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

package cmd

import (
	"fmt"
	"time"
)

// Footer returns the footer to be written into files.
func Footer() string {
	// Refer https://golang.org/src/time/format.go.
	dateFormatStr := "2006-01-02"
	formatStr := "\n---\n\nAuto-generated by [gaocegege/maintainer]" +
		"(https://github.com/maintainer-org/maintainer) on %s.\n"

	date := time.Now().Format(dateFormatStr)
	return fmt.Sprintf(formatStr, date)
}
