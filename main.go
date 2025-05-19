package main

import (
	"context"
	"flag"
	"net/http"
	"om/pkg/config"
	"om/pkg/db"
	"om/pkg/ngx"
	"om/pkg/task"
	"om/pkg/util"
	"om/pkg/web"
	"time"

	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
)

var e *echo.Echo

func main() {
	action := ""
	flag.StringVar(&action, "s", "run", "service control, such as \"start\", \"stop\", \"restart\", \"install\", \"uninstall\".")
	flag.Parse()

	serviceHandler(action)
}

func OMRun() {
	err := config.LoadOrCreate()
	if err != nil {
		log.Errorf("failed to load or create config: %v", err)
		return
	}

	err = util.InitLog(config.Cfg.LogLevel, util.DataDir+"oms.log")
	if err != nil {
		log.Errorf("failed to init log: %v", err)
		return
	}

	err = db.InitDB(config.Cfg.SqlDriver, config.Cfg.DSN)
	if err != nil {
		log.Errorf("failed to init db: %v", err)
		return
	}

	ngxCmd := ngx.NewNginxCommand()
	ngxCmd.Stop()
	out, err := ngxCmd.Start()
	if err != nil {
		if out != nil {
			println("failed to start OpenResty: " + string(out))
			log.Errorf("failed to start OpenResty: %s", string(out))
			return
		}
		println("failed to start OpenResty: " + err.Error())
		log.Errorf("failed to start OpenResty: %s", err.Error())
		return
	}

	err = task.InitTask()
	if err != nil {
		log.Errorf("failed to init task: %v", err)
		return
	}

	e = web.New(config.Cfg.JwtKey)
	if err = e.Start(config.Cfg.Listen); err != nil && err != http.ErrServerClosed {
		println("failed to listen: " + err.Error())
		log.Errorf("failed to listen: %v", err)
	}
}

func OMStop() {
	if e != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 7*time.Second)
		defer cancel()
		if err := e.Shutdown(ctx); err != nil {
			log.Errorf("failed to shutdown: %v", err)
		}
	}

	task.Stop()

	if err := db.Close(); err != nil {
		log.Errorf("failed to close db: %v", err)
	}

	ngxCmd := ngx.NewNginxCommand()
	out, err := ngxCmd.Stop()
	if err != nil {
		if out != nil {
			log.Errorf("failed to stop openresty: %s", string(out))
			return
		}
		log.Errorf("failed to stop openresty: %s", err.Error())
	}
}
