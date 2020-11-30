package main

import (
	"aoe2analysis/pkg/parser"
	"flag"
	"log"
	"os"
)

func main() {
	filePath := flag.String("file", "", "required - the file to parse")
	flag.Parse()

	if filePath == nil || *filePath == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	file, err := os.Open(*filePath)
	if err != nil {
		log.Fatalf("Could not open file %q, error: %v", *filePath, err)
	}
	defer func() {
		err := file.Close()
		if err != nil {
			log.Fatal("Failed to close file", err)
		}
	}()

	game, err := parser.Parse(file)
	if err != nil {
		log.Fatalf("Failed to parse game: %v", err)
	}

	log.Printf("Game: %#v", game)
	log.Printf("Game: %+v", game)

	log.Printf("Difficulty: %v", game.Header.DE.Difficulty)
	log.Printf("Victory Type: %v", game.Header.DE.VictoryType)
	log.Printf("Starting Resources: %v", game.Header.DE.StartingResources)
	log.Printf("Starting Age: %v - Ending Age: %v", game.Header.DE.StartingAge, game.Header.DE.EndingAge)

	for i, player := range game.Header.DE.Players {
		log.Printf("Player %v: %v", i, player)
	}

	log.Printf("ReplayHeader: %+v", game.Header.Replay)
}
