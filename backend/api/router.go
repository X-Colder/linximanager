package api

import (
	"github.com/gin-gonic/gin"
	"github.com/linximanager/backend/internal/config"
	"github.com/linximanager/backend/internal/handler"
	adminH "github.com/linximanager/backend/internal/handler/admin"
	consumerH "github.com/linximanager/backend/internal/handler/consumer"
	merchantH "github.com/linximanager/backend/internal/handler/merchant"
	"github.com/linximanager/backend/internal/middleware"
	"github.com/linximanager/backend/internal/repository"
	"github.com/linximanager/backend/internal/service"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB, rdb *redis.Client, cfg *config.Config) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.Logger())
	r.Use(middleware.CORS())
	r.Use(middleware.RateLimit(cfg.RateLimit.QPS, cfg.RateLimit.Burst))

	// 初始化 Repositories
	userRepo := repository.NewUserRepo(db)
	merchantRepo := repository.NewMerchantRepo(db)
	productRepo := repository.NewProductRepo(db)
	inventoryRepo := repository.NewInventoryRepo(db)
	orderRepo := repository.NewOrderRepo(db)

	// 初始化 Services
	authSvc := service.NewAuthService(userRepo, cfg)
	merchantSvc := service.NewMerchantService(merchantRepo, userRepo)
	productSvc := service.NewProductService(productRepo)
	inventorySvc := service.NewInventoryService(inventoryRepo)
	orderSvc := service.NewOrderService(orderRepo, inventoryRepo)
	aiSvc := service.NewAIService(inventoryRepo, productRepo)

	// 初始化 Handlers
	authHandler := handler.NewAuthHandler(authSvc)
	adminMerchantHandler := adminH.NewMerchantHandler(merchantSvc)
	adminDashboardHandler := adminH.NewDashboardHandler(merchantSvc)
	merchantDashboardHandler := merchantH.NewDashboardHandler(merchantSvc, orderSvc, inventorySvc)
	merchantProductHandler := merchantH.NewProductHandler(productSvc)
	merchantInventoryHandler := merchantH.NewInventoryHandler(inventorySvc)
	merchantAIHandler := merchantH.NewAIHandler(aiSvc)
	merchantOrderHandler := merchantH.NewOrderHandler(orderSvc)
	shopHandler := consumerH.NewShopHandler(merchantSvc, productSvc)
	cartHandler := consumerH.NewCartHandler(rdb)
	consumerOrderHandler := consumerH.NewOrderHandler(orderSvc)

	v1 := r.Group("/api/v1")

	// ==================== 认证（公开） ====================
	auth := v1.Group("/auth")
	{
		auth.POST("/login/phone", authHandler.LoginByPhone)
		auth.POST("/login/wechat", authHandler.LoginByWechat)
		auth.POST("/refresh", authHandler.RefreshToken)
		auth.POST("/logout", middleware.Auth(cfg.JWT.AccessSecret), authHandler.Logout)
	}

	// ==================== 平台管理端 ====================
	admin := v1.Group("/admin",
		middleware.Auth(cfg.JWT.AccessSecret),
		middleware.RequireRole("admin"),
	)
	{
		admin.GET("/dashboard", adminDashboardHandler.Overview)
		admin.GET("/merchants", adminMerchantHandler.List)
		admin.GET("/merchants/:id", adminMerchantHandler.GetByID)
		admin.PUT("/merchants/:id/audit", adminMerchantHandler.Audit)
		admin.PUT("/merchants/:id/freeze", adminMerchantHandler.Freeze)
	}

	// ==================== 商家端 ====================
	mc := v1.Group("/merchant",
		middleware.Auth(cfg.JWT.AccessSecret),
		middleware.RequireRole("merchant"),
	)
	{
		// 工作台
		mc.GET("/dashboard", merchantDashboardHandler.Overview)
		mc.GET("/alerts", merchantDashboardHandler.Alerts)
		mc.GET("/todos", merchantDashboardHandler.Todos)

		// 商品管理
		mc.GET("/products", merchantProductHandler.List)
		mc.POST("/products", merchantProductHandler.Create)
		mc.PUT("/products/:id", merchantProductHandler.Update)
		mc.DELETE("/products/:id", merchantProductHandler.Delete)
		mc.GET("/products/:id/bom", merchantProductHandler.GetBOM)
		mc.PUT("/products/:id/bom", merchantProductHandler.SetBOM)

		// 库存
		mc.GET("/inventory", merchantInventoryHandler.List)
		mc.POST("/inventory/purchase", merchantInventoryHandler.Purchase)
		mc.POST("/inventory/stocktake", merchantInventoryHandler.Stocktake)

		// AI决策
		mc.GET("/ai/replenishment", merchantAIHandler.Replenishment)
		mc.POST("/ai/replenishment/confirm", merchantAIHandler.ConfirmReplenishment)
		mc.GET("/ai/promotions", merchantAIHandler.PromotionSuggestions)
		mc.POST("/ai/promotions/execute", merchantAIHandler.ExecutePromotion)
		mc.POST("/ai/chat", merchantAIHandler.Chat)

		// 订单
		mc.GET("/orders", merchantOrderHandler.List)
		mc.GET("/orders/:id", merchantOrderHandler.GetByID)
		mc.POST("/orders/verify", merchantOrderHandler.Verify)
	}

	// ==================== C端（需要登录） ====================
	consumerAuth := v1.Group("",
		middleware.Auth(cfg.JWT.AccessSecret),
	)
	{
		// 购物车
		consumerAuth.GET("/cart", cartHandler.GetCart)
		consumerAuth.POST("/cart/items", cartHandler.AddItem)
		consumerAuth.PUT("/cart/items/:id", cartHandler.UpdateItem)
		consumerAuth.DELETE("/cart/items/:id", cartHandler.DeleteItem)

		// 订单
		consumerAuth.POST("/orders", consumerOrderHandler.Create)
		consumerAuth.GET("/orders", consumerOrderHandler.List)
		consumerAuth.GET("/orders/:id", consumerOrderHandler.GetByID)
		consumerAuth.POST("/orders/:id/pay", consumerOrderHandler.Pay)
		consumerAuth.POST("/orders/:id/cancel", consumerOrderHandler.Cancel)
	}

	// ==================== C端（公开） ====================
	v1.GET("/shop/:merchant_id", shopHandler.GetShop)
	v1.GET("/shop/:merchant_id/products", shopHandler.ListProducts)
	v1.GET("/shop/:merchant_id/products/:id", shopHandler.GetProduct)

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	return r
}
