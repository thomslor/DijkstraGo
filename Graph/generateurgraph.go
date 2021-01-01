package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
)

func GetArgs() [2]int {
	var output [2]int
	if len(os.Args) != 3 {
		fmt.Printf("Usage: go run generateurgraph.go <nbDeSommetsDuGraph:int> <numeroDuFichierVoulu:int>\n")
		os.Exit(1)
	} else {
		conv, err := strconv.Atoi(os.Args[1])
		conv2, err1 := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Printf("Le nombre de Sommets doit être un entier\n")
			os.Exit(1)
		} else {
			output[0] = conv
		}
		if err1 != nil {
			fmt.Printf("Le numéro de fichier doit être un entier\n")
			os.Exit(1)
		} else {
			output[1] = conv2
		}
	}
	return output
}

func main() {
	tab := GetArgs()
	nbSommets := tab[0]
	numFichier := tab[1]
	f, err := os.Open("graph5000.txt")
	defer f.Close()
	if err != nil {
		log.Fatal(err)
	}
	nomFichier := fmt.Sprintf("graph%d_n°%d.txt", nbSommets, numFichier)

	file, err1 := os.OpenFile(nomFichier, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
	defer file.Close()
	if err1 != nil {
		log.Fatal(err1)
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
		tsep := strings.Split(line, "\t")
		res1, _ := strconv.Atoi(tsep[0])
		res2, _ := strconv.Atoi(strings.TrimSuffix(tsep[1], "\n"))

		if res1 < nbSommets {
			if res2 < nbSommets {
				_, err = file.WriteString(fmt.Sprintf("%d;%d;%d\n", res1, res2, rand.Intn(100)))
				if err != nil {
					panic(err)
				}

			}
		}

	}
	fmt.Printf("Graph disponible dans le répertoire Graph sous le nom : %s", nomFichier)
}
