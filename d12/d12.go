
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

const isDebug = false

func Init(ver constants.VersionIndex) {
  lines, err := io.GetLinesFor(constants.Twelve, ver)
  if (err != nil) {
    panic(fmt.Sprintf("Error loading file for day %d, version %d: %v", constants.Twelve, ver, err))
  }
  land := getLand(lines)
  priceP1, priceP2 := solve(land)
  fmt.Printf("Price for fences (part 1): %d\n", priceP1)
  fmt.Printf("Price for fences (part 2): %d\n", priceP2)
}

func solve(land *Land) (int, int) {
  sumP1 := 0
  sumP2 := 0
  visited := make(map[int]struct{})
  for i := range len(land.plots) {
    area := 0
    perimeter := 0
    corners := 0
    scanRegion(land, i, &visited, &area, &perimeter, &corners)
    sumP1 += area * perimeter
    sumP2 += area * corners
    if isDebug { fmt.Printf("==== END REGION ====\n\n") }
  }
  return sumP1, sumP2
}

func scanRegion(land *Land, i int, visited *map[int]struct{}, area *int, perimeter *int, corners *int) {
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
      scanRegion(land, newI, visited, area, perimeter, corners)
    }
  }
  if downY < land.height {
    newI := land.coordToI(x, downY)
    if land.plots[i] == land.plots[newI] {
      *perimeter--
      scanRegion(land, newI, visited, area, perimeter, corners)
    }
  }
  if leftX >= 0 {
    newI := land.coordToI(leftX, y)
    if land.plots[i] == land.plots[newI] {
      *perimeter--
      scanRegion(land, newI, visited, area, perimeter, corners)
    }
  }
  if rightX < land.width {
    newI := land.coordToI(rightX, y)
    if land.plots[i] == land.plots[newI] {
      *perimeter--
      scanRegion(land, newI, visited, area, perimeter, corners)
    }
  }
  cornerFlags := []bool{
    // Outer corners when against one edge
    upY == -1 && leftX >= 0 && land.plots[i] != land.plots[land.coordToI(leftX, y)],
    upY == -1 && rightX < land.width && land.plots[i] != land.plots[land.coordToI(rightX, y)],
    downY == land.height && leftX >= 0 && land.plots[i] != land.plots[land.coordToI(leftX, y)],
    downY == land.height && rightX < land.width && land.plots[i] != land.plots[land.coordToI(rightX, y)],
    upY >= 0 && leftX == -1 && land.plots[i] != land.plots[land.coordToI(x, upY)],
    downY < land.height && leftX == -1 && land.plots[i] != land.plots[land.coordToI(x, downY)],
    upY >= 0 && rightX == land.width && land.plots[i] != land.plots[land.coordToI(x, upY)],
    downY < land.height && rightX == land.width && land.plots[i] != land.plots[land.coordToI(x, downY)],
    // Outer corners at the corners of the map
    upY == -1 && leftX == -1,
    upY == -1 && rightX == land.width,
    downY == land.height && rightX == land.width,
    downY == land.height && leftX == -1,
    // Outer corners away from edges
    upY >= 0 && leftX >= 0 && land.plots[i] != land.plots[land.coordToI(x, upY)] && land.plots[i] != land.plots[land.coordToI(leftX, y)],
    downY < land.height && leftX >= 0 && land.plots[i] != land.plots[land.coordToI(x, downY)] && land.plots[i] != land.plots[land.coordToI(leftX, y)],
    upY >= 0 && rightX < land.width && land.plots[i] != land.plots[land.coordToI(x, upY)] && land.plots[i] != land.plots[land.coordToI(rightX, y)],
    downY < land.height && rightX < land.width && land.plots[i] != land.plots[land.coordToI(x, downY)] && land.plots[i] != land.plots[land.coordToI(rightX, y)],
    // Inner corners
    upY >= 0 && leftX >= 0 && land.plots[i] == land.plots[land.coordToI(x, upY)] && land.plots[i] == land.plots[land.coordToI(leftX, y)] && land.plots[i] != land.plots[land.coordToI(leftX, upY)],
    upY >= 0 && rightX < land.width && land.plots[i] == land.plots[land.coordToI(x, upY)] && land.plots[i] == land.plots[land.coordToI(rightX, y)] && land.plots[i] != land.plots[land.coordToI(rightX, upY)],
    downY < land.height && rightX < land.width && land.plots[i] == land.plots[land.coordToI(x, downY)] && land.plots[i] == land.plots[land.coordToI(rightX, y)] && land.plots[i] != land.plots[land.coordToI(rightX, downY)],
    downY < land.height && leftX >= 0 && land.plots[i] == land.plots[land.coordToI(x, downY)] && land.plots[i] == land.plots[land.coordToI(leftX, y)] && land.plots[i] != land.plots[land.coordToI(leftX, downY)],
  }
  debugCorners := 0
  for _, isCorner := range cornerFlags {
    if isCorner { *corners++; debugCorners++ }
  }
  if isDebug { fmt.Printf("Tile at x %d, y %d (%c) has %d corners (%d corners so far)\n", x, y, land.plots[i], debugCorners, *corners) }
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
