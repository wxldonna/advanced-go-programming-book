package main

import (
	"fmt"
)

type Vertex struct {
	visited    bool
	value      string
	neighbours []*Vertex
}

func NewVertex(value string) *Vertex {
	return &Vertex{
		value: value,

		// the two following lines can be deleted, because the will be initialized with their null value
		visited:    false,
		neighbours: nil, // comment 5.
	}
}

func (v *Vertex) connect(vertex ...*Vertex) { // see comment 4.
	v.neighbours = append(v.neighbours, vertex...)
}

type Graph struct{}

func (g *Graph) dfs(vertex *Vertex) {
	if vertex.visited { // see comment 1.
		return // see comment 2.
	}
	vertex.visited = true
	fmt.Println(vertex.value)
	for _, v := range vertex.neighbours { // see comment 3.
		g.dfs(v)
	}
}

func (g *Graph) disconnected(vertices ...*Vertex) {
	for _, v := range vertices {
		g.dfs(v)
	}
}

func main() {
	v1 := NewVertex("A")
	v2 := NewVertex("B")
	v3 := NewVertex("C")
	v4 := NewVertex("D")
	v5 := NewVertex("E")
	g := Graph{}
	v1.connect(v2)
	v2.connect(v4, v5) // see comment 4.
	v3.connect(v4, v5) // see comment 4.
	g.dfs(v1)
}
