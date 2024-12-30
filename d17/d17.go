package d17

import (
	"aoc2k24/constants"
	"aoc2k24/io"
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

var reader = bufio.NewReader(os.Stdin)

type InstFn func(uint8)

var registers = map[rune]int {
  'A': -1, 'B': -1, 'C': -1,
}
var registerDefaults = map[rune]int {
  'A': -1, 'B': -1, 'C': -1,
}
var pointer int = 0

type Output []uint8
var output = Output{}

var instructions = map[uint8]InstFn {
  // adv - 0 = division between A reg and 2^combo
  0: func(combo uint8) {
    registers['A'] /= pow2(getComboVal(combo))
  },
  // bxl - 1 = bitwise XOR of B reg and literal
  1: func(literal uint8) {
    registers['B'] ^= int(literal)
  },
  // bst - 2 = combo operand modulo 8
  2: func(combo uint8) {
    registers['B'] = getComboVal(combo) % 8
  },
  // jnz - 3 = jump to literal index in program if register A > 0
  3: func(literal uint8) {
    if registers['A'] > 0 {
      // Subtract 2 from literal to cancel out the +2 applied to pointer at the end of each cycle
      pointer = int(literal) - 2
    }
  },
  // bxc - 4 = bitwise XOR of B reg and C reg
  4: func(_ uint8) {
    registers['B'] ^= registers['C']
  },
  // out - 5 = output combo modulo 8
  5: func(combo uint8) {
    output = append(output, uint8(getComboVal(combo) % 8))
  },
  // bdv - 6 = exactly like 0 but storing result in B reg
  6: func(combo uint8) {
    registers['B'] = registers['A'] / pow2(getComboVal(combo))
  },
  // cdv - 7 = exactly like 0 but storing result in C reg
  7: func(combo uint8) {
    registers['C'] = registers['A'] / pow2(getComboVal(combo))
  },
}

func analyze(program *[]uint8) {
  for pointer < len(*program) - 1 {
    opcode := (*program)[pointer]
    operand := (*program)[pointer + 1]
    instructions[opcode](operand)
    pointer += 2
  }
}

func Init(ver constants.VersionIndex) {
  lines, err := io.GetLinesFor(constants.Seventeen, ver)
  if (err != nil) {
    panic(fmt.Sprintf("Error loading file for day %d, version %d: %v", constants.Seventeen, ver, err))
  }
  program := parseInput(&lines)
  exploratoryMode := true
  expMin := int(math.Pow(8, float64(len(*program))))
  expMax := int(math.Pow(8, float64(len(*program) + 1)))

  // Part 1
  if !exploratoryMode {
    fixedValue, _ := registers['A']
    expMin, expMax = fixedValue, fixedValue
    analyze(program)
    fmt.Printf("Output: %+v\n", output)
    return
  }

  // Part 2
  // First I printed the resulting output by brute force from an initial value of Register A = 0, increasing by 1 on each loop
  // I let it run for a minute until I got to 3 or 4 digits output
  // That output allowed me to spot a pattern related to the Reg A values at which each new digit appeared in the output
  // Extrapolating that pattern, I came up with the min and max values for Reg A which would produce a 16-digit output (the program is a 16-digit array as well)
  // After doing that, I let the program generate and print a small portion of those 16-digit outputs
  // Analyzing those outputs allowed me to spot another pattern related to how long it took for the number to change initially on each output position (the five map),
  // as well as another pattern related to how long it took the digit to change in each position after those initial numbers changed (the num change map)
  // Armed with that knowledge, now I can instruct the program to only look for solutions within much narrower Reg A intervals,
  // hopefully taking a gazillion years less than a brute force solution would.
  // I should start from the end since the beginning digits are the ones that change most often
  fiveMap := []int{
    0, 0, 0, 0, 
    0, 0, 0, 0,
    0, 0, 0, 0, 
    0, 0, 0, 0,
  }
  numChangeMap := []int{
    0, 0, 0, 0,
    0, 0, 0, 0,
    0, 0, 0, 0,
    0, 0, 0, 0,
  }
  exp := 1
  for i := range len(fiveMap) {
    fiveMap[i] = int(math.Pow(2, float64(exp)))
    numChangeMap[i] = fiveMap[i] / 2
    exp += 3
  }
  regAValues := []int{}
  curPos := expMax
  curMapPos := len(*program) - 1
  for curPos >= expMin {
    regAValues = append(regAValues, curPos)
    curPos -= numChangeMap[curMapPos]
    if curPos == fiveMap[curMapPos] {
      curMapPos--  
    }
  }
  for _, regAValue := range regAValues {
    kInput := []byte{}
    reader.Read(kInput)
    for _, b := range kInput {
      if rune(b) == 'q' { break }
    }
    output = []uint8{}
    pointer = 0
    registers['A'] = regAValue
    registers['B'], _ = registerDefaults['B']
    registers['C'], _ = registerDefaults['C']
    analyze(program)
    fmt.Printf("%d: %+v\n", regAValue, output)
  }
  fmt.Printf("fiveMap: %+v\n", fiveMap)
  fmt.Printf("numChangeMap: %+v\n", numChangeMap)
}

func parseInput(lines *[]string) (*[]uint8) {
  program := []uint8{}
  for _, line := range *lines {
    if len(line) == 0 { continue }
    if strings.Contains(line, "Register") {
      reg := strings.Split(line, ": ")
      name := reg[0][len(reg[0]) - 1]
      value, _ := strconv.Atoi(reg[1])
      registers[rune(name)] = value
      registerDefaults[rune(name)] = value
      continue
    }
    programValues := strings.Split(strings.Split(line, ": ")[1], ",")
    for _, value := range programValues {
      v, _ := strconv.Atoi(value)
      program = append(program, uint8(v))
    }
  }
  return &program
}

func getComboVal(combo uint8) int {
  switch c := combo; c {
  case 4:
      return registers['A']
  case 5:
      return registers['B']
  case 6:
      return registers['C']
  case 7:
      panic("Illegal use of reserved combo operand 7! Program is invalid")
  default:
      return int(combo)
  }
}

func pow2(x int) int {
  return int(math.Pow(2, float64(x)))
}
