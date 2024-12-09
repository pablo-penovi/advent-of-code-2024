package d7

import (
	"aoc2k24/constants"
	"aoc2k24/io"
	"fmt"
	"strconv"
	"strings"
)

const isDebug = false

type Operator int

type OperandNode struct {
  next *OperandNode
  value int
}

type OperandList struct {
  first *OperandNode
  last *OperandNode
  length int
}

func (o *OperandList) append(operand int) {
  node := OperandNode{nil, operand}
  if o.length == 0 {
    o.first = &node
    o.last = &node
  } else {
    o.last.next = &node
    o.last = &node
  }
  o.length++
}

func newOperandList(operands *[]int) *OperandList {
  operandList := OperandList{nil, nil, 0}
  for _, o := range *operands {
    operandList.append(o)
  }
  return &operandList
}

const (
  Addition Operator = iota
  Product
  Concatenation
)

type Operation struct {
  operator Operator
  apply func(int, int) int
}

type Equation struct {
  result int
  operands *OperandList
}

func Init(ver constants.VersionIndex) {
  lines, err := io.GetLinesFor(constants.Seven, ver)
  if (err != nil) {
    panic(fmt.Sprintf("Error loading file for day %d, version %d: %v", constants.Seven, ver, err))
  }
  equations := getEquations(&lines)
  part1Operations := []Operation{
    {Addition, add},
    {Product, multiply},
  }
  part2Operations := []Operation{
    {Addition, add},
    {Product, multiply},
    {Concatenation, concatenate},
  }
  part1Result := solve(equations, &part1Operations)
  part2Result := solve(equations, &part2Operations)
  fmt.Printf("Part 1 result: %d\n", part1Result)
  fmt.Printf("Part 2 result: %d\n", part2Result)
}

func solve(equations *[]Equation, operations *[]Operation) int {
  sum := 0
  for i, equation := range *equations {
    if isDebug { fmt.Printf("Processing equation %d (expected %d)\n", i, equation.result) }
    if isValid(equation.operands.first, operations, equation.result, 0) {
      if isDebug { fmt.Printf("Equation %d is valid\n\n", i) }
      sum += equation.result
      continue
    }
    if isDebug { fmt.Printf("Equation %d is NOT valid\n\n", i) }
  }
  return sum
}

func isValid(operand *OperandNode, operations *[]Operation, expected int, soFar int) bool {
  if isDebug { fmt.Printf("Processing operand %d for operations %+v. Next -> %v | So far: %d\n", operand.value, operations, operand.next, soFar) }
  for _, operation := range *operations {
    if operand.next == nil {
      if operation.apply(soFar, operand.value) == expected { return true }
      continue
    } else {
      if isValid(operand.next, operations, expected, operation.apply(soFar, operand.value)) { return true }
    }
  }
  return false
}

func add(n1 int, n2 int) int {
  return n1 + n2
}

func multiply(n1 int, n2 int) int {
  return n1 * n2
}

func concatenate(n1 int, n2 int) int {
  n1str := strconv.Itoa(n1)
  n2str := strconv.Itoa(n2)
  result, _ := strconv.Atoi(n1str + n2str)
  return result
}

func getEquations(lines *[]string) *[]Equation {
  equations := make([]Equation, len(*lines))
  for i, line := range *lines {
    components := strings.Split(line, ": ")
    result, _ := strconv.Atoi(components[0])
    strOps := strings.Split(components[1], " ")
    operands := make([]int, len(strOps))
    for j, strOp := range strOps {
      op, _ := strconv.Atoi(strOp)
      operands[j] = op
    }
    equations[i] = Equation{result, newOperandList(&operands)}
  }
  return &equations
}

