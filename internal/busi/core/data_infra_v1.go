package core

import (
	"api-server/pkg/models/busi"
	"api-server/pkg/utils"
	"context"
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"

	log "github.com/sirupsen/logrus"
)

func ListCirculatingSupply(ctx context.Context, c *gin.Context, r *APIRequest) (interface{}, *utils.BuErrorResponse) {
	circulating_supply := make([]*busi.CirculatingSupply, 0)
	return circulating_supply, busiSQLExecute(r, &circulating_supply)
}

func busiSQLExecute(r *APIRequest, rowsSlicePtr interface{}) *utils.BuErrorResponse {
	v := reflect.ValueOf(rowsSlicePtr)
	if v.Kind() != reflect.Ptr || reflect.Indirect(v).Kind() != reflect.Slice {
		log.Errorf("needs a pointer to a slice, v.Kind() = %v, reflect.Indirect(v).Kind() = %v", v.Kind(), reflect.Indirect(v).Kind())
		return nil
	}

	if err := utils.EngineGroup[utils.DB].Where("stat_date between ? and ?", r.StartDate, r.EndDate).Find(rowsSlicePtr); err != nil {
		log.Errorf("Execute sql error: %v", err)
		return &utils.BuErrorResponse{HttpCode: http.StatusInternalServerError, Response: utils.ErrDataInfraAPIInternal}
	}
	return nil
}
