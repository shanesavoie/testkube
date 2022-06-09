package telemetry

import (
	"net/http"
	"sync"

	"strings"

	"github.com/spf13/cobra"

	"github.com/kubeshop/testkube/cmd/tools/commands"
	httpclient "github.com/kubeshop/testkube/pkg/http"
	"github.com/kubeshop/testkube/pkg/log"
)

var (
	client  = httpclient.NewClient()
	senders = []Sender{
		GoogleAnalyticsSender,
		TestkubeAnalyticsSender,
	}
)

type Sender func(client *http.Client, payload Payload) (out string, err error)

// SendServerStartEvent will send event to GA
func SendServerStartEvent() (string, error) {
	payload := NewCLIPayload(MachineID(), "testkube_api_start", commands.Version, "execution")
	return SendData(senders, payload)
}

// SendCmdEvent will send CLI event to GA
func SendCmdEvent(cmd *cobra.Command, version string) (string, error) {
	// get all sub-commands passed to cli
	command := strings.TrimPrefix(cmd.CommandPath(), "kubectl-testkube ")
	if command == "" {
		command = "root"
	}

	payload := NewCLIPayload(MachineID(), command, version, "execution")
	return SendData(senders, payload)
}

// SendCmdInitEvent will send CLI event to GA
func SendCmdInitEvent(cmd *cobra.Command, version string) (string, error) {
	payload := NewCLIPayload(MachineID(), "init", version, "execution")
	return SendData(senders, payload)
}

// SendHeartbeatEvent will send CLI event to GA
func SendHeartbeatEvent(host, version, clusterId string) (string, error) {
	payload := NewAPIPayload(clusterId, "testkube_api_heartbeat", version, host)
	return SendData(senders, payload)
}

// SendData sends data to all telemetry storages  in parallel and syncs sending
func SendData(senders []Sender, payload Payload) (out string, err error) {
	var wg sync.WaitGroup
	wg.Add(len(senders))
	for _, sender := range senders {
		go func(sender Sender) {
			defer wg.Done()
			o, err := sender(client, payload)
			if err != nil {
				log.DefaultLogger.Debugw("sending telemetry data error", "payload", payload, "error", err.Error())
				return
			}
			log.DefaultLogger.Debugw("sending telemetry data", "payload", payload, "output", o)
		}(sender)
	}

	wg.Wait()

	return out, nil
}