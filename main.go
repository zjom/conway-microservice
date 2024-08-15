package main

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		p, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var g game
		if err := json.Unmarshal(p, &g); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var next = &game{
			Board: step(g.Board),
		}

		r.Header.Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(next); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})
	http.ListenAndServe(":8080", nil)
}

type Error struct {
	Message string `json:"message"`
}

func (e Error) Error() string {
	return e.Message
}

var (
	Empty    = Error{"empty game"}
	EmptyCol = func(col int) Error {
		return Error{"empty column: " + strconv.Itoa(col)}
	}
	DifferentRowLength = func(row int) Error {
		return Error{"different row length: " + strconv.Itoa(row)}
	}
	InvalidValue = func(row, col int) Error {
		return Error{"invalid value at position [" + strconv.Itoa(row) + "][" + strconv.Itoa(col) + "]"}
	}
)

type game struct {
	Board board `json:"board"`
}

type board [][]int

func (b *board) UnmarshalJSON(data []byte) error {
	// Unmarshal the JSON data into a slice of slices
	var rawSlice [][]int
	if err := json.Unmarshal(data, &rawSlice); err != nil {
		return err
	}

	// Check if the length is greater than 0
	if len(rawSlice) == 0 {
		return Empty
	}

	// Check if all rows have the same length
	rowLength := len(rawSlice[0])
	if rowLength == 0 {
		return EmptyCol(0)
	}

	for i, row := range rawSlice {
		if len(row) != rowLength {
			return DifferentRowLength(i)
		}

		// Check if all values are 0 or 1
		for j, val := range row {
			if val != 0 && val != 1 {
				return InvalidValue(i, j)
			}
		}
	}

	// If all checks pass, assign the raw slice to the Board
	*b = rawSlice
	return nil
}

func step(g board) board {
	next := make([][]int, len(g))
	for i := range next {
		next[i] = make([]int, len(g[i]))
	}

	for i, row := range g {
		for j, c := range row {
			neighbours := getNeighbours(i, j, g)
			aliveNeighbours := countAliveNeighbours(neighbours)
			if c != 1 && aliveNeighbours == 3 {
				// Any dead cell with exactly three live neighbors becomes a live cell, as if by reproduction.
				next[i][j] = 1
				continue
			}

			switch aliveNeighbours {
			// 1. Any live cell with fewer than two live neighbors dies, as if by underpopulation.
			case 0, 1:
				next[i][j] = 0
			// 2. Any live cell with two or three live neighbors lives on to the next generation.
			case 2, 3:
				next[i][j] = 1
			// 3. Any live cell with more than three live neighbors dies, as if by overpopulation.
			default:
				next[i][j] = 0
			}
		}
	}

	return next
}

func getNeighbours(i int, j int, g board) []int {
	row := len(g[0])
	neighbours := make([]int, 0)
	// checks
	if i != 0 {
		// up
		neighbours = append(neighbours, g[i-1][j])
		// up left
		if j != 0 {
			neighbours = append(neighbours, g[i-1][j-1])
		}
		// up right
		if j < row-1 {
			neighbours = append(neighbours, g[i-1][j+1])
		}
	}

	// left
	if j != 0 {
		neighbours = append(neighbours, g[i][j-1])
	}

	// right
	if j < row-1 {
		neighbours = append(neighbours, g[i][j+1])
	}

	if i < len(g)-1 {
		// bot
		neighbours = append(neighbours, g[i+1][j])
		// bot left
		if j != 0 {
			neighbours = append(neighbours, g[i+1][j-1])
		}
		// bot right
		if j < row-1 {
			neighbours = append(neighbours, g[i+1][j+1])
		}
	}
	return neighbours
}

func countAliveNeighbours(neighbours []int) int {
	count := 0
	for _, n := range neighbours {
		if n == 1 {
			count++
		}
	}
	return count
}
