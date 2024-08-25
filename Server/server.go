package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	structs "server/m/v2/utils"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

var tickrate = 60
var connections = make(map[*structs.Connection]bool)

func main() {
	fmt.Println("Starting server...")

	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		websocket, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
			return
		}
		log.Println("Websocket Connected!")
		listen(websocket)
	})
	http.ListenAndServe(":8080", nil)
}

func listen(conn *websocket.Conn) {
	connection := &structs.Connection{
		Socket: conn,
		ID:     randInt(1, 1000),
	}

	fmt.Println("Connection ID:", connection.ID)

	connections[connection] = true
	disconnectChan := make(chan *structs.Connection)
	messageChan := make(chan string)

	// Goroutine to handle incoming messages
	go func() {
		for {
			messageType, messageContent, err := conn.ReadMessage()
			if err != nil {
				log.Println(err)
				disconnectChan <- connection
				break
			}

			// Process the message and send a response to the messageChan
			processMessage(messageType, messageContent, connection, messageChan)
		}
	}()

	// Goroutine to handle disconnects and broadcast messages
	go func() {
		for {
			select {
			case <-disconnectChan:
				//removePlayer(connection.ID)
				connections[connection] = false
				delete(connections, connection)
				broadcastMessage(fmt.Sprintf("Player leaved with ID: %d", connection.ID))
				return // Ends the goroutine
			case msg := <-messageChan:
				broadcastMessage(msg)
			}
		}
	}()
}

func broadcastMessage(message string) {
	for connection := range connections {
		if err := connection.Socket.WriteMessage(1, []byte(message)); err != nil {
			log.Println(err)
		}
	}
}

func processMessage(messageType int, messageContent []byte, connection *structs.Connection, messageChan chan string) {
	fmt.Println("Received message:", string(messageContent))
	commandType := strings.Split(string(messageContent), "{")[0]
	command := "{" + strings.SplitN(string(messageContent), "{", 2)[1]

	switch commandType {
	case "Join":
		join(command, connection)
	case "Create":
		create(command, messageType, connection.Socket)
	default:
		fmt.Println("Unknown command type:", commandType)
	}

	sync(messageType, connection)
	messageChan <- "Processed message for connection ID: " + fmt.Sprintf("%d", connection.ID)
}

func sync(messageType int, conn *structs.Connection) {
	fmt.Println("Syncing")

	playerValues := make([]structs.Player, 0, len(players))

	for _, v := range players {
		playerValues = append(playerValues, *v)
	}

	values, _ := json.Marshal(playerValues)

	fmt.Println("Values:", string(values))

	messageResponse := fmt.Sprintf("Sync: %s", string(values))

	for connection := range connections {
		if connections[conn] == false {
			fmt.Println("Connection closed")
			continue
		}

		if err := connection.Socket.WriteMessage(messageType, []byte(messageResponse)); err != nil {
			log.Println(err)
			return
		}
	}
}

// Random int fun
func randInt(min int, max int) int {
	rand.Seed(time.Now().UnixNano())
	return min + rand.Intn(max-min)
}

func join(command string, conn *structs.Connection) {
	var data structs.Player
	json.Unmarshal([]byte(string(command)), &data)
	data.ID = conn.ID
	fmt.Println("Player joined with ID:", data.ID)
	fmt.Println("Player joined with Name:", data.Name)

	data.PositionX = randInt(10, 1910)
	data.PositionY = randInt(10, 1070)

	fmt.Println("Player joined with Position:", data.PositionX, data.PositionY)

	fmt.Println("Player joined: %s", data)

	data.Collider = structs.BoxCollider{Position: structs.Vector2Int{data.PositionX, data.PositionY}, Size: structs.Vector2Int{128, 128}}

	//addPlayer(data)

	// Create an array only with values
	playersValues := make([]structs.Player, 0, len(players))

	for _, v := range players {
		playersValues = append(playersValues, *v)
	}

	dataJson, _ := json.Marshal(playersValues)

	messageResponse := fmt.Sprintf("Join: %s", dataJson)

	connections[conn] = true

	if err := conn.Socket.WriteMessage(1, []byte(messageResponse)); err != nil {
		log.Println(err)
		return
	}
}

func create(command string, messageType int, conn *websocket.Conn) {
	var data structs.Object
	json.Unmarshal([]byte(string(command)), &data)
	fmt.Println("Object created with ID:", data.ID)
	fmt.Println("Object created with Position:", data.PositionX, data.PositionY)
	//addObject(data)

	dataJson, _ := json.Marshal(data)

	messageResponse := fmt.Sprintf("Create: %s", dataJson)

	fmt.Println("Sending message: %s", dataJson)

	if err := conn.WriteMessage(messageType, []byte(messageResponse)); err != nil {
		log.Println(err)
		return
	}
}
