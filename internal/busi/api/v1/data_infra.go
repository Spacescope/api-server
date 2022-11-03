package v1

import (
	"api-server/internal/busi/core"
	"api-server/pkg/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ListCirculatingSupply godoc
// @Description List table circulating_supply
// @Tags DATA-INFRA-API-External-V1
// @Accept application/json,json
// @Produce application/json,json
// @Param user_id header string true "user_id"
// @Param authorization header string true "authorization"
// @Param APIRequest query core.APIRequest false "APIRequest"
// @Success 200 {object} busi.CirculatingSupply
// @Failure 400 {object} utils.ResponseWithRequestId
// @Failure 500 {object} utils.ResponseWithRequestId
// @Router /v1/network_core/circulating-supply/circulating_supply [get]
func ListCirculatingSupply(c *gin.Context) {
	app := utils.Gin{C: c}

	r := commonValidate(c)
	if r == nil {
		return
	}

	result, resp := core.ListCirculatingSupply(c.Request.Context(), c, r)
	if resp != nil {
		app.HTTPResponse(resp.HttpCode, resp.Response)
		return
	}

	app.HTTPResponseOK(result)
}

func commonValidate(c *gin.Context) *core.APIRequest {
	app := utils.Gin{C: c}

	var r core.APIRequest
	if err := c.ShouldBindQuery(&r); err != nil {
		app.HTTPResponse(http.StatusBadRequest, utils.NewResponse(utils.CodeDataInfraAPIParamsErr, err.Error(), nil))
		return nil
	}
	if err := r.Validate(); err != nil {
		app.HTTPResponse(http.StatusBadRequest, utils.NewResponse(utils.CodeDataInfraAPIParamsErr, err.Error(), nil))
		return nil
	}

	return &r
}
