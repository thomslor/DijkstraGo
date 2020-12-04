package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
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
	Sommet  Nd
}

//var graph map[Nd][]Lien = makeGraph()

//const Infinity = int(^uint(0) >> 1)

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

/*
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

//Recuperer distance entre 2 noeuds a partir graph
func GetDistance(dep Nd, fin Nd) (distance int) {
	for i := range graph[dep] {
		if graph[dep][i].dep == dep && graph[dep][i].fin == fin {
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

}*/

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

	for {
		fmt.Printf("#DEBUG MAIN Accepting next connection\n")
		conn, errconn := ln.Accept()

		if errconn != nil {
			fmt.Printf("DEBUG MAIN Error when accepting next connection\n")
			panic(errconn)

		}

		//If we're here, we did not panic and conn is a valid handler to the new connection

		go handleConnection(conn, connum)
		connum += 1

	}
}

func handleConnection(connection net.Conn, connum int) {
	//PFR !!!
	defer connection.Close()
	connReader := bufio.NewReader(connection)
	//    if err !=nil{
	//        fmt.Printf("#DEBUG %d handleConnection could not create reader\n", connum)
	//        return
	//    }

	for {
		//on lit la ligne recue du client
		inputLine, err := connReader.ReadString('\n')
		if err != nil {
			fmt.Printf("#DEBUG %d RCV ERROR no panic, just a client\n", connum)
			fmt.Printf("#DEBUG Error :|%s|\n", err.Error())
			break
		}

		//print la ligne recue
		inputLine = strings.TrimSuffix(inputLine, "\n")
		fmt.Printf("#DEBUG %d RCV |%s|\n", connum, inputLine)

		//Stocke la ligne recue

		//Construit un graphe grace a la ligne recue

		//Applique Dijkstra au graphe

		//ENVOI DU GRAPH AU CLIENT
		//print la ligne a envoyer au client

		//envoie les plus courts chemins au client

		//truc du prof pour exemple :
		splitLine := strings.Split(inputLine, " ")
		returnedString := splitLine[len(splitLine)-1]
		fmt.Printf("#DEBUG %d SND |%s|\n", connum, returnedString)

		io.WriteString(connection, fmt.Sprintf("%s\n", returnedString))

	}

}
