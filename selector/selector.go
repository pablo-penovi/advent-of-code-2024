package selector

import (
	"aoc2k24/constants"
	"aoc2k24/d1"
  "aoc2k24/d2"
  "aoc2k24/d3"
  "aoc2k24/d4"
  "aoc2k24/d5"
  "aoc2k24/d6"
  "aoc2k24/d7"
  "aoc2k24/d8"
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
    case constants.Four:
      d4.Init(ver)
    case constants.Five:
      d5.Init(ver)
    case constants.Six:
      d6.Init(ver)
    case constants.Seven:
      d7.Init(ver)
    case constants.Eight:
      d8.Init(ver)
    default:
      panic(fmt.Sprintf("Day %d is not present", day))
  }
}
