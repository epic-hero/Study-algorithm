package main

import (
	"fmt"
	"math"
	"sort"
	"sync"
)

const MaxInt = math.MaxInt32

type Graph struct {
	Graph   [][]int
	Edges   Edges
	GroupId int
}

func newGraph(graph [][]int) *Graph {
	edges := Edges{
		Edge:        make([]Edge, 0),
		VertexMutex: new(sync.Mutex),
	}
	return &Graph{
		Graph:   graph,
		Edges:   edges,
		GroupId: 0,
	}
}

type Edges struct {
	VertexMutex *sync.Mutex
	Edge        []Edge
}

type MST struct {
	Edge []Edge
}

type Edge struct {
	VertexStart  int // Начальная вершина ребра
	VertexFinish int // Конечная вершина ребра
	Weight       int // Вес ребра
	Group        int // Компонент связности
}

func newEdge() Edge {
	return Edge{}
}
func newEdges() Edges {
	return Edges{
		Edge:        make([]Edge, 0),
		VertexMutex: new(sync.Mutex),
	}
}

func (edges Edges) sort() {
	sort.Slice(edges.Edge, func(i, j int) bool {
		return edges.Edge[i].Weight < edges.Edge[j].Weight
	})
}

func (graph *Graph) loadMinEdges() {
	graph.Edges = newEdges()

	wg := new(sync.WaitGroup)
	lenGraph := len(graph.Graph)
	wg.Add(lenGraph)
	for i := 0; i < lenGraph; i++ {
		edge := newEdge()

		go graph.findAndSetMinElem(edge, i, wg)
	}
	wg.Wait()
}

func (graph *Graph) findAndSetMinElem(edge Edge, i int, wg *sync.WaitGroup) {

	defer wg.Done()
	graph.Edges.VertexMutex.Lock()
	defer graph.Edges.VertexMutex.Unlock()

	indexMin, minValue := graph.getMinEdge(graph.Graph[i])

	if minValue != MaxInt {
		edge.VertexStart = i
		edge.VertexFinish = indexMin
		edge.Weight = minValue
		// Изначально каждое ребро лежит в разных Компонентах связности
		edge.Group = graph.GroupId
		graph.GroupId++
		graph.Edges.Edge = append(graph.Edges.Edge, edge)
		graph.Graph[i][indexMin] = 0
	}
}

// Берем у вершины миниммальное ребро и позицию
func (graph *Graph) getMinEdge(array []int) (int, int) {
	var min int = MaxInt
	var i int
	for index, value := range array {
		if value != 0 && min > value {
			min = value
			i = index
		}
	}
	return i, min
}

func (graph *Graph) addVertexInMST(mst *MST) {
	for _, edge := range graph.Edges.Edge {
		mst.joinEdgeMST(edge)
	}
}

func (mst *MST) joinEdgeMST(edgeInput Edge) {
	isLoop, group1, group2 := mst.createLoop(edgeInput)
	if len(mst.Edge) == 0 || (mst.createDuplicate(edgeInput) == false && isLoop == false) {
		mst.joinGroup(group1, group2, edgeInput.Group)
		mst.Edge = append(mst.Edge, edgeInput)
	}
}

func (mst *MST) createDuplicate(edgeInput Edge) bool {
	for _, edgeMST := range mst.Edge {
		if (edgeMST.VertexStart == edgeInput.VertexStart && edgeMST.VertexFinish == edgeInput.VertexFinish) ||
			(edgeMST.VertexStart == edgeInput.VertexFinish && edgeMST.VertexFinish == edgeInput.VertexStart) {
			return true
		}
	}
	return false
}

func (mst *MST) createLoop(edgeInput Edge) (bool, int, int) {
	existInMstVertexStart := false
	existInMstVertexFinish := false
	group1 := -1
	group2 := -1
	for _, edgeMST := range mst.Edge {
		//Если начальная точка edgeInput есть в MST
		if edgeInput.VertexStart == edgeMST.VertexStart || edgeInput.VertexStart == edgeMST.VertexFinish {
			group1 = edgeMST.Group
			existInMstVertexStart = true
		}
		//Если конечная точка edgeInput есть в MST
		if edgeInput.VertexFinish == edgeMST.VertexStart || edgeInput.VertexFinish == edgeMST.VertexFinish {
			group2 = edgeMST.Group
			existInMstVertexFinish = true
		}
	}
	if existInMstVertexStart == true && existInMstVertexFinish ==true && group1 == group2{
		return true, group1, group2
	}
	return false, group1, group2
}

func (mst *MST) joinGroup(groupIdForUpdate int, group2 int, groupIdReference int) {
	for i, edge := range mst.Edge {
		if edge.Group == groupIdForUpdate || edge.Group == groupIdReference || edge.Group == group2 {
			mst.Edge[i].Group = groupIdReference
		}
	}
}

func getGraph() [][]int {
	return [][]int{
		{0, 5, 2, 0, 0, 0, 0, 0, 0},
		{5, 0, 2, 4, 7, 0, 0, 0, 0},
		{2, 2, 0, 3, 0, 0, 9, 0, 0},
		{0, 4, 3, 0, 2, 0, 6, 0, 0},
		{0, 7, 0, 2, 0, 8, 5, 7, 0},
		{0, 0, 0, 0, 8, 0, 0, 3, 4},
		{0, 0, 9, 6, 5, 0, 0, 2, 0},
		{0, 0, 0, 0, 7, 3, 2, 0, 0},
		{0, 0, 0, 0, 0, 4, 0, 0, 0},
	}

	//return [][]int{
	//	{0, 4, 0, 3, 0, 5},
	//	{4, 0, 3, 4, 0, 0},
	//	{0, 3, 0, 2, 0, 0},
	//	{3, 4, 2, 0, 3, 0},
	//	{0, 0, 0, 3, 0, 1},
	//	{5, 0, 0, 0, 1, 0},
	//}
	//return [][]int{
	//	{0, 3, 5, 0, 0, 0},
	//	{3, 0, 8, 0, 12, 4},
	//	{5, 8, 0, 10, 9, 0},
	//	{0, 0, 10, 0, 7, 0},
	//	{0, 12, 9, 7, 0, 15},
	//	{0, 4, 0, 0, 15, 0},
	//}
	//return [][]int{
	//	{0, 15, 1, 9, 0},
	//	{15, 0, 18, 0, 6},
	//	{1, 18, 0, 4, 11},
	//	{9, 0, 4, 0, 23},
	//	{0, 6, 11, 23, 0},
	//}
	//return [][]int{
	//	{0, 8, 5, 0, 0},
	//	{8, 0, 9, 11, 0},
	//	{5, 9, 0, 15, 10},
	//	{0, 11, 15, 0, 7},
	//	{0, 0, 10, 7, 0},
	//}

}

func main() {
	graphInput := getGraph()
	mst := MST{}
	graph := newGraph(graphInput)

	graph.loadMinEdges()

	graph.Edges.sort()

	graph.addVertexInMST(&mst)

	graph.loadMinEdges()

	graph.Edges.sort()

	graph.addVertexInMST(&mst)
	fmt.Println("MST: ",mst.Edge)
}