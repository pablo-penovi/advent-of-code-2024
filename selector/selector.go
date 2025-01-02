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
  "aoc2k24/d9"
  "aoc2k24/d10"
  "aoc2k24/d11"
  "aoc2k24/d12"
  "aoc2k24/d13"
  "aoc2k24/d14"
  "aoc2k24/d15"
  "aoc2k24/d16"
  "aoc2k24/d17"
  "aoc2k24/d18"
  "aoc2k24/d19"
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
    case constants.Nine:
      d9.Init(ver)
    case constants.Ten:
      d10.Init(ver)
    case constants.Eleven:
      d11.Init(ver)
    case constants.Twelve:
      d12.Init(ver)
    case constants.Thirteen:
      d13.Init(ver)
    case constants.Fourteen:
      d14.Init(ver)
    case constants.Fifteen:
      d15.Init(ver)
    case constants.Sixteen:
      d16.Init(ver)
    case constants.Seventeen:
      d17.Init(ver)
    case constants.Eighteen:
      d18.Init(ver)
    case constants.Nineteen:
      d19.Init(ver)
    default:
      panic(fmt.Sprintf("Day %d is not present", day))
  }
}
