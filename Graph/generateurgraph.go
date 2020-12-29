//GENERATION DE GRAPH PLUS PETIT GRACE A GROS GRAPH
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

func main() {
	f, err := os.Open("graph5000.txt")
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
		tsep := strings.Split(line, "\t")
		res1, _ := strconv.Atoi(tsep[0])
		res2, _ := strconv.Atoi(strings.TrimSuffix(tsep[1], "\n"))

		if res1 < 10 {
			if res2 < 10 {
				fmt.Printf("%d;%d;%d\n", res1, res2, rand.Intn(100))
			}
		}

	}
}
