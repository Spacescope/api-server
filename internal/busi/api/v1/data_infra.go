package v1

import (
	"net/http"
	"strconv"
	"strings"

	"api-server/internal/busi/core"
	"api-server/pkg/utils"

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
		app.HTTPResponse(http.StatusOK, utils.NewResponse(utils.CodeBadRequest, err.Error(), nil))
		return
	}

	if err := r.ListValidate(); err != nil {
		app.HTTPResponse(http.StatusOK, utils.NewResponse(utils.CodeBadRequest, err.Error(), nil))
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
		app.HTTPResponse(http.StatusOK, utils.NewResponse(utils.CodeBadRequest, err.Error(), nil))
	}

	result, resp := core.GetContract(c.Request.Context(), strings.ToLower(address))
	if resp != nil {
		app.HTTPResponse(resp.HttpCode, resp.Response)
		return
	}

	app.HTTPResponseOK(result)
}

// ListContractTXNs godoc
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
func ListContractTXNs(c *gin.Context) {
	app := utils.Gin{C: c}
	validate := validator.New()

	address := c.Param("address")
	if err := validate.Var(address, "required"); err != nil {
		app.HTTPResponse(http.StatusOK, utils.NewResponse(utils.CodeBadRequest, err.Error(), nil))
	}

	var r core.ListQuery
	if err := c.ShouldBindQuery(&r); err != nil {
		app.HTTPResponse(http.StatusOK, utils.NewResponse(utils.CodeBadRequest, err.Error(), nil))
		return
	}

	if err := r.ListValidate(); err != nil {
		app.HTTPResponse(http.StatusOK, utils.NewResponse(utils.CodeBadRequest, err.Error(), nil))
		return
	}

	result, resp := core.ListContractTXNs(c.Request.Context(), strings.ToLower(address), &r)
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
		app.HTTPResponse(http.StatusOK, utils.NewResponse(utils.CodeBadRequest, err.Error(), nil))
	}

	var r core.ListQuery
	if err := c.ShouldBindQuery(&r); err != nil {
		app.HTTPResponse(http.StatusOK, utils.NewResponse(utils.CodeBadRequest, err.Error(), nil))
		return
	}

	if err := r.ListValidate(); err != nil {
		app.HTTPResponse(http.StatusOK, utils.NewResponse(utils.CodeBadRequest, err.Error(), nil))
		return
	}

	result, resp := core.ListInternalTXNs(c.Request.Context(), strings.ToLower(address), &r)
	if resp != nil {
		app.HTTPResponse(resp.HttpCode, resp.Response)
		return
	}

	app.HTTPResponseOK(result)

}

// SubmitContractVerify godoc
// @Description submit contract verify
// @Tags DATA-INFRA-API-External-V1
// @Accept application/json,json
// @Produce application/json,json
// @Param SubmitContractVerifyRequest query core.SubmitContractVerifyRequest true "SubmitContractVerifyRequest"
// @Param address path string true "address"
// @Success 200 {object} busi.EVMContractVerify
// @Failure 400 {object} utils.ResponseWithRequestId
// @Failure 500 {object} utils.ResponseWithRequestId
// @Router /api/v1/contractverify/{address} [post]
func SubmitContractVerify(c *gin.Context) {
	app := utils.Gin{C: c}
	validate := validator.New()

	address := c.Param("address")
	if err := validate.Var(address, "required"); err != nil {
		app.HTTPResponse(http.StatusOK, utils.NewResponse(utils.CodeBadRequest, err.Error(), nil))
	}

	var r core.SubmitContractVerifyRequest
	if err := c.ShouldBindJSON(&r); err != nil {
		app.HTTPResponse(http.StatusOK, utils.NewResponse(utils.CodeBadRequest, err.Error(), nil))
		return
	}

	if err := r.Validate(); err != nil {
		app.HTTPResponse(http.StatusOK, utils.NewResponse(utils.CodeBadRequest, err.Error(), nil))
		return
	}

	result, resp := core.SubmitContractVerify(c.Request.Context(), strings.ToLower(address), &r)
	if resp != nil {
		app.HTTPResponse(resp.HttpCode, resp.Response)
		return
	}

	app.HTTPResponseOK(result)

}

// GetContractVerify godoc
// @Description check contract verify
// @Tags DATA-INFRA-API-External-V1
// @Accept application/json,json
// @Produce application/json,json
// @Param id path int true "id"
// @Success 200 {object} core.ContractVerify
// @Failure 400 {object} utils.ResponseWithRequestId
// @Failure 500 {object} utils.ResponseWithRequestId
// @Router /api/v1/contractverify/{id} [get]
func GetContractVerify(c *gin.Context) {
	app := utils.Gin{C: c}

	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)

	result, resp := core.GetContractVerifyByID(c.Request.Context(), id)
	if resp != nil {
		app.HTTPResponse(resp.HttpCode, resp.Response)
		return
	}

	app.HTTPResponseOK(result)
}

