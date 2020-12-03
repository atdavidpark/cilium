// Copyright 2020 Authors of Cilium
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

package inctimer

import "time"

// IncTimer wraps time.Timer
type IncTimer struct {
	t *time.Timer
}

// New creates a new IncTimer and a done function
func New() (*IncTimer, func() bool) {
	t := time.NewTimer(time.Nanosecond)
	return &IncTimer{
		t: t,
	}, t.Stop
}

// After returns a channel that will fire after
// the specified duration.
func (it *IncTimer) After(d time.Duration) <-chan time.Time {
	// We cannot call reset on an expired timer,
	// so we need to stop it and drain it first.
	// See https://golang.org/pkg/time/#Timer.Reset for more details.
	if !it.t.Stop() {
		// It could be that the channel was read already
		select {
		case <-it.t.C:
		default:
		}
	}
	it.t.Reset(d)
	return it.t.C
}
