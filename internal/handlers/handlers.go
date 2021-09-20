package handlers

import (
	"fmt"
	"net/http"

	"log"

	"github.com/CloudyKit/jet/v6"
	"github.com/gorilla/websocket"
)

var wsChan = make(chan WsPayload)

var clients = make(map[WsConnection]string)

var views = jet.NewSet(
	jet.NewOSFileSystemLoader("./html"),
	jet.InDevelopmentMode(),
)

var upgradeConnection = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func Home(rw http.ResponseWriter, r *http.Request) {
	err := renderPage(rw, "home.jet", nil)
	if err != nil {
		log.Println("Error rendering page : ", err)

	}
}

type WsConnection struct {
	*websocket.Conn
}

//WsJsonResponse defines the response sent back from websocket
type WsJsonResponse struct {
	Action      string `json:"action"`
	Message     string `json:"messsage"`
	MessageType string `json:"messsage_type"`
}

type WsPayload struct {
	Action   string       `json:"action"`
	Message  string       `json:"message"`
	Username string       `json:"username"`
	Conn     WsConnection `json:"-"`
}

//WsEndpoint upgrades connection to websocket.
func WsEndpoint(rw http.ResponseWriter, r *http.Request) {
	ws, err := upgradeConnection.Upgrade(rw, r, nil)
	if err != nil {
		log.Println("Upgrade to websocket error : ", err)
	}

	log.Println("Client connected to Websocket endpoint")

	var response WsJsonResponse
	response.Message = `<em><small>Connected to server</small></em>`
	conn := WsConnection{Conn: ws}
	clients[conn] = ""

	err = ws.WriteJSON(response)
	if err != nil {
		log.Println("Error parsing and sending json response : ", err)
	}

	go ListenForWs(&conn)
}

func ListenForWs(conn *WsConnection) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Error", fmt.Sprintf("%v", r))
		}
	}()

	var payload WsPayload

	err := conn.ReadJSON(&payload)
	if err != nil {
		//do nothing, there is no payload.
	} else {
		payload.Conn = *conn
		wsChan <- payload
	}
}

func ListenToWsChannel() {
	var response WsJsonResponse
	for {
		e := <-wsChan
		response.Action = "Got here"
		response.Message = fmt.Sprintf("Some message, action was %s", e.Action)
		broadcastToAll(response)
	}
}

func broadcastToAll(response WsJsonResponse) {
	for client := range clients {
		err := client.WriteJSON(response)
		if err != nil {
			log.Println("Websocket Error : ", err)
			_ = client.Close()
			delete(clients, client)
		}

	}
}

func renderPage(rw http.ResponseWriter, tmpl string, data jet.VarMap) error {
	view, err := views.GetTemplate(tmpl)
	if err != nil {
		log.Println("Error rendering page : ", err)
		return err
	}
	err = view.Execute(rw, data, nil)
	if err != nil {
		log.Println("Error rendering page : ", err)
		return err
	}
	return nil
}
