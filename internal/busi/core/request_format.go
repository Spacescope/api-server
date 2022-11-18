package core

import (
	"errors"

	"api-server/pkg/models/busi"
)

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

type SubmitContractVerifyRequest struct {
	// TODO need support more CompilerType, only implement single file now.
	CompilerType    int    `form:"compiler_type" json:"compiler_type" binding:"required,oneof=1 2 3" desc:"1-single file 2-multi part 3-jsoninput"`
	CompilerVersion string `from:"compiler_version" json:"compiler_version" binding:"required,semver"`
	LicenseType     string `from:"license_type" json:"license_type" binding:"required"`
	IsOptimization  bool   `from:"is_optimization" json:"is_optimization"`
	SourceCode      string `from:"source_code" json:"source_code"`
	// TODO need more field support more CompilerType
	Runs       int    `from:"runs" json:"runs" binding:"required"`
	EVMVersion string `from:"evm_version" json:"evm_version"`
}

func (s *SubmitContractVerifyRequest) Validate() error {
	if s.CompilerType == busi.CompilerTypeSingleFile && s.SourceCode == "" {
		return errors.New("source code can not empty")
	}
	return nil
}
