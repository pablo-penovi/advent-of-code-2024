package d13

import (
	"aoc2k24/constants"
	"aoc2k24/io"
	"fmt"
	"strconv"
	"strings"
)

const isDebug = false
const isPart2 = true
const part2Multiplier = 10000000000000
const aCost = 3
const bCost = 1

type Button struct {
  xInc int
  yInc int
}

type Machine struct {
  prizeX int
  prizeY int
  buttonA Button
  buttonB Button
  tokensToWin int
}

func (m Machine) toString() string {
  return fmt.Sprintf("Prize: X=%d, Y=%d | Button A: X+%d, Y+%d | Button B: X+%d, Y+%d\n", m.prizeX, m.prizeY, m.buttonA.xInc, m.buttonA.yInc, m.buttonB.xInc, m.buttonB.yInc)
}

func Init(ver constants.VersionIndex) {
  lines, err := io.GetLinesFor(constants.Thirteen, ver)
  if (err != nil) {
    panic(fmt.Sprintf("Error loading file for day %d, version %d: %v", constants.Thirteen, ver, err))
  }
  machines := getMachines(lines)
  solve(machines)
  winnable := 0
  tokens := 0
  for i := range len(*machines) {
    if (*machines)[i].tokensToWin >= 0 {
      winnable++
      tokens += (*machines)[i].tokensToWin
    }
    fmt.Printf("Machine %d: Winnable? %v | Tokens: %d\n", i + 1, (*machines)[i].tokensToWin >= 0, (*machines)[i].tokensToWin)
  }
  fmt.Printf("\nTotal winnable: %d | Total tokens: %d\n", winnable, tokens)
}

func solve(machines *[]Machine) {
  for i, m := range *machines {
    (*machines)[i].tokensToWin = solveMachine(&m)
  }
}

// Solved using Cramer's Rule. Thanks to Grant Riordan (https://dev.to/grantdotdev) for introducing me to this approach
func solveMachine(m *Machine) int {
  // Find determinants
  det := m.buttonA.xInc * m.buttonB.yInc - m.buttonB.xInc * m.buttonA.yInc
  detX := m.prizeX * m.buttonB.yInc - m.prizeY * m.buttonB.xInc
  detY := m.prizeY * m.buttonA.xInc - m.prizeX * m.buttonA.yInc
  // Only solvable is detX / det and detY / det (amount of button presses) are both integers
  isSolvable := detX % det == 0 && detY % det == 0
  if !isSolvable { return - 1 }
  aPresses := detX / det
  bPresses := detY / det
  return aPresses * aCost + bPresses * bCost
}

// Alternative recursive method with memoization used for part 1. Impractical for part 2
func solveRecursive(machines *[]Machine) {
  log := ""
  for i, m := range *machines {
    mem := make(map[string]int)
    (*machines)[i].tokensToWin = solveMachineRecursive(&m, 0, 0, 0, &log, &mem)
    log += "\n"
  }
  if isDebug { fmt.Print(log) }
}

func solveMachineRecursive(m *Machine, x, y, toks int, log *string, mem *map[string]int) int {
  key := fmt.Sprintf("x%d-y%d", x, y)
  memo, isMemo := (*mem)[key]; if isMemo { 
    return memo
  }
  if x > m.prizeX || y > m.prizeY || x == m.prizeX && y < m.prizeY || x < m.prizeX && y == m.prizeY {
    *log += "Target overshot, returning -1\n"
    return -1
  }
  if x == m.prizeX && y == m.prizeY {
    *log += fmt.Sprintf("Target reached! Returning cost of %d\n", toks)
    return toks
  }
  resB := solveMachineRecursive(m, x + m.buttonB.xInc, y + m.buttonB.yInc, toks + bCost, log, mem)
  resA := solveMachineRecursive(m, x + m.buttonA.xInc, y + m.buttonA.yInc, toks + aCost, log, mem)
  res := -1
  if resA > -1 && resB > -1 {
    if resA < resB {
      res = resA
    } else {
      res = resB
    }
  } else if resA > -1 || resB > -1 {
    if resB == -1 {
      res = resA
    } else {
      res = resB
    }
  }
  (*mem)[key] = res
  return res
}

func getMachines(lines []string) *[]Machine {
  machines := make([]Machine, 0)
  machine := Machine{0, 0, Button{0, 0}, Button{0, 0}, -1}
  for _, line := range lines {
    if len(line) == 0 { continue }
    parts := strings.Split(line, ": ")
    if parts[0] == "Button A" {
      comp := strings.Split(parts[1], ", ")
      xInc, _ := strconv.Atoi(comp[0][2:])
      yInc, _ := strconv.Atoi(comp[1][2:])
      machine.buttonA.xInc = xInc
      machine.buttonA.yInc = yInc
    } else if parts[0] == "Button B" {
      comp := strings.Split(parts[1], ", ")
      xInc, _ := strconv.Atoi(comp[0][2:])
      yInc, _ := strconv.Atoi(comp[1][2:])
      machine.buttonB.xInc = xInc
      machine.buttonB.yInc = yInc
    } else {
      comp := strings.Split(parts[1], ", ")
      x, _ := strconv.Atoi(comp[0][2:])
      y, _ := strconv.Atoi(comp[1][2:])
      machine.prizeX = x
      machine.prizeY = y
      if isPart2 {
        machine.prizeX += part2Multiplier
        machine.prizeY += part2Multiplier
      }
      machines = append(machines, machine)
      machine = Machine{0, 0, Button{0, 0}, Button{0, 0}, -1}
    }
  }
  return &machines
}
