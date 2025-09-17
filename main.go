package main

import (
	"context"
	"embed"

	"github.com/adrianpk/clio/internal/am"
	"github.com/adrianpk/clio/internal/core"
	"github.com/adrianpk/clio/internal/feat/auth"
	"github.com/adrianpk/clio/internal/feat/ssg"
	"github.com/adrianpk/clio/internal/repo/sqlite"
	appssg "github.com/adrianpk/clio/internal/web/ssg"
)

const (
	name      = "clio"
	version   = "v1"
	namespace = "clio"
	engine    = "sqlite"
)

//go:embed assets
var assetsFS embed.FS

func main() {
	ctx := context.Background()
	log := am.NewLogger("debug")
	cfg := am.LoadCfg(namespace, am.Flags)
	opts := am.DefOpts(log, cfg)

	fm := am.NewFlashManager()
	workspace := ssg.NewWorkspace(opts...)
	app := core.NewApp(name, version, assetsFS, opts...)
	queryManager := am.NewQueryManager(assetsFS, engine)
	templateManager := am.NewTemplateManager(assetsFS)
	repo := sqlite.NewClioRepo(queryManager)
	migrator := am.NewMigrator(assetsFS, engine)
	fileServer := am.NewFileServer(assetsFS)

	app.MountFileServer("/", fileServer)

	apiRouter := am.NewAPIRouter("api-router", opts...)

	// Auth feature
	authSeeder := auth.NewSeeder(assetsFS, engine, repo)
	authService := auth.NewService(repo, opts...)
	authAPIHandler := auth.NewAPIHandler("auth-api-handler", authService, opts...)
	authAPIRouter := auth.NewAPIRouter(authAPIHandler, nil) // No middleware for now
	apiRouter.Mount("/auth", authAPIRouter)

	// SSG feature
	ssgSeeder := ssg.NewSeeder(assetsFS, engine, repo)
	ssgService := ssg.NewService(repo, opts...)
	ssgAPIHandler := ssg.NewAPIHandler("ssg-api-handler", ssgService)
	ssgAPIRouter := ssg.NewAPIRouter(ssgAPIHandler, nil) // No middleware for now
	apiRouter.Mount("/ssg", ssgAPIRouter)

	app.MountAPI("v1", "/", apiRouter)

	// Web app
	ssgWebHandler := appssg.NewWebHandler(templateManager, fm, opts...)
	ssgWebRouter := appssg.NewWebRouter(ssgWebHandler, append(fm.Middlewares(), am.LogHeadersMw))

	app.MountWeb("/ssg", ssgWebRouter)

	// Add deps
	app.Add(workspace)
	app.Add(migrator)
	app.Add(fm)
	app.Add(fileServer)
	app.Add(queryManager)
	app.Add(templateManager)
	app.Add(repo)
	app.Add(ssgWebHandler)
	app.Add(ssgWebRouter)
	app.Add(ssgSeeder)
	app.Add(ssgService)
	app.Add(ssgAPIHandler)
	app.Add(ssgAPIRouter)
	app.Add(apiRouter)
	app.Add(authSeeder)

	err := app.Setup(ctx)
	if err != nil {
		log.Errorf("Cannot setup %s(%s): %v", name, version, err)
		return
	}

	err = app.Start(ctx)
	if err != nil {
		log.Errorf("Cannot start %s(%s): %v", name, version, err)
	}
}
