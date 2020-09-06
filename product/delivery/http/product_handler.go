package http

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	validator "gopkg.in/go-playground/validator.v9"

	"testPackage/domain"
)

// ResponseError represent the reseponse error struct
type ResponseError struct {
	Message string `json:"message"`
}

// ProductHandler  represent the httphandler for product
type ProductHandler struct {
	ProductUsecase domain.ProductUsecase
}

// NewProductHandler will initialize the products/ resources endpoint
func NewProductHandler(e *echo.Echo, us domain.ProductUsecase) {
	handler := &ProductHandler{
		ProductUsecase: us,
	}
	e.GET("/products", handler.FetchProduct)
	e.POST("/products", handler.Store)
	e.POST("/products/:id", handler.Update)
	e.GET("/products/:id", handler.GetByID)
	e.DELETE("/products/:id", handler.Delete)
}

// FetchProduct godoc
// @Summary Fetch all product
// @Description get string by ID
// @ID get-string-by-int
// @Accept  json
// @Produce  json
// @Param num body int false "Limit"
// @Param cursor body string false "next data"
// @Success 200 {object} domain.Product
// @Router /products [get]
func (a *ProductHandler) FetchProduct(c echo.Context) error {
	numS := c.QueryParam("num")
	num, _ := strconv.Atoi(numS)
	cursor := c.QueryParam("cursor")
	ctx := c.Request().Context()

	listAr, nextCursor, err := a.ProductUsecase.Fetch(ctx, cursor, int64(num))
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}

	c.Response().Header().Set(`X-Cursor`, nextCursor)
	return c.JSON(http.StatusOK, listAr)
}

// GetByID will get product by given id
// GBetById godoc
// @Summary Get product by id
// @Description get string by ID
// @ID get-string-by-int
// @Accept  json
// @Produce  json
// @Param id path int false "product ID"
// @Success 200 {object} domain.Product
// @Router /products/{id} [get]
func (a *ProductHandler) GetByID(c echo.Context) error {
	idP, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusNotFound, domain.ErrNotFound.Error())
	}

	id := int64(idP)
	ctx := c.Request().Context()

	art, err := a.ProductUsecase.GetByID(ctx, id)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, art)
}

func isRequestValid(m *domain.Product) (bool, error) {
	validate := validator.New()
	err := validate.Struct(m)
	if err != nil {
		return false, err
	}
	return true, nil
}

// Store will store the product by given request body
// Store godoc
// @Summary create a new product
// @Description create new product from request
// @ID get-string-by-int
// @Accept  json
// @Produce  json
// @Param title body string true "Title of the product"
// @Param content body string true "Content of the product"
// @Success 200 {object} domain.Product
// @Router /products [post]
func (a *ProductHandler) Store(c echo.Context) (err error) {
	var product domain.Product
	err = c.Bind(&product)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	var ok bool
	if ok, err = isRequestValid(&product); !ok {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	ctx := c.Request().Context()
	err = a.ProductUsecase.Store(ctx, &product)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}

	return c.JSON(http.StatusCreated, product)
}

// Delete will delete product by given param
// Delete godoc
// @Summary delete a new product
// @Description delete the created product
// @ID get-string-by-int
// @Accept  json
// @Produce  json
// @Param id path int true "ID of the product"
// @Success 200
// @Router /products [delete]
func (a *ProductHandler) Delete(c echo.Context) error {
	idP, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusNotFound, domain.ErrNotFound.Error())
	}

	id := int64(idP)
	ctx := c.Request().Context()

	err = a.ProductUsecase.Delete(ctx, id)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}

	return c.NoContent(http.StatusNoContent)
}

// Update will update the product by given request body
// Update godoc
// @Summary update a product
// @Description update the product from given request
// @ID get-string-by-int
// @Accept  json
// @Produce  json
// @Param title body string true "Title of the product"
// @Param content body string true "Content of the product"
// @Success 200 {object} domain.Product
// @Router /products/{id} [post]
func (a *ProductHandler) Update(c echo.Context) (err error) {
	idP, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusNotFound, domain.ErrNotFound.Error())
	}
	var product domain.Product
	id := int64(idP)

	err = c.Bind(&product)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	var ok bool
	if ok, err = isRequestValid(&product); !ok {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	ctx := c.Request().Context()
	product.ID = id
	err = a.ProductUsecase.Update(ctx, &product)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, product)
}

func getStatusCode(err error) int {
	if err == nil {
		return http.StatusOK
	}

	logrus.Error(err)
	switch err {
	case domain.ErrInternalServerError:
		return http.StatusInternalServerError
	case domain.ErrNotFound:
		return http.StatusNotFound
	case domain.ErrConflict:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}
