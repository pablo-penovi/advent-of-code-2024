package d16

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

func (c Coord) getPos(dir Direction) Coord {
  if dir == Up {
    return Coord{c.x, c.y - 1}
  }
  if dir == Right {
    return Coord{c.x + 1, c.y}
  }
  if dir == Down {
    return Coord{c.x, c.y + 1}
  }
  return Coord{c.x - 1, c.y}
}

type Path struct {
  nodes []PathNode
  score int
}

func (p Path) getScore() int {
  turns := 0
  for _, node := range p.nodes {
    if node.isTurn { turns++ }
  }
  return turns * 1000 + len(p.nodes)
}

type PathNode struct {
  nodeId int
  isTurn bool
}

type FrontierNode struct {
  next *FrontierNode
  prev *FrontierNode
  h int
  node *MazeNode
  isTurn bool
  dir Direction
  path Path
  visited map[int]*MazeNode
}

type FrontierMap struct {
  first *FrontierNode
  length int
}

func (fm *FrontierMap) add(node *FrontierNode) {
  examined := fm.first
  if examined == nil {
    fm.first = node
  } else {
    for examined != nil {
      if examined.h > node.h {
        if examined == fm.first {
          node.next = fm.first
          fm.first.prev = node
          fm.first = node
        } else {
          node.prev = examined.prev
          node.next = examined
          examined.prev.next = node
          examined.prev = node
        }
        break
      } else if examined.next == nil {
        examined.next = node
        node.prev = examined
        break
      }
      examined = examined.next
    }
  }
  fm.length++
}

func (fm *FrontierMap) pop() *FrontierNode {
  if fm.length == 0 { return nil }
  node := fm.first
  fm.first = fm.first.next
  fm.length--
  return node
}

func (fm *FrontierMap) clear() {
  fm.length = 0
  fm.first = nil
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
  solutions := []Path{}
  
  clearScr()
  renderPlan(maze)
  node := FrontierNode{nil, nil, computeHeuristic(&maze.start.pos, &maze.goal.pos), maze.start, false, Right, Path{[]PathNode{}, 0}, make(map[int]*MazeNode)}
  frontier := FrontierMap{&node, 1}
  next := frontier.pop()
 
  for next != nil {
    processNode(&frontier, maze, &bestScore, &solutions, next)
    next = frontier.pop()
  }

  fmt.Printf("\nSolutions: %d -->\n%+v\n\n", len(solutions), solutions)

  return bestScore, bestSeats
}

func processNode(frontier *FrontierMap, maze *Maze, bestScore *int, solutions *[]Path, fNode *FrontierNode) {
  fmt.Printf("Analyzing node ID %d at x %d, y %d\n", fNode.node.id, fNode.node.pos.x, fNode.node.pos.y)
  // Check if already visited for this solution 
  _, wasVisited := fNode.visited[fNode.node.id]; if wasVisited { fmt.Print("Node already visited\n"); return }

  // Add current node to visited and path
  fNode.visited[fNode.node.id] = fNode.node
  fNode.path.nodes = append(fNode.path.nodes, PathNode{fNode.node.id, fNode.isTurn})

  // Check if this is the goal
  if fNode.node.id == maze.goal.id {
    fmt.Print("\n************** Goal found! *****************\n")
    // Get score and update bestScore and solutions if necessary
    fNode.path.score = fNode.path.getScore()
    if fNode.path.score < *bestScore {
      fmt.Printf("New best score! Previous: %d | New: %d\n", *bestScore, fNode.path.score)
      *bestScore = fNode.path.score
      *solutions = []Path{fNode.path}
    } else if fNode.path.score == *bestScore {
      fmt.Print("Path matches current best score, adding to solutions\n")
      *solutions = append(*solutions, fNode.path)
    } else {
      fmt.Print("This path's score is greater than current best, so no best or equally good solution exist. Exiting\n")
      return
    }
    fmt.Printf("Current solutions: %+v\n", *solutions)
    return
  }

  // Now we get all possible exits from this node
  exits := []Direction{}
  for _, dir := range directions {
    // Count as an exit if node exists in this direction, if it was not visited, and if it is not in the direction whence we came
    if dir != fNode.dir.getOpposite() && fNode.node.links[dir] != nil && fNode.visited[fNode.node.links[dir].id] == nil {
      exits = append(exits, dir)
    }
  }
  // If node has no exits, it's a dead end
  fmt.Printf("Possible exits: %+v\n", exits)
  if len(exits) == 0 {
    fmt.Print("Dead end found!\n")
    return
  }
  fmt.Printf("Possible exits: %+v\n", exits)

  // If there are exits, compute if reindeer will turn to move to each of them. Also add new node to frontier
  for _, dir := range exits {
    isTurn := false; if dir != fNode.dir { isTurn = true }
    heuristic := computeHeuristic(&fNode.node.links[dir].pos, &maze.goal.pos)
    fmt.Printf("Adding %s node (x %d, y %d) to frontier. Is turn? %v. Heuristic: %d\n", dir.toStr(), fNode.node.links[dir].pos.x, fNode.node.links[dir].pos.y, isTurn, heuristic)
    frontier.add(&FrontierNode{nil, nil, heuristic, fNode.node.links[dir], isTurn, dir, fNode.path, fNode.visited})
  }
}

// My estimate will be the diff along x axis + the diff along y axis
func computeHeuristic(pos1 *Coord, pos2 *Coord) int {
  xDiff := pos1.x - pos2.x; if xDiff < 0 { xDiff = -xDiff }
  yDiff := pos1.y - pos2.y; if yDiff < 0 { yDiff = -yDiff }
  return xDiff + yDiff
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
          otherC := c.getPos(dir)
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

func render(m *Maze, dir Direction, currPos *Coord) {
  clearScr()
  for y := range m.height {
    for x := range m.width {
      c := Coord{x, y}
      if c.equals(currPos) {
        fmt.Print(string(Magenta) + dir.toStr()[:1] + string(Reset))
        continue
      }
      if m.start.pos.equals(&c) {
        fmt.Print(string(Green) + string(Start) + string(Reset))
        continue
      }
      if m.goal.pos.equals(&c) {
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
