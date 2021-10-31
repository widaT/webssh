package webssh

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"sync"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var (
	record     bool   = true
	recPath    string = "./rec"
	remoteAddr string = "localhost:22"
	user       string = "wida"
	password   string = "wida"
)
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024 * 10,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func ServeConn(c *gin.Context) {
	wsConn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.AbortWithStatusJSON(200, gin.H{"ok": false, "msg": err.Error()})
		return
	}
	defer wsConn.Close()
	client, err := NewSSHClient(
		//使用私钥登录
		/* 		SSHClientConfigPulicKey(
			"host:22",
			"user",
			"/home/user/ssh/user.id_rsa",
		), */
		//使用密码登录
		SSHClientConfigPassword(
			remoteAddr,
			user,
			password,
		),
	)

	if err != nil {
		wsConn.WriteControl(websocket.CloseMessage,
			[]byte(err.Error()), time.Now().Add(time.Second))
		return
	}
	defer client.Close()

	var recorder *Recorder
	if record {
		mask := syscall.Umask(0)
		defer syscall.Umask(mask)
		os.MkdirAll(recPath, 0766)
		fileName := path.Join(recPath, fmt.Sprintf("%s_%s_%s.cast", remoteAddr, user, time.Now().Format("20060102_150405")))
		f, err := os.OpenFile(fileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0766)
		if err != nil {
			c.AbortWithStatusJSON(200, gin.H{"ok": false, "msg": err.Error()})
		}
		defer f.Close()
		recorder = NewRecorder(f)
	}

	turn, err := NewTurn(wsConn, client, recorder)

	if err != nil {
		wsConn.WriteControl(websocket.CloseMessage,
			[]byte(err.Error()), time.Now().Add(time.Second))
		return
	}
	defer turn.Close()

	var logBuff = bufPool.Get().(*bytes.Buffer)
	logBuff.Reset()
	defer bufPool.Put(logBuff)

	ctx, cancel := context.WithCancel(context.Background())
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		err := turn.LoopRead(logBuff, ctx)
		if err != nil {
			log.Printf("%#v", err)
		}
	}()
	go func() {
		defer wg.Done()
		err := turn.SessionWait()
		if err != nil {
			log.Printf("%#v", err)
		}
		cancel()
	}()
	wg.Wait()
}
