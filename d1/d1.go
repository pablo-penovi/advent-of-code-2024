package d1

import (
	"aoc2k24/constants"
	"aoc2k24/io"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

func Init(ver constants.VersionIndex) {
  lines, err := io.GetLinesFor(constants.One, ver)
  if (err != nil) {
    panic(fmt.Sprint("Error loading file for day %d, version %d: %v", constants.One, ver, err))
  }
  seq1, seq2 := parseSequences(lines)
  sort.Sort(sort.IntSlice(seq1))
  sort.Sort(sort.IntSlice(seq2))
  sum := getSum(seq1, seq2)
  similarity := getSimilarity(seq1, seq2)
  fmt.Printf("The total sum is: %d\n", sum)
  fmt.Printf("The total similarity is: %d\n", similarity)
}

func getSum(seq1 []int, seq2 []int) int {
  sum := 0
  for i := range len(seq1) {
    if (seq2[i] > seq1[i]) {
      sum += seq2[i] - seq1[i]
    } else {
      sum += seq1[i] - seq2[i]
    }
  }
  return sum
}

func getSimilarity(seq1 []int, seq2 []int) int {
  itemCount := getItemCountMap(seq2)
  sim := 0
  for _, num := range seq1 {
    count, isInSecond := itemCount[num]
    if !isInSecond { continue }
    sim += num * count
  }
  return sim
}

func getItemCountMap(seq []int) map[int]int {
  m := make(map[int]int)
  for _, num := range seq {
    _, isInMap := m[num]
    if isInMap {
      m[num] += 1
    } else {
      m[num] = 1
    }
  }
  return m
}

func parseSequences(lines []string) ([]int, []int) {
  seq1 := make([]int, len(lines))
  seq2 := make([]int, len(lines))
  for _, line := range lines {
    strnums := strings.Split(line, "   ")
    n1, _ := strconv.Atoi(strnums[0])
    n2, _ := strconv.Atoi(strnums[1])
    seq1 = append(seq1, n1)
    seq2 = append(seq2, n2)
  }
  return seq1, seq2
}
