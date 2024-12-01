package io

import (
	"aoc2k24/constants"
	"bufio"
	"fmt"
	"os"
)

func GetLinesFor(day constants.DayIndex, ver constants.VersionIndex) ([]string, error) {
  file, err := os.Open(fmt.Sprintf("/home/pablo/projects/aoc/2024/go/files/%d-%d.txt", day, ver))
  if err != nil {
    return nil, err
  }
  defer file.Close()

  var lines []string
  scanner := bufio.NewScanner(file)
  for scanner.Scan() {
    lines = append(lines, scanner.Text())
  }
  return lines, scanner.Err()
}
