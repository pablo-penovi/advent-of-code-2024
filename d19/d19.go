package d19

import (
	"aoc2k24/constants"
	"aoc2k24/io"
	"fmt"
	"strings"
)

const isDebug = true

func Init(ver constants.VersionIndex) {
  lines, err := io.GetLinesFor(constants.Nineteen, ver)
  if (err != nil) {
    panic(fmt.Sprintf("Error loading file for day %d, version %d: %v", constants.Nineteen, ver, err))
  }
  available, desired := parsePatterns(&lines)
  possible := findPossible(available, desired)
  fmt.Printf("%d designs are possible\n", possible)
}

func findPossible(available *map[string]struct{}, desired *[]string) int {
  count := 0
  for _, d := range *desired {
    if isDebug { fmt.Printf("\n**** Desired pattern %s ****\n", d) }
    res := solve(d, available, 0, len(d))
    if isDebug { fmt.Printf("Matches for pattern %s: %d\n\n", d, res) }
    count += res
  }
  return count
}

func solve(d string, a *map[string]struct{}, start int, end int) int {
  // End condition 1: End cursor reached start so everything matched
  if end == 0 {
    if isDebug { fmt.Printf("Pattern %s is a match\n", d) }
    return 1
  }
  // End condition 2: Start cursor caught up with end cursor, so no match was found
  if start == end {
    if isDebug { fmt.Printf("Pattern %s is NOT a match\n", d) }
    return 0
  }
  
  // If there is a match, recurse from beginning of input to start of current match
  // (i.e. we start again but leaving out what has been matched so far)
  _, isMatch := (*a)[d[start:end]]
  if isMatch {
    if isDebug { fmt.Printf("Segment %s is a match\n", d[start:end]) }
    return solve(d, a, 0, start)
  }

  // If no match, advance start by 1 and recurse
  if isDebug { fmt.Printf("Segment %s is NOT a match\n", d[start:end]) }
  return solve(d, a, start + 1, end)
}

func parsePatterns(lines *[]string) (*map[string]struct{}, *[]string) {
  a := make(map[string]struct{})
  d := make([]string, len(*lines) - 2)

  for _, design := range strings.Split((*lines)[0], ", ") {
    a[design] = struct{}{}
  }

  for i, design := range (*lines)[2:] {
    d[i] = design
  }

  return &a, &d
}
