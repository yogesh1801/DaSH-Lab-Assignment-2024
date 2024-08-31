package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"time"
)

func main() {
	url := "ws://localhost:8080/ws"
	c := NewClient(url)
	go c.readMessages()

	file, err := os.Open("./input.txt")
	if err != nil {
		log.Fatal("Error opening input file:", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		message := fmt.Sprintf("%s:%s", c.id, scanner.Text())
		c.sendMessage(message)
		time.Sleep(100 * time.Millisecond)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal("Error reading input file:", err)
	}

	select {}
}
