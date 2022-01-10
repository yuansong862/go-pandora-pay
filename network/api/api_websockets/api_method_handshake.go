package api_websockets

import (
	"pandora-pay/config"
	"pandora-pay/network/websocks/connection"
)

func (api *APIWebsockets) getHandshake(conn *connection.AdvancedConnection, values []byte) (interface{}, error) {
	return &connection.ConnectionHandshake{config.NAME, config.VERSION, config.NETWORK_SELECTED, config.CONSENSUS, config.NETWORK_ADDRESS_URL_STRING}, nil
}