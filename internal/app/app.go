package app

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"slices"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/veapach/golang-wallet-flow/internal/closer"
	"github.com/veapach/golang-wallet-flow/internal/config"
)

type App struct {
	DIContainer *DIContainer
	httpServer  *http.Server
}

func New(cfg config.Config) *App {
	a := &App{
		DIContainer: NewDIContainer(&cfg),
	}
	a.initDeps()
	return a
}

func (a *App) initDeps() {
	inits := []func(){
		a.initHTTPServer,
	}
	for _, fn := range inits {
		fn()
	}
}

func (a *App) initHTTPServer() {
	gin.SetMode(gin.DebugMode)
	r := gin.New()

	err := r.SetTrustedProxies([]string{
		"127.0.0.1",
		"10.0.0.0/8",
		"172.17.0.2/12",
		"192.168.0.0/16",
	})

	if err != nil {
		slog.Error("не удалось установить trusted proxies", "err", err)
		os.Exit(1)
	}

	allowedOrigins := []string{
		"http://localhost:8080",
	}

	corsConfig := cors.Config{
		AllowOriginFunc: func(origin string) bool {
			if slices.Contains(allowedOrigins, origin) {
				return true
			}
			return false
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization", "Accept", "X-Requested-With", "Cookie", "X-Request-ID"},
		ExposeHeaders:    []string{"Content-Length", "Authorization", "X-Request-ID", "Content-Type", "Content-Disposition", "Cache-Control", "ETag", "Last-Modified", "Accept-Ranges"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}
	r.Use(cors.New(corsConfig))

	r.Use(gin.Recovery(), gin.Logger())

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	a.httpServer = &http.Server{
		Addr:    a.DIContainer.cfg.HTTPAddr,
		Handler: r,
	}
}

func (a *App) Run() error {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	slog.Info("сервер запущен", "addr", a.DIContainer.cfg.HTTPAddr)

	go func() {
		if err := a.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("ошибка сервера", "err", err)
		}
	}()

	<-ctx.Done()
	slog.Info("получен сигнал, завершаем...")

	// Паттер "двойной Ctrl + C": снимаем custom handler. Второй Ctrl + C сразу убъёт процесс
	stop()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer shutdownCancel()

	if err := a.httpServer.Shutdown(shutdownCtx); err != nil {
		slog.Error("ошибка при остановке сервера", "err", err)
	}

	slog.Info("сервер остановлен")

	closerCtx, closerCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer closerCancel()

	if err := closer.CloseAll(closerCtx); err != nil {
		slog.Error("ошибки при закрытии ресурсов", "err", err)
	}

	return nil
}
