package usecase

import (
	"context"
	"time"

	//"github.com/sirupsen/logrus"
	//"golang.org/x/sync/errgroup"

	"testPackage/domain"
)

type productUsecase struct {
	productRepo    domain.ProductRepository
	contextTimeout time.Duration
}

// NewProductUsecase will create new an productUsecase object representation of domain.ProductUsecase interface
func NewProductUsecase(a domain.ProductRepository, timeout time.Duration) domain.ProductUsecase {
	return &productUsecase{
		productRepo:    a,
		contextTimeout: timeout,
	}
}

func (a *productUsecase) Fetch(c context.Context, cursor string, num int64) (res []domain.Product, nextCursor string, err error) {
	if num == 0 {
		num = 10
	}

	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()

	res, nextCursor, err = a.productRepo.Fetch(ctx, cursor, num)
	if err != nil {
		return nil, "", err
	}

	return
}

func (a *productUsecase) GetByID(c context.Context, id int64) (res domain.Product, err error) {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()

	res, err = a.productRepo.GetByID(ctx, id)
	if err != nil {
		return domain.Product{}, err
	}

	return
}

func (a *productUsecase) Update(c context.Context, ar *domain.Product) (err error) {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()

	ar.UpdatedAt = time.Now()
	return a.productRepo.Update(ctx, ar)
}

func (a *productUsecase) GetByTitle(c context.Context, title string) (res domain.Product, err error) {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()
	res, err = a.productRepo.GetByTitle(ctx, title)
	if err != nil {
		return domain.Product{}, err
	}

	return
}

func (a *productUsecase) Store(c context.Context, m *domain.Product) (err error) {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()
	existedProduct, _ := a.GetByTitle(ctx, m.Title)
	if existedProduct != (domain.Product{}) {
		return domain.ErrConflict
	}

	err = a.productRepo.Store(ctx, m)
	return
}

func (a *productUsecase) Delete(c context.Context, id int64) (err error) {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()
	existedProduct, err := a.productRepo.GetByID(ctx, id)
	if err != nil {
		return
	}
	if existedProduct == (domain.Product{}) {
		return domain.ErrNotFound
	}
	return a.productRepo.Delete(ctx, id)
}
