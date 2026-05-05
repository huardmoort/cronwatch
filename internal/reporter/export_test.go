// export_test.go exposes internal symbols for white-box testing.
package reporter

// JobStatusFields returns the fields of a JobStatus for assertion in tests.
func JobStatusFields(s JobStatus) (name string, missed bool, hasLastRun bool) {
	return s.Name, s.Missed, s.LastRun != nil
}
