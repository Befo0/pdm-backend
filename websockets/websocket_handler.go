package websockets

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type BroadCastMessage struct {
	FinanzaId uint
	Payload   gin.H
}

var clientsFinanza = make(map[uint]map[*websocket.Conn]bool)
var mensajeBroadcast = make(chan BroadCastMessage)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func handleConnection(c *gin.Context) {

	idParam := c.Param("id")

	idUint, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		log.Println("Ocurrio un error convertir el id")
		return
	}
	idFinanza := uint(idUint)

	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("Ocurrio un error al convertir en websocket")
		return
	}

	defer ws.Close()

	if clientsFinanza[idFinanza] == nil {
		clientsFinanza[idFinanza] = make(map[*websocket.Conn]bool)
	}

	clientsFinanza[idFinanza][ws] = true

	for {
		var msg interface{}
		err := ws.ReadJSON(&msg)
		if err != nil {
			delete(clientsFinanza[idFinanza], ws)

			if len(clientsFinanza[idFinanza]) == 0 {
				delete(clientsFinanza, idFinanza)
			}
			break
		}
	}
}

func handleBroadCast() {
	for {
		msg := <-mensajeBroadcast

		clients := clientsFinanza[msg.FinanzaId]
		for client := range clients {
			if err := client.WriteJSON(msg.Payload); err != nil {
				log.Println("Error enviando mensaje")
				client.Close()
				delete(clients, client)
			}
		}
	}
}
