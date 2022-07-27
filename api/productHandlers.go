package api

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/machado-br/e-commerce/domain/dtos"
	"github.com/pborman/uuid"
)

func (a api) getAllProducts(c *gin.Context) {
	log.Println("GET /products")

	products, err := a.ProductsService.FindAll(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
	}
	c.JSON(http.StatusOK, products)
}

func (a api) getProductById(c *gin.Context) {
	log.Println("GET /products/:id")

	productId := uuid.Parse(c.Param("id"))
	if productId == nil {
		c.JSON(http.StatusUnprocessableEntity, "productId must be a uuid")
	} else if strings.Compare(productId.String(), uuid.NIL.String()) == 0 {
		c.JSON(http.StatusUnprocessableEntity, "productId cannot be an empty uuid")
	} else {
		product, err := a.ProductsService.Find(c, productId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
		}
		c.JSON(http.StatusOK, product)
	}
}

func (a api) createProduct(c *gin.Context) {
	log.Println("POST /products")

	var product dtos.Product
	err := c.ShouldBindJSON(&product)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
	} else if strings.Compare(product.ProductId.String(), uuid.NIL.String()) == 0 {
		c.JSON(http.StatusUnprocessableEntity, "productId must cannot be an empty uuid")
	} else {
		result, err := a.ProductsService.Create(c, product)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
		}
		c.JSON(http.StatusOK, result)
	}
}

func (a api) deleteProduct(c *gin.Context) {
	log.Println("DELETE /products")

	productId := uuid.Parse(c.Param("id"))
	if productId == nil {
		c.JSON(http.StatusUnprocessableEntity, "productId must be a uuid")
	} else if strings.Compare(productId.String(), uuid.NIL.String()) == 0 {
		c.JSON(http.StatusUnprocessableEntity, "productId must cannot be an empty uuid")
	} else {
		err := a.ProductsService.Delete(c, productId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
		} else {
			c.Status(http.StatusNoContent)
		}
	}
}

func (a api) updateProduct(c *gin.Context) {
	log.Println("PUT /products")

	var product dtos.Product
	err := c.ShouldBindJSON(&product)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, err)
	} else if strings.Compare(product.ProductId.String(), uuid.NIL.String()) == 0 {
		c.JSON(http.StatusUnprocessableEntity, "productId must cannot be an empty uuid")
	} else {
		err = a.ProductsService.Update(c, product)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
		} else {
			c.Status(http.StatusNoContent)
		}
	}
}
