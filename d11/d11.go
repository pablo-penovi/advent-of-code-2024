package d11

import (
	"aoc2k24/constants"
	"aoc2k24/io"
	"fmt"
	"strconv"
	"strings"
)

var isDebug = false
const isDebug2 = false

func Init(ver constants.VersionIndex) {
  lines, err := io.GetLinesFor(constants.Eleven, ver)
  if (err != nil) {
    panic(fmt.Sprintf("Error loading file for day %d, version %d: %v", constants.Eleven, ver, err))
  }
  stones := getStoneList(lines)
  blinksP1 := 25
  blinksP2 := 75
  if isDebug2 { fmt.Printf("\n********** INITIAL STONES: %+v *************\n\n", stones) }
  for i := range blinksP2 {
    stones =  blink(stones)
    if isDebug2 { fmt.Printf("\n********** STONES AFTER %d BLINKS: %+v *************\n\n", i + 1, stones) }
    if i == blinksP1 - 1 {
      fmt.Printf("[PART 1] Total stones after %d blinks: %d\n", blinksP1, countStones(stones))
    }
  }
  fmt.Printf("[PART 2] Total stones after %d blinks: %d\n", blinksP2, countStones(stones))
}

func countStones(numbers *map[int]int) int {
  count := 0
  for number := range *numbers {
    count += (*numbers)[number]
  }
  return count
}

func blink (numbers *map[int]int) *map[int]int {
  newNums := make(map[int]int)
  keys := make([]int, 0, len(*numbers))
  for number := range *numbers {
    keys = append(keys, number)
  }
  for _, number := range keys {
    count, _ := (*numbers)[number]
    delete(*numbers, number)
    if number == 0 {
      if isDebug { fmt.Print("Value is 0, so converting value to 1\n\n") }
      _, exists := newNums[1]
      if !exists {
        newNums[1] = count
      } else {
        newNums[1] += count
      }
      continue
    }
    str := strconv.Itoa(number)
    if len(str) % 2 == 0 {
      if isDebug { fmt.Print("Stone has even number of digits\n") }
      strn1 := str[:len(str) / 2]
      strn2 := str[len(str) / 2:]
      n1, _ := strconv.Atoi(strn1)
      n2, _ := strconv.Atoi(strn2)
      if isDebug { fmt.Printf("Value %d replaced by %d and %d\n\n", number, n1, n2) }
      _, exists := newNums[n1]
      if !exists { 
        newNums[n1] = count
      } else {
        newNums[n1] += count
      }
      _, exists = newNums[n2]
      if !exists { 
        newNums[n2] = count
      } else {
        newNums[n2] += count
      }
      continue
    }
    if isDebug { fmt.Printf("Value %d replaced by %d (x 2024)\n\n", number, number * 2024) }
    _, exists := newNums[number * 2024]
    if !exists {
      newNums[number * 2024] = count
    } else {
      newNums[number * 2024] += count
    }
  }
  return &newNums
}

func getStoneList(lines []string) *map[int]int {
  sl := make(map[int]int)
  // Input only has 1 line 
  values := strings.Split(lines[0], " ")
  for _, value := range values {
    num, _ := strconv.Atoi(value)
    _, exists := sl[num]
    if !exists {
      sl[num] = 1
    } else {
      sl[num] += 1
    }
  }
  return &sl
}
