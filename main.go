package main

import (
	"aurum/env"
	"aurum/handlers"
	"aurum/packages/gitinfo"
	"fmt"
	"os"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func init() {
	logrus.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})
}

func registerHandlers(app *fiber.App, handlers map[string]handlers.Route) {
	methodMap := map[string]func(path string, handler fiber.Handler, middleware ...fiber.Handler) fiber.Router{
		"GET":     app.Get,
		"POST":    app.Post,
		"PUT":     app.Put,
		"DELETE":  app.Delete,
		"PATCH":   app.Patch,
		"OPTIONS": app.Options,
		"HEAD":    app.Head,
	}

	for path, handler := range handlers {
		parts := strings.Split(strings.TrimSpace(path), " ")
		method, route := parts[0], parts[1]

		if registerFunc, ok := methodMap[method]; ok {
			var middlewareChain []fiber.Handler
			middlewareChain = append(middlewareChain, handler.Middlewares...)

			registerFunc(route, handler.Handler, middlewareChain...)
			logrus.Infof("Router: %s %s registered with %d middleware(s)", method, route, len(middlewareChain))
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
	logrus.Infof("Commit hash: %s", env.COMMIT_HASH)

	branch, err := gitinfo.GetBranch()
	if err != nil {
		logrus.Fatalf("Error getting git branch: %v", err)
	}
	env.BRANCH = branch
	logrus.Infof("Branch: %s", env.BRANCH)

	remote, err := gitinfo.GetRemote()
	if err != nil {
		logrus.Fatalf("Error getting git remote: %v", err)
	}
	env.REMOTE = remote
	logrus.Infof("Remote: %s", env.REMOTE)

	env.PORT = fmt.Sprintf(":%s", os.Getenv("PORT"))
}

func main() {
	logrus.Infoln("Starting server")

	err := godotenv.Load(".env")
	if err != nil {
		logrus.Fatalf("Error loading .env file: %v", err)
	}
	initEnv()

	app := fiber.New()
	registerHandlers(app, handlers.Handlers)

	logrus.Infof("Server started on %s", env.PORT)
	logrus.Fatal(app.Listen(env.PORT, fiber.ListenConfig{DisableStartupMessage: true}))
}
