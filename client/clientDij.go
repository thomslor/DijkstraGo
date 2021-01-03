package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

func getArgs() []string {
	//on créé notre tableau de sortie
	var res []string
	//On vérifie qu'on ait bien 2 arguments
	if len(os.Args) != 3 {
		fmt.Printf("Usage: go run clientDij.go <portnumber> <yourgraph.txt>\n")
		os.Exit(1)
	} else {

		//on vérifie que le 1er argument soit bien un int
		fmt.Printf("#DEBUG ARGS Port Number : %s\n", os.Args[1])
		_, err := strconv.Atoi(os.Args[1])
		if err != nil {
			fmt.Printf("Usage: go run clientDij.go <portnumber> <yourgraph.txt>\n")
			os.Exit(1)
		} else {
			//ajout du port au tableau de sortie
			res = append(res, os.Args[1])
		}
		fmt.Printf("#DEBUG ARGS Graph : %s\n", os.Args[2])

		//ajout du graph au tableau de sortie
		res = append(res, os.Args[2])
	}
	return res
}

func main() {
	//Récupération du numéro de port et du nom du graphe
	res := getArgs()
	port, _ := strconv.Atoi(res[0])
	graphName := res[1]
	fmt.Printf("#DEBUG DIALING TCP Server on port %d\n", port)

	//Création du string du port cible
	portString := fmt.Sprintf("127.0.0.1:%s", strconv.Itoa(port))
	fmt.Printf("Connexion sur le port |%s|\n", portString)

	//Connection
	conn, err := net.Dial("tcp", portString)
	if err != nil {
		//Quitter la connexion si elle ne fonctionne pas
		fmt.Printf("#DEBUG MAIN could not connect\n")
		os.Exit(1)
	} else {
		//Fermeture automatique à la fin de l'exécution
		defer conn.Close()

		//Création d'un reader
		reader := bufio.NewReader(conn)
		envoi := ""
		fmt.Printf("Connexion réussie\n")
		depart := time.Now()

		//Le client lit le graphe texte et l'envoie en format string ligne par ligne
		f, err := os.Open(graphName)
		defer f.Close()
		if err != nil {
			log.Fatal(err)
		}
		//Création d'un scanner
		scanner := bufio.NewScanner(f)
		nbLigne := 0
		for scanner.Scan() {
			envoi += fmt.Sprintf("%s\n", scanner.Text())
			nbLigne++
		}
		//Envoi des données au serveur
		io.WriteString(conn, fmt.Sprintf("%d\n", nbLigne))
		io.WriteString(conn, envoi)

		//Après l'envoi, le client attend une réponse du serveur avec les chemins les plus courts
		//Lecture jusqu'au symbole de fin ($) et stockage de l'information dans resultString
		resultString, err := reader.ReadString('$')
		if err != nil {
			fmt.Printf("Lecture impossible de la réponse du serveur")
			os.Exit(1)
		}
		//On enlève le symbole de fin du resultString
		resultString = strings.TrimSuffix(resultString, "$")

		//Stockage de l'ID du client dans la variable ID
		ID := resultString[0]
		time.Sleep(1000 * time.Millisecond)

		//Stockage des informations reçues dans un fichier texte
		file, err := os.OpenFile(fmt.Sprintf("Dijkstra%d.txt", ID), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
		//CREATE pour créer fichier s'il n'existe pas déjà
		//WR ONLY pour rendre le fichier (dans le programme) accessible en écriture seulement
		//APPEND pour rajouter les informations dans le fichier
		//0600 : permission -rw pour le fichier

		//Fermeture automatique à la fin de l'exécution
		defer file.Close()
		if err != nil {
			panic(err)
		}

		//Ecriture de l'id du client et du graphe dans le fichier texte
		_, err = file.WriteString(fmt.Sprintf("ID CLIENT : %s\n", resultString))
		if err != nil {
			panic(err)
		}

		fmt.Printf("Résultat disponible dans le fichier : Dijkstra%d.txt\n", ID)
		arrivee := time.Now()
		duree := arrivee.Sub(depart)
		fmt.Printf("Ca a pris %v millisecond", duree.Milliseconds())

	}

}
