package d3

import (
	"aoc2k24/constants"
	"aoc2k24/io"
	"fmt"
	"strconv"
	"strings"
)

type State int

const (
  Do State = iota
  Dont
)

func (s *State) flip() {
  if *s == Do {
    *s = Dont
  } else {
    *s = Do
  }
}

func (s State) ToString() string {
  if s == Do {
    return "do()"
  }
  return "don't()"
}

const isDebug = false

func Init(ver constants.VersionIndex) {
  lines, err := io.GetLinesFor(constants.Three, ver)
  if (err != nil) {
    panic(fmt.Sprintf("Error loading file for day %d, version %d: %v", constants.Three, ver, err))
  }
  part1Result := getMulSums(&lines, false)
  fmt.Printf("Part 1 result: %d\n\n", part1Result)

  part2Result := getMulSums(&lines, true)
  fmt.Printf("Part 2 result: %d\n\n", part2Result)
}

func getMulSums(lines *[]string, isPart2 bool) int {
  currState := Do
  muls := 0
  for _, line := range *lines {
    for len(line) >= 8 {
      start := -1
      candidate := ""
      if isPart2 {
        start, candidate = getCandidatePart2(&line, &currState)
      } else {
        start, candidate = getCandidate(&line)
      }
      if start < 0 { continue }
      n1, n2 := getOperands(&line, start, candidate)
      if n1 < 0 { continue }
      muls += n1 * n2
      line = line[start + len(candidate):]
    }
  }
  return muls
}

func getOperands(line *string, start int, candidate string) (int, int) {
  strnums := strings.Split(candidate[4:len(candidate)-1], ",")
  if len(strnums) < 2 {
    updateLine(line, start + len(candidate)) 
    return -1, -1
  }
  n1, err := strconv.Atoi(strnums[0])
  if err != nil {
    updateLine(line, start + len(candidate)) 
    return -1, -1
  }
  n2, err := strconv.Atoi(strnums[1])
  if err != nil {
    updateLine(line, start + len(candidate))
    return -1, -1
  }
  if isDebug { fmt.Printf("Numbers: %d, %d\n", n1, n2) }
  return n1, n2
}

func getCandidatePart2(line *string, currState *State) (int, string) {
  nextStateChangeIndex := getNextStateChangeIndex(line, currState)
  start := strings.Index(*line, "mul(")
  if isDebug { fmt.Printf("Next state change: %d | Next candidate: %d\n", nextStateChangeIndex, start) }
  if nextStateChangeIndex >= 0 && nextStateChangeIndex < start {
    currState.flip()
    if isDebug { fmt.Print("Switching state\n") }
  }
  if *currState == Dont { 
    if isDebug { fmt.Print("State is DON'T so skipping this candidate\n") }
    updateLine(line, start+4)
    return -1, "" 
  }
  end := start + 13
  if end > len(*line) {
    end = len(*line)
  }
  candidate := (*line)[start:end]
  if !strings.Contains(candidate, ",") || !strings.Contains(candidate, ")") {
    updateLine(line, start+4)
    return -1, ""
  }
  end = strings.Index(candidate, ")") + 1
  candidate = candidate[:end]
  if isDebug { fmt.Printf("Candidate: %s ", candidate) }
  return start, candidate
}

func getNextStateChangeIndex(line *string, currState *State) int {
  toSearch := Do
  if *currState == Do {
    toSearch = Dont
  }
  return strings.Index(*line, toSearch.ToString())
}

func getCandidate(line *string) (int, string) {
  start := strings.Index(*line, "mul(")
  end := start + 13
  if end > len(*line) {
    end = len(*line)
  }
  candidate := (*line)[start:end]
  if !strings.Contains(candidate, ",") || !strings.Contains(candidate, ")") {
    updateLine(line, start+4)
    return -1, ""
  }
  end = strings.Index(candidate, ")") + 1
  candidate = candidate[:end]
  if isDebug { fmt.Printf("Candidate: %s ", candidate) }
  return start, candidate
}

func updateLine(line *string, positionsToTrim int) {
  *line = (*line)[positionsToTrim:]
}
