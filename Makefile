dfs:
	go run . -file=maze2.txt -search=dfs -animate=true

bfs:
	go run . -file=maze2.txt -search=bfs -animate=true

dijkstra:
	go run . -file=maze2.txt -search=dijkstra -animate=true	

gbfs:
	go run . -file=maze2.txt -search=gbfs -animate=true		


astar:
	go run . -file=maze2.txt -search=astar -animate=true		

astar2:
	go run . -file=maze-100-steps.txt -search=astar -animate=true	