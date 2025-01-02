package d19

import (
	"aoc2k24/constants"
	"aoc2k24/io"
	"fmt"
	"strings"
)

const isDebug = false

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
  memo := make(map[string]int)
  for _, d := range *desired {
    if isDebug { fmt.Printf("\n**** Desired pattern %s ****\n", d) }
    res := solve(d, available, 0, 1, &memo)
    if isDebug { fmt.Printf("Matches for pattern %s: %d\n\n", d, res) }
    count += res
  }
  return count
}

func solve(d string, a *map[string]struct{}, start int, end int, memo *map[string]int) int {
  // End condition 1: End of desired pattern reached with no match
  if end > len(d) {
    if isDebug { fmt.Print("End reached and no match found\n") }
    return 0
  }

  if isDebug { fmt.Printf("Analyzing %s (%d, %d)\n", d[start:end], start, end) }

  // End condition 2: Match found for last segment of desired pattern
  _, isMatch := (*a)[d[start:end]]
  if end == len(d) && isMatch {
    if isDebug { fmt.Printf("Match found for %s. Also this is the last segment so this is a match!\n", d[start:end]) }
    (*memo)[d] = 1
    return 1
  }

  memoVal, isRestInMemo := (*memo)[d[start:]]
  if isRestInMemo {
    if isDebug { fmt.Printf("Match found in memo for %s: %d\n", d[start:], memoVal) }
    return memoVal
  }

  // First possibility: Match found
  // In this case we branch, looking for the next char (begin at current end) and also expanding the current search (end + 1)
  if isMatch {
    if isDebug { fmt.Printf("Match found for %s. Keep looking for %d, %d (%s)\n", d[start:end], end, end + 1, d[end:end+1]) }
    b1 := solve(d, a, end, end + 1, memo)
    b2 := solve(d, a, start, end + 1, memo)
    if b1 > 0 || b2 > 0 {
      return 1
    }
    return 0
  }
  if isDebug { fmt.Printf("Match not found for %s. Keep looking for %d, %d\n", d[start:end], start, end + 1) }
  // Second possibility: Match not found
  // In this case we increase end by 1 and look again
  return solve(d, a, start, end + 1, memo)
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
