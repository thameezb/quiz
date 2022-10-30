package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"time"
)

func getFlags() (string, int) {
	filename := flag.String("filename", "problems.csv", "source quiz csv file - defaults to problems.csv")
	maxTime := flag.Int("maxTime", 30, "set max quiz time (in seconds) - defaults to 30 seconds")
	flag.Parse()
	return *filename, *maxTime
}

func readCSV(filename string) ([]string, []string, error) {
	s, err := os.Open(filename)
	if err != nil {
		return []string{}, []string{}, err
	}
	defer s.Close()

	lines, err := csv.NewReader(s).ReadAll()
	if err != nil {
		return []string{}, []string{}, err
	}

	quest := make([]string, len(lines))
	sol := make([]string, len(lines))

	for i, l := range lines {
		quest[i] = l[0]
		sol[i] = l[1]
	}
	return quest, sol, nil
}

func runGame(quest []string, sol []string, maxTime int) int {
	var (
		a     string
		start string
		score int
	)
	done := make(chan bool)

	runGame := func() {
		for i, v := range quest {
			fmt.Printf("%s? ", v)
			fmt.Scan(&a)
			if a == sol[i] {
				score++
			}
		}
		done <- true
	}

	for {
		fmt.Print("start game (type y to begin)? ")
		fmt.Scan(&start)
		if start == "y" {
			break
		}
	}
	timer := time.NewTimer(time.Duration(maxTime * int(time.Second)))
	go runGame()
	for {
		select {
		case <-timer.C:
			return score
		case <-done:
			return score
		}
	}
}

func main() {
	filename, maxTime := getFlags()

	quest, sol, err := readCSV(filename)
	if err != nil {
		log.Fatalf("error reading file, %s", err)
	}

	score := runGame(quest, sol, maxTime)

	fmt.Printf("you got %d answers correct and %d answers incorrect \n", score, len(quest)-score)
	fmt.Printf("from a total of %d questions \n", len(quest))
}
