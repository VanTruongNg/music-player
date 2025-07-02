package main

import (
	"fmt"
	"log"

	"auth-service/configs"
)

func main() {
	appCfg := configs.LoadAppConfig()
	dbCfg := configs.LoadDBConfig()
	redisCfg := configs.LoadRedisConfig()

	app, err := InitializeApp(appCfg, dbCfg, redisCfg)
	if err != nil {
		log.Fatalf("[FATAL] Failed to initialize app: %v", err)
	}

	addr := fmt.Sprintf(":%s", appCfg.Port)
	log.Printf("[INFO] App running in %s mode on port %s", appCfg.Env, appCfg.Port)
	if err := app.Router.Run(addr); err != nil {
		log.Fatalf("[FATAL] Failed to start Gin server: %v", err)
	}
}
