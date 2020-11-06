package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

type Nd struct {
	nom string
}

type Lien struct {
	dep   *Nd
	fin   *Nd
	poids int
}

graph := makeGraph()

const Infinity = int(^uint(0) >> 1)

func getGraph(graph map[*Nd][]Lien) {
	for i := range graph {
		fmt.Println(graph[i])
	}
}

func makeGraph() map[*Nd][]Lien {
	graph := make(map[*Nd][]Lien)
	f, err := os.Open("graph.txt")
	defer f.Close()
	if err != nil {
		log.Fatal(err)
	}

	rd := bufio.NewReader(f)

	for {

		line, err := rd.ReadString('\n')
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatal(err)
		}
		tsep := strings.Split(line, ";")
		res1 := &Nd{nom: tsep[0]}
		res2 := &Nd{nom: tsep[1]}
		resq := strings.TrimSuffix(tsep[2], "\r\n")
		res3, err := strconv.Atoi(resq)
		if err == nil {
			//fmt.Println(res3)
		}

		lien := Lien{res1, res2, res3}

		graph[res1] = append(graph[res1], lien)

		//fmt.Println("Depart:", res1, "|", "Fin:", res2, "|", "ID: ", res3, "|", "Distance:")
	}
	return graph

}

func ListeNd(graph map[*Nd][]Lien) []Nd {
	keys := make([]Nd, 0, len(graph))
	for k := range graph {
		keys = append(keys, k)
	}
}

func NewDistTab(NdInit *Nd) map[*Nd]int {
	DistTab := make(map[*Nd]int)
	DistTab[NdInit] = 0

	for _, nd := range ListeNd(graph) {
		if nd != NdInit {
			DistTab[nd] = Infinity
		}
	}

	return DistTab
}

func main() {
	getGraph(makeGraph())

}
