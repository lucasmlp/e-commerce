package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/machado-br/e-commerce/domain/orders"
	"github.com/machado-br/e-commerce/domain/products"
)

type api struct {
	Debug           bool
	OrdersService   orders.Service
	ProductsService products.Service
}

func NewApi(
	debug bool,
	orderService orders.Service,
	productsService products.Service,
) (api, error) {
	return api{
		Debug:           debug,
		OrdersService:   orderService,
		ProductsService: productsService,
	}, nil
}

func (a api) Engine() *gin.Engine {
	router := gin.Default()

	root := router.Group("")
	{
		root.GET("/ping", func(c *gin.Context) {
			c.JSON(http.StatusOK, "pong")
		})
	}
	ordersRoot := router.Group("/orders")
	{
		ordersRoot.GET("", a.getAllOrders)
		ordersRoot.GET("/:id", a.getOrderById)
		ordersRoot.POST("", a.createOrder)
		ordersRoot.DELETE(":id", a.deleteOrder)
		ordersRoot.PUT("", a.updateOrder)
	}

	productsRoot := router.Group("/products")
	{
		productsRoot.GET("", a.getAllProducts)
		productsRoot.GET("/:id", a.getProductById)
		productsRoot.POST("", a.createProduct)
		productsRoot.DELETE(":id", a.deleteProduct)
		productsRoot.PUT("", a.updateProduct)
	}

	return router
}

func (a api) Run() {

	router := a.Engine()
	router.Run()
}
