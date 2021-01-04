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

//Création de la constante infinie
const Infinity = int(^uint(0) >> 1)

func getArgs() int {
	//On vérifie qu'on ait bien 2 arguments
	if len(os.Args) != 2 {
		fmt.Printf("Usage: go run serverDij.go <portnumber>\n")
		os.Exit(1)
	} else {
		//on vérifie que le 1er argument soit bien un int et on le retourne
		fmt.Printf("#DEBUG ARGS Port Number : %s\n", os.Args[1])
		portNumber, err := strconv.Atoi(os.Args[1])
		if err != nil {
			fmt.Printf("Usage: go run serverDij.go <portnumber>\n")
			os.Exit(1)
		} else {
			return portNumber
		}

	}
	//PFR should never be reached
	return -1
}

//Fonction qui donne la liste des noeuds du graphe en entrée
func ListeNd(graph map[Nd][]Lien) []Nd {
	keys := make([]Nd, 0, len(graph))
	for k := range graph {
		keys = append(keys, k)
	}
	return keys
}

//Fonction qui créé le tableau initial de distances à partir du graphe en entrée
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

//Fonction qui donne le noeud non visité avec la plus petite distance
func getBestNonVisitedNode(distTab map[Nd]int, visited []Nd) Nd {
	type DistTabATrier struct {
		Node     Nd
		Distance int
	}
	var triOK []DistTabATrier
	//Pour voir si le noeud a déjà été visité
	for nd, distance := range distTab {
		var visiteOK bool
		for _, ndVisiteOK := range visited {
			if nd == ndVisiteOK {
				visiteOK = true
			}
		}
		//Si le noeud n'a pas été visité, on l'ajoute au slice triOK
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

//ALGORITHME DE DIJKSTRA : la fonction renvoie le chemin le plus court du noeud source a tous les autres noeuds
func Dijkstra(initNd Nd, graph map[Nd][]Lien) (plusCourtChemin string) {

	//Création du tableau de distances
	distTab := NewDistTab(initNd, graph)

	//Création du tableau de résultat
	ResTab := make(map[Nd]Lien)

	//Création d'une liste vide des noeuds visités. Dès qu'un noeud est visité, il est ajouté à la liste
	var visiteOK []Nd

	//Création d'une boucle pour visiter tous les noeuds
	for len(visiteOK) != len(ListeNd(graph)) {

		//On prend le noeud non visité le plus proche a partir de distTab
		nd := getBestNonVisitedNode(distTab, visiteOK)

		//On marque le noeud comme etant visité
		visiteOK = append(visiteOK, nd)

		//On prend les voisins du noeud visité (liste de liens)
		voisins := graph[nd]

		//On calcule les nouvelles distances et on met à jour le distTab
		for _, lien := range voisins {
			distanceVoisin := distTab[nd] + lien.poids
			//si distanceVoisin est plus petite que la distance dans le distTab pour ce voisin
			if distanceVoisin < distTab[lien.fin] {
				//On met à jour la distTab pour ce voisin
				distTab[lien.fin] = distanceVoisin
				ResTab[lien.fin] = lien

			}
		}
	}
	plusCourtChemin += fmt.Sprintf("Dijkstra pour le Sommet %s \n", initNd)
	for dest, lien := range ResTab {
		plusCourtChemin += fmt.Sprintf("%s --> %s, %d \n", lien.dep, dest, lien.poids)
	}
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
			results <- ResWorker{Dijkstra(f.Sommet, f.graph), f.idGraph}
		}
	}
}

func main() {
	//Récupération du port
	port := getArgs()
	fmt.Printf("#DEBUG MAIN Creating TCP Server on port %d\n", port)

	//Création d'un string portString nous laissant accepter les connexions sur toutes les interfaces de l'hôte
	portString := fmt.Sprintf(":%s", strconv.Itoa(port))
	fmt.Printf("Connexion sur le port |%s|\n", portString)

	//Création d'un listener
	ln, err := net.Listen("tcp", portString)
	if err != nil {
		fmt.Printf("#DEBUG MAIN Could not create listener\n")
		panic(err)
	}

	//Si nous sommes arrivés ici, alors le listener est valide
	//connum = id de connexion
	connum := 1

	//Création des channels pour faire communiquer worker et main
	//jobs : channel des datas
	//results : channel de wait group (s'assurer que toutes les go routines ont fini) et transmission des résultats
	jobs := make(chan GraphSommet, 1000)
	results := make(chan ResWorker, 1000)

	//Initialisation des Workers
	for i := 1; i <= 20; i++ {
		go worker(jobs, results)
	}

	for {
		fmt.Printf("#DEBUG MAIN Accepting next connection\n")
		conn, errconn := ln.Accept()

		if errconn != nil {
			fmt.Printf("DEBUG MAIN Error when accepting next connection\n")
			panic(errconn)
		}

		//Si nous sommes arrivés ici, alors conn est un handler valide pour la nouvelle connection
		go handleConnection(conn, connum, jobs, results)
		connum += 1
	}
}

