package main

import (
	"log"
	"net/url"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
	"github.com/swaggo/echo-swagger"

	_productHttpDelivery "testPackage/product/delivery/http"
	_productHttpDeliveryMiddleware "testPackage/product/delivery/http/middleware"
	_productRepo "testPackage/product/repository/memory"
	_productUcase "testPackage/product/usecase"

	_ "testPacakge/docs"
	"testPackage/domain"
)

func init() {
	viper.SetConfigFile(`config.json`)
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	if viper.GetBool(`debug`) {
		log.Println("Service RUN on DEBUG mode")
	}
}

// @title Swagger Example API
// @version 1.0
// @description This is a sample server Petstore server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:9090
// @BasePath /
func main() {
	val := url.Values{}
	val.Add("parseTime", "1")
	val.Add("loc", "Asia/Jakarta")

	e := echo.New()
	middL := _productHttpDeliveryMiddleware.InitMiddleware()
	e.Use(middL.CORS)

	var products []domain.Product
	ar := _productRepo.NewMemoryProductRepository(products, 0)

	timeoutContext := time.Duration(viper.GetInt("context.timeout")) * time.Second
	au := _productUcase.NewProductUsecase(ar, timeoutContext)
	_productHttpDelivery.NewProductHandler(e, au)

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	log.Fatal(e.Start(viper.GetString("server.address")))
}
