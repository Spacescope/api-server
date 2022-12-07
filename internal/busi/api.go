package busi

import (
	"context"
	"fmt"
	"time"

	v1 "api-server/internal/busi/api/v1"
	"api-server/pkg/models/busi"
	"api-server/pkg/utils"

	log "github.com/sirupsen/logrus"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func registerV1(r *gin.Engine) {
	apiv1 := r.Group("/api/v1")
	{
		apiv1.GET("/ping", v1.Ping)

		{
			apiv1.GET("/contracts", v1.ListContracts)                          // list contracts
			apiv1.GET("/contract/:address", v1.GetContract)                    // contract detail
			apiv1.GET("/contract/:address/txns", v1.ListContractTXNs)          // list contract's txns
			apiv1.GET("/contract/:address/internal_txns", v1.ListInternalTXNs) // list contract's internal txns
			apiv1.GET("/contract/:address/is_verify", v1.ContractIsVerify)     // contract is verify
		}

		{
			apiv1.GET("/txns", v1.ListTXNs)
			apiv1.GET("/txn/:txnHash", v1.GetTXN)
		}

		{
			apiv1.GET("/block/:height", v1.GetBlock)
		}

		{
			apiv1.POST("/contractverify/:address", v1.SubmitContractVerify) // submit contract verify
			apiv1.GET("/contractverify/:id", v1.GetContractVerify)          // get contract verify
			apiv1.GET("/complieversions", v1.ListCompileVersion)            // list contract compile cersion
		}

		{
			apiv1.GET("/address/:address", v1.GetAddress)
			apiv1.GET("/address/:address/txns", v1.ListAddressTXNs)                  // list address's txns
			apiv1.GET("/address/:address/internal_txns", v1.ListAddressInternalTXNs) // list address's internal txns
		}
	}
}

func RegisterRoutes(r *gin.Engine) {
	// r.Use(utils.Cors())
	r.Use(cors.Default())
	r.GET("/api-server/swagger/*any", swagHandler)

	registerV1(r)
}

func initconfig(ctx context.Context, cf *utils.TomlConfig) {
	if err := utils.InitConfFile(Flags.Config, cf); err != nil {
		log.Fatalf("Load configuration file err: %v", err)
	}

	utils.EngineGroup = utils.NewEngineGroup(ctx, &[]utils.EngineInfo{
		{utils.TaskDB, cf.APIServer.DB, nil},
		{utils.APIDB, cf.APIServer.BusiDB, busi.Tables},
	})
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
