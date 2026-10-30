package product

import (
	"context"
	"errors"
	"go-pet-shop/internal/models"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi"
)

// Get Product - Ready
func TestGetAllProducts_Success(t *testing.T) {
	// Мокаем storage — он вернёт один продукт.
	mock := &ProductsMock{
		GetAllProductsFunc: func(ctx context.Context) ([]models.Product, error) {
			return []models.Product{
				{ID: 1, Name: "Dog Food"},
			}, nil
		},
	}

	// Создаем HTTP-запрос GET /products
	req := httptest.NewRequest(http.MethodGet, "/products", nil)
	w := httptest.NewRecorder()

	// Создаем хендлер с мок-хранилищем
	handler := New(slog.Default(), mock)

	// Вызываем метод GetAllProducts, который является http.HandlerFunc
	handler.GetAllProducts(w, req)

	// Проверяем HTTP-код
	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}
}
func TestGetAllProducts_Error(t *testing.T) {
	// Мокаем storage — он будет возвращать ошибку
	mock := &ProductsMock{
		GetAllProductsFunc: func(ctx context.Context) ([]models.Product, error) {
			return nil, errors.New("DB error")
		},
	}

	// Создаем запрос
	req := httptest.NewRequest(http.MethodGet, "/products", nil)
	w := httptest.NewRecorder()

	handler := New(slog.Default(), mock)
	handler.GetAllProducts(w, req)

	// Ожидаем HTTP 500
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", w.Code)
	}
}

// =======================
// Create Product
// =======================

func TestCreateProduct_Success(t *testing.T) { //200
	mock := &ProductsMock{
		CreateProductFunc: func(ctx context.Context, product models.Product) (int, error) {
			return 1, nil
		},
	}

	body := strings.NewReader(`{"Name": "Cat food", "Price": 149, "Stock": 67}`)
	req := httptest.NewRequest(http.MethodPost, "/products", body)
	w := httptest.NewRecorder()
	
	handler := New(slog.Default(), mock)
	handler.CreateProduct(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}
}

func TestCreateProduct_BadRequest(t *testing.T) { //невалидным JSON (400 Bad Request)
	mock := &ProductsMock{
		CreateProductFunc: func(ctx context.Context, product models.Product) (int, error){
			return 1, nil
		},
	}

	body := strings.NewReader(`{"Name": 'Cat shampoo', "Price': 987, 'Stock': 1}`)
	req := httptest.NewRequest(http.MethodPost, "/products", body)
	w := httptest.NewRecorder()

	handler := New(slog.Default(), mock)
	handler.CreateProduct(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestCreateProduct_Fail(t *testing.T) { // (500 Internal Server Error)
	mock := &ProductsMock{
		CreateProductFunc: func(ctx context.Context, product models.Product) (int, error) {
			return 1, errors.New("DB error")
		},
	}

	body := strings.NewReader(`{"Name": "Dog Shampoo", "Price": 67, "Stock": 23}`)
	req := httptest.NewRequest(http.MethodPost, "/products", body)
	w := httptest.NewRecorder()

	handler := New(slog.Default(), mock)
	handler.CreateProduct(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", w.Code)
	}
}

// =======================
// Update Product
// =======================

func TestUpdateProduct_Success(t *testing.T) { // (200 OK)
	mock := &ProductsMock{
		UpdateProductFunc: func(ctx context.Context, product models.Product) error {
			return nil
		},
	}

	body := strings.NewReader(`{"Name": "Kitten_Shampoo", "Price": 67, "Stock": 12}`)
	req := httptest.NewRequest(http.MethodPut, "/products/{id}", body)
	w := httptest.NewRecorder()

	// записать в конспект
	rctx := chi.NewRouteContext() // создаем карман для айдишки
	rctx.URLParams.Add("id", "1") //кладет айди
	req = req.WithContext( //через withcontext прикрепляем к запросу
		context.WithValue(req.Context(), chi.RouteCtxKey, rctx),
	)

	handler := New(slog.Default(), mock)
	handler.UpdateProduct(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestUpdateProduct_BadRequest(t *testing.T) { // невалидным JSON (400 Bad Request)
	mock := &ProductsMock{
		UpdateProductFunc: func(ctx context.Context, product models.Product) error {
			return nil
		},
	}

	body := strings.NewReader(`{"Name": 'Shampoo_cat', "Price': 76, 'Stock': 67}`)
	req := httptest.NewRequest(http.MethodPut, "/products/{id}", body)
	w := httptest.NewRecorder()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	req = req.WithContext(
		context.WithValue(req.Context(), chi.RouteCtxKey, rctx),
	)

	handler := New(slog.Default(), mock)
	handler.UpdateProduct(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestUpdateProduct_Fail(t *testing.T) {//(500 Internal Server Error)
	mock := &ProductsMock{
		UpdateProductFunc: func(ctx context.Context, product models.Product) error {
			return errors.New("Db err")
		},
	}

	body := strings.NewReader(`{"Name": "Shampoo_for_kittens", "Price": 67, "Stock": 3}`)
	req := httptest.NewRequest(http.MethodPut, "/products/{id}", body)
	w := httptest.NewRecorder()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	req = req.WithContext(
		context.WithValue(req.Context(), chi.RouteCtxKey, rctx),
	)

	handler := New(slog.Default(), mock)
	handler.UpdateProduct(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", w.Code)
	}
}

// =======================
// Delete Product
// =======================

func TestDeleteProduct_Success(t *testing.T) { // (200 OK)
	mock := &ProductsMock{
		DeleteProductFunc: func(ctx context.Context, id int) error {
			return nil
		},
	}

	body := strings.NewReader(`{"Name": "Shampoo_for_kittens", "Price": 67, "Stock": 3}`)
	req := httptest.NewRequest(http.MethodDelete, "/products/{id}", body)
	w := httptest.NewRecorder()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	req = req.WithContext(
		context.WithValue(req.Context(), chi.RouteCtxKey, rctx),
	)

	handler := New(slog.Default(), mock)
	handler.DeleteProduct(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestDeleteProduct_BadRequest(t *testing.T) { //(400 Bad Request)
	mock := &ProductsMock{
		DeleteProductFunc: func(ctx context.Context, id int) error {
			return nil
		},
	}

	body := strings.NewReader(`{"Name": "Shampoo_for_kittens", "Price": 67, "Stock": 3}`)
	req := httptest.NewRequest(http.MethodDelete, "/products/{id}", body)
	w := httptest.NewRecorder()

	handler := New(slog.Default(), mock)
	handler.DeleteProduct(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestDeleteProduct_Fail(t *testing.T) { // (500 Internal Server Error)
	mock := &ProductsMock{
		DeleteProductFunc: func(ctx context.Context, id int) error {
			return errors.New("Db err")
		},
	}

	body := strings.NewReader(`{"Name": "Shampoo_for_kittens", "Price": 67, "Stock": 3}`)
	req := httptest.NewRequest(http.MethodDelete, "/products/{id}", body)
	w := httptest.NewRecorder()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	req = req.WithContext(
		context.WithValue(req.Context(), chi.RouteCtxKey, rctx),
	)

	handler := New(slog.Default(), mock)
	handler.DeleteProduct(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", w.Code)
	}
}