package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"sync"
	"time"

	"github.com/oguzbilgic/pusher"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(2)
	accountKey := os.Getenv("SCOUT_KEY")
	scoutGemBinPath := os.Getenv("SCOUT_GEM_BIN_PATH") + "/scout"
	go listenForRealtime(accountKey, scoutGemBinPath, &wg)
	go reportLoop(accountKey, scoutGemBinPath, &wg)
	wg.Wait()
}

func reportLoop(accountKey string, scoutGemBinPath string, wg *sync.WaitGroup) {
	c := time.Tick(1 * time.Second)
	for _ = range c {
		// fmt.Printf("report loop\n")
		cmd := exec.Command(scoutGemBinPath, accountKey)
		err := cmd.Run()
		if err != nil {
			log.Fatal(err)
		}
	}
	wg.Done()
}

func listenForRealtime(accountKey string, scoutGemBinPath string, wg *sync.WaitGroup) {
	conn, err := pusher.New("f07eaa39898f3c36c8cf")
	if err != nil {
		panic(err)
	}

	hostName := os.Getenv("SCOUT_HOSTNAME")
	commandChan := conn.Channel(accountKey + "-" + hostName)
	messages := commandChan.Bind("streamer_command") // a go channel is returned

	for {
		msg := <-messages
		fmt.Printf(scoutGemBinPath + " realtime " + msg.(string))
		cmd := exec.Command(scoutGemBinPath, "realtime", msg.(string))
		err := cmd.Run()
		if err != nil {
			log.Fatal(err)
		}
	}
	wg.Done()
}
