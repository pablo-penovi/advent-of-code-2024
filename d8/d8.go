package d8

import (
	"aoc2k24/constants"
	"fmt"
  "aoc2k24/io"
)

const isDebugP1 = false
const isDebugP2 = false

type AntennaMap struct {
  antennae map[rune][]int
  width int
  height int
}

func newAntennaMap(lines *[]string) *AntennaMap {
  am := AntennaMap{make(map[rune][]int), len((*lines)[0]), len(*lines)}
  for y, line := range *lines {
    for x, char := range line {
      if char == '.' { continue }
      i := y * am.width + x
      _, exists := am.antennae[rune(char)]; if exists {
        am.antennae[rune(char)] = append(am.antennae[rune(char)], i)
        continue
      }
      am.antennae[rune(char)] = []int{i}
    }
  }
  return &am
}

func Init(ver constants.VersionIndex) {
  lines, err := io.GetLinesFor(constants.Eight, ver)
  if (err != nil) {
    panic(fmt.Sprintf("Error loading file for day %d, version %d: %v", constants.Eight, ver, err))
  }
  am := newAntennaMap(&lines)
  antinodeCount := solvePart1(am)
  antinodeCount2 := solvePart2(am)
  fmt.Printf("Unique antinode locations (part 1): %d\n", antinodeCount)
  fmt.Printf("Unique antinode locations (part 2): %d\n", antinodeCount2)
}

func solvePart2(am *AntennaMap) int {
  antinodes := make(map[int]struct{})
  for antennaType := range am.antennae {
    if isDebugP2 { fmt.Printf("Processing antennae of type '%c'\n", antennaType) }
    for _, position := range am.antennae[antennaType] {
      _, exists := antinodes[position]; if !exists { antinodes[position] = struct{}{} }
      for _, otherPosition := range am.antennae[antennaType] {
        if position == otherPosition { continue }
        _, exists := antinodes[otherPosition]; if !exists { antinodes[otherPosition] = struct{}{} }
        x := position % am.width
        y := position / am.width
        otherX := otherPosition % am.width
        otherY := otherPosition / am.width
        if isDebugP2 { fmt.Printf("Computing antinodes for antennae '%c' in positions %d (x %d, y %d) and %d (x %d, y %d)\n", antennaType, position, x, y, otherPosition, otherX, otherY) }
        xDiff := x - otherX; if xDiff < 0 { xDiff = -xDiff }
        yDiff := y - otherY; if yDiff < 0 { yDiff = -yDiff }
        addAntinodes(&antinodes, x, y, otherX, otherY, xDiff, yDiff, am.width, am.height)
      }
      if isDebugP2 { fmt.Print("\n") }
    }
    if isDebugP2 { fmt.Print("\n") }
  }
  return len(antinodes)
}

func addAntinodes(antinodes *map[int]struct{}, x, y, otherX, otherY, xDiff, yDiff, width, height int) {
  if (x < 0 || x >= width || y < 0 || y >= height) && (otherX < 0 || otherX >= width || otherY < 0 || otherY >= height) {
    return
  }
  antinode1X := x - xDiff
  antinode2X := otherX + xDiff
  if otherX < x {
    antinode1X = x + xDiff
    antinode2X = otherX - xDiff
  }
  antinode1Y := y - yDiff
  antinode2Y := otherY + yDiff
  if otherY < y {
    antinode1Y = y + yDiff
    antinode2Y = otherY - yDiff
  }
  if isDebugP2 { fmt.Printf("First antinode coordinates: x %d, y %d. Second antinode coordinates: x %d, y %d\n", antinode1X, antinode1Y, antinode2X, antinode2Y) }
  if antinode1X >= 0 && antinode1X < width && antinode1Y >= 0 && antinode1Y < height {
    if isDebugP2 { fmt.Print("Antinode 1 position is inside field. Adding to map if it doesn't already exist\n") }
    antinodePos := antinode1Y * width + antinode1X
    _, exists := (*antinodes)[antinodePos]; if !exists { (*antinodes)[antinodePos] = struct{}{} }
  }
  if antinode2X >= 0 && antinode2X < width && antinode2Y >= 0 && antinode2Y < height {
    if isDebugP2 { fmt.Print("Antinode 2 position is inside field. Adding to map if it doesn't already exist\n") }
    antinodePos := antinode2Y * width + antinode2X
    _, exists := (*antinodes)[antinodePos]; if !exists { (*antinodes)[antinodePos] = struct{}{} }
  }
  addAntinodes(antinodes, antinode1X, antinode1Y, antinode2X, antinode2Y, xDiff, yDiff, width, height)
}

func solvePart1(am *AntennaMap) int {
  antinodes := make(map[int]struct{})
  for antennaType := range am.antennae {
    if isDebugP1 { fmt.Printf("Processing antennae of type '%c'\n", antennaType) }
    for _, position := range am.antennae[antennaType] {
      for _, otherPosition := range am.antennae[antennaType] {
        if position == otherPosition { continue }
        x := position % am.width
        y := position / am.width
        otherX := otherPosition % am.width
        otherY := otherPosition / am.width
        if isDebugP1 { fmt.Printf("Computing antinodes for antennae '%c' in positions %d (x %d, y %d) and %d (x %d, y %d)\n", antennaType, position, x, y, otherPosition, otherX, otherY) }
        xDiff := x - otherX; if xDiff < 0 { xDiff = -xDiff }
        yDiff := y - otherY; if yDiff < 0 { yDiff = -yDiff }
        antinode1X := x - xDiff
        antinode2X := otherX + xDiff
        if otherX < x {
          antinode1X = x + xDiff
          antinode2X = otherX - xDiff
        }
        antinode1Y := y - yDiff
        antinode2Y := otherY + yDiff
        if otherY < y {
          antinode1Y = y + yDiff
          antinode2Y = otherY - yDiff
        }
        if isDebugP1 { fmt.Printf("First antinode coordinates: x %d, y %d. Second antinode coordinates: x %d, y %d\n", antinode1X, antinode1Y, antinode2X, antinode2Y) }
        if antinode1X >= 0 && antinode1X < am.width && antinode1Y >= 0 && antinode1Y < am.height {
          if isDebugP1 { fmt.Print("Antinode 1 position is inside field. Adding to map if it doesn't already exist\n") }
          antinodePos := antinode1Y * am.width + antinode1X
          _, exists := antinodes[antinodePos]; if !exists { antinodes[antinodePos] = struct{}{} }
        }
        if antinode2X >= 0 && antinode2X < am.width && antinode2Y >= 0 && antinode2Y < am.height {
          if isDebugP1 { fmt.Print("Antinode 2 position is inside field. Adding to map if it doesn't already exist\n") }
          antinodePos := antinode2Y * am.width + antinode2X
          _, exists := antinodes[antinodePos]; if !exists { antinodes[antinodePos] = struct{}{} }
        }
      }
      if isDebugP1 { fmt.Print("\n") }
    }
    if isDebugP1 { fmt.Print("\n") }
  }
  return len(antinodes)
}
