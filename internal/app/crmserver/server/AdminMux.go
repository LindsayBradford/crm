// Copyright (c) 2018 Australian Rivers Institute.

package server

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"

	"golang.org/x/net/context"
)

type Status struct {
	Name    string
	Version string
	Message string
	Time    string
}

type AdminMux struct {
	RestMux
	Status Status

	doneChannel chan bool
}

func (am *AdminMux) Initialise() *AdminMux {
	am.RestMux.Initialise()

	am.doneChannel = make(chan bool)
	am.handlerMap["/status"] = am.statusHandler
	am.handlerMap["/shutdown"] = am.shutdownHandler

	return am
}

func (am *AdminMux) WithType(muxType string) *AdminMux {
	am.RestMux.WithType(muxType)
	return am
}

func (am *AdminMux) setStatus(statusMessage string) {
	am.Logger.Info("Changed server Status to [" + statusMessage + "]")
	am.Status.Message = statusMessage
	am.UpdateStatusTime()
}

func (am *AdminMux) Start(address string) {
	am.setStatus("RUNNING")
	am.RestMux.Start(address)
}

func (am *AdminMux) WaitForShutdownSignal() {
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, os.Kill)
		<-sigint

		am.Logger.Warn("Received Operating System Interrupt/Kill signal -- triggering graceful shutdown")

		close(am.doneChannel)
	}()

	<-am.doneChannel
}

func (am *AdminMux) statusHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		am.ServeMethodNotAllowedError(w, r)
		return
	}

	am.Logger.Debug("Responding with status [" + am.Status.Message + "]")
	am.UpdateStatusTime()

	statusJson, encodeError := json.MarshalIndent(am.Status, "", "  ")
	if encodeError != nil {
		am.Logger.Error(encodeError)
	}

	fmt.Fprintf(w, string(statusJson))
}

func (am *AdminMux) shutdownHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		am.ServeMethodNotAllowedError(w, r)
		return
	}

	am.Status.Message = "SHUTTING_DOWN"
	am.Logger.Debug("Responding with status [" + am.Status.Message + "]")
	am.UpdateStatusTime()

	statusJson, encodeError := json.MarshalIndent(am.Status, "", "  ")
	if encodeError != nil {
		am.Logger.Error(encodeError)
	}

	bufferedWriter := bufio.NewWriter(w)

	fmt.Fprintf(bufferedWriter, string(statusJson))
	bufferedWriter.Flush()

	am.doneChannel <- true
}

func (am *AdminMux) UpdateStatusTime() {
	am.Status.Time = FormattedTimestamp()
}

func (am *AdminMux) Shutdown() {
	am.server.Shutdown(context.Background())
	am.setStatus("DEAD")
}
