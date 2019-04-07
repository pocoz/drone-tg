package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/go-kit/kit/transport/http/jsonrpc"
	"github.com/kelseyhightower/envconfig"

	"github.com/pocoz/drone-tg/types"
)

func main() {
	const (
		exitCodeSuccess = 0
		exitCodeFailure = 1
	)

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	logger := log.NewLogfmtLogger(log.NewSyncWriter(os.Stdout))
	logger = log.With(logger, "ts", log.DefaultTimestampUTC, "caller", log.DefaultCaller)

	var (
		cfgPlugin types.ConfigurationPlugin
		cfgDrone  types.ConfigurationDrone
	)
	if err := envconfig.Process("", &cfgPlugin); err != nil {
		level.Error(logger).Log("msg", "failed to load configuration plugin", "err", err)
		os.Exit(exitCodeFailure)
	}
	if err := envconfig.Process("", &cfgDrone); err != nil {
		level.Error(logger).Log("msg", "failed to load configuration drone", "err", err)
		os.Exit(exitCodeFailure)
	}

	if strings.TrimSpace(cfgPlugin.Token) == "" {
		level.Error(logger).Log("msg", "empty bot token")
		os.Exit(exitCodeFailure)
	}
	if cfgPlugin.ChatID == 0 {
		level.Error(logger).Log("msg", "empty chat id")
		os.Exit(exitCodeFailure)
	}

	// Create URL
	url := cfgPlugin.ProxyURL + "/bot" + cfgPlugin.Token + "/sendMessage"

	// Create body message
	text := "Build status: " + cfgDrone.BuildStatus + "\nBuild link: " + cfgDrone.BuildLink + "\nRepo: " + cfgDrone.RepoName + "\nCommit: " + cfgDrone.CommitMessage
	msg := &types.MessageBody{
		ChatID: cfgPlugin.ChatID,
		Text:   text,
	}

	msgByte, err := json.Marshal(msg)
	if err != nil {
		level.Info(logger).Log("msg", "json marshal failed", "err", err)
		os.Exit(exitCodeFailure)
	}

	// SendMessage method call
	// Called Method Documentation: https://core.telegram.org/bots/api#sendmessage
	r, err := http.Post(url, jsonrpc.ContentType, bytes.NewBuffer(msgByte))
	if err != nil {
		level.Info(logger).Log("msg", "the request failed", "err", err)
		os.Exit(exitCodeFailure)
	}
	defer r.Body.Close()
	if cfgPlugin.Debug {
		contents, err := ioutil.ReadAll(r.Body)
		if err != nil {
			fmt.Printf("%s", err)
			os.Exit(exitCodeFailure)
		}
		level.Info(logger).Log("msg", "tg api", "response", string(contents))
	}

	errc := make(chan error, 1)
	donec := make(chan struct{})
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, syscall.SIGTERM, os.Interrupt)
	defer func() {
		signal.Stop(sigc)
		cancel()
	}()

	close(donec)

	go func() {
		select {
		case sig := <-sigc:
			level.Info(logger).Log("msg", "received signal, exiting", "signal", sig)
			signal.Stop(sigc)
			close(donec)
		case <-errc:
			level.Info(logger).Log("msg", "now exiting with error", "error code", exitCodeFailure)
			os.Exit(exitCodeFailure)
		}
	}()

	<-donec
	level.Info(logger).Log("msg", "all done")
	os.Exit(exitCodeSuccess)
}
