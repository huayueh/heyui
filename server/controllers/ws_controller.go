package controllers

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"net/http"
	"sync"
)

var wsMap = sync.Map{}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func reader(conn *websocket.Conn) {
	for {
		// read in a message
		_, p, err := conn.ReadMessage()
		if err != nil {
			fmt.Print(err)
			return
		}
		// print out that message for clarity
		fmt.Println(string(p))

		//Testing by echo client message
		//if err := conn.WriteMessage(messageType, p); err != nil {
		//	fmt.Print(err)
		//	return
		//}
	}
}

func WSconnection(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	vars := mux.Vars(r)
	acct := vars["acct"]
	msg := fmt.Sprintf("Hi %v from server!", acct)

	if con, ok := GetWSConn(acct); ok {
		con.Close()
	}

	ws, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		fmt.Print(err)
	}

	errW := ws.WriteMessage(1, []byte(msg))
	if errW != nil {
		fmt.Print(errW)
	}

	wsMap.Store(acct, ws)
	reader(ws)
}

func GetWSConn(acct string) (*websocket.Conn, bool) {
	if cValue, ok := wsMap.Load(acct); ok {
		if conn, ok := cValue.(*websocket.Conn); ok {
			return conn, ok
		}
	}
	return nil, false
}
