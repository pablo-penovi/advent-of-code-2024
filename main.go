package main

import (
	"aoc2k24/constants"
	"aoc2k24/selector"
	"flag"
)
 
func main() {
  dayParam := flag.Int("day", 10, "The Advent of Code 2024 day you wish to see")
  versionParam := flag.Int("v", 0, "The version. 0 is full puzzle input, successive ones are test data")
  flag.Parse()
  selector.RunDay(constants.DayIndex(*dayParam), constants.VersionIndex(*versionParam))
}
