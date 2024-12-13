
package d12

import (
	"aoc2k24/constants"
	"aoc2k24/io"
	"fmt"
)

var isDebug = false

func Init(ver constants.VersionIndex) {
  lines, err := io.GetLinesFor(constants.Twelve, ver)
  if (err != nil) {
    panic(fmt.Sprintf("Error loading file for day %d, version %d: %v", constants.Twelve, ver, err))
  }
}
