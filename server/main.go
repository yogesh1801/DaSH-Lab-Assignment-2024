package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
	"github.com/jpoz/groq"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func getOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading the env file.")
	}

	groq_api_key := os.Getenv("GROQ_API_KEY")
	client := groq.NewClient(groq.WithAPIKey(groq_api_key))

	s := NewServer(client)
	go s.broadcastMessage()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Server is ready!")
	})
	http.HandleFunc("/ws", s.handleConnection)
	http.HandleFunc("/save", s.handleSaveStatus)

	port := ":8080"
	ip := getOutboundIP()

	fmt.Printf("Server is running on:\n")
	fmt.Printf("- Local:   http://localhost%s\n", port)
	fmt.Printf("- Network: http://%s%s\n", ip.String(), port)

	err = http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatal("Error starting the server: ", err)
	}
}
