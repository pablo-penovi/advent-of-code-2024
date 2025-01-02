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
const isRender = true


type Color string

const (
  Reset Color = "\033[0m"
  Green = "\033[32m"
  Red = "\033[31m"
  Magenta = "\033[35m"
  White = "\033[97m"
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

type Memory struct {
  width int
  height int
  corrupted map[Coord]struct{}
}

func (m Memory) getHeuristic(c Coord) int {
  xDiff := m.width - 1 - c.x; if xDiff < 0 { xDiff = -xDiff }
  yDiff := m.height - 1 - c.y; if yDiff < 0 { yDiff = -yDiff }
  return xDiff + yDiff
}

type NextStep struct {
  dir Direction
  destination Coord
  h int
  next *NextStep
}

type Frontier struct {
  first *NextStep
  length int
}

func (f *Frontier) add(ns *NextStep) {
  if f.length == 0 {
    f.first = ns
    f.length++
    return
  }
  s := f.first
  for s != nil {
    if ns.h > s.h {
      if s.next == nil {
        s.next = ns
        break
      }
      s = s.next
    }
    ns.next = s.next
    s.next = ns
    break
  }
  f.length++
}

func (f *Frontier) pop() *NextStep {
  if f.length == 0 {
    return nil
  }
  s := f.first
  f.first = s.next
  f.length--
  return s
}

func newFrontier(start Coord) *Frontier {
  f := Frontier{nil, 0}
  f.add(&NextStep{None, start, 999_999_999_999, nil})
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
  fmt.Printf("Shortest path: %d steps\n%s\n", len(*path) - 1, path.toStr())
}

func findShortestPath(m *Memory) *Path {
  start := Coord{0, 0}
  end := Coord{m.width - 1, m.height - 1}
  shortest := Path{}
  frontier := newFrontier(start)
  seen := map[Coord]struct{}{start: {}}
  for frontier.length > 0 {
    nextStep(frontier.pop(), Path{}, frontier, &seen, m, &end, &shortest)
  }
  return &shortest
}

func nextStep(step *NextStep, path Path, frontier *Frontier, seen *map[Coord]struct{}, m *Memory, end *Coord, shortest *Path) {
  // For part 1 at least, since we only need one answer, stop recursion once a shortest path has been found
  if len(*shortest) > 0 { return }

  if isDebug { fmt.Printf("Analyzing node %s\n", step.destination.toStr()) }
  path = append(path, step.destination)

  // First exit condition: current path exceeds or matches length of shortest
  if len(*shortest) > 0 && len(path) >= len(*shortest) {
    if isDebug { fmt.Print("\nPath exceeds shortest found! Aborted\n") }
    return
  }
  // Second exit condition: Found the exit. If new path is shortest, set it to shortest variable
  if step.destination.equals(end) {
    if isDebug { fmt.Printf("Exit found [%d steps]!\n%s\n", len(path) - 1, path.toStr()) }
    if len(*shortest) == 0 || len(path) < len(*shortest) {
      if isDebug { fmt.Print("This is the new shortest path\n") }
      *shortest = path
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
    // For each valid direction we explore that path
    if isDebug { fmt.Print("Yes\n") }
    (*seen)[step.destination] = struct{}{}
    frontier.add(&NextStep{dir, newNode, m.getHeuristic(newNode), nil})

    if isRender {
      render(m, step.destination, step.dir, frontier)
      time.Sleep(20 * time.Millisecond)
    }
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

func render(m *Memory, pos Coord, dir Direction, frontier *Frontier) {
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
      if isCorrupted {
        if isFrontier {
          line += string(Red) + Wall + string(Reset)
        } else {
          line += string(White) + Wall + string(Reset)
        }
      } else if c.equals(&pos) {
        line += string(Magenta) + dir.toStr() + string(Reset)
      } else {
        if isFrontier {
          line += string(Red) + Ground + string(Reset)
        } else {
          line += string(Green) + Ground + string(Reset)
        }
      }
    }
    fmt.Printf("%s\n", line)
  }
}

func clearScr() {
  cmd := exec.Command("clear")
  cmd.Stdout = os.Stdout
  cmd.Run()
}

