package d15

import (
	"aoc2k24/constants"
	"aoc2k24/io"
	"fmt"
)

const isDebug = false

type Direction int

const (
  Up Direction = iota
  Right
  Down
  Left
)

type Warehouse struct {
  width int
  height int
  boxes map[int]struct{}
  walls map[int]struct{}
}

type Robot struct {
  position int
  moves []Direction
}

func Init(ver constants.VersionIndex) {
  lines, err := io.GetLinesFor(constants.Fifteen, ver)
  if (err != nil) {
    panic(fmt.Sprintf("Error loading file for day %d, version %d: %v", constants.Fifteen, ver, err))
  }
  warehouse, robot := parseInput(lines)
  if isDebug {
    fmt.Printf("Warehouse: width %d, height %d\n", warehouse.width, warehouse.height)
    fmt.Printf("Boxes: %+v\n", warehouse.boxes)
    fmt.Printf("Walls: %+v\n", warehouse.walls)
    fmt.Printf("Robot: %+v\n", *robot)
    render(warehouse, robot)
  }
  gpsSum := solve(warehouse, robot)
  fmt.Printf("GPS coordinate sum (part 1): %d", gpsSum)
}

func move(i int, w *Warehouse, dir Direction, r *Robot) bool {
  j := i + 1
  switch dir {
  case Up:
    j = i - w.width
  case Down:
    j = i + w.width
  case Left:
    j = i - 1
  }
  
  _, isNextWall := w.walls[j]; if isNextWall { return false }
  _, isNextBox := w.boxes[j];
  // If next is box (else is empty space)
  if isNextBox {
    if !move(j, w, dir, r) { return false }
  }
  // If current is robot (else is box)
  if r.position == i {
    r.position = j
  } else {
    delete(w.boxes, i)
    w.boxes[j] = struct{}{}
  }
  return true
}

func solve(warehouse *Warehouse, robot *Robot) int {
  sum := 0
  for _, direction := range robot.moves {
    move(robot.position, warehouse, direction, robot)
  }
  if isDebug { render(warehouse, robot) }
  for box := range warehouse.boxes {
    x := box % warehouse.width
    y := box / warehouse.width
    sum += 100 * y + x
  }
  return sum
}

func parseInput(lines []string) (*Warehouse, *Robot) {
  w := Warehouse{len(lines[0]), 0, make(map[int]struct{}), make(map[int]struct{})}
  for y, line := range lines {
    if len(line) == 0 {
      w.height = y
    }
  }
  r := Robot{0, make([]Direction, 0)}
  parseWarehouse := true
  for y, line := range lines {
    if len(line) == 0 {
      parseWarehouse = false
    }
    if parseWarehouse {
      for x, char := range line {
        i := y * w.width + x
        if char == '.' { continue }
        if char == '#' {
          w.walls[i] = struct{}{}
        } else if char == 'O' {
          w.boxes[i] = struct{}{}
        } else {
          fmt.Print("is robot")
          r.position = i
        }
      }
      continue
    }
    for _, char := range line {
      toAppend := Up
      if char == '>' {
        toAppend = Right
      } else if char == 'v' {
        toAppend = Down
      } else if char == '<'{
        toAppend = Left
      }
      r.moves = append(r.moves, toAppend)
    }
  }
  return &w, &r
}

func render(w *Warehouse, r *Robot) {
  fmt.Print("\n\n")
  for y := range w.height {
    for x := range w.width {
      i := y * w.width + x
      _, isWall := w.walls[i]
      _, isBox := w.boxes[i]
      isRobot := i == r.position
      if isWall {
        fmt.Print("#")
      } else if isBox {
        fmt.Print("O")
      } else if isRobot {
        fmt.Print("@")
      } else {
        fmt.Print(".")
      }
    }
    fmt.Print("\n")
  }
  fmt.Print("\n")
}
