package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	entities "github.com/ant0nix/GoCloudCampPart1/pkg"
)

func main() {
	ch := make(chan string)
	pauseCh := make(chan bool)
	sc := bufio.NewReader(os.Stdin)
	playlist := entities.NewPlaylist()
	fmt.Println("Add song. Example: add [id] [duration]")

	go func(ch chan string, pauseCh chan bool) {
		for {
			playlist.Start(ch, pauseCh)

		}
	}(ch, pauseCh)

	for {
		input, err := sc.ReadString('\n')
		if err != nil {
			log.Printf("Error with input:%s", err.Error())
		}
		input = strings.TrimSpace(input)
		ch <- input
	}
}
