// @title           MaqamDetector API
// @version         1.0.0
// @description     API pendeteksi maqam musik Arab & Timur Tengah untuk komunitas banjari.

// @contact.name   Anas (Cypress Consulting)

// @license.name  ISC

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.apikey  ApiKeyAuth
// @in                          header
// @name                        Authorization

package main

import (
	"api/config"
	_ "api/docs"
	"api/handler"
	"api/middleware"
	"api/model"
	"api/repository"
	"api/seeders"
	"api/service"
	"api/utility"
	"os"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func main() {
	config.LoadEnv()
	db := config.DBInit()

	// ── Auto-Migrate ────────────────────────────
	if err := db.AutoMigrate(
		&model.User{},
		&model.Session{},
		&model.Maqam{},
		&model.Analysis{},
		&model.AnalysisCandidate{},
	); err != nil {
		panic("Gagal auto-migrate tabel: " + err.Error())
	}

	// ── Seed Data ───────────────────────────────
	seeders.SeedMaqamat(db)

	// ── Repository ──────────────────────────────
	userRepo := repository.NewUserRepository(db)
	sessionRepo := repository.NewSessionRepository(db)
	analysisRepo := repository.NewAnalysisRepository(db)
	maqamRepo := repository.NewMaqamRepository(db)

	// ── Service ─────────────────────────────────
	userService := service.NewUserService(userRepo)
	claudeService := service.NewClaudeService()
	analyzeService := service.NewAnalyzeService(analysisRepo, sessionRepo, claudeService)
	historyService := service.NewHistoryService(analysisRepo)
	maqamService := service.NewMaqamService(maqamRepo)

	// ── Handler ─────────────────────────────────
	userHandler := handler.NewUserHandler(userService)
	analyzeHandler := handler.NewAnalyzeHandler(analyzeService)
	historyHandler := handler.NewHistoryHandler(historyService)
	maqamHandler := handler.NewMaqamHandler(maqamService)

	// ── Echo Setup ──────────────────────────────
	e := echo.New()
	e.Validator = &utility.CustomValidator{Validator: validator.New()}
	e.Use(middleware.MiddlewareLogging)
	e.Use(echoMiddleware.CORSWithConfig(echoMiddleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE, echo.OPTIONS},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, "X-Session-ID"},
	}))

	e.HTTPErrorHandler = handler.CustomHTTPErrorHandler

	// ── Rate Limiter ────────────────────────────
	rateLimiter := middleware.NewRateLimiter(100, 1*time.Hour)

	// ── Routes: Legacy (User/Auth — dari starter kit) ──
	api := e.Group("/api")
	api.POST("/users/register", userHandler.CreateUser)
	api.POST("/users/login", userHandler.LoginUser)

	auth := api.Group("")
	auth.Use(middleware.JWTAuth)
	auth.GET("/users", userHandler.GetUser)

	// ── Routes: MaqamDetector v1 ────────────────
	v1 := e.Group("/api/v1")
	v1.Use(middleware.SessionMiddleware)

	// Analyze (rate limited)
	analyze := v1.Group("")
	analyze.Use(rateLimiter.Middleware())
	analyze.POST("/analyze/youtube", analyzeHandler.AnalyzeYoutube)
	analyze.POST("/analyze/upload", analyzeHandler.AnalyzeUpload)
	analyze.POST("/analyze/record", analyzeHandler.AnalyzeRecord)

	// Analysis polling (no rate limit)
	v1.GET("/analyses/:id", analyzeHandler.GetAnalysis)

	// History
	v1.GET("/history", historyHandler.GetHistory)
	v1.DELETE("/history/:id", historyHandler.DeleteHistory)

	// Maqam (public — no session required)
	maqamGroup := e.Group("/api/v1")
	maqamGroup.GET("/maqamat", maqamHandler.GetMaqamat)
	maqamGroup.GET("/maqamat/:id", maqamHandler.GetMaqamByID)

	// ── Swagger ─────────────────────────────────
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// ── Static Files (Uploaded Audio) ────────────
	e.Static("/uploads", "apps/api/uploads")

	// ── Start ───────────────────────────────────
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	e.Logger.Fatal(e.Start(":" + port))
}
