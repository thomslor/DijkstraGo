package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Nd struct {
	nom string
}

type Lien struct {
	dep   Nd
	fin   Nd
	poids int
}

type GraphSommet struct {
	Job bool
	idGraph int
	Sommet Nd
}

var graph map[Nd][]Lien = makeGraph()

const Infinity = int(^uint(0) >> 1)

//fonction qui permet d'afficher le graphe en entrée
func getGraph(graph map[Nd][]Lien) {
	for i := range graph {
		fmt.Println(graph[i])
	}
}

//fonction qui permet de créer un graphe avec un fichier .txt
func makeGraph() map[Nd][]Lien {
	graph := make(map[Nd][]Lien)
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
		res1 := Nd{nom: tsep[0]}
		res2 := Nd{nom: tsep[1]}
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

//fonction qui nous donne la liste des noeuds du graphe en entrée
func ListeNd(graph map[Nd][]Lien) []Nd {
	keys := make([]Nd, 0, len(graph))
	for k := range graph {
		keys = append(keys, k)
	}
	return keys
}

//fonction qui créé le tableau initial de distances à partir du graphe en entrée
//le noeud source se voit attribuer la valeur 0 et tous les autres noeuds la valeur infinie
func NewDistTab(NdInit Nd) map[Nd]int {
	DistTab := make(map[Nd]int)
	DistTab[NdInit] = 0

	for _, nd := range ListeNd(graph) {
		if nd != NdInit {
			DistTab[nd] = Infinity
		}
	}

	return DistTab
}

//fonction qui donne le noeud non visité avec la plus petite distance
func getBestNonVisitedNode(distTab map[Nd]int, visited []Nd) (Nd) {
	type DistTabATrier struct {
		Node     Nd
		Distance int
	}
	var triOK []DistTabATrier
	//Pour voir si le noeud a deja ete visite
	for nd, distance := range distTab {
		var visiteOK bool
		for _, ndVisiteOK := range visited {
			if nd == ndVisiteOK {
				visiteOK = true
			}
		}
		//Si le noeud n'a pas ete visite, on l'ajoute au slice triOK
		if !visiteOK {
			triOK = append(triOK, DistTabATrier{nd, distance})
		}
	}
	//Pour avoir la plus petite distance il faut trier le slice triOK et prendre la première valeur
	sort.Slice(triOK, func(i, j int) bool {
		return triOK[i].Distance < triOK[j].Distance
	})

	return triOK[0].Node
}

//Recuperer distance entre 2 noeuds a partir graph
func GetDistance(dep Nd, fin Nd) (distance int){
	for i := range graph[dep]{
		if graph[dep][i].dep==dep && graph[dep][i].fin==fin {
			distance = graph[dep][i].poids
		}
	}
	return distance
}

//ALGORITHME DE DIJKSTRA : la fonction renvoie le chemin le plus court du noeud source a tous les autres noeuds
func Djikstra(initNd Nd) (plusCourtChemin string) {


	//Creation du tableau de distances
	distTab := NewDistTab(initNd)
	//fmt.Println(distTab)
	ResTab := make(map[Nd]Lien)

	//Creation d'une liste vide des noeuds visites. Des qu'un noeud est visite, il est ajoute a la liste
	var visiteOK []Nd


	//Creation d'une boucle pour visiter tous les noeuds
	for len(visiteOK) != len(ListeNd(graph)) {

		//On prend le noeud non visité le plus proche a partir de distTab
		nd := getBestNonVisitedNode(distTab, visiteOK)

		//On marque le noeud comme etant visite
		visiteOK = append(visiteOK, nd)



		//On prend les voisins du noeud visite (liste de liens)
		voisins := graph[nd]

		//On calcule les nouvelles distances et met a jour le distTab
		for _, lien := range voisins {
			distanceVoisin := distTab[nd] + lien.poids
			//si distanceVoisin plus petite que la distance dans le distTab pour ce voisin
			if distanceVoisin < distTab[lien.fin] {
				//On met a jour la distTab pour ce voisin
				distTab[lien.fin] = distanceVoisin
				ResTab[lien.fin] = lien

			}
			//Rajouter la condition "2 chemins égaux"
		}

	}
	//for nd, distance := range distTab {
		//plusCourtChemin += fmt.Sprintf("La distance de %s à %s est %d \n", initNd, nd.nom, distance)
	//}
	plusCourtChemin += fmt.Sprintf("Djikstra pour le Sommet %s \n", initNd)
	for dest, lien := range ResTab {
		plusCourtChemin += fmt.Sprintf("%s --> %s, %d \n", lien.dep, dest, lien.poids)
		}

	return plusCourtChemin

}

func worker(id int, work chan GraphSommet, results chan string){
	for f := range work {
		if f.Job {
			fmt.Printf("Go routine %d réalise un Djikstra \n", id)
			fmt.Println(Djikstra(f.Sommet))
			fmt.Printf("Go routine %d a fini le Djikstra \n", id)
			results <- Djikstra(f.Sommet)

		}
	}
}

func main() {
	getGraph(makeGraph())

	nbSommets := len(ListeNd(graph))
	listSommet := ListeNd(graph)
	listGraphSommet := make([]GraphSommet, 0, nbSommets)


	//remplit GraphSommet avec les Sommets d'un graph donné
	for sommet := range ListeNd(graph){
		f := GraphSommet{true, 1, listSommet[sommet]}
		listGraphSommet = append(listGraphSommet, f)
	}

	//création des chan pour faire communiquer worker et main
	//jobs : channel des datas
	//results : channel de wait group (s'assurer que toutes les go routines ont fini)
	jobs := make(chan GraphSommet, nbSommets)
	results := make(chan string, nbSommets)

	//Initialisation des Workers
	for i := 1; i <= 5; i++{
		go worker(i, jobs, results)
	}

	//Envoi de la Data
	for j := 0; j < nbSommets; j++ {
		jobs <- listGraphSommet[j]
	}
	close(jobs) //Fermeture du chan

	//Vide le channel résultats
	for a :=0; a < nbSommets; a++ {
		<-results

	}

}
