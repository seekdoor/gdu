package web

import (
	"encoding/json"
	"log"
	"time"

	"golang.org/x/net/websocket"
)

type progressMsg struct {
	MsgType   string
	Done      bool
	ItemCount int
	TotalSize int64
}

type commandMsg struct {
	MsgType string
}

func (ui *UI) handleWs(conn *websocket.Conn) {
	jsonE := json.NewEncoder(conn)
	jsonD := json.NewDecoder(conn)

	progress := ui.analyzer.GetProgress()

	for {
		progress.Mutex.Lock()

		if progress.Done {
			jsonE.Encode(&progressMsg{
				MsgType:   "progress",
				Done:      progress.Done,
				ItemCount: progress.ItemCount,
				TotalSize: progress.TotalSize,
			})
			progress.Mutex.Unlock()
			break
		}

		jsonE.Encode(&progressMsg{
			MsgType:   "progress",
			Done:      progress.Done,
			ItemCount: progress.ItemCount,
			TotalSize: progress.TotalSize,
		})

		progress.Mutex.Unlock()

		time.Sleep(100 * time.Millisecond)
	}

	for {
		msg := &commandMsg{}
		err := jsonD.Decode(msg)
		if err != nil {
			log.Printf("Error parsing websocket message: %s", err.Error())
		}

		println(msg.MsgType)
	}
}
