package apierror

import (
	"fmt"
	"testing"
)

func TestError(t *testing.T) {
	//ers := "{\"code\":\"FeatureTemporaryDisabled\",\"message\":\"feature upgrading\",\"requestId\":\"0bc13b0116660546641984077e4b93\",\"resultCode\":\"FeatureTemporaryDisabled\",\"display_message\":\"功能维护中，预计10月底前维护完成\"}"
	ers := "{\"code\":\"FeatureTemporaryDisabled\",\"message\":\"feature upgrading\",\"requestId\":\"0bc13b0116660546641984077e4b93\"}"

	e := ParseCommonApiError([]byte(ers))
	fmt.Println(e)
}
