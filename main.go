package main

import (
	"aurum/env"
	"aurum/handlers"
	"aurum/packages/gitinfo"
	"fmt"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func init() {
	logrus.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})
}

func registerHandlers(app *fiber.App, handlers map[string]handlers.Route) {
	methodMap := map[string]func(path string, handlers ...fiber.Handler) fiber.Router{
		"GET":     app.Get,
		"POST":    app.Post,
		"PUT":     app.Put,
		"DELETE":  app.Delete,
		"PATCH":   app.Patch,
		"OPTIONS": app.Options,
		"HEAD":    app.Head,
		"SOCKET":  app.Get,
	}

	for path, handler := range handlers {
		parts := strings.Split(strings.TrimSpace(path), " ")
		method, route := parts[0], parts[1]

		if method == "SOCKET" {
			app.Use(route, handler.Handler)
			app.Get(route, websocket.New(handler.Endpoint))
			logrus.Infof("Router: %s %s registered as WebSocket", method, route)
			continue
		}

		if registerFunc, ok := methodMap[method]; ok {
			var handlerPipeline []fiber.Handler
			handlerPipeline = append(handlerPipeline, handler.Handler)
			handlerPipeline = append(handlerPipeline, handler.Middlewares...)

			registerFunc(route, handlerPipeline...)
			logrus.Infof("Router: %s %s registered with %d middleware(s)", method, route, len(handlerPipeline)-1)
		} else {
			logrus.Fatalf("Invalid method: %s", method)
		}
	}
}

func initEnv() {
	commit, err := gitinfo.GetCommit()
	if err != nil {
		logrus.Fatalf("Error getting git commit: %v", err)
	}
	env.COMMIT_HASH = commit

	branch, err := gitinfo.GetBranch()
	if err != nil {
		logrus.Fatalf("Error getting git branch: %v", err)
	}
	env.BRANCH = branch

	remote, err := gitinfo.GetRemote()
	if err != nil {
		logrus.Fatalf("Error getting git remote: %v", err)
	}
	env.REMOTE = remote

	env.PORT = fmt.Sprintf(":%s", os.Getenv("PORT"))
	env.ETHERSCAN_API_KEY = os.Getenv("ETHERSCAN_API_KEY")
}

func main() {
	logrus.Infoln("Starting server")

	err := godotenv.Load(".env")
	if err != nil {
		logrus.Fatalf("Error loading .env file: %v", err)
	}
	initEnv()

	app := fiber.New(fiber.Config{DisableStartupMessage: true})
	registerHandlers(app, handlers.Handlers)

	logrus.Infof("Server started on %s", env.PORT)
	logrus.Fatal(app.Listen(env.PORT))
}
