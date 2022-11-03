package busi

import (
	v1 "api-server/internal/busi/api/v1"
	"api-server/internal/busi/core"
	"api-server/pkg/utils"
	"context"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func registerV1(r *gin.Engine) {
	apiv1 := r.Group("/v1")
	tokenGroup_apiv1 := apiv1.Group("", core.AuthFromGateway(), core.GinUidLogMiddleware())
	{
		tokenGroup_apiv1.GET("/network_core/circulating-supply/circulating_supply", v1.ListCirculatingSupply)
	}
}

func RegisterRoutes(r *gin.Engine) {
	// r.Use(utils.Cors())
	r.Use(cors.Default())
	r.GET("/api-server/swagger/*any", swagHandler)
	r.GET("/api/v1/ping", v1.Ping)

	registerV1(r)
}

func initconfig(ctx context.Context, cf *utils.TomlConfig) {
	if err := utils.InitConfFile(Flags.Config, cf); err != nil {
		log.Fatalf("Load configuration file err: %v", err)
	}

	utils.EngineGroup = utils.NewEngineGroup(ctx, &[]utils.EngineInfo{{utils.DB, cf.APIServer.DB, nil}})

	// utils.InitKVEngine(ctx, cf.DataInfra.KV, "", 0)
}

func Start() {
	initconfig(context.Background(), &utils.CNF)

	// if Flags.Mode == "prod" {
	gin.SetMode(gin.ReleaseMode)
	// }

	// r := gin.Default()
	r := gin.New()
	// LoggerWithFormatter middleware will write the logs to gin.DefaultWriter
	// By default gin.DefaultWriter = os.Stdout
	r.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		// your custom format
		return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
			param.ClientIP,
			param.TimeStamp.Format(time.RFC1123),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	}))
	r.Use(gin.Recovery())

	RegisterRoutes(r)

	r.Run(utils.CNF.APIServer.Addr)
}
