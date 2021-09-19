package handlers

import (
	"net/http"

	"log"

	"github.com/CloudyKit/jet/v6"
	"github.com/gorilla/websocket"
)

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

//WsJsonResponse defines the response sent back from websocket
type WsJsonResponse struct {
	Action      string `json:"action"`
	Message     string `json:"messsage"`
	MessageType string `json:"messsage_type"`
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
	err = ws.WriteJSON(response)
	if err != nil {
		log.Println("Error parsing and sending json response : ", err)
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
