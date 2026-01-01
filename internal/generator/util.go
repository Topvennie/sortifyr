package generator

import (
	"slices"
	"time"
)

func hasBurst(times []time.Time, minPlays int, interval time.Duration) bool {
	if len(times) < minPlays {
		return false
	}

	slices.SortFunc(times, func(a, b time.Time) int { return a.Compare(b) })

	left := 0
	for right := range times {
		for times[right].Sub(times[left]) > interval {
			left++
		}
		if right-left+1 >= minPlays {
			return true
		}
	}

	return false
}
