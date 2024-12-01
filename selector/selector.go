package selector

import (
	"aoc2k24/constants"
	"aoc2k24/d1"
	"fmt"
)

func RunDay(day constants.DayIndex, ver constants.VersionIndex) {
  switch day {
    case constants.One:
      d1.Init(ver)
    default:
      panic(fmt.Sprintf("Day %d is not present", day))
  }
}
