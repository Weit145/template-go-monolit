package infrastructure

import "github.com/Weit145/template-go-monolit/internal/infrastructure/websocket"

type InfrastructureInterfase interface {
	websocket.HubInfastructure
}

type Infrastructure struct {
	*websocket.Hub
}

func NewInfrastructure() *Infrastructure {
	return &Infrastructure{
		Hub: websocket.NewHub(),
	}
}
