package main

import (
	"aoe2analysis/pkg/parser"
	"log"
	"os"
	"path"
)

func main() {

	saveGameFolder := "C:\\Users\\caffe\\Games\\Age of Empires 2 DE\\76561198025300479\\savegame"
	saveGameFile   := "MP Replay v101.101.43210.0 @2020.11.28 220107 (2).aoe2record"
	fullPath := path.Join(saveGameFolder, saveGameFile)

	file, err := os.Open(fullPath)
	if err != nil {
		log.Fatalf("Could not open file %q, error: %v", fullPath, err)
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
}
