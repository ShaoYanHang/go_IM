package models

import (
	"crypto/ecdsa"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"sync"

	"github.com/gorilla/websocket"
	"gopkg.in/fatih/set.v0"
	"gorm.io/gorm"
)

type Message struct {
	gorm.Model
	FormId   string // 发送者
	TargetId string // 接收者
	Type     string // 群聊 私聊 广播
	Media    string // 文字 图片 音频
	Content  string // 内容
	Pic      string
	Url      string
	Desc     string
	Amount   int // 其他数字统计
}

func (table *Message) TableName() string {
	return "message"
}

type Node struct {
	Conn      *websocket.Conn
	DataQueue chan []byte
	GroupSets set.Interface
}

// 关系
var clientMap map[int64]*Node = make(map[int64]*Node, 0)

var rwLocker sync.RWMutex

// 需要: 发送者ID 接收者ID 消息类型 发送的内容 发送类型
func Chat(writer http.ResponseWriter, request *http.Request) {
	// 1. 获取参数 并 校验token 等合法性
	// token := query.Get("token")

	query := request.URL.Query()
	Id := query.Get("userId")
	userId, _ = strconv.ParseInt(Id, 10, 64)
	msgType := query.Get("type")
	targetId := query.Get("TargetId")
	context := query.Get("context")

	isvalida := true // checktoken

	conn, err := (&websocket.Upgrader{
		// token 校验
		CheckOrigin: func(r *http.Request) bool {
			return isvalida
		},
	}).Upgrade(writer, request, nil)

	if err != nil {
		fmt.Println(err)
		return
	}

	// 2. 获取conn
	node := &Node{
		Conn:      conn,
		DataQueue: make(chan []byte, 50),
		GroupSets: set.New(set.ThreadSafe),
	}

	// 3. 用户关系
	// 4. userid 跟 node绑定 并加锁
	rwLocker.Lock()
	clientMap[userId] = node
	rwLocker.Unlock()
	// 5. 完成发送逻辑
	go sendProc(node)
	// 6. 完成接受逻辑
	go recvProc(node)

}

func sendProc(node *Node) {
	for {
		select {
		case data := <- node.DataQueue:
			err := node.Conn.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}
}

func recvProc(node *Node) {
	for {
		_, data, err := node.Conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			return 
		}
		broadMsg(data)	// 广播
		fmt.Println("[ws] <<<< ", data)
	}
}

var udpsendChan chan []byte = make(chan []byte, 1024)

func broadMsg(data []byte) {
	udpsendChan <- data
}

func init() {
	go udpSendProc()
	go udpRecvProc()
}

// 完成udp数据发送协程
func udpSendProc() {
	conn, err := net.DialUDP("udp", nil, &net.UDPAddr{
		IP: net.IPv4(192,168,3,255),
		Port: 3000,
	})
	defer conn.Close()
	if err != nil {
		fmt.Println(err)
		return 
	}

	for {
		select {
		case data := <- udpsendChan:
			_, err := conn.Write(data)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}
}

// 完成udp数据接收协程
func udpRecvProc() {
	conn, err := net.ListenUDP("udp", &net.UDPAddr{
		IP: net.IPv4zero,
		Port: 3000,
	})
	if err !=nil {
		fmt.Println(err)
	}
	defer conn.Close()
	for {
		var buf [512]byte
		n, err := conn.Read(buf[0:])
		if err != nil {
			fmt.Println(err)
			return
		}
		dispatch(buf[0:n])
	}
}

// 后端调度逻辑处理
func dispatch() {
	
}