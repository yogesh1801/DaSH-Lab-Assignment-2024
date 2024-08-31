package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/jpoz/groq"
)

type Server struct {
	connections  []*websocket.Conn
	connectionmu sync.Mutex
	messages     chan string
	status       map[string]bool
	statusmu     sync.Mutex
	groqclient   *groq.Client
}

func NewServer(client *groq.Client) *Server {
	return &Server{
		connections: make([]*websocket.Conn, 0),
		messages:    make(chan string, 20),
		status:      make(map[string]bool),
		groqclient:  client,
	}
}

func (s *Server) handleConnection(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading the Http connection: ", err)
		return
	}

	s.connectionmu.Lock()
	s.connections = append(s.connections, conn)
	s.connectionmu.Unlock()
	log.Println("New websocket client connected")
	s.handleMessage(conn)
}

func (s *Server) handleMessage(conn *websocket.Conn) {
	for {
		_, mess, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Error reading message from the connection: %v", err)
			conn.Close()
			return
		}

		message := string(mess)
		parts := strings.SplitN(message, ":", 2)

		if len(parts) != 2 {
			log.Printf("Invalid message format: %s", message)
			continue
		}

		clientID := parts[0]
		msgContent := parts[1]

		res, err := groqClient(s.groqclient, msgContent)
		if err != nil {
			log.Printf("Error processing message with groqClient: %v", err)
			continue
		}

		timestamp := time.Now().Format(time.RFC3339)
		messageID := uuid.New().String()

		formattedMessage := fmt.Sprintf("SR@%s@%s@%s@%s@%s@groq", clientID, messageID, timestamp, msgContent, res)
		s.status[messageID] = false
		s.messages <- formattedMessage
	}
}

func (s *Server) handleSaveStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var reqData struct {
		ID string `json:"id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&reqData); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	result := s.saveStatus(reqData.ID)
	response := map[string]interface{}{
		"status": result,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (s *Server) saveStatus(id string) int {
	s.statusmu.Lock()
	check := s.status[id]
	s.statusmu.Unlock()
	if check {
		return 0
	} else {
		s.status[id] = true
		return 1
	}
}

func (s *Server) broadcastMessage() {
	for {
		message := <-s.messages
		s.connectionmu.Lock()
		var activeConnections []*websocket.Conn
		for _, conn := range s.connections {
			err := conn.WriteMessage(websocket.TextMessage, []byte(message))
			if err != nil {
				log.Printf("Websocket broadcast error: %v", err)
				conn.Close()
			} else {
				activeConnections = append(activeConnections, conn)
			}
		}

		s.connections = activeConnections
		s.connectionmu.Unlock()
		time.Sleep(500 * time.Millisecond)
	}
}