//Fonction qui va gérer la connexion et appeler les différentes fonctions de Dijkstra
func handleConnection(connection net.Conn, connum int, jobs chan GraphSommet, results chan ResWorker) {
	//Fermeture automatique de la connection en fin  d'exécution
	defer connection.Close()

	//Création d'un reader
	connReader := bufio.NewReader(connection)

	//Création de différentes variables
	graph := make(map[Nd][]Lien)
	nbLigneS, _ := connReader.ReadString('\n')
	nbLigneS = strings.TrimSuffix(nbLigneS, "\n")
	nbLigne, errC := strconv.Atoi(nbLigneS)
	if errC != nil {
		fmt.Println(errC)
	}
	c := 0

	for c < nbLigne {

		//Lecture de la ligne recue du client
		inputLine, err := connReader.ReadString('\n')
		if err != nil {
			fmt.Printf("#DEBUG %d RCV ERROR no panic, just a client\n", connum)
			fmt.Printf("#DEBUG Error :|%s|\n", err.Error())
			break
		}

		//On print la ligne reçue
		inputLine = strings.TrimSuffix(inputLine, "\n")
		fmt.Printf("Client %d Réception |%s|\n", connum, inputLine)

		//Découpage de la ligne reçue en 3 parties
		tsep := strings.Split(inputLine, ";")

		//Noeud de départ
		res1 := Nd{nom: tsep[0]}

		//Noeud de fin
		res2 := Nd{nom: tsep[1]}

		//Distance
		resq := strings.TrimSuffix(tsep[2], "\n")
		//convertir le res3 de string vers int
		res3, erreur := strconv.Atoi(resq)
		if erreur != nil {
			fmt.Println(erreur)
		}

		//Stockage de la ligne reçue
		lien := Lien{res1, res2, res3}

		//Construction d'un graphe grâce à la ligne reçue
		graph[res1] = append(graph[res1], lien)

		c++

	}

	//Nombre de sommets du graphe
	nbSommets := len(ListeNd(graph))

	//Liste des sommets du graphe
	listSommet := ListeNd(graph)

	//Liste des sommets avec l'id du graph (qui est l'id du client) et un boolean Job
	listGraphSommet := make([]GraphSommet, 0, nbSommets)

	//Pour chaque sommet du graphe, on définit un GraphSommet qui a un Job=true, l'id de connexion du client et la liste de sommets du graphe
	//Puis on remplit notre slice listGraphSommet
	for sommet := range ListeNd(graph) {
		f := GraphSommet{true, connum, graph, listSommet[sommet]}
		listGraphSommet = append(listGraphSommet, f)
	}

	//Permet d'envoyer les noeuds sur lesquels on va appliquer dijkstra via le channel jobs
	for j := 0; j < nbSommets; j++ {
		jobs <- listGraphSommet[j]
	}

	returnString := fmt.Sprintf("%d\n", connum)
	compteur := 0

	//Permet de synchroniser les go routines
	//Principe : On attend que toutes les goroutines qui travaillent sur notre graph finissent leur travail et,
	//pour chaque résultat transmis, on garde le résultat en mémoire pour le client
	for compteur != nbSommets {
		t := <-results
		if t.id == connum {
			returnString += t.res
			compteur += 1

		} else {
			results <- t
		}
	}
	//Ajout d'un symbole de fin au graphe
	returnString += "$"

	//Affichage du graphe à envoyer sur le terminal
	fmt.Println(returnString)

	//Envoi du graphe au client
	io.WriteString(connection, fmt.Sprintf("%s\n", returnString))

}
