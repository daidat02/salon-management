package middlewares

import (
	"context"
	"errors"
	"strings"

	"github.com/daidat02/server/pkg/apperror"
	"github.com/daidat02/server/pkg/response"
	"github.com/daidat02/server/pkg/utils"
	"github.com/gin-gonic/gin"
)

type contextKey string
const OrgIdKey contextKey = contextKey("organization_id")

func VerifyToken() gin.HandlerFunc{
	return func(c *gin.Context){
		// Lấy token từ header Authorization
		authHeader := c.GetHeader("Authorization")

		if authHeader == ""{
			response.Error(c,apperror.New(401,"Thiếu token xác thực, đăng nhập để tiếp tục",nil))
			c.Abort()
			return	
		}

		paths := strings.Split(authHeader," ")
		if len(paths) != 2 || paths[0] != "Bearer"{
			response.Error(c,apperror.New(401,"Token xác thực không hợp lệ",nil))
			c.Abort()
			return
		}

		tokenString := paths[1]
		claims ,err := utils.ValidateToken(tokenString)
		if err != nil{
			response.Error(c,apperror.New(401,"Token xác thực không hợp lệ",nil))
			c.Abort()
			return
		}
		c.Set("userId", claims.UserID)
		c.Set("orgId", claims.OrgID)
		ctx := context.WithValue(c.Request.Context(), OrgIdKey, claims.OrgID)

		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}



func GetOrgID(ctx context.Context) (string, error) {
    orgID, ok := ctx.Value(OrgIdKey).(string)
    
    if !ok || orgID == "" {
        return "", errors.New("missing organization_id in context")
    }
    return orgID, nil
}