package main

import (
	"container/heap"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"slices"
)

// Chooses the neighbor with the lowest manhattan cost
type AstarSearch struct {
	Frontier PriorityQueueAstar
	Game     *Maze
}

func (d *AstarSearch) GetFrontier() []*Node {
	return d.Frontier
}

func (d *AstarSearch) Add(i *Node) {
	i.CostToGoal = i.ManhattanDistance(d.Game.Start)
	i.EstimatedCostToGoal = eculideanDist(i.State, d.Game.Goal) + float64(i.CostToGoal)
	d.Frontier.Push(i)
	heap.Init(&d.Frontier)
}

func (d *AstarSearch) ContainsState(i *Node) bool {
	for _, x := range d.Frontier {
		if x.State == i.State {
			return true
		}
	}
	return false
}

func (d *AstarSearch) Empty() bool {
	return len(d.Frontier) == 0
}

func (d *AstarSearch) Remove() (*Node, error) {
	if len(d.Frontier) > 0 {
		if d.Game.Debug {
			fmt.Println("Frontier before remove:")
			for _, x := range d.Frontier {
				fmt.Println("Node:", x.State)
			}
		}
		return heap.Pop(&d.Frontier).(*Node), nil
	}
	return nil, errors.New("frontier is empty")
}

func (d *AstarSearch) Solve() {
	fmt.Println("Starting to solve maze using astar search...")
	d.Game.NumExplored = 0

	start := Node{
		State:  d.Game.Start,
		Parent: nil,
		Action: "",
	}

	d.Add(&start)
	d.Game.CurrentNode = &start

	for {
		if d.Empty() {
			return
		}

		currentNode, err := d.Remove()
		if err != nil {
			log.Println(err)
			return
		}

		if d.Game.Debug {
			fmt.Println("Removed", currentNode.State)
			fmt.Println("-------")
			fmt.Println()
		}

		d.Game.CurrentNode = currentNode
		d.Game.NumExplored += 1

		// Have we found the solution?
		if d.Game.Goal == currentNode.State {
			var actions []string
			var cells []Point

			for {
				if currentNode.Parent != nil {
					actions = append(actions, currentNode.Action)
					cells = append(cells, currentNode.State)
					currentNode = currentNode.Parent
				} else {
					break
				}
			}

			slices.Reverse(actions)
			slices.Reverse(cells)

			d.Game.Solution = Solution{
				Actions: actions,
				Cells:   cells,
			}
			d.Game.Explored = append(d.Game.Explored, currentNode.State)
			break
		}

		d.Game.Explored = append(d.Game.Explored, currentNode.State)

		if d.Game.Animate {
			d.Game.OutputImage(fmt.Sprintf("tmp/%06d.png", d.Game.NumExplored))
		}

		for _, x := range d.Neighbors(currentNode) {
			if !d.ContainsState(x) {
				if !inExplored(x.State, d.Game.Explored) {
					d.Add(&Node{
						State:  x.State,
						Parent: currentNode,
						Action: x.Action,
					})
				}
			}
		}
	}
}

func (d *AstarSearch) Neighbors(node *Node) []*Node {
	row := node.State.Row
	col := node.State.Col

	candidates := []*Node{
		{State: Point{Row: row - 1, Col: col}, Parent: node, Action: "up"},
		{State: Point{Row: row + 1, Col: col}, Parent: node, Action: "down"},
		{State: Point{Row: row, Col: col - 1}, Parent: node, Action: "left"},
		{State: Point{Row: row, Col: col + 1}, Parent: node, Action: "right"},
	}

	var neighbors []*Node
	for _, x := range candidates {
		if 0 <= x.State.Row && x.State.Row < d.Game.Height {
			if 0 <= x.State.Col && x.State.Col < d.Game.Width {
				if !d.Game.Walls[x.State.Row][x.State.Col].wall {
					neighbors = append(neighbors, x)
				}
			}
		}
	}

	for i := range neighbors {
		j := rand.Intn(i + 1)
		neighbors[i], neighbors[j] = neighbors[j], neighbors[i]
	}

	return neighbors
}
