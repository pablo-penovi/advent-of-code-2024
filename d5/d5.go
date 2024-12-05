package d5

import (
	"aoc2k24/constants"
	"aoc2k24/io"
	"fmt"
	"slices"
	"sort"
	"strconv"
	"strings"
)

type RuleMap map[int][]int

func Init(ver constants.VersionIndex) {
  lines, err := io.GetLinesFor(constants.Five, ver)
  if (err != nil) {
    panic(fmt.Sprintf("Error loading file for day %d, version %d: %v", constants.Five, ver, err))
  }
  rules, updates := parseEntry(&lines)
  part1Result, part2Result := solve(rules, updates)
  fmt.Printf("Part 1 result: %d\n", part1Result)
  fmt.Printf("Part 2 result: %d\n", part2Result)
}

func solve(rules *RuleMap, updates *[][]int) (int, int) {
  part1Sum := 0
  part2Sum := 0
  for _, update := range *updates {
    if !isUpdateInOrder(&update, rules) {
      sortUpdate(&update, rules)
      part2Sum += update[len(update) / 2]
      continue
    }
    part1Sum += update[len(update) / 2]
  }
  return part1Sum, part2Sum
}

func sortUpdate(update *[]int, rules *RuleMap) {
  sort.Slice(*update, func(i, j int) bool {
    n1 := (*update)[i]
    n2 := (*update)[j]
    rules1, exists := (*rules)[n1]
    if exists && slices.Contains(rules1, n2) {
      return true
    }
    rules2, exists := (*rules)[n2]
    if exists && slices.Contains(rules2, n1) {
      return false
    }
    return true
  })
}

func isUpdateInOrder(update *[]int, rm *RuleMap) bool {
  for i := len(*update) - 1; i > 0; i-- {
    subject := (*update)[i]
    rules, hasRules := (*rm)[subject]; if !hasRules { continue }
    for j := i - 1; j >= 0; j-- {
      toCheck := (*update)[j]
      if slices.Contains(rules, toCheck) {
        return false
      }
    }
  }
  return true
}

func parseEntry(lines *[]string) (*RuleMap, *[][]int) {
  rules := RuleMap(make(map[int][]int))
  updates := make([][]int, 0)
  for _, line := range *lines {
    if len(line) == 0 { continue }
    if strings.Contains(line, "|") {
      addRule(&line, &rules)
    } else {
      addUpdate(&line, &updates)
    }
  }
  return &rules, &updates
}

func addRule(line *string, rules *RuleMap) {
  ruleComponents := strings.Split(*line, "|")
  n1, _ := strconv.Atoi(ruleComponents[0])
  n2, _ := strconv.Atoi(ruleComponents[1])
  _, exists := (*rules)[n1]
  if exists {
    (*rules)[n1] = append((*rules)[n1], n2)
  } else {
    (*rules)[n1] = make([]int, 1)
    (*rules)[n1][0] = n2
  }
}

func addUpdate(line *string, updates *[][]int) {
  updComponents := strings.Split(*line, ",")
  update := make([]int, len(updComponents))
  for i, strnum := range updComponents {
    num, _ := strconv.Atoi(strnum)
    update[i] = num
  }
  (*updates) = append(*updates, update)
}
