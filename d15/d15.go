package d15

import (
	"aoc2k24/constants"
	"aoc2k24/io"
	"fmt"
	"os"
	"os/exec"
	"slices"
	"sort"
	// "time"
)

const isDebug = false
const useDebugData = false
const isAnimated = false

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

type Warehouse2 struct {
  width int
  height int
  boxes map[int]*Box
  walls map[int]struct{}
}

type Robot struct {
  position int
  moves []Direction
}

type Box struct {
  start int
  end int
}

func Init(ver constants.VersionIndex) {
  lines, err := io.GetLinesFor(constants.Fifteen, ver)
  if (err != nil) {
    panic(fmt.Sprintf("Error loading file for day %d, version %d: %v", constants.Fifteen, ver, err))
  }
  warehouse, robot := parseInput(lines)
  warehouse2, robot2 := parseInputPart2(lines)
  if isDebug {
    fmt.Printf("Warehouse: width %d, height %d\n", warehouse2.width, warehouse2.height)
    fmt.Printf("Boxes: %+v\n", warehouse2.boxes)
    fmt.Printf("Walls: %+v\n", warehouse2.walls)
    fmt.Printf("Robot: %+v\n", *robot2)
    renderPart2(warehouse2, robot2)
  }
  gpsSum := solve(warehouse, robot)
  gpsSumPart2 := solvePart2(warehouse2, robot2)
  fmt.Printf("GPS coordinate sum (part 1): %d\n", gpsSum)
  fmt.Printf("GPS coordinate sum (part 2): %d\n", gpsSumPart2)
}

