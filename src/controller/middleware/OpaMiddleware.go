package middleware

import (
	"encoding/json"
	"go-app/src/logger"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Input struct {
	Method string   `json:"method"`
	Path   []string `json:"path"`
	Role   string   `json:"role"`
}
type OPAInput struct {
	Input Input `json:"input"`
}

type OPAResponse struct {
	Result bool `json:"result"`
}

type opaMiddlewareFactory struct {
	OPAInput
	OPAResponse
	Input
}

func NewOpaMiddlewareFactory() *opaMiddlewareFactory {
	return &opaMiddlewareFactory{}
}

func (s *opaMiddlewareFactory) OPAMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		input := OPAInput{
			Input: Input{
				Method: c.Request.Method,
				Path:   splitPath(c.Request.URL.Path),
				Role:   c.GetHeader("X-User"),
			},
		}
		logger.Info("input", zap.Any("input", input))
		inputBytes, err := json.Marshal(input)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		opaReq, err := http.NewRequest("POST", "http://localhost:8181/v1/data/authz/allow", strings.NewReader(string(inputBytes)))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		opaReq.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(opaReq)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer resp.Body.Close()

		var opaResp OPAResponse
		err = json.NewDecoder(resp.Body).Decode(&opaResp)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		logger.Info("opaResp", zap.Any("opaResp", opaResp))
		if !opaResp.Result {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		c.Next()
	}
}
func splitPath(path string) []string {
	cleanedPath := strings.Trim(path, "/")
	if cleanedPath == "" {
		return []string{}
	}
	return strings.Split(cleanedPath, "/")
}
