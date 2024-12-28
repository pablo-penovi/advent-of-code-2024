package d17

import (
	"aoc2k24/constants"
	"aoc2k24/io"
	"fmt"
	"math"
	"strconv"
	"strings"
)

type InstFn func(uint8)

var registers = map[rune]int {
  'A': -1, 'B': -1, 'C': -1,
}
var pointer int = 0
var output = ""

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
    if len(output) > 0 { output += "," }
    output += fmt.Sprintf("%d", getComboVal(combo) % 8)
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

func Init(ver constants.VersionIndex) {
  lines, err := io.GetLinesFor(constants.Seventeen, ver)
  if (err != nil) {
    panic(fmt.Sprintf("Error loading file for day %d, version %d: %v", constants.Seventeen, ver, err))
  }
  program := parseInput(&lines)
  for pointer < len(*program) - 1 {
    opcode := (*program)[pointer]
    operand := (*program)[pointer + 1]
    instructions[opcode](operand)
    pointer += 2
  }
  fmt.Print("\nEnd state:\n\n")
  fmt.Printf("Registers - A: %d, B: %d, C: %d\n", registers['A'], registers['B'], registers['C'])
  fmt.Printf("Output: %s\n", output)
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