func generateDebugData(isReverse bool) []string {
  warehouse := []string{
    "##################",
    "##..............##",
    "##..............##",
    "##..[][]..[][]..##",
    "##...[]....[]...##",
    "##..[]..[][]....##",
    "##...[]..[].....##",
    "##....[][]......##",
    "##.....[].......##",
    "##......@.......##",
    "##..............##",
    "##..............##",
    "##################",
  }
  instructions := "^^"
  if isReverse {
    slices.Reverse(warehouse)
    instructions = "vv"
  }
  warehouse = append(warehouse, "")
  return append(warehouse, instructions)
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

func canMove(pos []int, w *Warehouse2, dir Direction, r *Robot, boxesToMove *map[int]struct{}) bool {
  result := true
  for _, i := range pos {
    if isDebug { fmt.Printf("Checking if pos %d can move\n", i) }
    nextPos := make([]int, 1)
    box, isCurrentBox := w.boxes[i]
    if isCurrentBox {
      nextPos[0] = box.end + 1
    } else {
      nextPos[0] = i + 1
    }
    switch dir {
    case Up:
      if isCurrentBox {
        nextPos[0] = box.start - w.width
        nextPos = append(nextPos, box.end - w.width)
      } else {
        nextPos[0] = i - w.width
      }
    case Down:
      if isCurrentBox {
        nextPos[0] = box.start + w.width
        nextPos = append(nextPos, box.end + w.width)
      } else {
        nextPos[0] = i + w.width
      }
    case Left:
      if isCurrentBox {
        nextPos[0] = box.start - 1
      } else {
        nextPos[0] = i - 1
      }
    }
    if isCurrentBox && isDebug { fmt.Printf("This is a box starting at %d and ending at %d\n", box.start, box.end) }
    if isDebug { fmt.Printf("Next position(s) are: %+v\n", nextPos) }

    isNextWall := false
    for _, j := range nextPos {
      _, isNextWall = w.walls[j]
      if isNextWall {
        if isDebug { fmt.Printf("Next position %d is a wall\n\n", j) }
        return false
      }
    }
    
    nextBox1, isNextBox1 := w.boxes[nextPos[0]]
    var nextBox2 *Box 
    isNextBox2 := false
    if len(nextPos) > 1 {
      nextBox2, isNextBox2 = w.boxes[nextPos[1]]
    }
    if isNextBox1 {
      if isDebug { fmt.Printf("Will check if box starting at %d and ending at %d can be moved\n", nextBox1.start, nextBox1.end) }
      (*boxesToMove)[nextBox1.start] = struct{}{}
      nextPos[0] = nextBox1.start
    }
    if isNextBox1 && isNextBox2 && nextBox1.start != nextBox2.start || !isNextBox1 && isNextBox2 {
      if isDebug { fmt.Printf("Will check if box starting at %d and ending at %d can be moved\n", nextBox2.start, nextBox2.end) }
      (*boxesToMove)[nextBox2.start] = struct{}{}
      nextPos[1] = nextBox2.start
    }
    if isDebug { fmt.Print("\n") }
    if isNextBox1 || isNextBox2 { result = canMove(nextPos, w, dir, r, boxesToMove) }
  }

  // If next is not wall and not box, then it must be open space
  if isDebug { fmt.Print("All next positions are open space\n\n") }
  return result
}

func movePart2(w *Warehouse2, dir Direction, r *Robot) bool {
  if isDebug { fmt.Printf("Trying to move robot in position %d in direction %d. ", r.position, dir) }
  boxesToMove := make(map[int]struct{})
  couldMove := canMove([]int{r.position}, w, dir, r, &boxesToMove)
  if isDebug { fmt.Printf("Can move? %v\n", couldMove) }
  if couldMove {
    if isDebug {
      fmt.Printf("Boxes to move: %+v\n", boxesToMove)
      fmt.Printf("Boxes before moving: %+v\n", w.boxes)
    }
    boxPositions := make([]int, len(boxesToMove))
    i := 0
    for pos := range boxesToMove {
      boxPositions[i] = pos
      i++
    }
    if dir == Up || dir == Left {
      sort.Sort(sort.IntSlice(boxPositions))
    } else {
      sort.Sort(sort.Reverse(sort.IntSlice(boxPositions)))
    }
    for _, pos := range boxPositions {
      b, _ := w.boxes[pos]
      if isDebug { fmt.Printf("Moving box at position %d (start %d, end %d)\n", pos, b.start, b.end) }
      delete(w.boxes, b.start)
      delete(w.boxes, b.end)
      if isDebug { fmt.Printf("Deleted box from map: %+v\n", w.boxes) }
      if dir == Up {
        b.start = b.start - w.width
        b.end = b.end - w.width
      } else if dir == Down {
        b.start = b.start + w.width
        b.end = b.end + w.width
      } else if dir == Left {
        b.start = b.start - 1
        b.end = b.end - 1
      } else {
        b.start = b.start + 1
        b.end = b.end + 1
      }
      w.boxes[b.start] = b
      w.boxes[b.end] = b
      if isDebug { fmt.Printf("Added new box references to map: %+v\n", w.boxes) }
    }
    if isDebug { fmt.Printf("Boxes after moving: %+v\n", w.boxes) }
    if dir == Up {
      r.position = r.position - w.width
    } else if dir == Down {
      r.position = r.position + w.width
    } else if dir == Left {
      r.position = r.position - 1
    } else {
      r.position = r.position + 1
    }
    if isDebug { fmt.Printf("Robot is now at position %d\n", r.position) }
  }
  return couldMove
}

func solve(warehouse *Warehouse, robot *Robot) int {
  sum := 0
  for _, direction := range robot.moves {
    move(robot.position, warehouse, direction, robot)
  }
  for box := range warehouse.boxes {
    x := box % warehouse.width
    y := box / warehouse.width
    sum += 100 * y + x
  }
  return sum
}

func solvePart2(warehouse *Warehouse2, robot *Robot) int {
  sum := 0
  for i, direction := range robot.moves {
    if isAnimated {
      renderAnimated(i, direction, warehouse, robot)
    }
    movePart2(warehouse, direction, robot)
    if isDebug { renderPart2(warehouse, robot) }
  }
  boxes := make(map[int]struct{})
  for box := range warehouse.boxes {
    boxes[warehouse.boxes[box].start] = struct{}{}
  }
  for box := range boxes {
    x := box % warehouse.width
    y := box / warehouse.width
    sum += 100 * y + x
  }
  return sum
}

func clearScr() {
  c := exec.Command("clear")
  c.Stdout = os.Stdout
  c.Run()
}

func renderAnimated(frame int, dir Direction, warehouse *Warehouse2, robot *Robot) {
  clearScr()
  heading := '^'
  if dir == Right {
    heading = '>'
  } else if dir == Down {
    heading = 'v'
  } else if dir == Left {
    heading = '<'
  }
  fmt.Printf("\n\nFrame: %d || Robot moving %c || Moves pending: %d\n\n", frame + 1, heading, len(robot.moves) - frame)
  i := 0
  for range warehouse.height {
    for range warehouse.width {
      _, isWall := warehouse.walls[i]
      box, isBox := warehouse.boxes[i]
      if isWall {
        fmt.Print("▓")
      } else if isBox {
        if box.start == i {
          fmt.Print("\033[34m[")
        } else {
          fmt.Print("]\033[0m")
        }
      } else if i == robot.position {
        fmt.Print("§")
      } else {
        fmt.Print("\033[32m░\033[0m")
      }
      i++
    }
    fmt.Print("\n")
  }
  fmt.Print("\n\n")
  //time.Sleep(30 * time.Millisecond)
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

func convertInputToPart2(lines *[]string) {
  for i, line := range *lines {
    if len(line) == 0 {
      break
    }
    newLine := ""
    for _, char := range line {
      if char == 'O' {
        newLine += "[]"
      } else if char == '@' {
        newLine += "@."
      } else {
        newLine += fmt.Sprintf("%c%c", char, char)
      }
    }
    (*lines)[i] = newLine
  }
}

func parseInputPart2(lines []string) (*Warehouse2, *Robot) {
  convertInputToPart2(&lines)
  if isDebug && useDebugData { lines = generateDebugData(true) }
  w := Warehouse2{len(lines[0]), 0, make(map[int]*Box), make(map[int]struct{})}
  for y, line := range lines {
    if len(line) == 0 {
      w.height = y
      break
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
        i := y * len(line) + x
        if char == '.' || char == ']' { continue }
        if char == '#' {
          w.walls[i] = struct{}{}
        } else if char == '[' {
          box := Box{i, i + 1}
          w.boxes[i] = &box
          w.boxes[i + 1] = &box
        } else {
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

func renderPart2(w *Warehouse2, r *Robot) {
  fmt.Print("\n\n")
  for y := range w.height {
    for x := 0; x < w.width; x++ {
      i := y * w.width + x
      _, isWall := w.walls[i]
      box, isBox := w.boxes[i]; if isBox { isBox = box.start == i }
      isRobot := i == r.position
      if isWall {
        fmt.Print("#")
      } else if isBox {
        fmt.Print("[]")
        x++
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
