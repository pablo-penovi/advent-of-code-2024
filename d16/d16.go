package d16

import (
	"aoc2k24/constants"
	"aoc2k24/io"
	"fmt"
	"os"
	"os/exec"
	"time"
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

func (d Direction) getNext(c *Coord) *Coord {
  if d == Up { return &Coord{c.x, c.y - 1} }
  if d == Right { return &Coord{c.x + 1, c.y} }
  if d == Down { return &Coord{c.x, c.y + 1} }
  return &Coord{c.x - 1, c.y}
}

type Coord struct {
  x int
  y int
}

func (c Coord) equals(o *Coord) bool {
  return c.x == o.x && c.y == o.y
}

func (c Coord) toKey() string {
  return fmt.Sprintf("%d-%d", c.x, c.y)
}

func (c Coord) getDirTo(o *Coord) Direction {
  difference := c.x - o.x + c.y - o.y
  if difference == 0 {
    panic("Trying to get direction from a coordinate and itself!")
  }
  if difference > 1 || difference < -1 {
    panic("Trying to get direction from two coordinates that are not adjacent!")
  }
  if c.x > o.x {
    return Left
  } else if c.x < o.x {
    return Right
  } else if c.y > o.y {
    return Up
  }
  return Down
}

type Maze struct {
  height int
  width int
  start Coord
  goal Coord
  walls map[Coord]struct{}
}

type Path struct {
  start *Node
  score int
}

type Node struct {
  pos Coord
  up *Node
  right *Node
  down *Node
  left *Node
}

type MemoItem struct {
  node *Node
  score int
}

func Init(ver constants.VersionIndex) {
  lines, err := io.GetLinesFor(constants.Sixteen, ver)
  if (err != nil) {
    panic(fmt.Sprintf("Error loading file for day %d, version %d: %v", constants.Sixteen, ver, err))
  }
  maze := parseInput(lines)
  bestScore := solve(maze)
  fmt.Printf("\nBest path score: %d\n", bestScore)
}

func solve(maze *Maze) int {
  currPos := maze.start
  prevPos := Coord{-1, -1}
  path := Path{&Node{maze.start, nil, nil, nil, nil}, 0}
  memo := make(map[string]MemoItem)
  memo[maze.start.toKey()] = MemoItem{path.start, 0}
  findBestScoreToGoal(&currPos, &prevPos, Right, maze, &path, path.start, 0, &memo)
  return path.score
}

func findBestScoreToGoal(pos *Coord, prevPos *Coord, dir Direction, maze *Maze, path *Path, current *Node, accumulator int, memo *map[string]MemoItem) {
  if isDebug {
    fmt.Printf("Analyzing for pos x %d, y %d. Accumulator: %d\n", pos.x, pos.y, accumulator)
    render(maze, dir, pos)
  }
  
  // End condition: Goal reached
  if maze.goal.equals(pos) {
    if path.score == 0 || path.score > accumulator + 1 {
      path.score = accumulator + 1
    }
    return
  }
  // Get possible next positions (not previous position && not wall)
  nextPositions := getNextPositions(pos, prevPos, maze)
  // End condition: If no next positions, this is a dead end
  if len(*nextPositions) == 0 { return }

  for _, nextPos := range *nextPositions {
    currNodeScore := 1; if pos.equals(&maze.start) { currNodeScore = 0 }
    newDir := pos.getDirTo(&nextPos)

    // If next node hasn't been visited before, create it and add pointer to memo
    next, isInMemo := (*memo)[nextPos.toKey()];
    if !isInMemo {
      next = MemoItem{&Node{nextPos, nil, nil, nil, nil}, 0}
    }

    // Now link new next node to correct direction in current node
    if newDir == Up {
      current.up = next.node
    } else if newDir == Right {
      current.right = next.node
    } else if newDir == Down {
      current.down = next.node
    } else {
      current.left = next.node
    }

    // If the reindeer turns at this point, add 1000 score to current node
    if newDir != dir { currNodeScore += 1000 }

    // If next node is new, add current score to it in memory
    // If next node is in memory and its memory score is lower or equal than its curent score, skip
    if isDebug { fmt.Printf("Node x%d, y %d || In memo? %v || Score before: %d || Accumulator: %d\n", next.node.pos.x, next.node.pos.y, isInMemo, next.score, accumulator + currNodeScore) }
    if !isInMemo {
      next.score = accumulator + currNodeScore
      (*memo)[nextPos.toKey()] = next
    } else if next.score <= accumulator + currNodeScore {
      continue
    } else {
      next.score = accumulator + currNodeScore
      (*memo)[nextPos.toKey()] = next
    }
    if isDebug { fmt.Printf("Score after: %d\n", next.score) }
    // Finally, recurse for new node
    findBestScoreToGoal(&nextPos, pos, newDir, maze, path, next.node, accumulator + currNodeScore, memo)
  }
}

func getNextPositions(pos *Coord, prevPos *Coord, maze *Maze) *[]Coord {
  coords := make([]Coord, 0)
  possibleDirs := []Direction{Up, Right, Down, Left}
  for _, dir := range possibleDirs {
    nextPos := dir.getNext(pos)
    if nextPos.equals(prevPos) { continue }
    _, isWall := maze.walls[*nextPos]
    if isWall { continue }
    coords = append(coords, *nextPos)
  }
  return &coords
}

func parseInput(lines []string) *Maze {
  m := Maze{len(lines), len(lines[0]), Coord{-1, -1}, Coord{-1, -1}, make(map[Coord]struct{})}
  for y, line := range lines {
    for x, char := range line {
      if char == '.' { continue }
      if char == '#' {
        m.walls[Coord{x, y}] = struct{}{}
      } else if char == 'S' {
        m.start.x = x
        m.start.y = y
      } else {
        m.goal.x = x
        m.goal.y = y
      }
    }
  }
  return &m
}

func render(m *Maze, dir Direction, currPos *Coord) {
  clearScr()
  for y := range m.height {
    for x := range m.width {
      c := Coord{x, y}
      if c.equals(currPos) {
        fmt.Print(string(Magenta) + dir.toStr()[:1] + string(Reset))
        continue
      }
      if m.start.equals(&c) {
        fmt.Print(string(Green) + string(Start) + string(Reset))
        continue
      }
      if m.goal.equals(&c) {
        fmt.Print(string(Red) + string(Goal) + string(Reset))
        continue
      }
      _, isWall := m.walls[c]
      if isWall {
        fmt.Print(string(White) + string(Wall) + string(Reset))
        continue
      }
      fmt.Print(string(Green) + string(Ground) + string(Reset))
    }
    fmt.Print("\n")
  }
  time.Sleep(100 * time.Millisecond)
}

func clearScr() {
  c := exec.Command("clear")
  c.Stdout = os.Stdout
  c.Run()
}
