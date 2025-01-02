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
  // memo := make(map[string]int)
  for _, d := range *desired {
    if isDebug { fmt.Printf("\n**** Desired pattern %s ****\n", d) }
    // res := solve(d, available, 0, 1, &memo)
    res := solveV2(d, 0, available)
    if isDebug { fmt.Printf("Matches for pattern %s: %d\n\n", d, res) }
    count += res
  }
  return count
}

func solveV2(d string, start int, a *map[string]struct{}) int {
  count := 0
  for aPattern := range *a {
    if count > 0 { break }
    end := start + len(aPattern)
    // If this pattern doesn't fit the rest of the input, disregard
    if end > len(d) {
      if isDebug { fmt.Printf("Pattern length exceeds input (required length %d, actual length %d)\n", end, len(d)) }
      continue
    }
    // If this pattern does not match, disregard
    if d[start:end] != aPattern {
      if isDebug { fmt.Printf("Pattern doesn't match! (pattern %s, input segment %s)\n", aPattern, d[start:end]) }
      continue
    }
    // If the pattern matches, keep going for next input segment unless we're done
    if isDebug { fmt.Printf("Pattern matches! (pattern %s, input segment %s)\n", aPattern, d[start:end]) }
    if end == len(d) {
      return 1
    }
    count += solveV2(d, end, a)
  }

  res := 0
  if count > 0 { res = 1 }
  return res
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
