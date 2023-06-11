package model

import (
	"fmt"
	"github.com/gorilla/websocket"
	"sync"
	"time"
)

var (
	group   = &sync.Map{}
	lock    = &sync.Mutex{}
	clients = map[string]*Client{}
)

const (
	maxLostTime = 12 * time.Second
	heartbeat   = 10 * time.Second
)

const (
	TextMsg = iota
	CommandMsg
	PingMsg
	PongMsg
	ErrorMsg
	CloseMsg
)

type Client struct {
	Name          string
	Room          string
	conn          *websocket.Conn
	inChan        chan Msg
	outChan       chan Msg
	isClosed      bool
	Error         error
	isLost        bool
	CommandChan   chan string
	lastHeartbeat time.Time
}

func (c *Client) String() string {
	return "lastHeartbeat " + fmt.Sprintf("%v\n", c.lastHeartbeat) +
		"Name " + fmt.Sprintf("%v\n", c.Name) +
		"isLost " + fmt.Sprintf("%v\n", c.isLost) +
		"isClosed " + fmt.Sprintf("%v\n", c.isClosed) +
		"Error " + fmt.Sprintf("%v\n", c.Error) +
		"HasHeartbeat " + fmt.Sprintf("%v\n", c.lastHeartbeat.Add(heartbeat).After(time.Now())) +
		"isClosed " + fmt.Sprintf("%v\n", c.lastHeartbeat.Add(maxLostTime).Before(time.Now()))
}

func ShowClients() {
	for {
		time.Sleep(1 * time.Second)
		go group.Range(func(key, value any) bool {
			fmt.Printf("now:%v-->%T\n", key, value)
			_, ok := value.(chan Msg)
			if !ok {
				fmt.Println("err : kind is not [chan Msg]")
			}
			return true
		})
		for _, client := range clients {
			fmt.Println(client)
		}
	}
}
func NewClient(name, room string, conn *websocket.Conn) *Client {
	cli, ok := clients[name+room]
	if !ok {
		cli = &Client{
			Name:        name,
			Room:        room,
			inChan:      make(chan Msg, 10),
			outChan:     make(chan Msg, 10),
			CommandChan: make(chan string, 10),
		}
	}
	cli.conn = conn
	go cli.connect()
	cli.joinRoom()
	cli.joinClients()
	go cli.writeProc()
	go cli.readProc()
	return cli
}

func (c *Client) joinClients() {
	clients[c.Name+c.Room] = c
}

func (c *Client) joinRoom() {
	lock.Lock()
	group.Store(c.Name+c.Room, c.outChan)
	lock.Unlock()
}

func (c *Client) tick() {
	c.lastHeartbeat = time.Now()
	for {
		if c.isClosed || c.isLost || c.Error != nil {
			return
		}
		if time.Now().After(c.lastHeartbeat.Add(heartbeat)) {
			c.isLost = true
			go c.reconnect()
		}
	}
}

func (c *Client) Exit() {
	if c.isClosed {
		return
	}
	if c.Error != nil {
		_ = c.conn.WriteJSON(Msg{
			Type:    ErrorMsg,
			Time:    time.Now(),
			Content: c.Error.Error(),
		})
	} else {
		fmt.Println("*********", c.Name, "exit", c.Room, "***********")
		err := c.conn.WriteJSON(Msg{
			Type:    CloseMsg,
			Time:    time.Now(),
			Content: "[" + c.Name + "] exit the room",
		})
		if err != nil {
			fmt.Println("********" + err.Error() + "***********")
		}
	}
	_ = c.conn.Close()
	group.Delete(c.Name + c.Room)
	c.isClosed = true
}

func (c *Client) reconnect() {
	fmt.Println(c.Name + " reconnecting ...")
	for {
		switch {
		case c.lastHeartbeat.Add(maxLostTime).After(time.Now()):
			c.Exit()
			return
		case c.Error != nil:
			c.Exit()
			return
		case c.isClosed:
			return
		case !c.isLost:
			c.connect()
			go c.tick()
			go c.readProc()
			go c.writeProc()
			return
		}

	}
}

func (c *Client) connect() {
	c.Error = nil
	c.isClosed = false
	c.isLost = false
	c.lastHeartbeat = time.Now()
}

func (c *Client) handleError(err error) bool {
	if err != nil {
		c.Error = err
		c.Exit()
		return true
	}
	return false
}

// 不断地读
func (c *Client) readProc() {
	defer c.Exit()
	for {
		if c.Error != nil || c.isLost || c.isClosed {
			fmt.Println(c.Name, "read stopped")
			return
		}
		var msg Msg
		err := c.conn.ReadJSON(&msg)
		msg.Time = time.Now()
		if c.handleError(err) {
			return
		}
		myLog.Info("msg in:", msg)
		switch msg.Type {
		case TextMsg:
			group.Range(func(key, value any) bool {
				if c.Name+c.Room == key.(string) {
					fmt.Println("same room")
					return true
				}
				fmt.Println(key, fmt.Sprintf("%T", key), "--->", value, fmt.Sprintf("%T", value))
				value.(chan Msg) <- msg
				return true
			})
		case PingMsg:
			c.pong()
		case CloseMsg:
			c.Exit()
			return
		case CommandMsg:
			c.CommandChan <- msg.Content
		default:
		}
	}
}

func (c *Client) pong() {
	if c.Error != nil {
		c.Exit()
		return
	}
	_ = c.conn.WriteJSON(Msg{
		Type:    PongMsg,
		Time:    time.Now(),
		Content: "pong" + c.lastHeartbeat.Format("15:04:05"),
	})
	c.isClosed = false
	c.lastHeartbeat = time.Now()
}

// 不断的写
func (c *Client) writeProc() {
	defer c.Exit()
	for {
		if c.Error != nil || c.isLost || c.isClosed {
			return
		}
		select {
		case msg := <-c.outChan:
			if c.handleError(c.conn.WriteJSON(msg)) {
				return
			}
		default:
		}
	}
}
