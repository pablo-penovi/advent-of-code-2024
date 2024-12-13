
package d12

import (
	"aoc2k24/constants"
	"aoc2k24/io"
	"fmt"
)

type Land struct {
  height int
  width int
  plots []rune
}

type RegionData struct {
  perimeter int
  area int
}

type PosRegion map[int]int
type RegionDataMap map[int]RegionData

func (pr *PosRegion) add(pos int, rId int) bool {
  _, exists := (*pr)[pos]
  if !exists {
    (*pr)[pos] = rId
    return true
  }
  return false
}

func (pr PosRegion) get(pos int) int {
  rId, exists := pr[pos]
  if !exists { return -1 }
  return rId
}

func (rdm *RegionDataMap) update(rId, perimeter int) {
  data, exists := (*rdm)[rId]
  if !exists {
    (*rdm)[rId] = RegionData{perimeter, 1}
  } else {
    data.perimeter += perimeter
    data.area += 1
    (*rdm)[rId] = data
  }
}

func (l Land) iToCoord(i int) (int, int) {
  return i % l.width, i / l.width
}

func (l Land) coordToI(x, y int) int {
  return y * l.width + x
}

func (l Land) isXOutOfBound(x int) bool {
  return x < 0 || x >= l.width
}

func (l Land) isYOutOfBound(y int) bool {
  return y < 0 || y >= l.height
}

const isDebug = true

func Init(ver constants.VersionIndex) {
  lines, err := io.GetLinesFor(constants.Twelve, ver)
  if (err != nil) {
    panic(fmt.Sprintf("Error loading file for day %d, version %d: %v", constants.Twelve, ver, err))
  }
  land := getLand(lines)
  price := solve(land)
  fmt.Printf("Price for fences (part 1): %d\n", price)
}

func solve(land *Land) int {
  sum := 0
  visited := make(map[int]struct{})
  for i := range len(land.plots) {
    area := 0
    perimeter := 0
    scanRegion(land, i, &visited, &area, &perimeter)
    sum += area * perimeter
  }
  return sum
}

func scanRegion(land *Land, i int, visited *map[int]struct{}, area *int, perimeter *int) {
  _, isVisited := (*visited)[i]
  if isVisited { return }
  (*visited)[i] = struct{}{}
  x, y := land.iToCoord(i)
  *area += 1
  *perimeter += 4
  upY := y - 1
  downY := y + 1
  leftX := x - 1
  rightX := x + 1
  if upY >= 0 {
    newI := land.coordToI(x, upY)
    if land.plots[i] == land.plots[newI] {
      *perimeter--
      scanRegion(land, newI, visited, area, perimeter)
    }
  }
  if downY < land.height {
    newI := land.coordToI(x, downY)
    if land.plots[i] == land.plots[newI] {
      *perimeter--
      scanRegion(land, newI, visited, area, perimeter)
    }
  }
  if leftX >= 0 {
    newI := land.coordToI(leftX, y)
    if land.plots[i] == land.plots[newI] {
      *perimeter--
      scanRegion(land, newI, visited, area, perimeter)
    }
  }
  if rightX < land.width {
    newI := land.coordToI(rightX, y)
    if land.plots[i] == land.plots[newI] {
      *perimeter--
      scanRegion(land, newI, visited, area, perimeter)
    }
  }
}

func getLand(lines []string) *Land {
  l := Land{len(lines), len(lines[0]), make([]rune, len(lines) * len(lines[0]))}
  for y, line := range lines {
    for x, char := range line {
      l.plots[l.coordToI(x, y)] = rune(char)
    }
  }
  return &l
}
