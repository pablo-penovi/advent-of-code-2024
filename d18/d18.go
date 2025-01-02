package d18

import (
	"aoc2k24/constants"
	"aoc2k24/io"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

const isDebug = false
const isRender = false

type Color string

const (
  Reset Color = "\033[0m"
  Green = "\033[32m"
  Red = "\033[31m"
  Magenta = "\033[35m"
  White = "\033[97m"
  Blue = "\033[34m"
)

const (
  Wall = "▓"
  Ground = "░"
)

type Direction int 

func (d Direction) getOpposite() Direction {
  if d == None { return None }
  if d == Up { return Down }
  if d == Right { return Left }
  if d == Down { return Up }
  return Right
}

func (d Direction) getNext(c *Coord) Coord {
  newC := Coord{c.x, c.y}
  if d == None {
    panic("You shouldn't do getNext over a NONE direction!")
  } else if d == Up {
    newC.y--
  } else if d == Down {
    newC.y++
  } else if d == Right {
    newC.x++
  } else {
    newC.x--
  }
  return newC
}

func (d Direction) toStr() string {
  if d == None { return "N" }
  if d == Up { return "U" }
  if d == Right { return "R" }
  if d == Down { return "D" }
  return "L"
}

const (
  Up Direction = iota
  Right
  Down
  Left
  None
)

var directions = [4]Direction{Up, Right, Down, Left}

type Coord struct {
  x int
  y int
}

func (c Coord) toStr() string {
  return fmt.Sprintf("(%d, %d)", c.x, c.y)
}

func (c Coord) equals(o *Coord) bool {
  return c.x == o.x && c.y == o.y
}

func (c Coord) isOutside(m *Memory) bool {
  return c.x < 0 || c.y < 0 || c.x >= m.width || c.y >= m.height
}

type Path []Coord

func (p Path) toStr() string {
  output := ""
  for i, node := range p {
    output += node.toStr()
    if i < len(p) - 1 {
      output += " -> "
    }
  }
  return output
}

func (p Path) has(c *Coord) bool {
  for _, node := range p {
    if node.equals(c) { return true }
  }
  return false
}

type Memory struct {
  width int
  height int
  corrupted map[Coord]struct{}
}

func (m Memory) getHeuristic(dir Direction, c Coord) int {
  xDiff := m.width - 1 - c.x; if xDiff < 0 { xDiff = -xDiff }
  yDiff := m.height - 1 - c.y; if yDiff < 0 { yDiff = -yDiff }
  h := xDiff + yDiff
  if xDiff > yDiff && dir == Right || yDiff > xDiff && dir == Down {
    h -= 1
  } else if (xDiff > yDiff && dir == Down || dir == Up) || yDiff > xDiff && dir == Right || dir == Left {
    h += 1
  }
  return h
}

type NextStep struct {
  dir Direction
  destination Coord
  path Path
  next *NextStep
  prev *NextStep
}

type Frontier struct {
  first *NextStep
  last *NextStep
  length int
}

func (f *Frontier) append(ns *NextStep) {
  if f.length == 0 {
    f.first = ns
  }
  ns.prev = f.last
  if f.last != nil {
    f.last.next = ns
  }
  f.last = ns
  f.length++
}

func (f *Frontier) pop() *NextStep {
  if f.length == 0 {
    return nil
  }
  s := f.first
  f.length--
  if f.length == 0 {
    f.first = nil
    f.last = nil
  } else {
    f.first = s.next
    f.first.prev = nil
  }
  return s
}

func newFrontier(start Coord) *Frontier {
  f := Frontier{nil, nil, 0}
  f.append(&NextStep{None, start, Path{}, nil, nil})
  return &f
}

func Init(ver constants.VersionIndex) {
  lines, err := io.GetLinesFor(constants.Eighteen, ver)
  if (err != nil) {
    panic(fmt.Sprintf("Error loading file for day %d, version %d: %v", constants.Eighteen, ver, err))
  }
  fallenBytes := 1024
  memory := getMemory(&lines, fallenBytes)
  path := findShortestPath(memory)
  renderPath(memory, path)
}

func findShortestPath(m *Memory) *Path {
  start := Coord{0, 0}
  end := Coord{m.width - 1, m.height - 1}
  shortest := Path{}
  frontier := newFrontier(start)
  seen := map[Coord]struct{}{start: {}}
  for frontier.length > 0 {
    nextStep(frontier.pop(), frontier, &seen, m, &end, &shortest)
    if len(shortest) > 0 {
      break
    }
  }
  return &shortest
}

func nextStep(step *NextStep, frontier *Frontier, seen *map[Coord]struct{}, m *Memory, end *Coord, shortest *Path) {
  if isDebug { fmt.Printf("Analyzing node %s\n", step.destination.toStr()) }
  step.path = append(step.path, step.destination)
  log := fmt.Sprintf("Analyzing %s: path so far: %s\n", step.destination.toStr(), step.path.toStr())

  // Second exit condition: Found the exit. If new path is shortest, set it to shortest variable
  if step.destination.equals(end) {
    if isDebug { fmt.Printf("Exit found [%d steps]!\n%s\n", len(step.path) - 1, step.path.toStr()) }
    if len(*shortest) == 0 || len(step.path) < len(*shortest) {
      if isDebug { fmt.Print("This is the new shortest path\n") }
      *shortest = step.path
    }
    return
  }
  
  for _, dir := range directions {
    if isDebug { fmt.Printf("Explore direction %s from %s? ", dir.toStr(), step.destination.toStr()) }
    // Skip direction conditions: Has already been visited, it's the direction we came from, next byte is outside memory space, or next byte is corrupted
    if dir.getOpposite() == step.dir { if isDebug { fmt.Print("No, we came from here\n") }; continue }
    newNode := dir.getNext(&step.destination) 
    if newNode.isOutside(m) { if isDebug { fmt.Print("No, this steps out of memory\n") }; continue }
    _, isCorrupted := m.corrupted[newNode]
    if isCorrupted { if isDebug { fmt.Print("No, the next byte is corrupted\n") }; continue }
    _, wasSeen := (*seen)[newNode]
    if wasSeen { if isDebug { fmt.Print("No, it was visited before\n") }; continue }
    // For each valid direction we add those new nodes to the back of the frontier queue
    if isDebug { fmt.Print("Yes\n") }
    (*seen)[newNode] = struct{}{}
    newPath := make(Path, len(step.path))
    for i := range len(step.path) {
      newPath[i] = step.path[i]
    }
    frontier.append(&NextStep{dir, newNode, newPath, nil, nil})
  }

  if isRender {
    render(m, step, frontier)
    fmt.Print(log)
    time.Sleep(2000 * time.Millisecond)

  }
}

func getMemory(lines *[]string, fallenBytes int) *Memory {
  m := Memory{0, 0, make(map[Coord]struct{})}
  for i, byte := range *lines {
    coordVals := strings.Split(byte, ",")
    x, _ := strconv.Atoi(coordVals[0])
    y, _ := strconv.Atoi(coordVals[1])
    if x > m.width {
      m.width = x
    }
    if y > m.height {
      m.height = y
    }
    if i < fallenBytes {
      c := Coord{x, y}
      m.corrupted[c] = struct{}{}
    }
  }
  m.width++
  m.height++
  return &m
}

func render(m *Memory, step *NextStep, frontier *Frontier) {
  fMap := make(map[Coord]struct{})
  s := frontier.first
  for s != nil {
    fMap[s.destination] = struct{}{}
    s = s.next
  }
  clearScr()
  for y := range m.height {
    line := ""
    for x := range m.width {
      c := Coord{x, y}
      _, isCorrupted := m.corrupted[c]
      _, isFrontier := fMap[c]
      isPath := step.path.has(&c)
      if isCorrupted {
        if isFrontier {
          line += string(Red) + Wall + string(Reset)
        } else {
          line += string(White) + Wall + string(Reset)
        }
      } else if c.equals(&step.destination) {
        line += string(Magenta) + step.dir.toStr() + string(Reset)
      } else {
        if isFrontier {
          line += string(Red) + Ground + string(Reset)
        } else if isPath {
          line += string(Blue) + Ground + string(Reset)
        } else {
          line += string(Green) + Ground + string(Reset)
        }
      }
    }
    fmt.Printf("%s\n", line)
  }
}

func renderPath(m *Memory, path *Path) {
  clearScr()
  for y := range m.height {
    line := ""
    for x := range m.width {
      c := Coord{x, y}
      _, isCorrupted := m.corrupted[c]
      isPath := path.has(&c)
      if isCorrupted {
        line += string(White) + Wall + string(Reset)
      } else {
        if isPath {
          line += string(Red) + Ground + string(Reset)
        } else {
          line += string(Green) + Ground + string(Reset)
        }
      }
    }
    fmt.Printf("%s\n", line)
  }
  fmt.Printf("Path: %d steps\n%s\n", len(*path) - 1, path.toStr())
}

func clearScr() {
  cmd := exec.Command("clear")
  cmd.Stdout = os.Stdout
  cmd.Run()
}

