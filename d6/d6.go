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
  obstructions map[int]struct{}
  width int
  height int
}

func newFloorplan(lines []string) *Floorplan {
  floorplan := Floorplan{-1, make(map[int]struct{}), len(lines[0]), len(lines)}
  for y, line := range lines {
    for x, char := range line {
      if char == '#' {
        i := y * floorplan.width + x
        floorplan.obstructions[i] = struct{}{}
      } else if char == '^' {
        i := y * floorplan.width + x
        floorplan.guardPos = i
      }
    }
  }
  return &floorplan
}

func Init(ver constants.VersionIndex) {
  lines, err := io.GetLinesFor(constants.Six, ver)
  if (err != nil) {
    panic(fmt.Sprintf("Error loading file for day %d, version %d: %v", constants.Six, ver, err))
  }
  floorplan := newFloorplan(lines)
  part1Result := solvePart1(floorplan)
  fmt.Printf("Part 1 result: %d\n", part1Result)
}

func solvePart1(floorplan *Floorplan) int {
  visited := make(map[int]struct{})
  dir := North
  y := floorplan.guardPos / floorplan.width
  x := floorplan.guardPos % floorplan.width
  visited[floorplan.guardPos] = struct{}{}
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
    if nextX < 0 || nextY < 0 || nextX >= floorplan.width || nextY >= floorplan.height { break }
    next := nextY * floorplan.width + nextX
    _, isObstructed := floorplan.obstructions[next]
    if isObstructed {
      dir = dir.turnRight()
      continue
    }
    x = nextX
    y = nextY
    _, wasVisited := visited[next]
    if wasVisited { continue }
    visited[next] = struct{}{}
  }
  return len(visited)
}
