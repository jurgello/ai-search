package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

// Set ids for search types.
const (
	DFS      = iota // Depth first search.
	BFS             // Breadth first search.
	GBFS            // Greedy best first search.
	ASTAR           // A* search.
	DIJKSTRA        // Dijkstra algorithm.
)

// Point is a simple struct to store XY coordinates of a node.
type Point struct {
	Row int
	Col int
}

// Wall is the type used to keep track of potential nodes that
// are walls, and cannot be explored.
type Wall struct {
	State Point
	wall  bool
}

// Maze is the type for our game. It keeps track of all the information we need to complete the
// maze, if possible.
type Maze struct {
	Height int      // How tall is the maze.
	Width  int      // How wide is the maze.
	Start  Point    // The start location.
	Goal   Point    // The end location.
	Walls  [][]Wall // A slice of slices of Wall type; one per row of the maze.
}

// main is the entry point to our application.
func main() {
	// Declare some variables.
	var m Maze
	var maze, searchType string

	// Read command line flags, and set some sensible defaults.
	flag.StringVar(&maze, "file", "maze.txt", "maze file")
	flag.StringVar(&searchType, "search", "dfs", "search type")
	flag.Parse()

	// Load and parse the maze file.
	err := m.Load(maze)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("maze height/width:", m.Height, m.Width)
}

func (g *Maze) Load(fileName string) error {
	f, err := os.Open(fileName)
	if err != nil {
		fmt.Printf("error opening %s: %s\n", fileName, err)
	}
	defer f.Close()

	var fileContents []string

	reader := bufio.NewReader(f)
	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		} else if err != nil {
			return fmt.Errorf("cannot open file %s: %s", fileName, err)
		}
		fileContents = append(fileContents, line)
	}

	foundStart, foundEnd := false, false
	for _, line := range fileContents {
		if strings.Contains(line, "A") {
			foundStart = true
		}
		if strings.Contains(line, "B") {
			foundEnd = true
		}
	}

	if !foundStart {
		return errors.New("starting location not found")
	}

	if !foundEnd {
		return errors.New("ending location not found")
	}

	g.Height = len(fileContents)
	g.Width = len(fileContents[0])

	var rows [][]Wall

	for i, row := range fileContents {
		var cols []Wall

		for j, col := range row {
			curLetter := fmt.Sprintf("%c", col)
			var wall Wall
			switch curLetter {
			case "A":
				g.Start = Point{Row: i, Col: j}
				wall.State.Row = i
				wall.State.Col = j
				wall.wall = false
			case "B":
				g.Goal = Point{Row: i, Col: j}
				wall.State.Row = i
				wall.State.Col = j
				wall.wall = false
			case " ":
				wall.State.Row = i
				wall.State.Col = j
				wall.wall = false
			case "#":
				wall.State.Row = i
				wall.State.Col = j
				wall.wall = true
			default:
				continue
			}
			cols = append(cols, wall)
		}
		rows = append(rows, cols)
	}

	g.Walls = rows
	return nil
}
