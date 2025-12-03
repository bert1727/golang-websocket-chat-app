package main

import (
	"log"

	// "github.com/gofiber/contrib/v3/websocket"

	"github.com/bert1727/ChatApp/internal/db"
	"github.com/bert1727/ChatApp/internal/handlers"
	httphandlers "github.com/bert1727/ChatApp/internal/http_handlers"
	"github.com/bert1727/ChatApp/internal/repository"
	"github.com/bert1727/ChatApp/internal/service"
	"github.com/bert1727/ChatApp/internal/ws"
	"github.com/gofiber/contrib/v3/websocket"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
)

func main() {
	app := fiber.New()
	// hub := ws.NewHub()
	// go hub.Run()

	app.Use(cors.New())

	// TODO: do this in separate function
	app.Use("/ws", func(c fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			return c.Next()
		}
		log.Println("User connected without ws")
		return fiber.NewError(fiber.StatusUpgradeRequired, "WebSocket upgrade required")
	})

	db := db.SetupDB()

	messageRepo := repository.NewMessageRepository(db)
	userRepo := repository.NewUserRepository(db)

	userService := service.NewUserService(userRepo)
	messageService := service.NewMessageService(messageRepo)

	hub := ws.NewHub(userService, messageService)
	wsHandler := handlers.NewWSHandler(hub, userService)

	go hub.Run()

	userHandler := httphandlers.NewUserHanlder(userService)
	app.Post("/login", userHandler.Register)
	app.Get("/ws", wsHandler.HandleWSConnection)

	// id := "1"
	// app.Get("/ws", websocket.New(func(conn *websocket.Conn) {
	// 	id += "1"
	// 	client := &ws.Client{
	// 		ID:       id,
	// 		User:     ws.User{},
	// 		Conn:     conn,
	// 		Hub:      hub,
	// 		SendChan: make(chan ws.Message, 32), // limit the size
	// 	}
	//
	// 	log.Println("client connected:", client.ID)
	// 	hub.RegisterClient(client)
	//
	// 	defer func() {
	// 		log.Println("client disconnected:", client.ID)
	// 		// client.Hub.Unregister <- client
	// 		hub.UnregisterClient(client)
	// 		conn.Close()
	// 	}()
	//
	// 	go client.WritePump()
	// 	client.ReadPump()
	// }))

	log.Println("Server running on http://localhost:6969")
	log.Fatal(app.Listen(":6969"))
}
