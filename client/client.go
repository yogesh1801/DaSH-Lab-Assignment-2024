package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type client struct {
	id   string
	conn *websocket.Conn
}

func NewClient(url string) *client {
	return &client{
		id:   uuid.New().String(),
		conn: wsconnect(url),
	}
}

func wsconnect(url string) *websocket.Conn {
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Fatal("dial: ", err)
	}
	return conn
}

func (c *client) sendMessage(message string) {
	err := c.conn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		log.Println("Failed to send msg to the server: ", err)
		c.conn.Close()
	}
}

func (c *client) readMessages() {
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			log.Printf("Error reading message from the connection: %v", err)
			c.conn.Close()
			return
		}

		TimeRecvd := time.Now().Format(time.RFC3339)

		fmt.Println(string(message))
		tokens := strings.Split(string(message), "@")

		messageid := tokens[2]

		data := map[string]interface{}{
			"ID": messageid,
		}

		postData, _ := json.Marshal(data)
		res, _ := http.Post("http://localhost:8080/save", "application/json", bytes.NewBuffer(postData))

		body, _ := io.ReadAll(res.Body)
		fmt.Println(string(body))

		type Response struct {
			Status int `json:"status"`
		}

		var response Response
		err = json.Unmarshal(body, &response)
		if err != nil {
			log.Fatal("Error unmarshalling response body:", err)
		}

		if response.Status == 1 {
			Prompt := tokens[4]
			Message := tokens[5]
			TimeSent := tokens[3]
			Source := tokens[6]
			if tokens[1] != c.id {
				Source = "user"
			}

			output := map[string]interface{}{
				"Prompt":    Prompt,
				"Message":   Message,
				"TimeSent":  TimeSent,
				"TimeRecvd": TimeRecvd,
				"Source":    Source,
			}

			jsonOutput, err := json.Marshal(output)
			if err != nil {
				log.Printf("Error marshalling output to JSON: %v", err)
				continue
			}

			f, err := os.OpenFile("../output.json", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
			if err != nil {
				log.Printf("Error opening output.json: %v", err)
				continue
			}
			defer f.Close()
			if _, err = f.WriteString("\n"); err != nil {
				log.Printf("Error writing newline to output.json: %v", err)
			}

			if _, err = f.Write(jsonOutput); err != nil {
				log.Printf("Error appending to output.json: %v", err)
			}
		}
	}
}
