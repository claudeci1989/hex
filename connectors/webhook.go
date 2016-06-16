package connectors

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/projectjane/jane/models"
)

// Webhook Struct for manipulating the webhook connector
type Webhook struct {
	CommandMsgs chan<- models.Message
	Connector   models.Connector
}

var webhook Webhook

// Listen Webhook listener
func (x Webhook) Listen(commandMsgs chan<- models.Message, connector models.Connector) {
	defer Recovery(connector)

	webhook = x
	webhook.CommandMsgs = commandMsgs
	webhook.Connector = connector

	port, _ := strconv.Atoi(connector.Port)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: nil,
	}

	if connector.Debug {
		log.Printf("Starting Webhook connector at: %s", server.Addr)
	}

	http.HandleFunc("/webhook/", webhookHandler)

	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}

// Command Webhook command parser
func (x Webhook) Command(message models.Message, publishMsgs chan<- models.Message, connector models.Connector) {
	if message.In.Process && connector.Type == "webhook" {
		if connector.Debug {
			log.Printf("Processing command: %s", message.In.Text)
		}

		for _, c := range connector.Commands {
			if strings.HasPrefix(strings.ToLower(message.In.Text), strings.ToLower(c.Match)) {
				msg := strings.TrimSpace(strings.Replace(message.In.Text, c.Match, "", 1))

				message.Out.Text = msg
				publishMsgs <- message
			}
		}
	}
}

// Publish Webhook does not publish
func (x Webhook) Publish(connector models.Connector, message models.Message, target string) {
	return
}

// Help Webhook help information
func (x Webhook) Help(connector models.Connector) (help string) {
	return
}

func webhookHandler(w http.ResponseWriter, r *http.Request) {
	var segs []string

	// get everything past /webhook/
	webhookString := r.URL.Path[9:]

	segs = strings.Split(webhookString, "+")

	if len(segs) < 2 || segs[1] == "" {
		segs = strings.Split(webhookString, "/")

		if len(segs) < 1 {
			w.WriteHeader(http.StatusNotFound)
			log.Println("Route not found")
			fmt.Fprintf(w, "Route not found")
			return
		}
	}

	if segs[0] == "" {
		w.WriteHeader(http.StatusNotFound)
		log.Println("Empty webhook command")
		fmt.Fprintf(w, "Empty webhook command")
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
	}
	defer r.Body.Close()

	command := strings.Join(segs[0:], " ")

	if webhook.Connector.Debug {
		log.Printf("Command received: %s", command)
	}

	var m models.Message
	m.Routes = webhook.Connector.Routes
	m.In.Source = webhook.Connector.ID
	m.In.Text = command
	m.In.Process = true
	m.Out.Detail = string(body)
	webhook.CommandMsgs <- m

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Request received."))
}
