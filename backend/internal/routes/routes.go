package router

import (
	"ecommerce-platform/internal/database"
	"ecommerce-platform/internal/handlers"
	"ecommerce-platform/internal/middleware"
	"net/http"
	"strconv"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	httpRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests.",
		},
		[]string{"method", "route", "status_code"},
	)

	httpDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_duration_seconds",
			Help:    "Histogram of HTTP request durations.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "route"},
	)
)

func init() {
	prometheus.MustRegister(httpRequests)
	prometheus.MustRegister(httpDuration)
}

func RegisterRoutes(client database.Client, r *gin.Engine) {
	// Роут для сбора метрик
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	api := r.Group("/api")
	{
		api.POST("/login", handlers.Login(client.DB))
		api.POST("/logout/:session_id", handlers.Logout(client.DB))

		api.POST("/users", handlers.CreateUser(client))            // Регистрация обычного пользователя
		api.POST("/admin/users", handlers.CreateAdminUser(client)) // Создание пользователя с ролью "admin"

		api.GET("/users/:id", func(c *gin.Context) {
			userIDStr := c.Param("id")
			userID, err := strconv.ParseUint(userIDStr, 10, 32)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
				return
			}

			user, err := handlers.GetUserByID(client, uint(userID))
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "User not found"})
				return
			}

			c.JSON(http.StatusOK, user)
		}) // Получить пользователя

		api.PUT("/users/:id", middleware.AuthMiddleware(&client), handlers.UpdateUser(client))    // Обновить пользователя
		api.DELETE("/users/:id", middleware.AuthMiddleware(&client), handlers.DeleteUser(client)) // Удалить пользователя
		api.GET("/users", middleware.AuthMiddleware(&client), middleware.RoleMiddleware("admin"), handlers.GetAllUsers(client))

		// Роуты для изображений
		api.GET("/products/:product_id/images", handlers.GetProductImages(client))    // Получить все изображения
		api.GET("/reviews/product/:product_id", handlers.GetReviewsByProduct(client)) // Получить отзывы по продукту

		// Роуты для продуктов
		api.GET("/products", handlers.GetAllProducts(client))             // Получить все продукты
		api.GET("/products/:product_id", handlers.GetProductByID(client)) // Получить продукт по ID

		// Роуты для категорий
		api.GET("/categories", handlers.GetAllCategories(client))    // Получить все категории
		api.GET("/categories/:id", handlers.GetCategoryByID(client)) // Получить категорию по ID

		protected := api.Group("")
		protected.Use(middleware.AuthMiddleware(&client))
		{
			protected.GET("/profile", handlers.GetProfile)

			// Роуты для корзины
			protected.POST("/cart", handlers.CreateShoppingCart(client))                 // Создать корзину
			protected.GET("/cart/:user_id", handlers.GetCartByUserID(client))            // Получить корзину по cart_id
			protected.GET("/cart/:user_id/items", handlers.GetCartItemsByUserID(client)) // Получить товары в корзине
			protected.POST("/cart/:cart_id/items", handlers.AddItemToCart(client))       // Добавить товар в корзину
			protected.DELETE("/cart/items/:item_id", handlers.DeleteCartItem(client))    // Удалить товар из корзины

			// Роуты для заказов
			protected.POST("/orders", handlers.CreateOrder(client))                   // Создать заказ
			protected.GET("/orders/:id", handlers.GetOrderByID(client))               // Получить заказ по ID
			protected.GET("/orders/user/:user_id", handlers.GetOrdersByUser(client))  // Заказы пользователя
			protected.POST("/orders/items", handlers.CreateOrderItem(client))         // Создание позиции заказа
			protected.PUT("/orders/items/:item_id", handlers.UpdateOrderItem(client)) // Обновление позиции заказа
			protected.GET("/orders/:id/items", handlers.GetOrderItems(client))        // Получение всех позиций заказа

			// Роуты для отзывов
			protected.POST("/reviews", handlers.CreateReview(client))       // Создать отзыв
			protected.DELETE("/reviews/:id", handlers.DeleteReview(client)) // Удалить отзыв

			// Роуты для кэширования
			protected.GET("/cache/:cache_key", handlers.GetProductFromCache(client)) // Получить продукт из кэша или базы данных
			protected.DELETE("/cache/:key", handlers.DeleteCache(client.DB))

			protected.POST("/payments", handlers.CreatePayment(client))       // Создать новый платеж
			protected.GET("/payments", handlers.GetPayments(client))          // Получить список всех платежей
			protected.GET("/payments/:id", handlers.GetPaymentByID(client))   // Получить платеж по ID
			protected.DELETE("/payments/:id", handlers.DeletePayment(client)) // Удалить платеж

			protected.POST("/addresses", handlers.CreateUserAddress(client))       // Создать адрес
			protected.GET("/addresses", handlers.GetUserAddresses(client))         // Получить все адреса пользователя
			protected.PUT("/addresses/:id", handlers.UpdateUserAddress(client))    // Обновить адрес
			protected.DELETE("/addresses/:id", handlers.DeleteUserAddress(client)) // Удалить адрес

			// Только администраторы
			admin := protected.Group("")
			admin.Use(middleware.RoleMiddleware("admin"))
			{
				// Роуты для продуктов
				admin.POST("/products", handlers.CreateProduct(client))       // Добавить продукт
				admin.PUT("/products/:id", handlers.UpdateProduct(client))    // Обновить продукт
				admin.DELETE("/products/:id", handlers.DeleteProduct(client)) // Удалить продукт

				// Роуты для изображений
				admin.POST("/products/:product_id/images", handlers.AddProductImage(client)) // Добавить изображение
				admin.PUT("/products/images/:id", handlers.UpdateProductImage(client))       // Обновить изображение
				admin.DELETE("/products/images/:id", handlers.DeleteProductImage(client))    // Удалить изображение

				// Роуты для категорий
				admin.POST("/categories", handlers.CreateCategory(client))       // Создать категорию
				admin.PUT("/categories/:id", handlers.UpdateCategory(client))    // Обновить категорию
				admin.DELETE("/categories/:id", handlers.DeleteCategory(client)) // Удалить категорию

				// Роуты для сессий и журналов аудита
				admin.GET("/sessions/:user_id", handlers.GetUserSessions(client)) // Получить сессии пользователя
				admin.GET("/audit-logs/:user_id", handlers.GetAuditLogs(client))  // Получить журнал аудита пользователя

			}
		}
		api.Use(func(c *gin.Context) {
			timer := prometheus.NewTimer(httpDuration.WithLabelValues(c.Request.Method, c.FullPath()))
			defer timer.ObserveDuration()

			c.Next()

			httpRequests.WithLabelValues(c.Request.Method, c.FullPath(), strconv.Itoa(c.Writer.Status())).Inc()
		})

	}
}

func SetupRouter(client database.Client) *gin.Engine {
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	RegisterRoutes(client, router)

	return router
}
