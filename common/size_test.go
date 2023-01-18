// Copyright 2021 The nbn Authors
// This file is part of the nbn library.
//
// The nbn library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The nbn library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the nbn library. If not, see <http://www.gnu.org/licenses/>.

package common

import (
	"testing"
)

func TestStorageSizeString(t *testing.T) {
	tests := []struct {
		size StorageSize
		str  string
	}{
		{2381273, "2.27 MiB"},
		{2192, "2.14 KiB"},
		{12, "12.00 B"},
	}

	for _, test := range tests {
		if test.size.String() != test.str {
			t.Errorf("%f: got %q, want %q", float64(test.size), test.size.String(), test.str)
		}
	}
}
