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

package downloader

import "github.com/token/core/types"

type DoneEvent struct {
	Latest *types.Header
}
type StartEvent struct{}
type FailedEvent struct{ Err error }