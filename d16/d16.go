package d16

import (
	"aoc2k24/constants"
	"aoc2k24/io"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const isDebug = false

type Tile rune
type Color string
type Direction int

const (
  Reset Color = "\033[0m"
  Green = "\033[32m"
  Red = "\033[31m"
  Magenta = "\033[35m"
  White = "\033[97m"
)

const (
  Start Tile = 'S'
  Goal = 'G'
  Wall = '▓'
  Ground = '░'
  Reindeer = '¥'
)

const (
  Up Direction = iota
  Right
  Down
  Left
  NullDir
)

func (d Direction) toStr() string {
  if d == Up { return "Up" }
  if d == Right { return "Right" }
  if d == Down { return "Down" }
  return "Left"
}

func (d Direction) getOpposite() Direction {
  if d == Up { return Down }
  if d == Right { return Left }
  if d == Down { return Up }
  return Right
}

var directions = []Direction{Up, Right, Down, Left}

type Coord struct {
  x int
  y int
}

func (c Coord) equals(o *Coord) bool {
  return c.x == o.x && c.y == o.y
}

func (c Coord) getNextInDir(dir Direction) Coord {
  if dir == Up {
    return Coord{c.x, c.y - 1}
  }
  if dir == Down {
    return Coord{c.x, c.y + 1}
  }
  if dir == Right {
    return Coord{c.x + 1, c.y}
  }
  return Coord{c.x - 1, c.y}
}

type Maze struct {
  width int
  height int
  walls map[Coord]struct{}
  start *MazeNode
  goal *MazeNode
  nodes map[Coord]*MazeNode
}

type MazeNode struct {
  id int
  pos Coord
  links map[Direction]*MazeNode
}

func Init(ver constants.VersionIndex) {
  lines, err := io.GetLinesFor(constants.Sixteen, ver)
  if (err != nil) {
    panic(fmt.Sprintf("Error loading file for day %d, version %d: %v", constants.Sixteen, ver, err))
  }
  maze := parse(lines)
  bestScore, bestSeats := solve(maze)
  fmt.Printf("\nBest Score: %d || Best seats: %d\n", bestScore, bestSeats)
}

func solve(maze *Maze) (int, int) {
  // OK here's the idea. I'll use Dijkstra to find the best path, but won't stop searching when goal is reached
  // Instead, I'll keep the score from that first path that reaches the goal (the best path)
  // Then I'll keep finding new solutions until the solution score is greater than the best score
  // At that point I know I have exhausted all best paths (ie. all paths that share the best score)
  // For each solution I'll keep its nodes in a set. That way when I count that set at the end, I'll
  // get the unique nodes belonging to each best path (the best seats asked for in part 2)
  // Sounds neat right?
  bestSeats := 0
  bestScore := 9223372036854775807

  return bestScore, bestSeats
}

func parse(lines []string) *Maze {
  m := Maze{len(lines[0]), len(lines), make(map[Coord]struct{}), nil, nil, make(map[Coord]*MazeNode)}
  nodeId := 0
  for y, line := range lines {
    for x, char := range line {
      c := Coord{x, y}
      if char == '#' {
        m.walls[c] = struct{}{}
      } else {
        node := MazeNode{nodeId, c, make(map[Direction]*MazeNode)}
        nodeId++
        for _, dir := range directions {
          otherC := c.getNextInDir(dir)
          otherNode, isMapped := m.nodes[otherC]; if !isMapped { continue }
          node.links[dir] = otherNode
          otherNode.links[dir.getOpposite()] = &node
        }
        m.nodes[c] = &node
        if char == 'S' { m.start = &node }
        if char == 'E' { m.goal = &node }
      }
    }
  }
  return &m
}

func renderPlan(m *Maze) {
  maxId := 0
  for coord := range m.nodes {
    nId := m.nodes[coord].id
    if maxId < nId { maxId = nId }
  }
  cellSize := len(strconv.Itoa(maxId))
  horLine := strings.Repeat("─", cellSize)
  wallLine := strings.Repeat("#", cellSize)
  horFrameTop := string(White) + "┌" + horLine + strings.Repeat("┬" + horLine, m.width - 1) + "┐" + string(Reset)
  horFrameBetween := string(White) + "├" + horLine + strings.Repeat("┼" + horLine, m.width - 1) + "┤" + string(Reset)
  horFrameBottom := string(White) + "└" + horLine + strings.Repeat("┴" + horLine, m.width - 1) + "┘" + string(Reset)
  fmt.Printf("%s\n", horFrameTop)
  for y := range m.height {
    for x := range m.width {
      fmt.Print(string(White) + "│" + string(Reset))
      c := Coord{x, y}
      node := m.nodes[c]
      if m.start.pos.equals(&c) {
        buf := strings.Repeat(" ", cellSize - len(strconv.Itoa(node.id)))
        fmt.Printf("%s%s%d%s", string(Magenta), buf, node.id, string(Reset))
        continue
      }
      if m.goal.pos.equals(&c) {
        buf := strings.Repeat(" ", cellSize - len(strconv.Itoa(node.id)))
        fmt.Printf("%s%s%d%s", string(Red), buf, node.id, string(Reset))
        continue
      }
      if node == nil {
        fmt.Print(string(White) + wallLine + string(Reset))
        continue
      }
      buf := strings.Repeat(" ", cellSize - len(strconv.Itoa(node.id)))
      fmt.Printf("%s%s%d%s", string(Green), buf, node.id, string(Reset))
    }
    fmt.Print(string(White) + "│" + string(Reset) + "\n")
    if y < m.height - 1 { 
      fmt.Print(horFrameBetween + "\n")
    } else {
      fmt.Print(horFrameBottom + "\n\n")
    }
  }
}

func clearScr() {
  c := exec.Command("clear")
  c.Stdout = os.Stdout
  c.Run()
}
