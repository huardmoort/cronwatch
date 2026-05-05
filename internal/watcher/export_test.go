package watcher

import "time"

// CheckAllAt exposes the private checkAll method for white-box testing.
func (w *Watcher) CheckAllAt(now time.Time) {
	w.checkAll(now)
}
