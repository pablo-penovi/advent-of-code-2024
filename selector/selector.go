package selector

import (
	"aoc2k24/constants"
	"aoc2k24/d1"
  "aoc2k24/d2"
  "aoc2k24/d3"
	"fmt"
)

func RunDay(day constants.DayIndex, ver constants.VersionIndex) {
  switch day {
    case constants.One:
      d1.Init(ver)
    case constants.Two:
      d2.Init(ver)
    case constants.Three:
      d3.Init(ver)
    default:
      panic(fmt.Sprintf("Day %d is not present", day))
  }
}
