package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
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
	Job     bool
	idGraph int
	graph   map[Nd][]Lien
	Sommet  Nd
}

type ResWorker struct {
	res string
	id  int
}

const Infinity = int(^uint(0) >> 1)

func getArgs() int {
	//On verifie qu'on a bien 2 arguments
	if len(os.Args) != 2 {
		fmt.Printf("Usage: go run server.go <portnumber>\n")
		os.Exit(1)
	} else {
		//on verifie que l'argument est bien un int
		fmt.Printf("#DEBUG ARGS Port Number : %s\n", os.Args[1])
		portNumber, err := strconv.Atoi(os.Args[1])
		if err != nil {
			fmt.Printf("Usage: go run server.go <portnumber>\n")
			os.Exit(1)
		} else {
			return portNumber
		}

	}
	//PFR should never be reached
	return -1
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
func NewDistTab(NdInit Nd, graph map[Nd][]Lien) map[Nd]int {
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
func getBestNonVisitedNode(distTab map[Nd]int, visited []Nd) Nd {
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

/*
//Recuperer distance entre 2 noeuds a partir graph
func GetDistance(dep Nd, fin Nd, graph map[Nd][]Lien) (distance int) {
	for i := range graph[dep] {
		if graph[dep][i].dep == dep && graph[dep][i].fin == fin {
			distance = graph[dep][i].poids
		}
	}
	return distance
}
*/
//ALGORITHME DE DIJKSTRA : la fonction renvoie le chemin le plus court du noeud source a tous les autres noeuds
func Dijkstra(initNd Nd, graph map[Nd][]Lien) (plusCourtChemin string) {

	//Creation du tableau de distances
	distTab := NewDistTab(initNd, graph)
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
	//fmt.Println(plusCourtChemin)
	return plusCourtChemin
}

//Fonction qui permet d'afficher le graph dans le Terminal, très utile pour débug
/*
func getGraph(graph map[Nd][]Lien) {
	for i := range graph {
		fmt.Println(graph[i])
	}
}
*/

func worker(work chan GraphSommet, results chan ResWorker) {
	for f := range work {
		if f.Job {
			//sortie[f.idGraph] += Dijkstra(f.Sommet, f.graph)
			//fmt.Println(sortie)
			results <- ResWorker{Dijkstra(f.Sommet, f.graph), f.idGraph}
		}
	}
}

func main() {
	port := getArgs()
	fmt.Printf("#DEBUG MAIN Creating TCP Server on port %d\n", port)
	//Create a port string that lets us accept connection on all interfaces of the host
	portString := fmt.Sprintf(":%s", strconv.Itoa(port))
	fmt.Printf("#DEBUG MAIN PORT STRING |%s|\n", portString)

	ln, err := net.Listen("tcp", portString)
	if err != nil {
		fmt.Printf("#DEBUG MAIN Could not create listener\n")
		panic(err)
	}

	//If we're here, we did not panic and ln is a valid listener
	//connum = id de connexion
	connum := 1

	//création des chan pour faire communiquer worker et main
	//jobs : channel des datas
	//results : channel de wait group (s'assurer que toutes les go routines ont fini)
	jobs := make(chan GraphSommet, 1000)
	results := make(chan ResWorker, 1000)

	//Initialisation des Workers
	for i := 1; i <= 5; i++ {
		go worker(jobs, results)
	}

	for {
		fmt.Printf("#DEBUG MAIN Accepting next connection\n")
		conn, errconn := ln.Accept()

		if errconn != nil {
			fmt.Printf("DEBUG MAIN Error when accepting next connection\n")
			panic(errconn)

		}

		//If we're here, we did not panic and conn is a valid handler to the new connection

		go handleConnection(conn, connum, jobs, results)
		connum += 1

	}
}

func handleConnection(connection net.Conn, connum int, jobs chan GraphSommet, results chan ResWorker) {
	//PFR !!!
	defer connection.Close()
	connReader := bufio.NewReader(connection)
	graph := make(map[Nd][]Lien)
	nbLigneS, _ := connReader.ReadString('\n')
	nbLigneS = strings.TrimSuffix(nbLigneS, "\n")
	nbLigne, errC := strconv.Atoi(nbLigneS)
	if errC != nil {
		fmt.Println(errC)
	}
	c := 0
	fmt.Println(nbLigne)

	for c < nbLigne {

		//on lit la ligne recue du client
		inputLine, err := connReader.ReadString('\n')

		if err != nil {
			fmt.Printf("#DEBUG %d RCV ERROR no panic, just a client\n", connum)
			fmt.Printf("#DEBUG Error :|%s|\n", err.Error())
			break
		}

		//print la ligne recue
		inputLine = strings.TrimSuffix(inputLine, "\n")
		/*
			if inputLine == "EOF" {
				break
			}
		*/
		fmt.Printf("#DEBUG %d RCV |%s|\n", connum, inputLine)

		//convertir le res3 de string vers int
		//Stocke la ligne recue
		//Construit un graphe grace a la ligne recue

		tsep := strings.Split(inputLine, ";")
		res1 := Nd{nom: tsep[0]}
		res2 := Nd{nom: tsep[1]}
		resq := strings.TrimSuffix(tsep[2], "\n")
		res3, erreur := strconv.Atoi(resq)
		if erreur != nil {
			fmt.Println(erreur)
		}

		lien := Lien{res1, res2, res3}

		graph[res1] = append(graph[res1], lien)

		c++

	}

	//nombre de sommets du graphe :
	nbSommets := len(ListeNd(graph))

	//liste des sommets du graphe :
	listSommet := ListeNd(graph)

	//liste des sommets avec l'id du graph (qui est l'id du client) et un boolean Job :
	listGraphSommet := make([]GraphSommet, 0, nbSommets)

	//pour chaque sommet du graphe, on definit un GraphSommet qui a un Job=true, l'id de connexion du client et la liste de sommets du graphe
	//puis on remplit notre slice listGraphSommet
	for sommet := range ListeNd(graph) {
		f := GraphSommet{true, connum, graph, listSommet[sommet]}
		listGraphSommet = append(listGraphSommet, f)
	}

	//permet d'envoyer les noeuds sur lesquels on va appliquer dijkstra via le channel jobs
	for j := 0; j < nbSommets; j++ {
		jobs <- listGraphSommet[j]
	}

	returnString := fmt.Sprintf("%d\n", connum)
	compteur := 0

	//permet de synchroniser les go routines
	//Principe : On attend que toutes les goroutines qui travaillent sur notre graph finissent leur travail et, pour chaque résultat transmis, on garde le résultat en mémoire pour le client
	for compteur != nbSommets {
		t := <-results
		if t.id == connum {
			returnString += t.res
			compteur += 1

		} else {
			results <- t
		}
	}

	//fmt.Println(sortie)
	//fmt.Println(sortie[connum])

	returnString += "$"
	fmt.Println(returnString)
	io.WriteString(connection, fmt.Sprintf("%s\n", returnString))

}
