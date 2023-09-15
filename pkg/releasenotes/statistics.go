package releasenotes

// Statistics represents counters about the merged in PRs.
type Statistics struct {
	total     int64
	totalNone int64
	authors   map[string]int64
}
