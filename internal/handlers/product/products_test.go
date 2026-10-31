package product

import (
	"context"
	"encoding/json"
	"errors"
	"go-pet-shop/internal/models"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"
)

// Get Product - Ready
func TestGetAllProducts_Success(t *testing.T) { //200
	mock := &ProductsMock{
		GetAllProductsFunc: func(ctx context.Context) ([]models.Product, error) {
			return []models.Product{
				{ID: 1, Name: "Dog Food"},
			}, nil
		},
	}

	req := httptest.NewRequest(http.MethodGet, "/products", nil)
	w := httptest.NewRecorder()

	handler := New(slog.Default(), mock)
	handler.GetAllProducts(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	var got []models.Product
	if err := json.NewDecoder(w.Body).Decode(&got); err != nil {
		t.Fatalf("failed to decode body: %v", err)
	}

	assert.Equal(t, []models.Product{{ID: 1, Name: "Dog Food"}}, got)
}

func TestGetAllProducts_Error(t *testing.T) { //500
	mock := &ProductsMock{
		GetAllProductsFunc: func(ctx context.Context) ([]models.Product, error) {
			return nil, errors.New("DB error")
		},
	}

	req := httptest.NewRequest(http.MethodGet, "/products", nil)
	w := httptest.NewRecorder()

	handler := New(slog.Default(), mock)
	handler.GetAllProducts(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", w.Code)
	}

	var resp struct {
		Error   string `json:"error"`
		Message string `json:"message"`
	}
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode body: %v", err)
	}

	assert.Equal(t, "Internal server error", resp.Error)
	assert.Equal(t, "Failed to retrieve products", resp.Message)
}

// =======================
// Create Product
// =======================

func TestCreateProduct_Success(t *testing.T) { //200
	mock := &ProductsMock{
		CreateProductFunc: func(ctx context.Context, product models.Product) (int, error) {
			return 0, nil
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

	var resp struct {
		Status  string         `json:"status"`
		ID      int            `json:"id"`
		Product models.Product `json:"product"`
	}
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode body: %v", err)
	}

	assert.Equal(t, "Product created successfully", resp.Status)
	assert.Equal(t, 0, resp.ID)
	assert.Equal(t, models.Product{Name: "Cat food", Price: 149, Stock: 67}, resp.Product)
}

func TestCreateProduct_BadRequest(t *testing.T) { //невалидным JSON (400 Bad Request)
	mock := &ProductsMock{
		CreateProductFunc: func(ctx context.Context, product models.Product) (int, error) {
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

	var resp struct {
		Error   string `json:"error"`
		Message string `json:"message"`
	}
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode body: %v", err)
	}

	assert.Equal(t, "Bad request", resp.Error)
	assert.Equal(t, "Invalid JSON payload", resp.Message)
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

	var resp struct {
		Error   string `json:"error"`
		Message string `json:"message"`
	}
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode body: %v", err)
	}

	assert.Equal(t, "Internal server error", resp.Error)
	assert.Equal(t, "Failed to create product", resp.Message)
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

	rctx := chi.NewRouteContext() // создаем карман для айдишки
	rctx.URLParams.Add("id", "0") //кладет айди
	req = req.WithContext(        //через withcontext прикрепляем к запросу
		context.WithValue(req.Context(), chi.RouteCtxKey, rctx),
	)

	handler := New(slog.Default(), mock)
	handler.UpdateProduct(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var resp struct {
		Status  string         `json:"status"`
		ID      int            `json:"id"`
		Product models.Product `json:"product"`
	}
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode body: %v", err)
	}

	assert.Equal(t, "Product updated successfully", resp.Status)
	assert.Equal(t, 0, resp.ID)
	assert.Equal(t, models.Product{Name: "Kitten_Shampoo", Price: 67, Stock: 12}, resp.Product)
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
	rctx.URLParams.Add("id", "0")
	req = req.WithContext(
		context.WithValue(req.Context(), chi.RouteCtxKey, rctx),
	)

	handler := New(slog.Default(), mock)
	handler.UpdateProduct(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}

	var resp struct {
		Error   string `json:"error"`
		Message string `json:"message"`
	}
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode body: %v", err)
	}

	assert.Equal(t, "Bad request", resp.Error)
	assert.Equal(t, "Invalid JSON payload", resp.Message)
}

func TestUpdateProduct_Fail(t *testing.T) { //(500 Internal Server Error)
	mock := &ProductsMock{
		UpdateProductFunc: func(ctx context.Context, product models.Product) error {
			return errors.New("Db err")
		},
	}

	body := strings.NewReader(`{"Name": "Shampoo_for_kittens", "Price": 67, "Stock": 3}`)
	req := httptest.NewRequest(http.MethodPut, "/products/{id}", body)
	w := httptest.NewRecorder()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "0")
	req = req.WithContext(
		context.WithValue(req.Context(), chi.RouteCtxKey, rctx),
	)

	handler := New(slog.Default(), mock)
	handler.UpdateProduct(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", w.Code)
	}

	var resp struct {
		Error   string `json:"error"`
		Message string `json:"message"`
	}
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode body: %v", err)
	}

	assert.Equal(t, "Internal server error", resp.Error)
	assert.Equal(t, "Failed to update product", resp.Message)
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

	var resp struct {
		Status  string `json:"status"`
		ID      int    `json:"id"`
		Message string `json:"message"`
	}
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode body: %v", err)
	}

	assert.Equal(t, "Product deleted successfully", resp.Status)
	assert.Equal(t, 1, resp.ID)
	assert.Equal(t, "Product with ID 1 has been deleted", resp.Message)
}

func TestDeleteProduct_BadRequest(t *testing.T) { //(400 Bad Request) // пустое id
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

	var resp struct {
		Error   string `json:"error"`
		Message string `json:"message"`
	}
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode body: %v", err)
	}

	assert.Equal(t, "Bad request", resp.Error)
	assert.Equal(t, "Product ID is required", resp.Message)
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

	var resp struct {
		Error   string `json:"error"`
		Message string `json:"message"`
	}
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode body: %v", err)
	}

	assert.Equal(t, "Internal server error", resp.Error)
	assert.Equal(t, "Failed to delete product", resp.Message)
}
