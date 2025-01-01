package d18

import (
	"aoc2k24/constants"
	"aoc2k24/io"
	"fmt"
	"strconv"
	"strings"
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

func Init(ver constants.VersionIndex) {
  lines, err := io.GetLinesFor(constants.Eighteen, ver)
  if (err != nil) {
    panic(fmt.Sprintf("Error loading file for day %d, version %d: %v", constants.Eighteen, ver, err))
  }
  fallenBytes := 12
  memory := getMemory(&lines, fallenBytes)
  render(memory, Coord{0, 0}, None)
  path := findShortestPath(memory)
  fmt.Printf("Shortest path: %d steps\n%s", len(*path), path.toStr())
}

func findShortestPath(m *Memory) *Path {
  start := Coord{0, 0}
  end := Coord{m.width - 1, m.height - 1}
  var shortest *Path
  nextStep(start, None, []Coord{}, m, &end, shortest)
  return shortest
}

func nextStep(node Coord, prevDir Direction, path Path, m *Memory, end *Coord, shortest *Path) {
  fmt.Printf("Analyzing node %s\n", node.toStr())
  path = append(path, node)

  // First exit condition: current path exceeds or matches length of shortest
  if shortest != nil && len(path) >= len(*shortest) {
    fmt.Print("Path exceeds shortest found! Aborted\n")
    return
  }
  // Second exit condition: Found the exit. If new path is shortest, set it to shortest variable
  if node.equals(end) {
    fmt.Printf("Exit found!\n%s\n", path.toStr())
    if shortest == nil || len(path) < len(*shortest) { fmt.Print("(is new shortest)\n"); shortest = &path }
    return
  }

  for _, dir := range directions {
    fmt.Printf("Exploring direction %s from %s\n", dir.toStr(), node.toStr())
    // Skip direction conditions: It's the direction we came from, next byte is outside memory space, or next byte is corrupted
    if dir.getOpposite() == prevDir { fmt.Print("This is the direction we came from. Aborting\n"); continue }
    newNode := dir.getNext(&node) 
    if newNode.isOutside(m) { fmt.Print("Going in this direction steps out of memory space. Aborting\n"); continue }
    _, isCorrupted := m.corrupted[newNode]
    if isCorrupted { fmt.Print("This byte is corrupted. Aborting\n"); continue }
    // For each valid direction we explore that path
    nextStep(newNode, dir, path, m, end, shortest)
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

func render(m *Memory, pos Coord, dir Direction) {
  for y := range m.height {
    line := ""
    for x := range m.width {
      c := Coord{x, y}
      _, isCorrupted := m.corrupted[c];
      if isCorrupted {
        line += "#"
      } else if c.equals(&pos) {
        line += dir.toStr()
      } else {
        line += "."
      }
    }
    fmt.Printf("%s\n", line)
  }
}
