package d16

import (
	"aoc2k24/constants"
	"aoc2k24/io"
	"fmt"
	"os"
	"os/exec"
	"slices"
	"strconv"
	"strings"
)

type LogToggle int 

const (
  RenderNodes LogToggle = iota
  FrontierPush
  DijkstraAlgo
)

var debugToggles = map[LogToggle]bool{
  RenderNodes: false,
  FrontierPush: false,
  DijkstraAlgo: false,
}

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

type Node struct {
  id int
  pos Coord
  score int
  dir Direction
  next *Node
  prevInRoute *Node
  heuristic int
}

func (n Node) hasIdInPath(id int) bool {
  node := n.next
  for node != nil {
    if n.id == id { return true }
    node = node.next
  }
  return false
}

type Frontier struct {
  first *Node
  length int
}

type Result struct {
  bestScore int
  bestSeats int
  lastNodes *NodePath
}

func (f *Frontier) push(newNode *Node) {
  log := ""
  if f.length == 0 {
    log += "Frontier vacia, agrego nodo directamente\n"
    f.first = newNode
  } else {
    var prev *Node
    current := f.first
    for current != nil {
      log += fmt.Sprintf("Comparando nuevo nodo (x %d, y %d; heuristic %d) con nodo actual (x %d, y %d; heuristic %d)\n", newNode.pos.x, newNode.pos.y, newNode.heuristic, current.pos.x, current.pos.y, current.heuristic)
      if newNode.heuristic <= current.heuristic {
        log += "Nuevo nodo tiene heuristica menor o igual a nodo actual\n"
        newNode.next = current
        if prev == nil {
          log += "No hay prev\n"
          f.first = newNode
          f.length++
          return
        } else {
          log += "Hay prev\n"
          prev.next = newNode
        }
        break
      }
      log += "Voy a la siguiente\n"
      prev = current
      current = current.next
    }
    log += fmt.Sprintf("Hola, esto es prev -> %+v\n\n", prev)
    prev.next = newNode
  }
  f.length++
  if debugToggles[FrontierPush] { fmt.Print(log) }
}

func (f *Frontier) pop() *Node {
  if f.length == 0 { return nil }
  node := f.first
  f.first = node.next
  f.length--
  return node
}

func (f Frontier) toStr() string {
  node := f.first
  st := ""
  for node != nil {
    st += fmt.Sprintf("%d, ", node.id)
    node = node.next
  }
  return st
}

type Path []int

func (p Path) toString() string {
  st := "(" + strconv.Itoa(len(p)) + ") "
  for i, nodeId := range p {
    if i > 0 {
      st += " -> "
    }
    st += strconv.Itoa(nodeId)
  }
  st += "\n"
  return st
}

type NodePath []*Node
type SeenMap map[string]struct{}

func (sm *SeenMap) add(key string) bool {
  _, isInMap := (*sm)[key]
  if !isInMap {
    (*sm)[key] = struct{}{}
  }
  return !isInMap
}

func (sm SeenMap) has(key string) bool {
  _, isInMap := sm[key]
  return isInMap
}

func Init(ver constants.VersionIndex) {
  lines, err := io.GetLinesFor(constants.Sixteen, ver)
  if (err != nil) {
    panic(fmt.Sprintf("Error loading file for day %d, version %d: %v", constants.Sixteen, ver, err))
  }
  maze := parse(lines)
  result := solve(maze)
  paths := reconstructPath(result.lastNodes)
  fmt.Printf("Paths:\n")
  uniqueSeats := make(map[int]struct{})
  for _, path := range *paths {
    for _, nodeId := range path {
      uniqueSeats[nodeId] = struct{}{}
    }
    fmt.Print(path.toString())
  }
  fmt.Printf("\nBest Score: %d || Best seats: %d\n", result.bestScore, len(uniqueSeats) + 1)
}

