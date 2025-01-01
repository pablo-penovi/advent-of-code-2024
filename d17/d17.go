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
  isPart2 := true
  expMin := int(math.Pow(8, float64(len(*program) - 1)))
  expMax := int(math.Pow(8, float64(len(*program))))

  // Part 1
  if !isPart2 {
    expMin, expMax = registers['A'], registers['A']
    analyze(program)
    fmt.Printf("Output: %+v\n", output)
    return
  }

  // Part 2
  // OK, by increasing the register A initial value by 1 on every loop during brute force exploration, I've found that this resembles a Googol machine, in that
  // there is a fixed ratio between the number of times reg A has to increase to advance digit 0, digit 1, digit 2, etc.
  // It all begins with digit 0. The thing is, digit 0 does not advance a number for every increase of reg A, the number varies randomly each time
  // So the first thing I have to do is run a brute force to get the first, say, 500.000 digit 0 advancements and how many reg A numbers each advancement took
  // What I mean is, maybe when regA is 0, digit 0 is 1. Then when regA is 1, digit 0 is 2. So that took 1 advancement of reg A
  // But then, maybe for when regA is 2, digit 0 is 2 again, meaning there is no advancement, and only when regA is 3 does digit 0 change again. Which means the second advancement took 2 reg A numbers to occur

  // So when I get the first 500.000 advancements of digit 0, I can translate those advancements to their equivalents for each digit following this formula:
  // next advancement for digit x = digit0Advancement * 8^x

  // Knowing that, I can produce outputs only for the reg A values in which each digit changes, and then only keep those outputs that match the input, and recurse for the next digit within the confines of those reg A values
  // I start by exploring digit 15 (the last one). In case it matches the desired output, I keep that interval of reg A values that produce that desired output, and I explore digit 14
  // within the confines of said interval. When I find digit 14 candidates, I explore digit 13 within the confines of those reg A values, and so on and so forth,
  // allowing me to gradually home in on the solution. When I reach a candidate for digit 0, I've solved the problem

  d0Changes := getFirstChanges(500.000, expMin, program)
  result := -1
  narrowDown(15, expMin, d0Changes, expMax, program, &result)
  fmt.Printf("\nSolution found: %d\n", result)
}

func narrowDown(digit int, initial int, d0Changes *[]int, final int, program *[]uint8, result *int) {
  if *result != -1 { return }
  dMult := int(math.Pow(8, float64(digit)))
  regA := initial
  loopCount := 0
  for i := range len(*d0Changes) {
    if regA >= final { break }
    kInput := []byte{}
    reader.Read(kInput)
    for _, b := range kInput {
      if rune(b) == 'q' { break }
    }
    newOutput(regA, program)
    if output[digit] == (*program)[digit] {
      fmt.Print(" [C]\n")
      top := final; if i < len(*d0Changes) - 1 { top = regA + (*d0Changes)[i + 1] * dMult }
      fmt.Printf("Candidate found for digit i %d: %d - %d. Exploring now\n\n", digit, regA, top)
      if digit == 0 {
        fmt.Printf("\n\n*** SOLVED ****\n\n")
        *result = regA
        break
      }
      narrowDown(digit - 1, regA, d0Changes, top, program, result)
    }
    fmt.Print("OK this was a red herring, moving on\n")
    regA += (*d0Changes)[i + 1] * dMult
    loopCount++
  }
  return
}

func getFirstChanges(limit int, expMin int, program *[]uint8) *[]int {
  count := -1
  d0Changes := []int{}
  d0Value := -1
  for i := 0; i < limit; i++ {
    count++
    kInput := []byte{}
    reader.Read(kInput)
    for _, b := range kInput {
      if rune(b) == 'q' { break }
    }
    output = []uint8{}
    pointer = 0
    registers['A'] = expMin + i
    registers['B'], _ = registerDefaults['B']
    registers['C'], _ = registerDefaults['C']
    analyze(program)
    if int(output[0]) != d0Value {
      d0Value = int(output[0])
      if count > 0 {
        d0Changes = append(d0Changes, count)
        count = 0
      }
    }
  }
  return &d0Changes
}

func newOutput(regA int, program *[]uint8) {
  output = Output{}
  pointer = 0
  registers['A'] = regA
  registers['B'], _ = registerDefaults['B']
  registers['C'], _ = registerDefaults['C']
  analyze(program)
  fmt.Printf("Output for %d: %+v", regA, output)
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
