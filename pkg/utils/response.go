package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
)

type Response struct {
	Code    ResponseCode `json:"code"`
	Message string       `json:"message,omitempty"`
	Data    interface{}  `json:"data,omitempty"`
}

type ResponseWithRequestId struct {
	RequestId string `json:"request_id"`
	Response
}

type BuErrorResponse struct {
	HttpCode int
	*Response
}

func (res Response) Error() string {
	t, _ := json.Marshal(res)
	return string(t)
}

func NewResponse(code ResponseCode, msg string, data interface{}) *Response {
	return &Response{
		Code:    code,
		Message: msg,
		Data:    data,
	}
}

func NewResponseWithRequestId(requestId string, r *Response) *ResponseWithRequestId {
	return &ResponseWithRequestId{
		RequestId: requestId,
		Response: Response{
			Code:    r.Code,
			Message: r.Message,
			Data:    r.Data,
		},
	}
}

func ToResponse(err error) *Response {
	if err == nil {
		return OK
	}
	switch t := err.(type) {
	case *Response:
		return t
	default:
		//logrus.Error(err)
	}
	return NewResponse(CodeError, err.Error(), nil)
}

func ToOK(data interface{}) *Response {
	return NewResponse(OK.Code, OK.Message, data)
}

type ResponseCode int

const (
	CodeOk ResponseCode = iota
	CodeUnKnownReasonErr

	CodeError
	CodeInternalServer
	CodeBadRequest
	CodeNotFound
	CodeUserNotFound
	CodeForbidSendEmail

	CodeAccountRegister             = 10000
	CodeAccountRegisterParamsErr    = 10001
	CodeAccountInternalErr          = 10002
	CodeAccountNotFound             = 10003
	CodeAccountPasswordErr          = 10004
	CodeAccountUserNotFound         = 10005
	CodeAccountEmailFormatErr       = 10006
	CodeAccountEmailNotVerified     = 10007
	CodeAccountEmailRateLimit       = 10008
	CodeAccountCodeOrFingerPrintErr = 10009
	CodeAccountAlreadyExistsErr     = 10010

	CodeMailCode                          = 11000
	CodeMailCodeNotSend                   = 11001
	CodeMailCodeVerifiedFailed            = 11002
	CodeMailCodeVerifiedOK                = 11003
	CodeMailBrowserFingerPrintMetadataErr = 11004
	CodeMailSendFailed                    = 11005

	CodeSPRD                 = 20000
	CodeSPRDParamsErr        = 20001
	CodeSPRDInternalErr      = 20002
	CodeSPRDForbidden        = 20003
	CodeSPRDSPNotFound       = 20004
	CodeWatchRecordHasExists = 20005
	CodeSPRDRegionNotFound   = 20006

	CodeDataInfraAPI                = 30000
	CodeDataInfraAPIForbidden       = 30001
	CodeDataInfraAPIParamsErr       = 30002
	CodeDataInfraAPIInternalErr     = 30003
	CodeDataInfraAPIDataNotReadyErr = 30004
)

var (
	OK               = &Response{Code: CodeOk, Message: "success."}
	ErrUnKnownReason = &Response{Code: CodeUnKnownReasonErr, Message: "unknown reason."}

	ErrInternalServer = &Response{Code: CodeInternalServer, Message: "server internal error."}
	ErrBadRequest     = &Response{Code: CodeBadRequest, Message: "bad request."}
	ErrNotFound       = &Response{Code: CodeNotFound, Message: "object not found."}

	ErrAccountRegister          = &Response{Code: CodeAccountRegister, Message: ""}
	ErrAccountRegisterParams    = &Response{Code: CodeAccountRegisterParamsErr, Message: "register parameters error."}
	ErrAccountInternal          = &Response{Code: CodeAccountInternalErr, Message: "account internal server error."}
	ErrAccountNotFound          = &Response{Code: CodeAccountNotFound, Message: "account not found."}
	ErrAccountPassword          = &Response{Code: CodeAccountPasswordErr, Message: "password error."}
	ErrAccountUserNotFound      = &Response{Code: CodeAccountUserNotFound, Message: "user not found."}
	ErrAccountEmailFormat       = &Response{Code: CodeAccountEmailFormatErr, Message: "email format error."}
	ErrAccountEmailNotVerified  = &Response{Code: CodeAccountEmailNotVerified, Message: "email has not been verified."}
	ErrAccountEmailRateLimit    = &Response{Code: CodeAccountEmailRateLimit, Message: "email rate limit."}
	ErrAccountCodeOrFingerPrint = &Response{Code: CodeAccountCodeOrFingerPrintErr, Message: "finger print error."}
	ErrAccountAlreadyExists     = &Response{Code: CodeAccountAlreadyExistsErr, Message: "account already exists."}

	ErrMailCode                       = &Response{Code: CodeMailCode, Message: ""}
	ErrMailCodeNotSend                = &Response{Code: CodeMailCodeNotSend, Message: "code has not been sent."}
	ErrMailCodeVerifiedFailed         = &Response{Code: CodeMailCodeVerifiedFailed, Message: "code verified failed."}
	ErrMailCodeVerifiedOK             = &Response{Code: CodeMailCodeVerifiedOK, Message: "code verified ok."}
	ErrMailBrowserFingerPrintMetadata = &Response{Code: CodeMailBrowserFingerPrintMetadataErr, Message: "fp metadata error."}
	ErrMailSendFailed                 = &Response{Code: CodeMailSendFailed, Message: "Failed to send email."}

	ErrSPRDInternal             = &Response{Code: CodeSPRDInternalErr, Message: "sprd internal server error."}
	ErrSPRDParamsErr            = &Response{Code: CodeSPRDParamsErr, Message: "sprd parameters error."}
	ErrSPRDForBidden            = &Response{Code: CodeSPRDForbidden, Message: "sprd forbidden."}
	ErrSPRDSPNotFound           = &Response{Code: CodeSPRDSPNotFound, Message: "sp not found."}
	ErrSPRDWatchRecordHasExists = &Response{Code: CodeWatchRecordHasExists, Message: "record has exists."}
	ErrSPRDRegionNotFound       = &Response{Code: CodeSPRDRegionNotFound, Message: "region not found exists."}

	ErrDataInfraAPIForbidden    = &Response{Code: CodeDataInfraAPIForbidden, Message: "SpaceScope API forbidden."}
	ErrDataInfraAPIParamsErr    = &Response{Code: CodeDataInfraAPIParamsErr, Message: "SpaceScope API parameters error."}
	ErrDataInfraAPIInternal     = &Response{Code: CodeDataInfraAPIInternalErr, Message: "SpaceScope API internal server error."}
	ErrDataInfraAPIDataNotReady = &Response{Code: CodeDataInfraAPIDataNotReadyErr, Message: "SpaceScope API data is not available."}
)

type Gin struct {
	C *gin.Context
}

func (g *Gin) HTTPResponseOK(data interface{}) {
	requestId := g.C.Request.Header.Get("Kong-Request-ID")
	g.C.JSON(http.StatusOK, NewResponseWithRequestId(requestId, NewResponse(OK.Code, OK.Message, data)))
}

func (g *Gin) HTTPResponse204() {
	g.C.JSON(http.StatusNoContent, NewResponse(OK.Code, OK.Message, nil))
}

func (g *Gin) HTTPResponse(httpCode int, r *Response) {
	requestId := g.C.Request.Header.Get("Kong-Request-ID")
	g.C.JSON(httpCode, NewResponseWithRequestId(requestId, r))
}
