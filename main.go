package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
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

type Node struct {
	index  int
	State  Point
	Parent *Node
	Action string
}

type Solution struct {
	Actions []string
	Cells   []Point
}

// Maze is the type for our game. It keeps track of all the information we need to complete the
// maze, if possible.
type Maze struct {
	Height      int      // How tall is the maze.
	Width       int      // How wide is the maze.
	Start       Point    // The start location.
	Goal        Point    // The end location.
	Walls       [][]Wall // A slice of slices of Wall type; one per row of the maze.
	CurrentNode *Node
	Solution    Solution
	Explored    []Point
	Steps       int
	NumExplored int
	Debug       bool
	SearchType  int
	Animate     bool
}

func init() {
	_ = os.Mkdir("./tmp", os.ModePerm)
	emptyTmp()
}

// main is the entry point to our application.
func main() {
	// Declare some variables.
	var m Maze
	var maze, searchType string

	// Read command line flags, and set some sensible defaults.
	flag.StringVar(&maze, "file", "maze.txt", "maze file")
	flag.StringVar(&searchType, "search", "dfs", "search type")
	flag.BoolVar(&m.Debug, "debug", false, "write debugging info")
	flag.BoolVar(&m.Animate, "animate", false, "produce animation")

	flag.Parse()

	// Load and parse the maze file.
	err := m.Load(maze)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	startTime := time.Now()
	switch searchType {
	case "dfs":
		m.SearchType = DFS
		solveDFS(&m)
	case "bfs":
		m.SearchType = BFS
		solveBFS(&m)
	default:
		fmt.Println("Invalid search type.")
		os.Exit(1)
	}
	if len(m.Solution.Actions) > 0 {
		fmt.Println("Solution:")
		//m.printMaze()

		fmt.Println("Solutions is", len(m.Solution.Cells), "steps")
		fmt.Println("Time to solve:", time.Since(startTime))
		m.OutputImage("image.png")
	} else {
		fmt.Println("No solution.")
	}
	fmt.Println("Explored", len(m.Explored), "nodes.")
	if m.Animate {
		fmt.Println("Building animation...")
		m.OutputAnimatedImage()
		fmt.Println("Done!")
	}
}

func solveDFS(m *Maze) {
	var s DepthFirstSearch
	s.Game = m
	fmt.Println("Goal is", s.Game.Goal)
	s.Solve()
}

func solveBFS(m *Maze) {
	var s BreadthFirstSearch
	s.Game = m
	fmt.Println("Goal is", s.Game.Goal)
	s.Solve()
}

func (g *Maze) printMaze() {
	for r, row := range g.Walls {
		for c, col := range row {
			if col.wall {
				fmt.Print("█")
			} else if g.Start.Row == col.State.Row && g.Start.Col == col.State.Col {
				fmt.Print("A")
			} else if g.Goal.Row == col.State.Row && g.Goal.Col == col.State.Col {
				fmt.Print("B")
			} else if g.inSolution(Point{r, c}) {
				fmt.Print("*")
			} else {
				fmt.Print(" ")
			}
		}
		fmt.Println()
	}
}

func (g *Maze) inSolution(x Point) bool {
	for _, step := range g.Solution.Cells {
		if step.Row == x.Row && step.Col == x.Col {
			return true
		}
	}
	return false
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
