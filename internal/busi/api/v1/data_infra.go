package v1

import (
	"api-server/internal/busi/core"
	"api-server/pkg/utils"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/go-playground/validator/v10"
)

// ListContracts godoc
// @Description List contracts
// @Tags DATA-INFRA-API-External-V1
// @Accept application/json,json
// @Produce application/json,json
// @Param ListQuery query core.ListQuery true "ListQuery"
// @Success 200 {object} core.Contract
// @Failure 400 {object} utils.ResponseWithRequestId
// @Failure 500 {object} utils.ResponseWithRequestId
// @Router /api/v1/contracts [get]
func ListContracts(c *gin.Context) {
	app := utils.Gin{C: c}

	var r core.ListQuery
	if err := c.ShouldBindQuery(&r); err != nil {
		app.HTTPResponse(http.StatusBadRequest, utils.NewResponse(utils.CodeBadRequest, err.Error(), nil))
		return
	}

	if err := r.ListValidate(); err != nil {
		app.HTTPResponse(http.StatusBadRequest, utils.NewResponse(utils.CodeBadRequest, err.Error(), nil))
		return
	}

	result, resp := core.ListContracts(c.Request.Context(), &r)
	if resp != nil {
		app.HTTPResponse(resp.HttpCode, resp.Response)
		return
	}

	app.HTTPResponseOK(result)
}

// GetContract godoc
// @Description Get contract detail
// @Tags DATA-INFRA-API-External-V1
// @Accept application/json,json
// @Produce application/json,json
// @Param address path string true "address"
// @Success 200 {object} core.Contract
// @Failure 400 {object} utils.ResponseWithRequestId
// @Failure 500 {object} utils.ResponseWithRequestId
// @Router /api/v1/contract/{address} [get]
func GetContract(c *gin.Context) {
	app := utils.Gin{C: c}
	validate := validator.New()

	address := c.Param("address")
	if err := validate.Var(address, "required"); err != nil {
		app.HTTPResponse(http.StatusBadRequest, utils.NewResponse(utils.CodeBadRequest, err.Error(), nil))
	}

	result, resp := core.Getcontract(c.Request.Context(), address)
	if resp != nil {
		app.HTTPResponse(resp.HttpCode, resp.Response)
		return
	}

	app.HTTPResponseOK(result)
}

// ListTXNs godoc
// @Description List contract's transactions
// @Tags DATA-INFRA-API-External-V1
// @Accept application/json,json
// @Produce application/json,json
// @Param ListQuery query core.ListQuery true "ListQuery"
// @Param address path string true "address"
// @Success 200 {object} busi.EVMTransaction
// @Failure 400 {object} utils.ResponseWithRequestId
// @Failure 500 {object} utils.ResponseWithRequestId
// @Router /api/v1/contract/{address}/txns [get]
func ListTXNs(c *gin.Context) {
	app := utils.Gin{C: c}
	validate := validator.New()

	address := c.Param("address")
	if err := validate.Var(address, "required"); err != nil {
		app.HTTPResponse(http.StatusBadRequest, utils.NewResponse(utils.CodeBadRequest, err.Error(), nil))
	}

	var r core.ListQuery
	if err := c.ShouldBindQuery(&r); err != nil {
		app.HTTPResponse(http.StatusBadRequest, utils.NewResponse(utils.CodeBadRequest, err.Error(), nil))
		return
	}

	if err := r.ListValidate(); err != nil {
		app.HTTPResponse(http.StatusBadRequest, utils.NewResponse(utils.CodeBadRequest, err.Error(), nil))
		return
	}

	result, resp := core.ListTXNs(c.Request.Context(), address, &r)
	if resp != nil {
		app.HTTPResponse(resp.HttpCode, resp.Response)
		return
	}

	app.HTTPResponseOK(result)
}

// ListInternalTXNs godoc
// @Description List contract's internal transactions
// @Tags DATA-INFRA-API-External-V1
// @Accept application/json,json
// @Produce application/json,json
// @Param ListQuery query core.ListQuery true "ListQuery"
// @Param address path string true "address"
// @Success 200 {object} busi.EVMInternalTX
// @Failure 400 {object} utils.ResponseWithRequestId
// @Failure 500 {object} utils.ResponseWithRequestId
// @Router /api/v1/contract/{address}/internal_txns [get]
func ListInternalTXNs(c *gin.Context) {
	app := utils.Gin{C: c}
	validate := validator.New()

	address := c.Param("address")
	if err := validate.Var(address, "required"); err != nil {
		app.HTTPResponse(http.StatusBadRequest, utils.NewResponse(utils.CodeBadRequest, err.Error(), nil))
	}

	var r core.ListQuery
	if err := c.ShouldBindQuery(&r); err != nil {
		app.HTTPResponse(http.StatusBadRequest, utils.NewResponse(utils.CodeBadRequest, err.Error(), nil))
		return
	}

	if err := r.ListValidate(); err != nil {
		app.HTTPResponse(http.StatusBadRequest, utils.NewResponse(utils.CodeBadRequest, err.Error(), nil))
		return
	}

	result, resp := core.ListInternalTXNs(c.Request.Context(), address, &r)
	if resp != nil {
		app.HTTPResponse(resp.HttpCode, resp.Response)
		return
	}

	app.HTTPResponseOK(result)

}
