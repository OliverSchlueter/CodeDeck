package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	_ "github.com/gorilla/websocket"
	"log"
	"net/http"
	"os/exec"
	"strings"
	"time"
)

var paths = map[string]string{
	"vscode":      `D:\Programme\Microsoft VS Code\Code.exe`,
	"calculator":  `calc`,
	"taskmanager": `taskmgr`,
	"editor":      `notepad`,
	"explorer":    `explorer`,
}
var authPassword string
var allowedClients []string

func startProgram(program string, args []string) {
	cmd := exec.Command(program, args...)

	if err := cmd.Run(); err != nil {
		log.Println("Error:", err)
	}
}

func setupWebsocketServer() {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}

	http.HandleFunc("/codedeck", func(writer http.ResponseWriter, request *http.Request) {
		ws, err := upgrader.Upgrade(writer, request, nil)

		if err != nil {
			log.Println(err)
		}

		log.Println("Client connected with address: " + ws.RemoteAddr().String())

		websocketReader(ws)
	})
}

func websocketReader(con *websocket.Conn) {
	for {
		_, msgBytes, err := con.ReadMessage()
		if err != nil {
			log.Println(err)
			break
		}

		msg := string(msgBytes)
		words := strings.Split(msg, " ")

		log.Println("Message from " + con.RemoteAddr().String() + " -> " + msg)

		if len(words) != 2 {
			continue
		}

		switch words[0] {
		case "auth":
			password := words[1]
			if password == authPassword {
				allowedClients = append(allowedClients, con.RemoteAddr().String())
				con.WriteMessage(websocket.TextMessage, []byte("successfully authenticated"))
			} else {
				con.WriteMessage(websocket.TextMessage, []byte("wrong password"))
				time.Sleep(500 * time.Millisecond)
				con.Close()
				continue
			}

		case "run":
			isAllowed := false
			for _, addr := range allowedClients {
				if con.RemoteAddr().String() == addr {
					isAllowed = true
					break
				}
			}

			if !isAllowed {
				con.WriteMessage(websocket.TextMessage, []byte("you're not allowed to run programs"))
				time.Sleep(500 * time.Millisecond)
				con.Close()
				continue
			}

			program := strings.ToLower(words[1])
			path := paths[program]

			startProgram(path, []string{})

		default:
			con.WriteMessage(websocket.TextMessage, []byte("unknown command"))
			time.Sleep(500 * time.Millisecond)
			con.Close()
		}
	}
}

func main() {
	log.Print("Please enter the password: ")
	var pwd string
	fmt.Scanf("%s", &pwd)
	authPassword = pwd
	log.Println("Password is set to: " + authPassword)

	log.Println("Starting websocket server")
	setupWebsocketServer()
	http.ListenAndServe(":1337", nil)
}
