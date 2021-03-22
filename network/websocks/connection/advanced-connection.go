package connection

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/gorilla/websocket"
	"pandora-pay/config"
	"pandora-pay/helpers"
	"sync"
	"time"
)

type AdvancedConnectionMessage struct {
	Answer uint32
	Reply  bool
	Name   []byte
	Data   []byte
}

type AdvancedConnectionAnswer struct {
	Out []byte
	Err error
}

type AdvancedConnection struct {
	Conn          *websocket.Conn
	send          chan *AdvancedConnectionMessage
	answerCounter uint32
	answerMap     map[uint32]chan *AdvancedConnectionAnswer
	Closed        chan struct{}
	getMap        map[string]func(conn *AdvancedConnection, values []byte) interface{}
	sync.RWMutex  `json:"-"`
}

func (c *AdvancedConnection) sendNow(replyBackId uint32, name []byte, data interface{}, await, reply bool) *AdvancedConnectionAnswer {

	if await && replyBackId == 0 {
		c.Lock()
		c.answerCounter += 1
		if c.answerCounter == 0 {
			c.answerCounter = 1
		}
		for c.answerMap[c.answerCounter] != nil {
			c.answerCounter += 1
			if c.answerCounter == 0 {
				c.answerCounter = 1
			}
		}
		replyBackId = c.answerCounter
		c.answerMap[replyBackId] = make(chan *AdvancedConnectionAnswer)
		c.Unlock()
	}

	marshal, _ := json.Marshal(data)

	message := &AdvancedConnectionMessage{
		replyBackId,
		reply,
		name,
		marshal,
	}
	c.send <- message
	if await {
		select {
		case out, ok := <-c.answerMap[replyBackId]:
			if ok == false {
				return &AdvancedConnectionAnswer{Err: errors.New("Timeout - Closed channel")}
			}
			return out
		case <-time.NewTicker(config.WEBSOCKETS_TIMEOUT).C:
			return &AdvancedConnectionAnswer{Err: errors.New("Timeout")}
		}
	}
	return nil
}

func (c *AdvancedConnection) Send(name []byte, data interface{}) {
	c.sendNow(0, name, data, false, false)
}

func (c *AdvancedConnection) SendAwaitAnswer(name []byte, data interface{}) *AdvancedConnectionAnswer {
	return c.sendNow(0, name, data, true, false)
}

func (c *AdvancedConnection) get(message *AdvancedConnectionMessage) (out interface{}, err error) {

	defer func() {
		err = helpers.ConvertRecoverError(recover())
	}()

	route := string(message.Name)
	var callback func(conn *AdvancedConnection, values []byte) interface{}
	if callback = c.getMap[route]; callback != nil {
		out = callback(c, message.Data)
		return
	}

	err = errors.New("Unknown GET request")
	return
}

func (c *AdvancedConnection) ReadPump() {

	defer func() {
		close(c.Closed)
		c.Conn.Close()
	}()

	for {
		c.Conn.SetReadLimit(int64(config.WEBSOCKETS_MAX_READ))

		_, read, err := c.Conn.ReadMessage()
		if err != nil {
			break
		}

		message := new(AdvancedConnectionMessage)
		if err = json.Unmarshal(read, &message); err != nil {
			continue
		}

		if message.Answer == 0 || !message.Reply {

			var out interface{}
			out, err = c.get(message)

			if message.Answer != 0 {
				if err != nil {
					c.sendNow(message.Answer, []byte{0}, out, false, true)
				} else {
					c.sendNow(message.Answer, []byte{1}, out, false, true)
				}
			}

		} else {

			output := &AdvancedConnectionAnswer{}
			if bytes.Equal(message.Name, []byte{1}) {
				output.Out = message.Data
			} else {
				if err = json.Unmarshal(message.Data, &output.Err); err != nil {
					output.Err = errors.New("Error decoding received error")
				}
			}

			c.RLock()
			cn := c.answerMap[message.Answer]
			if cn != nil {
				delete(c.answerMap, message.Answer)
			}
			c.RUnlock()

			if cn != nil {
				cn <- output
			}
		}
	}

}

func (c *AdvancedConnection) WritePump() {
	defer func() {
		c.Conn.Close()
	}()

	var err error
	for {
		select {
		case message, ok := <-c.send:
			if !ok { // Closed the channel.
				return
			}
			if err = c.Conn.SetWriteDeadline(time.Now().Add(config.WEBSOCKETS_TIMEOUT)); err != nil {
				return
			}
			if err = c.Conn.WriteJSON(message); err != nil {
				return
			}
		}
	}

}

func CreateAdvancedConnection(conn *websocket.Conn, getMap map[string]func(conn *AdvancedConnection, values []byte) interface{}) *AdvancedConnection {
	return &AdvancedConnection{
		Conn:          conn,
		send:          make(chan *AdvancedConnectionMessage),
		Closed:        make(chan struct{}),
		answerCounter: 0,
		answerMap:     make(map[uint32]chan *AdvancedConnectionAnswer),
		getMap:        getMap,
	}
}