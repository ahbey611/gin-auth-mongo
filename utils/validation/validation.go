package validation

import (
	"gin-auth-mongo/utils/response"
	"regexp"

	"github.com/gin-gonic/gin"
)

// bind and validate request
// use for rest api
func BindAndValidate(c *gin.Context, request interface{}) error {
	if err := c.ShouldBind(request); err != nil {
		response.BadRequestWithMessage(c, err.Error())
		return err
	}
	if err := request.(interface{ Validate() error }).Validate(); err != nil {
		response.BadRequestWithMessage(c, err.Error())
		return err
	}
	return nil
}

func CheckDateFormat(date string) bool {
	return regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`).MatchString(date)
}