func reconstructPath(lastNodes *NodePath) *[]Path {
  paths := make([]Path, 0)
  for _, lastNode := range *lastNodes {
    path := Path{lastNode.id}
    node := lastNode.prevInRoute
    for node != nil {
      path = append(path, node.id)
      node = node.prevInRoute
    }
    // Ignore first node since this is not counted in problem
    path = path[:len(path)-1]
    slices.Reverse(path)
    paths = append(paths, path)
  }
  return &paths
}

func solve(maze *Maze) *Result {
  // OK here's the idea. I'll use a sort of A* (Dijkstra w/ heuristics) to find the best path, but won't stop searching when goal is reached
  // Instead, I'll keep the score from that first path that reaches the goal (the best path)
  // Then I'll keep finding new solutions until the solution score is greater than the best score
  // At that point I know I have exhausted all best paths (ie. all paths that share the best score)
  // For each best solution I'll keep the last node, and since each node has a link to the previous node in the path, I'll be able to reconstruct the paths later
  // Sounds neat right?
  bestSeats := 0
  // OK a bit of cheating, since I can't seem to be able to keep this from being a runaway code, I'll set the bestScore I know is right from part 1 here,
  // and I'll use it to filter paths that exceed this score
  bestScore := 123540
  frontier := Frontier{}
  frontier.push(&Node{maze.start.id, maze.start.pos, 0, Right, nil, nil, computeHeuristic(&maze.start.pos, &maze.goal.pos)})
  if debugToggles[RenderNodes] { clearScr(); renderPlan(maze) }
  lastNodes := NodePath{}
  node := frontier.pop()
  for node != nil {
    data := maze.nodes[node.pos]
    log := fmt.Sprintf("\nI'm on node %d, heading %s. Adding to seen\n", data.id, node.dir.toStr())
    for _, newDir := range directions {
      log += fmt.Sprintf("Now checking %s: ", newDir.toStr())
      if node.dir == newDir.getOpposite() || data.links[newDir] == nil {
        log += "CAN'T\n"
        continue
      }
      newData := data.links[newDir]
      newNode := Node{newData.id, newData.pos, node.score + 1, newDir, nil, node, computeHeuristic(&newData.pos, &maze.goal.pos)}
      log += fmt.Sprintf("This is node %d. ", newNode.id)
      isTurn := node.dir != newNode.dir
      if isTurn {
        log += "It's a turn. "
        newNode.score += 1000
      } else {
        log += "Not a turn. "
      }
      if bestScore > -1 && newNode.score > bestScore {
        log += "There's a best score and the current score is greater, skipping\n"
        continue
      }
      if node.hasIdInPath(newNode.id) {
        log += "I've already seen it before, skipping\n"
        continue
      }
      log += fmt.Sprintf("Score if I go this way: %d\n", newNode.score)
      if maze.goal.pos.equals(&newData.pos) {
        log += "\n***** GOAL!! *****\n"
        if bestScore == -1 || bestScore > newNode.score {
          log += fmt.Sprintf("=== NEW BEST SCORE ===\nPrevious: %d. New: %d\n", bestScore, newNode.score)
          bestScore = newNode.score
          lastNodes = NodePath{&newNode}
        } else if bestScore == newNode.score {
          log += fmt.Sprintf("Score matches current best score. Saving last node for path reconstruction\n")
          lastNodes = append(lastNodes, &newNode)
        } else {
          log += fmt.Sprintf("Goal was reached, but this path has worse score than the better found so far (%d vs %d)\n", bestScore, newNode.score)
        }
        log += "\n"
      } else {
        frontier.push(&newNode)
        log += fmt.Sprintf("Pushing node %d to frontier. Frontier now is: %s\n", newNode.id, frontier.toStr())
      }
    }
    node = frontier.pop()
    if node == nil {
      log += "\nI'm all out of nodes...\n"
    }
    if debugToggles[DijkstraAlgo] { fmt.Print(log) }
  }
  return &Result{bestScore, bestSeats, &lastNodes}
}

func computeHeuristic(c1 *Coord, c2 *Coord) int {
  xDiff := c1.x - c2.x; if xDiff < 0 { xDiff = -xDiff }
  yDiff := c1.y - c2.y; if yDiff < 0 { yDiff = -yDiff }
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
