package main

import (
	"fmt"
	"github.com/labstack/echo"
	nsqlib "github.com/nsqio/go-nsq"
	"golang-learn/examples/64_manager"
	"golang-learn/examples/64_manager/nsq"
	"golang-learn/examples/64_manager/sqlcon"
	"golang-learn/examples/64_manager/web"
	"net/http"
)

// EXAMPKE NSQ HANDLER
type DummyNSQHandler struct{}

func (instance *DummyNSQHandler) HandleMessage(msg *nsqlib.Message) error {
	return nil
}

// EXAMPLE WEB SERVER HANDLER
func exampleWebServerHandler(c echo.Context) error {
	// User ID from path `users/:id`
	id := c.Param("id")
	return c.String(http.StatusOK, id)
}

func main() {
	//
	// MANAGER
	//
	manager, _ := pm.NewManager()

	//
	// SQL CONNECTION
	//
	sqlConfig := sqlcon.NewConfig("localhost", "postgres", 1, 2)
	sqlConnection, _ := manager.NewSQLConnection(sqlConfig)
	_ = manager.AddConnection("conn_1", sqlConnection)

	//
	// NSQ
	//
	nsqConfig := &nsq.Config{
		Topic:   "topic_1",
		Channel: "channel_2",
		Lookupd: []string{"http://localhost:4151"},
	}

	// Consumer
	nsqConsumer, _ := manager.NewNSQConsumer(nsqConfig, &DummyNSQHandler{})
	manager.AddProcess("teste_1", nsqConsumer)

	// Producer
	nsqProducer, _ := manager.NewNSQProducer(nsqConfig)
	manager.AddProcess("teste_2", nsqProducer)

	//
	// CONFIGURATION
	//
	simpleConfig, _ := manager.NewSimpleConfig("/Users/joaoribeiro/workspace/go/src/golang-learn/examples/64_manager/run/", "config", "json")
	manager.AddConfig("teste_3", simpleConfig)

	// Get configuration by path
	fmt.Println("a: ", manager.GetConfig("teste_3").Get("a"))
	fmt.Println("caa: ", manager.GetConfig("teste_3").Get("c.ca.caa"))

	// Get configuration by tag
	fmt.Println("a: ", manager.GetConfig("teste_3").Get("a"))
	fmt.Println("caa: ", manager.GetConfig("teste_3").Get("c.ca.caa"))

	//
	// HTTP SERVER
	//
	configWebServer := web.NewConfig("8080")
	webServer, _ := manager.NewWEBServer(configWebServer)
	manager.AddWebServer("web_server_1", webServer)
	manager.AddWebServerRoute("web_server_1", http.MethodGet, "/example", exampleWebServerHandler)
	manager.StartWebServer("web_server_1")

	//
	// ELASTIC SEARCH CLIENT
	//

	manager.Start()
}
