package gf

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gofly/utils/results"
)

func ValidateRequiredFields(c *gin.Context, param map[string]interface{}, fields []string) bool {
	for _, field := range fields {
		val, ok := param[field]
		if !ok || val == nil || val == "" {
			results.Failed(c, fmt.Sprintf("字段 '%s' 不能为空", field), nil)
			return false
		}
	}
	return true
}
