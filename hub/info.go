package hub

import (
	"time"
)

type connectionInfo struct {
    IPAddress string `json:"ip_address"`
}

type eventInfo struct {
    Name string `json:"name"`
    Connections []connectionInfo `json:"connections"`
}

type infoMessage struct {
    Events []eventInfo `json:"events"`
    CreatedAt time.Time `json:"created_at"`
}
