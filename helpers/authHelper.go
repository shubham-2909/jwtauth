package helpers

import (
	"errors"

	"github.com/gin-gonic/gin"
)

func VerifyUserType(c *gin.Context, userType string) (err error) {
	claimedUserType := c.GetString("usertype")
	err = nil
	if userType != claimedUserType {
		err := errors.New("unauthorized access")
		return err
	}

	return err
}

func MapUserTypetoID(c *gin.Context, uid string) (err error) {
	err = nil
	claimedId := c.GetString("uid")
	claimedUserType := c.GetString("usertype")
	if claimedUserType == "USER" && claimedId != uid {
		err = errors.New("unauthorized access")
		return err
	}
	return err
}
