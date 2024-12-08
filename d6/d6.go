package d6

import (
	"aoc2k24/constants"
	"aoc2k24/io"
	"fmt"
)

type Direction int

const (
  North Direction = iota
  East
  South
  West
)

func (d Direction) turnRight() Direction {
  if d == North {
    return East
  }
  if d == East {
    return South
  }
  if d == South {
    return West
  }
  return North
}

type Floorplan struct {
  guardPos int
  width int
  height int
}

func newFloorplan(lines []string) (*Floorplan, *map[int]struct{}) {
  floorplan := Floorplan{-1, len(lines[0]), len(lines)}
  obstructions := make(map[int]struct{})
  for y, line := range lines {
    for x, char := range line {
      if char == '#' {
        i := y * floorplan.width + x
        obstructions[i] = struct{}{}
      } else if char == '^' {
        i := y * floorplan.width + x
        floorplan.guardPos = i
      }
    }
  }
  return &floorplan, &obstructions
}

func Init(ver constants.VersionIndex) {
  lines, err := io.GetLinesFor(constants.Six, ver)
  if (err != nil) {
    panic(fmt.Sprintf("Error loading file for day %d, version %d: %v", constants.Six, ver, err))
  }
  floorplan, obstructions := newFloorplan(lines)
  part1Result, path := solvePart1(floorplan, *obstructions)
  part2Result := solvePart2(floorplan, *obstructions, path)
  fmt.Printf("Part 1 result: %d\n", part1Result)
  fmt.Printf("Part 2 result: %d\n", part2Result)
}

func solvePart2(floorplan *Floorplan, obstructions map[int]struct{}, path *[]int) int {
  loopCount := 0
  usedObstructions := make(map[int]struct{})
  var obsIn int
  for i := 1; i < len(*path); i++ {
    newObs := make(map[int]struct{})
    obsIn = (*path)[i]
    newObs[obsIn] = struct{}{}
    _, alreadyUsed := usedObstructions[obsIn]; if alreadyUsed {
      continue 
    }
    usedObstructions[obsIn] = struct{}{}
    for obs := range obstructions {
      newObs[obs] = struct{}{}
    }
    visited, _ := solvePart1(floorplan, newObs)
    if visited == -1 {
      loopCount++
    }
  }
  return loopCount
}

func solvePart1(floorplan *Floorplan, obstructions map[int]struct{}) (int, *[]int) {
  path := make([]int, 0)
  visited := make(map[int]struct{})
  turnMap := make(map[string]struct{})
  dir := North
  y := floorplan.guardPos / floorplan.width
  x := floorplan.guardPos % floorplan.width
  for true {
    nextX := x
    nextY := y - 1
    if dir == East {
      nextX = x + 1
      nextY = y
    } else if dir == South {
      nextX = x
      nextY = y + 1
    } else if dir == West {
      nextX = x - 1
      nextY = y
    }
    current := y * floorplan.width + x
    next := nextY * floorplan.width + nextX
    _, isObstructed := obstructions[next]
    if isObstructed {
      _, hasTurnedHereBefore := turnMap[fmt.Sprintf("%d-%d", current, dir)]
      if hasTurnedHereBefore {
        return -1, &path
      }
      turnMap[fmt.Sprintf("%d-%d", current, dir)] = struct{}{}
      dir = dir.turnRight()
      continue
    }
    path = append(path, current)
    _, wasVisited := visited[current]
    if !wasVisited { 
      visited[current] = struct{}{}
    }
    if nextX < 0 || nextY < 0 || nextX >= floorplan.width || nextY >= floorplan.height { break }
    x = nextX
    y = nextY
  }
  return len(visited), &path
}
