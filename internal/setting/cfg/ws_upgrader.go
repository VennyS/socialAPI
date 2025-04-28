package cfg

import (
	"net/http"

	"github.com/gorilla/websocket"
)

type Upgrader interface {
	Upgrade(w http.ResponseWriter, r *http.Request, responseHeader http.Header) (*websocket.Conn, error)
}

type websocketUpgrader struct {
	upgrader *websocket.Upgrader
}

func NewUpgrader(allowedOrigins []string) Upgrader {
	return &websocketUpgrader{
		upgrader: &websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				origin := r.Header.Get("Origin")
				for _, allowedOrigin := range allowedOrigins {
					if origin == allowedOrigin {
						return true
					}
				}
				return false
			},
		},
	}
}

func (wu *websocketUpgrader) Upgrade(w http.ResponseWriter, r *http.Request, responseHeader http.Header) (*websocket.Conn, error) {
	return wu.upgrader.Upgrade(w, r, responseHeader)
}
