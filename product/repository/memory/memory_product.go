package memory

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"

	"testPackage/domain"
	"testPackage/product/repository"
)

type memoryProductRepository struct {
	products []domain.Product
	lastId   int64
}

type logProductData struct {
	OldData domain.Product
	NewData domain.Product
}

// NewMemoryProductRepository will create an object that represent the product.Repository interface
func NewMemoryProductRepository(products []domain.Product, lastId int64) domain.ProductRepository {
	return &memoryProductRepository{products, lastId}
}

func (m *memoryProductRepository) log(oldData domain.Product, newData domain.Product) error {
	// Log the data
	logProduct := logProductData{oldData, newData}
	jsonLogProduct, err := json.Marshal(logProduct)
	if err != nil {
		return err
	}

	logrus.Info(fmt.Println(string(jsonLogProduct)))

	return nil
}

func (m *memoryProductRepository) fetch(ctx context.Context, args ...interface{}) (result []domain.Product, err error) {
	result = make([]domain.Product, 0)
	for _, product := range m.products {
		result = append(result, product)
	}

	return result, nil
}

func (m *memoryProductRepository) Fetch(ctx context.Context, cursor string, num int64) (res []domain.Product, nextCursor string, err error) {
	decodedCursor, err := repository.DecodeCursor(cursor)
	if err != nil && cursor != "" {
		return nil, "", domain.ErrBadParamInput
	}

	res, err = m.fetch(ctx, decodedCursor, num)
	if err != nil {
		return nil, "", err
	}

	if len(res) == int(num) {
		nextCursor, err = repository.EncodeCursor(res[len(res)-1].ID)
		if err != nil {
			return nil, "", err
		}
	}

	return
}
func (m *memoryProductRepository) GetByID(ctx context.Context, id int64) (res domain.Product, err error) {
	found := false
	for _, product := range m.products {
		if product.ID == id {
			res = product
			found = true
			break;
		}
	}

	if found != true {
		return domain.Product{}, domain.ErrNotFound
	}

	return
}

func (m *memoryProductRepository) GetByTitle(ctx context.Context, title string) (res domain.Product, err error) {
	found := false
	for _, product := range m.products {
		if product.Title == title {
			res = product
			found = true
			break;
		}
	}

	if found != true {
		return domain.Product{}, domain.ErrNotFound
	}

	return
}

func (m *memoryProductRepository) Store(ctx context.Context, a *domain.Product) (err error) {
	m.lastId += 1
	now := time.Now()
	oldData := domain.Product{}
	product := domain.Product{
		ID:        m.lastId,
		Title:     a.Title,
		Content:   a.Content,
		CreatedAt: now,
		UpdatedAt: now,
	}
	m.products = append(m.products, product)
	a.ID = m.lastId

	// Log the created product
	err = m.log(oldData, product)

	return
}

func (m *memoryProductRepository) Delete(ctx context.Context, id int64) (err error) {
	oldData := domain.Product{}
	for i, product := range m.products {
		if product.ID == id {
			oldData = product
			m.products = append(m.products[:i], m.products[i+1:]...)
			break
		}
	}

	// Log the deleted product
	err = m.log(oldData, domain.Product{})
	return
}
func (m *memoryProductRepository) Update(ctx context.Context, ar *domain.Product) (err error) {
	now := time.Now()
	oldData := domain.Product{}
	newData := domain.Product{}
	for _, product := range m.products {
		if product.ID == ar.ID {
			oldData = product
			product.Title = ar.Title
			product.Content = ar.Content
			product.UpdatedAt = now
			newData = product
			break
		}
	}
	// Log the updated product
	err = m.log(oldData, newData)
	return
}
