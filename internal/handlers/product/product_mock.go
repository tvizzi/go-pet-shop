package product

import (
	"context"
	"go-pet-shop/internal/models"
)

type ProductsMock struct {
	GetAllProductsFunc func(ctx context.Context) ([]models.Product, error)
	CreateProductFunc  func(ctx context.Context, product models.Product) (int, error)
	DeleteProductFunc  func(ctx context.Context, id int) error
	UpdateProductFunc  func(ctx context.Context, product models.Product) error
}

func (m *ProductsMock) GetAllProducts(ctx context.Context) ([]models.Product, error) {
	if m.GetAllProductsFunc != nil {
		return m.GetAllProductsFunc(ctx)
	}
	return []models.Product{}, nil
}

func (m *ProductsMock) CreateProduct(ctx context.Context, product models.Product) (int, error) {
	if m.CreateProductFunc != nil {
		return m.CreateProductFunc(ctx, product)
	}
	return 0, nil
}

func (m *ProductsMock) DeleteProduct(ctx context.Context, id int) error {
	if m.DeleteProductFunc != nil {
		return m.DeleteProductFunc(ctx, id)
	}
	return nil
}

func (m *ProductsMock) UpdateProduct(ctx context.Context, product models.Product) error {
	if m.UpdateProductFunc != nil {
		return m.UpdateProductFunc(ctx, product)
	}
	return nil
}
