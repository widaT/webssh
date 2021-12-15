package webssh

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type WebSSHConfig struct {
	Record     bool
	RecPath    string
	RemoteAddr string
	User       string
	Password   string
	AuthModel  AuthModel
	PkPath     string
}

type WebSSH struct {
	*WebSSHConfig
}

func NewWebSSH(conf *WebSSHConfig) *WebSSH {
	return &WebSSH{
		WebSSHConfig: conf,
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024 * 10,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (w WebSSH) ServeConn(c *gin.Context) {
	wsConn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.AbortWithStatusJSON(200, gin.H{"ok": false, "msg": err.Error()})
		return
	}
	defer wsConn.Close()
	var config *SSHClientConfig
	switch w.AuthModel {

	case PASSWORD:
		config = SSHClientConfigPassword(
			w.RemoteAddr,
			w.User,
			w.Password,
		)
	case PUBLICKEY:
		config = SSHClientConfigPulicKey(
			w.RemoteAddr,
			w.User,
			w.PkPath,
		)
	}

	client, err := NewSSHClient(config)
	if err != nil {
		wsConn.WriteControl(websocket.CloseMessage,
			[]byte(err.Error()), time.Now().Add(time.Second))
		return
	}
	defer client.Close()

	var recorder *Recorder
	if w.Record {
		mask := syscall.Umask(0)
		defer syscall.Umask(mask)
		os.MkdirAll(w.RecPath, 0766)
		fileName := path.Join(w.RecPath, fmt.Sprintf("%s_%s_%s.cast", w.RemoteAddr, w.User, time.Now().Format("20060102_150405")))
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

func (w WebSSH) RecoderList(c *gin.Context) {
	files, err := ioutil.ReadDir(w.RecPath)
	if err != nil {
		c.AbortWithStatusJSON(200, gin.H{"ok": false, "msg": err.Error()})
		return
	}
	var filesName []string
	for _, f := range files {
		if f.IsDir() {
			continue
		}

		if strings.HasSuffix(f.Name(), ".cast") {
			filesName = append(filesName, f.Name())
		}
	}
	c.JSON(200, filesName)
}
