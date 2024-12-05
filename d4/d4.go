package d4

import (
	"aoc2k24/constants"
	"fmt"
	"aoc2k24/io"
)

const isDebug = false

type Direction int

const (
  U Direction = iota
  UR
  R
  DR
  D
  DL
  L
  UL
)

var dirList = [8]Direction{U, UR, R, DR, D, DL, L, UL}

type Salad struct {
  height int
  width int
  letters []rune
}

func newSalad(lines []string) (*Salad, *map[int]struct{}, *map[int]struct{}) {
  height := len(lines)
  width := len(lines[0])
  salad := Salad{height, width, make([]rune, height * width)}
  xMap := make(map[int]struct{})
  aMap := make(map[int]struct{})
  for y, line := range lines {
    for x, char := range line {
      salad.letters[width * y + x] = rune(char)
      if char == 'X' {
        xMap[width * y + x] = struct{}{}
      } else if char == 'A' {
        aMap[width * y + x] = struct{}{}
      }
    }
  }
  return &salad, &xMap, &aMap
}

func Init(ver constants.VersionIndex) {
  lines, err := io.GetLinesFor(constants.Four, ver)
  if (err != nil) {
    panic(fmt.Sprintf("Error loading file for day %d, version %d: %v", constants.Four, ver, err))
  }
  salad, xMap, aMap := newSalad(lines)
  p1Result := getP1Result(salad, xMap)
  p2Result := getP2Result(salad, aMap)
  fmt.Printf("Part 1 result: %d\n", p1Result)
  fmt.Printf("Part 2 result: %d\n", p2Result)
}

func getP1Result(salad *Salad, xMap *map[int]struct{}) int {
  w := "MAS"
  height := salad.height
  width := salad.width
  count := 0
  for i := range *xMap {
    y := i / width
    x := i % width
    hasSpaceUp := i - width * 3 >= 0
    hasSpaceDown := i + width * 3 < height * width
    hasSpaceRight := i % width < width - 3
    hasSpaceLeft := i % width > 2
    cb := count
    // Look up
    if hasSpaceUp && w == string(salad.letters[i-width]) + string(salad.letters[i-width*2]) + string(salad.letters[i-width*3]) {
      count++
      if isDebug { fmt.Printf("Found UP match for X in x %d, y %d\n", x, y) }
    }
    // Look right
    if hasSpaceRight && w == string(salad.letters[i+1]) + string(salad.letters[i+2]) + string(salad.letters[i+3]) {
      count++
      if isDebug { fmt.Printf("Found RIGHT match for X in x %d, y %d\n", x, y) }
    }
    // Look down
    if hasSpaceDown && w == string(salad.letters[i+width]) + string(salad.letters[i+width*2]) + string(salad.letters[i+width*3]) {
      count++
      if isDebug { fmt.Printf("Found DOWN match for X in x %d, y %d\n", x, y) }
    }
    // Look left
    if hasSpaceLeft && w == string(salad.letters[i-1]) + string(salad.letters[i-2]) + string(salad.letters[i-3]) {
      count++
      if isDebug { fmt.Printf("Found LEFT match for X in x %d, y %d\n", x, y) }
    }
    // Look up/right
    if hasSpaceUp && hasSpaceRight && w == string(salad.letters[i-width+1]) + string(salad.letters[i-width*2+2]) + string(salad.letters[i-width*3+3]) {
      count++
      if isDebug { fmt.Printf("Found UP/RIGHT match for X in x %d, y %d\n", x, y) }
    }
    // Look down/right
    if hasSpaceDown && hasSpaceRight && w == string(salad.letters[i+width+1]) + string(salad.letters[i+width*2+2]) + string(salad.letters[i+width*3+3]) {
      count++
      if isDebug { fmt.Printf("Found DOWN/RIGHT match for X in x %d, y %d\n", x, y) }
    }
    // Look down/left
    if hasSpaceDown && hasSpaceLeft && w == string(salad.letters[i+width-1]) + string(salad.letters[i+width*2-2]) + string(salad.letters[i+width*3-3]) {
      count++
      if isDebug { fmt.Printf("Found DOWN/LEFT match for X in x %d, y %d\n", x, y) }
    }
    // Look up/left
    if hasSpaceUp && hasSpaceLeft && w == string(salad.letters[i-width-1]) + string(salad.letters[i-width*2-2]) + string(salad.letters[i-width*3-3]) {
      count++
      if isDebug { fmt.Printf("Found UP/LEFT match for X in x %d, y %d\n", x, y) }
    }
    if isDebug && cb == count {
      fmt.Printf("No match found for X in x %d, y %d\n", x, y)
    }
  }
  return count
}

func getP2Result(salad *Salad, aMap *map[int]struct{}) int {
  height := salad.height
  width := salad.width
  count := 0
  for i := range *aMap {
    cb := count
    x := i % width
    y := i / width
    hasSpace := i - width >= 0 && i + width < height * width && x - 1 >= 0 && x + 1 < width
    if !hasSpace { 
      if isDebug { fmt.Printf("No space for A in x %d, y %d\n", x, y) }
      continue 
    }
    topLeft := i - width - 1
    topRight := i - width + 1
    bottomRight := i + width + 1
    bottomLeft := i + width - 1
    if (salad.letters[topLeft] == 'M' && salad.letters[bottomRight] == 'S' || salad.letters[topLeft] == 'S' && salad.letters[bottomRight] == 'M') &&
      (salad.letters[topRight] == 'M' && salad.letters[bottomLeft] == 'S' || salad.letters[topRight] == 'S' && salad.letters[bottomLeft] == 'M') {
      if isDebug { fmt.Printf("Match found for A in x %d, y %d\n", x, y) }
      count++
    }
    if cb == count {
      if isDebug { fmt.Printf("No match found for A in x %d, y %d\n", x, y) }
    }
  }
  return count
}
