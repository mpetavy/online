package main

import (
	"embed"
	"os/exec"
	"time"

	"github.com/mpetavy/common"
)

//go:embed go.mod
var resources embed.FS
var (
	lastOnline bool
	quit       = make(chan struct{})
)

func init() {
	common.Init("", "", "", "", "", "", "", "", &resources, start, stop, nil, 0)
}

func online() bool {
	cmd := exec.Command("ping", "-c", "3", "heise.de")
	ba,err := cmd.CombinedOutput()
	online := err == nil

	if online == lastOnline {
		return online
	}

	lastOnline = online

	n := time.Now().Format(time.RFC3339)

	if online {
		common.Info("%v Online", n)
	} else {
		common.Info("%v Offline - %s\n%s", n,err.Error(),string(ba))
	}

	return online
}

func run() error {
	ticker := time.NewTicker(time.Second * 5)
	defer func() {
		ticker.Stop()
	}()

loop:
	for {
		select {
		case <-quit:
			break loop
		case <-ticker.C:
			online()
		}
	}

	return nil
}

func start() error {
	go run()

	return nil
}

func stop() error {
	close(quit)

	return nil
}

func main() {
	common.Run(nil)
}
