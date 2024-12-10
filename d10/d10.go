package d10

import (
	"aoc2k24/constants"
	"aoc2k24/io"
	"fmt"
	"strconv"
)

const isDebug = false

type Path []int

func (p Path) equals(other Path) bool {
  if len(p) != len(other) { return false }
  for i := range len(p) {
    if p[i] != other[i] { return false }
  }
  return true
}

func (p Path) toString() string {
  result := ""
  for i := range len(p) {
    result += fmt.Sprintf("%d,", p[i])
  }
  return result
}

type Terrain struct {
  tiles []int
  width int
  height int
  heads map[int]struct{}
}

func newTerrain(lines *[]string) *Terrain {
  height := len(*lines)
  width := len((*lines)[0])
  tiles := make([]int, height * width)
  heads := make(map[int]struct{})
  for y, line := range *lines {
    for x, char := range line {
      i := y * width + x
      value, _ := strconv.Atoi(string(char))
      tiles[i] = value
      if value == 0 {
        heads[i] = struct{}{}
      }
    }
  }
  return &Terrain{tiles, width, height, heads}
}

func Init(ver constants.VersionIndex) {
  lines, err := io.GetLinesFor(constants.Ten, ver)
  if (err != nil) {
    panic(fmt.Sprintf("Error loading file for day %d, version %d: %v", constants.Ten, ver, err))
  }
  terrain := newTerrain(&lines)
  scoreSum, uniqueTrailsSum := solve(terrain)
  fmt.Printf("Sum of all trail scores (part 1): %d\n", scoreSum)
  fmt.Printf("Sum of all unique trails (part 2): %d\n", uniqueTrailsSum)
}

func solve(terrain *Terrain) (int, int) {
  scoreSum := 0
  uniqueTrailsSum := 0
  for headIndex := range terrain.heads {
    log := ""
    trailEnds := make(map[int]struct{})
    trailPath := Path(make([]int, 1))
    uniqueTrails := make(map[string]struct{})
    trailPath[0] = headIndex
    exploreTrail(terrain, headIndex, headIndex, &log, &trailEnds, &uniqueTrails, trailPath)
    if isDebug {
      x := headIndex % terrain.width
      y := headIndex / terrain.width
      fmt.Printf("==================== Exploring trail starting at x %d, y %d ====================\n\n", x, y)
      fmt.Print(log)
      fmt.Printf("==================== Trail starting at x %d, y %d EXPLORED! Ends reached: %d, Unique trails: %d ====================\n\n", x, y, len(trailEnds), len(uniqueTrails))
    }
    scoreSum += len(trailEnds)
    uniqueTrailsSum += len(uniqueTrails)
  }
  return scoreSum, uniqueTrailsSum
}

func exploreTrail(terrain *Terrain, head int, currPos int, log *string, ends *map[int]struct{}, unique *map[string]struct{}, path Path) {
  currVal := terrain.tiles[currPos]
  x := currPos % terrain.width
  y := currPos / terrain.width
  left := x - 1
  right := x + 1
  up := y - 1
  down := y + 1
  if left >= 0 {
    newPos := y * terrain.width + left
    if terrain.tiles[newPos] == currVal + 1 {
      *log += fmt.Sprintf("Going left from x %d, y %d to x %d, y %d (value %d).\n", x, y, left, y, currVal + 1)
      path = append(path, newPos)
      if terrain.tiles[newPos] == 9 { 
        *log += "Trail end reached!\n\n"
        _, alreadyMapped := (*ends)[newPos]; if !alreadyMapped { (*ends)[newPos] = struct{}{} }
        pathKey := path.toString()
        _, pathAlreadyRegistered := (*unique)[pathKey]; if !pathAlreadyRegistered { (*unique)[pathKey] = struct{}{} }
      } else {
        exploreTrail(terrain, head, newPos, log, ends, unique, path)
      }
    }
  }
  if right < terrain.width {
    newPos := y * terrain.width + right
    if terrain.tiles[newPos] == currVal + 1 {
      *log += fmt.Sprintf("Going right from x %d, y %d to x %d, y %d (value %d).\n", x, y, right, y, currVal + 1)
      path = append(path, newPos)
      if terrain.tiles[newPos] == 9 {
        *log += "Trail end reached!\n\n"
        _, alreadyMapped := (*ends)[newPos]; if !alreadyMapped { (*ends)[newPos] = struct{}{} }
        pathKey := path.toString()
        _, pathAlreadyRegistered := (*unique)[pathKey]; if !pathAlreadyRegistered { (*unique)[pathKey] = struct{}{} }
      } else {
        exploreTrail(terrain, head, newPos, log, ends, unique, path)
      }
    }
  }
  if up >= 0 {
    newPos := up * terrain.width + x
    if terrain.tiles[newPos] == currVal + 1 {
      *log += fmt.Sprintf("Going up from x %d, y %d to x %d, y %d (value %d).\n", x, y, x, up, currVal + 1)
      path = append(path, newPos)
      if terrain.tiles[newPos] == 9 {
        *log += "Trail end reached!\n\n"
        _, alreadyMapped := (*ends)[newPos]; if !alreadyMapped { (*ends)[newPos] = struct{}{} }
        pathKey := path.toString()
        _, pathAlreadyRegistered := (*unique)[pathKey]; if !pathAlreadyRegistered { (*unique)[pathKey] = struct{}{} }
      } else {
        exploreTrail(terrain, head, newPos, log, ends, unique, path)
      }
    }
  }
  if down < terrain.height {
    newPos := down * terrain.width + x
    if terrain.tiles[newPos] == currVal + 1 {
      *log += fmt.Sprintf("Going down from x %d, y %d to x %d, y %d (value %d).\n", x, y, x, down, currVal + 1)
      path = append(path, newPos)
      if terrain.tiles[newPos] == 9 {
        *log += "Trail end reached!\n\n"
        _, alreadyMapped := (*ends)[newPos]; if !alreadyMapped { (*ends)[newPos] = struct{}{} }
        pathKey := path.toString()
        _, pathAlreadyRegistered := (*unique)[pathKey]; if !pathAlreadyRegistered { (*unique)[pathKey] = struct{}{} }
      } else {
        exploreTrail(terrain, head, newPos, log, ends, unique, path)
      }
    }
  }
}
