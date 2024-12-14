package d14

import (
	"aoc2k24/constants"
	"aoc2k24/io"
	"fmt"
)

func Init(ver constants.VersionIndex) {
  lines, err := io.GetLinesFor(constants.Fourteen, ver)
  if (err != nil) {
    panic(fmt.Sprintf("Error loading file for day %d, version %d: %v", constants.Fourteen, ver, err))
  }
}
