package domain

import (
	"context"
	"time"
)

// Product ...
type Product struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title" validate:"required"`
	Content   string    `json:"content" validate:"required"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedAt time.Time `json:"created_at"`
}

// ProductUsecase represent the product's usecases
type ProductUsecase interface {
	Fetch(ctx context.Context, cursor string, num int64) ([]Product, string, error)
	GetByID(ctx context.Context, id int64) (Product, error)
	Update(ctx context.Context, ar *Product) error
	GetByTitle(ctx context.Context, title string) (Product, error)
	Store(context.Context, *Product) error
	Delete(ctx context.Context, id int64) error
}

// ProductRepository represent the product's repository contract
type ProductRepository interface {
	Fetch(ctx context.Context, cursor string, num int64) (res []Product, nextCursor string, err error)
	GetByID(ctx context.Context, id int64) (Product, error)
	GetByTitle(ctx context.Context, title string) (Product, error)
	Update(ctx context.Context, ar *Product) error
	Store(ctx context.Context, a *Product) error
	Delete(ctx context.Context, id int64) error
}
