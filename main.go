package main

import (
	"crypto/tls"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/olongfen/gorm-gin-admin/docs"
	"github.com/olongfen/gorm-gin-admin/src/models"
	"github.com/olongfen/gorm-gin-admin/src/pkg/gredis"
	"github.com/olongfen/gorm-gin-admin/src/pkg/setting"
	"github.com/olongfen/gorm-gin-admin/src/router"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"
)

var (
	engine *gin.Engine
)

func init() {
	// 初始化配置文件
	setting.InitConfig()
	// 初始化模型
	models.InitModel()
	// 初始化redis
	gredis.InitRedisInstance()
	// 初始化路由
	engine = router.InitRouter()
}

func main() {

	go func() {
		// 开启服务
		s := &http.Server{
			Addr:           setting.Setting.ServerAddr + ":" + setting.Setting.ServerPort,
			Handler:        engine,
			ReadTimeout:    60 * time.Second,
			WriteTimeout:   60 * time.Second,
			MaxHeaderBytes: 1 << 20, // 10M
		}
		logrus.Println("server listen on: ", s.Addr)
		if setting.Setting.IsTLS { // 开启tls
			TLSConfig := &tls.Config{
				MinVersion:               tls.VersionTLS11,
				CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
				PreferServerCipherSuites: true,
				CipherSuites: []uint16{
					tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
					tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
					tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
					tls.TLS_RSA_WITH_AES_256_CBC_SHA,
				},
			}

			TLSProto := make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0)

			s.TLSConfig = TLSConfig
			s.TLSNextProto = TLSProto

			if err := s.ListenAndServeTLS(setting.Setting.TLS.Cert, setting.Setting.TLS.Key); err != nil {
				logrus.Fatal(err)
			}
		} else {
			if err := s.ListenAndServe(); err != nil {
				logrus.Fatal(err)
			}
		}
	}()

	var state int32 = 1
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
EXIT:
	for {
		sig := <-sc
		fmt.Println("获取到信号[%s]", sig.String())
		switch sig {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			atomic.StoreInt32(&state, 0)
			break EXIT
		case syscall.SIGHUP:
		default:
			break EXIT
		}
	}

	fmt.Println("服务退出")
	time.Sleep(time.Second)
	os.Exit(int(atomic.LoadInt32(&state)))
}
