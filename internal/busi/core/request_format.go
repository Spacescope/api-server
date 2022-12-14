package core

import (
	"errors"

	"api-server/pkg/models/busi"
)

type ListContractsParams struct {
	ListQuery
	Verified uint8 `form:"verified" json:"verified" desc:"0-all contracts, 1-verified contracts"`
}

var (
	ContractBreakDownOrderByTxnsAsc  = 1
	ContractBreakDownOrderByTxnsDesc = 2

	// TODO next iter
	ContractBreakDownOrderByInternalTxnsAsc  = 3
	ContractBreakDownOrderByInternalTxnsDesc = 4

	ContractBreakDownOrderByFilBurnedAsc = 5
	ContractBreakDownOrderByFilBurneDesc = 6

	ContractBreakDownOrderByUserCountAsc  = 7
	ContractBreakDownOrderByUserCountDesc = 8

	ContractBreakDownOrderByCallInAsc  = 9
	ContractBreakDownOrderByCallInDesc = 10

	ContractBreakDownOrderByCallOutAsc  = 11
	ContractBreakDownOrderByCallOutDesc = 12
)

type ListStatContractBreakdownParams struct {
	ListQuery
	OrderBy int `form:"order_by" json:"order_by" binding:"oneof=0 1 2 3 4 5 6 7 8 9 10 11 12"`
}

type ListQuery struct {
	Offset int `form:"o" json:"o"`
	Limit  int `form:"l" json:"l"`
}

func (r *ListQuery) ListValidate() error {
	if r.Offset < 0 {
		return errors.New("the o(ffset) should be greater than or equal 0")
	}

	switch r.Limit {
	case 25, 50, 100:
	default:
		return errors.New("the l(imit) should be one of 25/50/100")
	}

	return nil
}

type SourceCodePart struct {
	Filename      string `json:"filename"`
	SourceCodeUrl string `json:"source_code_url"`
}

type JsonInput struct {
	Url string `json:"url"`
}

type SubmitContractVerifyRequest struct {
	CompilerType    int               `form:"compiler_type" json:"compiler_type" binding:"required,oneof=1 2 3" desc:"1-single file 2-multi part 3-jsoninput"`
	CompilerVersion string            `from:"compiler_version" json:"compiler_version" binding:"required,semver"`
	LicenseType     string            `from:"license_type" json:"license_type" binding:"required"`
	IsOptimization  bool              `from:"is_optimization" json:"is_optimization"`
	SourceCode      string            `from:"source_code" json:"source_code"`
	SourceCodeParts []*SourceCodePart `form:"source_code_parts" json:"source_code_parts"`
	Runs            int               `from:"runs" json:"runs" binding:"required"`
	EVMVersion      string            `from:"evm_version" json:"evm_version"`
}

func (s *SubmitContractVerifyRequest) Validate() error {
	if s.CompilerType == busi.CompilerTypeSingleFile && s.SourceCode == "" {
		return errors.New("source code can not empty")
	}
	return nil
}
