package d9

import (
	"aoc2k24/constants"
	"aoc2k24/io"
	"fmt"
	"strconv"
)

const isDebug = false

type Sparse []int

func (s Sparse) Print() {
  fmt.Print("\n")
  for _, block := range s {
    if block == -1 {
      fmt.Print(".")
    } else {
      fmt.Print(block)
    }
  }
  fmt.Print("\n")
}

func Init(ver constants.VersionIndex) {
  lines, err := io.GetLinesFor(constants.Nine, ver)
  if (err != nil) {
    panic(fmt.Sprintf("Error loading file for day %d, version %d: %v", constants.Nine, ver, err))
  }
  // Puzzle entry has only 1 line
  checksum := solvePart1(getSparse(&lines[0]))
  checksum2 := solvePart2(getSparse(&lines[0]))
  fmt.Printf("Checksum Part 1: %d\n", checksum)
  fmt.Printf("Checksum Part 2: %d\n", checksum2)
}

func solvePart2(blocks *Sparse) int {
  checksum := 0
  if isDebug { fmt.Print("\nBefore defragged consolidation of free blocks: "); blocks.Print() }
  defragConsolidateFreeBlocks(blocks)
  if isDebug { fmt.Print("\nAfter defragged consolidation of free blocks: "); blocks.Print() }
  for i, block := range *blocks {
    // Skip empty blocks
    if block == -1 { continue }
    checksum += i * block
  }
  return checksum
}

func defragConsolidateFreeBlocks(blocks *Sparse) {
  for i := len(*blocks) - 1; i >= 0; i-- {
    // Skip free blocks
    if (*blocks)[i] == -1 { continue }
    start := getStartOfFile(blocks, i)
    length := i - start + 1
    if isDebug { fmt.Printf("File of length %d found. Starts at %d and ends at %d\n", length, start, start + length - 1) }
    freeBlocksStart := getNextFreeBlocks(blocks, start, length)
    // If no free blocks of this size at the left of current file, skip file
    if freeBlocksStart == -1 { 
      i = start
      continue 
    }
    if isDebug { fmt.Printf("Free space found starting at %d and ending at %d\n", freeBlocksStart, freeBlocksStart + length - 1) }
    for j := range length {
      (*blocks)[freeBlocksStart + j] = (*blocks)[start + j]
      (*blocks)[start + j] = -1
    }
  }
}

func getNextFreeBlocks(blocks *Sparse, startOfFile int, length int) int {
  for i := range startOfFile {
    if (*blocks)[i] != -1 { continue }
    hasEnoughSize := true
    for j := i + 1; j < i + length; j++ {
      if (*blocks)[j] != -1 {
        hasEnoughSize = false
        break
      }
    }
    if !hasEnoughSize {
      i = i + length - 1
      continue
    }
    return i
  }
  return -1
}

func getStartOfFile(blocks *Sparse, end int) int {
  fileId := (*blocks)[end]
  for i := end; i >= 0; i-- {
    if (*blocks)[i] != fileId { return i + 1 }
  }
  return 0
}

func solvePart1(blocks *Sparse) int {
  if isDebug { fmt.Print("\nBefore consolidation of free blocks: "); blocks.Print() }
  consolidateFreeBlocks(blocks)
  if isDebug { fmt.Print("\nAfter consolidation of free blocks: "); blocks.Print(); fmt.Print("\n") }
  checksum := 0
  for i, block := range *blocks {
    // If we reached the empty slots, process is over
    if block == -1 { break }
    checksum += i * block
  }
  return checksum
}

func consolidateFreeBlocks(blocks *Sparse) {
  for i := len(*blocks) - 1; i >= 0; i-- {
    // Skip free blocks
    if (*blocks)[i] == -1 { continue }
    nextFree := getNextFreeBlock(blocks, i)
    if nextFree == -1 {
      // No more free blocks before current
      return
    }
    (*blocks)[nextFree] = (*blocks)[i]
    (*blocks)[i] = -1
  }
}

func getNextFreeBlock(blocks *Sparse, end int) int {
  for i := 0; i < end; i++ {
    if (*blocks)[i] == -1 {
      return i
    }
  }
  return -1
}

func getSparse(dense *string) *Sparse {
  blocks := Sparse(make([]int, 0))
  isFile := false
  fileId := 0
  if isDebug { fmt.Printf("Dense input: %s\n", *dense) }
  for _, char := range *dense {
    if isDebug { fmt.Printf("Analyzing dense unit '%c'\n", char) }
    isFile = !isFile
    if char == '0' { continue }
    if isDebug { fmt.Printf("Is this file? %v\n", isFile) }
    length, _ := strconv.Atoi(string(char))
    id := fileId; if !isFile { id = -1 }
    if isDebug { fmt.Printf("Adding %d blocks of id %d\n", length, id) }
    for range length {
      blocks = append(blocks, id)
    }
    if isFile { fileId++ }
  }
  return &blocks
}
