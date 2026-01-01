package generator

import (
	"slices"
	"time"

	"github.com/topvennie/sortifyr/internal/database/model"
)

func hasBurst(times []time.Time, window model.GeneratorWindow) bool {
	if len(times) < window.MinPlays {
		return false
	}

	slices.SortFunc(times, func(a, b time.Time) int { return a.Compare(b) })

	left := 0
	for right := range times {
		for times[right].Sub(times[left]) > window.BurstInterval {
			left++
		}
		if right-left+1 >= window.MinPlays {
			return true
		}
	}

	return false
}
