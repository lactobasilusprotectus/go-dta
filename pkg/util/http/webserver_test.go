package http

import (
	"github.com/gin-gonic/gin"
	"go-dts/pkg/common/config"
	"testing"
	"time"
)

type testHandler struct{}

func (t *testHandler) Register(g *gin.Engine) {}

func TestHttpServer_ok(t *testing.T) {
	srv := NewServer(config.HttpConfig{
		Port:    "0",
		TimeOut: 0,
	})

	srv.RegisterHandler(&testHandler{})
	srv.Run("local")
	time.Sleep(time.Millisecond * 300)
	srv.Stop()
}