// ListCompileVersion godoc
// @Description list compile version
// @Tags DATA-INFRA-API-External-V1
// @Accept application/json,json
// @Produce application/json,json
// @Param id path int true "id"
// @Success 200 {object} core.CompileVersionList
// @Failure 400 {object} utils.ResponseWithRequestId
// @Failure 500 {object} utils.ResponseWithRequestId
// @Router /api/v1/complieversions [get]
func ListCompileVersion(c *gin.Context) {
	app := utils.Gin{C: c}

	result, resp := core.ListCompileVersion(c.Request.Context())
	if resp != nil {
		app.HTTPResponse(resp.HttpCode, resp.Response)
		return
	}

	app.HTTPResponseOK(result)
}

// ContractIsVerify godoc
// @Description get contract is verify
// @Tags DATA-INFRA-API-External-V1
// @Accept application/json,json
// @Produce application/json,json
// @Param id path int true "id"
// @Success 200 {object} core.ContractIsVerify
// @Failure 400 {object} utils.ResponseWithRequestId
// @Failure 500 {object} utils.ResponseWithRequestId
// @Router /api/v1/contract/{address}/is_verify [get]
func ContractIsVerify(c *gin.Context) {
	app := utils.Gin{C: c}
	validate := validator.New()

	address := c.Param("address")
	if err := validate.Var(address, "required"); err != nil {
		app.HTTPResponse(http.StatusOK, utils.NewResponse(utils.CodeBadRequest, err.Error(), nil))
	}

	result, resp := core.GetContractIsVerify(c.Request.Context(), strings.ToLower(address))
	if resp != nil {
		app.HTTPResponse(resp.HttpCode, resp.Response)
		return
	}

	app.HTTPResponseOK(result)
}

// ListTXNs godoc
// @Description List transactions
// @Tags DATA-INFRA-API-External-V1
// @Accept application/json,json
// @Produce application/json,json
// @Param ListQuery query core.ListQuery true "ListQuery"
// @Success 200 {object} nil
// @Failure 400 {object} utils.ResponseWithRequestId
// @Failure 500 {object} utils.ResponseWithRequestId
// @Router /api/v1/txns [get]
func ListTXNs(c *gin.Context) {
	app := utils.Gin{C: c}

	var r core.ListQuery
	if err := c.ShouldBindQuery(&r); err != nil {
		app.HTTPResponse(http.StatusOK, utils.NewResponse(utils.CodeBadRequest, err.Error(), nil))
		return
	}

	if err := r.ListValidate(); err != nil {
		app.HTTPResponse(http.StatusOK, utils.NewResponse(utils.CodeBadRequest, err.Error(), nil))
		return
	}

	result, resp := core.ListTXNs(c.Request.Context(), &r)
	if resp != nil {
		app.HTTPResponse(resp.HttpCode, resp.Response)
		return
	}

	app.HTTPResponseOK(result)
}

// GetTXN godoc
// @Description Get transaction
// @Tags DATA-INFRA-API-External-V1
// @Accept application/json,json
// @Produce application/json,json
// @Param txnHash path string true "txnHash"
// @Success 200 {object} nil
// @Failure 400 {object} utils.ResponseWithRequestId
// @Failure 500 {object} utils.ResponseWithRequestId
// @Router /api/v1/txn/{txnHash} [get]
func GetTXN(c *gin.Context) {
	app := utils.Gin{C: c}
	validate := validator.New()

	txnHash := c.Param("txnHash")
	if err := validate.Var(txnHash, "required"); err != nil {
		app.HTTPResponse(http.StatusOK, utils.NewResponse(utils.CodeBadRequest, err.Error(), nil))
	}

	result, resp := core.GetTXN(c.Request.Context(), strings.ToLower(txnHash))
	if resp != nil {
		app.HTTPResponse(resp.HttpCode, resp.Response)
		return
	}

	app.HTTPResponseOK(result)
}

// GetBlock godoc
// @Description Get block detail
// @Tags DATA-INFRA-API-External-V1
// @Accept application/json,json
// @Produce application/json,json
// @Param height path string true "height"
// @Success 200 {object} busi.EVMBlockHeader
// @Failure 400 {object} utils.ResponseWithRequestId
// @Failure 500 {object} utils.ResponseWithRequestId
// @Router /api/v1/block/{height} [get]
func GetBlock(c *gin.Context) {
	app := utils.Gin{C: c}
	validate := validator.New()

	height := c.Param("height")
	if err := validate.Var(height, "required"); err != nil {
		app.HTTPResponse(http.StatusOK, utils.NewResponse(utils.CodeBadRequest, err.Error(), nil))
	}

	result, resp := core.GetBlock(c.Request.Context(), height)
	if resp != nil {
		app.HTTPResponse(resp.HttpCode, resp.Response)
		return
	}

	app.HTTPResponseOK(result)
}
