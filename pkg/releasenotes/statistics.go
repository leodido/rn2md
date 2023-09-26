package releasenotes

// Statistics represents counters about the merged in PRs.
type Statistics struct {
	total     int64
	nonFacing int64
	authors   map[string]int64
}
