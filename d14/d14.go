package d14

import (
	"aoc2k24/constants"
	"aoc2k24/io"
	"fmt"
	"strconv"
	"strings"
)

const isDebug = false
const isDebugAdd = false
const isDebugRobot = -1

type VelocityMap map[int][2]int
type PositionMap map[string][]int

func (p *PositionMap) add(id, x, y int) {
  if isDebugAdd { fmt.Printf("Attempting to add robot %d at position x %d, y %d\n", id, x, y) }
  key := coordToKey(x, y)
  ids, exists := (*p)[key]
  if !exists {
    if isDebugAdd { fmt.Print("Position not found in map, adding new slice\n") }
    ids = []int{id}
  } else {
    if isDebugAdd { fmt.Print("Position found in map, adding robot ID to slice\n") }
    ids = append(ids, id)
  }
  (*p)[key] = ids
  if isDebugAdd { fmt.Printf("New position array: %+v\n\n", ids) }
}

func (p PositionMap) get(x, y int) []int {
  key := coordToKey(x, y)
  ids, exists := p[key]
  if !exists { return []int{} }
  return ids
}

func Init(ver constants.VersionIndex) {
  lines, err := io.GetLinesFor(constants.Fourteen, ver)
  if (err != nil) {
    panic(fmt.Sprintf("Error loading file for day %d, version %d: %v", constants.Fourteen, ver, err))
  }
  steps := 100
  areaData := strings.Split(lines[0], ",")
  width, _ := strconv.Atoi(areaData[0])
  height, _ := strconv.Atoi(areaData[1])
  lines = lines[1:]
  vels, posits := parseInput(&lines)
  if isDebug { 
    fmt.Print("At start:\n")
    render(posits, width, height) 
  }
  for range steps {
    if isDebug { fmt.Print("\n\n ************** \n\n") }
    posits = move(vels, posits, width, height)
  }
  if isDebug { 
    fmt.Print("At end:\n")
    render(posits, width, height) 
  }
  safetyFactor := computeSafetyFactor(posits, width, height)
  fmt.Printf("Safety Factor (part 1): %d\n\n", safetyFactor)
  if isDebug { fmt.Printf("Position map: %+v", posits) }
}

func computeSafetyFactor(posits *PositionMap, width, height int) int {
  middleX := width / 2
  middleY := height / 2
  q1 := 0
  q2 := 0
  q3 := 0
  q4 := 0
  for key := range *posits {
    x, y := keyToCoord(key)
    if y == middleY || x == middleX { continue }
    ids, _ := (*posits)[key]
    if y < middleY {
      if x < middleX {
        q1 += len(ids)
      } else {
        q2 += len(ids)
      }
    } else {
      if x < middleX {
        q3 += len(ids)
      } else {
        q4 += len(ids)
      }
    }
  }
  return q1 * q2 * q3 * q4
}

func move(vels *VelocityMap, posits *PositionMap, width, height int) *PositionMap {
  newPos := make(PositionMap)
  keys := make([]string, len(*posits))
  i := 0
  for key := range *posits {
    keys[i] = key
    i++
  }
  for _, key := range keys {
    oldX, oldY := keyToCoord(key)
    if isDebug { fmt.Printf("Computing moves for x %d, y %d\n", oldX, oldY) }
    ids, _ := (*posits)[key]
    for _, robotId := range ids {
      deb := isDebug && isDebugRobot == robotId
      vels := (*vels)[robotId]
      velX := vels[0]
      velY := vels[1]
      if deb { fmt.Printf("Computing move robot %d (pos: %d, %d | vel: %d, %d): ", robotId, oldX, oldY, velX, velY) }
      x := oldX + velX
      y := oldY + velY
      if x >= width {
        x = x - width
      } else if x < 0 {
        x = width + x
      }
      if y >= height {
        y = y - height
      } else if y < 0 {
        y = height + y
      }
      if deb { fmt.Printf("New position %d, %d\n", x, y) }
      newPos.add(robotId, x, y)
    }
  }
  return &newPos
}

func parseInput(lines *[]string) (*VelocityMap, *PositionMap) {
  vels := make(VelocityMap)
  posits := make(PositionMap)
  for id, line := range *lines {
    parts := strings.Split(line, " ")
    pos := strings.Split(parts[0][2:], ",")
    vel := strings.Split(parts[1][2:], ",")
    posX, _ := strconv.Atoi(pos[0])
    posY, _ := strconv.Atoi(pos[1])
    velX, _ := strconv.Atoi(vel[0])
    velY, _ := strconv.Atoi(vel[1])
    vels[id] = [2]int{velX, velY}
    posits.add(id, posX, posY)
  }
  return &vels, &posits
}

func coordToKey(x, y int) string {
  return fmt.Sprintf("%d-%d", x, y)
}

func keyToCoord(key string) (int, int) {
  comps := strings.Split(key, "-")
  x, _ := strconv.Atoi(comps[0])
  y, _ := strconv.Atoi(comps[1])
  return x, y
}

func render(posits *PositionMap, width, height int) {
  fmt.Print("\n")
  for y := range height {
    for x := range width {
      key := coordToKey(x, y)
      ids, exists := (*posits)[key]
      if !exists {
        fmt.Print(".")
        continue
      }
      fmt.Print(len(ids))
    }
    fmt.Print("\n")
  }
  fmt.Print("\n\n")
}
