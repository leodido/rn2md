package releasenotes

// statistics represents counters about the merged in PRs.
type statistics struct {
	total     int
	totalNone int
	authors   map[string]int
}
