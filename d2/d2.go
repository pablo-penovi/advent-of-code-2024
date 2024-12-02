package d2

import (
	"aoc2k24/constants"
	"aoc2k24/io"
	"fmt"
	"strconv"
	"strings"
)

const isDebug = false

func Init(ver constants.VersionIndex) {
  lines, err := io.GetLinesFor(constants.Two, ver)
  if (err != nil) {
    panic(fmt.Sprintf("Error loading file for day %d, version %d: %v", constants.Two, ver, err))
  }
  fmt.Printf("Total safe reports (part 1): %d\n", processReports(lines, false))
  fmt.Printf("Total safe reports (part 2): %d\n", processReports(lines, true))
}

func processReports(lines []string, isPart2 bool) int {
  safeReports := 0
  for i, line := range lines {
    nums := getNums(strings.Split(line, " "))
    if isDebug { fmt.Printf("Processing report #%d: %s - part 2? %v\n", i, line, isPart2) }
    isSafe := isReportSafe(nums, getNew(nums, -1), isPart2)
    if isSafe { safeReports++ }
  }
  return safeReports
}

func isReportSafe(orig []int, nums []int, isPart2 bool) bool {
  isAscending := findTrend(nums)
  unsafeCount := 0
  updateUnsafeCount(nums, isAscending, &unsafeCount)
  if unsafeCount == 0 { 
    return true 
  }
  if isPart2 && len(orig) == len(nums) {
    for skip := 0; skip < len(orig); skip++ {
      newNums := getNew(orig, skip)
      if isDebug { fmt.Printf("Analyzing sub %v\n", newNums) }
      isSafe := isReportSafe(orig, newNums, isPart2)
      if isSafe { return true }
    }
  }
  return false
}

func getNew(nums []int, skip int) []int {
  newNums := make([]int, 0)
  for i, num := range nums {
    if i != skip {
      newNums = append(newNums, num)
    }
  }
  return newNums
}

func findTrend(nums []int) bool {
  for i := 0; i < len(nums) - 1; i++ {
    if nums[i] == nums[i+1] { continue }
    return nums[i] < nums[i+1]
  }
  return false
}

func updateUnsafeCount(nums []int, isAscending bool, unsafeCount *int) {
  if len(nums) == 1 { 
    if isDebug { fmt.Printf("Finished examining report. Unsafe count: %d\n\n", *unsafeCount) }
    return
  }
  safe := isValidProgression(nums[0], nums[1], isAscending)
  if !safe {
    (*unsafeCount)++
  }
  if isDebug { fmt.Printf("Is pair %d, %d safe? %v\n", nums[0], nums[1], safe) }
  updateUnsafeCount(nums[1:], isAscending, unsafeCount)
}

func isValidProgression(v1 int, v2 int, isAscending bool) bool {
  diff := v1 - v2
  if diff < 0 { diff = -diff }
  return diff > 0 && diff <= 3 && (!isAscending && v1 > v2 || isAscending && v1 < v2)
}

func getNums(strVals []string) []int {
  numVals := make([]int, len(strVals))
  for i := range len(strVals) {
    numVals[i], _ = strconv.Atoi(strVals[i])
  }
  return numVals
}
