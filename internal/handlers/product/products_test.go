package product

import (
	"context"
	"encoding/json"
	"errors"
	"go-pet-shop/internal/handlers/product/mocks"
	"go-pet-shop/internal/models"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Get Product - Ready
func TestGetAllProducts_Success(t *testing.T) { //200
	target := []models.Product{{ID: 1, Name: "Dog Food"}}

	m := mocks.NewProducts(t)
	m.On("GetAllProducts", mock.Anything).Return(target, nil)

	req := httptest.NewRequest(http.MethodGet, "/products", nil)
	w := httptest.NewRecorder()

	handler := New(slog.Default(), m)
	handler.GetAllProducts(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	var got []models.Product
	if err := json.NewDecoder(w.Body).Decode(&got); err != nil {
		t.Fatalf("failed to decode body: %v", err)
	}

	assert.Equal(t, target, got)
}
func TestGetAllProducts_Error(t *testing.T) { //500
	m := mocks.NewProducts(t)
	m.On("GetAllProducts", mock.Anything).Return([]models.Product{}, errors.New("DB error"))

	req := httptest.NewRequest(http.MethodGet, "/products", nil)
	w := httptest.NewRecorder()

	handler := New(slog.Default(), m)
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
	target := models.Product{Name: "Cat food", Price: 149, Stock: 67}

	m := mocks.NewProducts(t)
	m.On("CreateProduct", mock.Anything, target).Return(1, nil)

	body := strings.NewReader(`{"Name": "Cat food", "Price": 149, "Stock": 67}`)
	req := httptest.NewRequest(http.MethodPost, "/products", body)
	w := httptest.NewRecorder()

	handler := New(slog.Default(), m)
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
	assert.Equal(t, 1, resp.ID)
	assert.Equal(t, models.Product{ID: 1, Name: "Cat food", Price: 149, Stock: 67}, resp.Product)
}

func TestCreateProduct_BadRequest(t *testing.T) { //невалидным JSON (400 Bad Request)
	m := mocks.NewProducts(t)

	body := strings.NewReader(`{"Name": 'Cat shampoo', "Price': 987, 'Stock': 1}`)
	req := httptest.NewRequest(http.MethodPost, "/products", body)
	w := httptest.NewRecorder()

	handler := New(slog.Default(), m)
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
	target := models.Product{Name: "Dog Shampoo", Price: 67, Stock: 23}

	m := mocks.NewProducts(t)
	m.On("CreateProduct", mock.Anything, target).Return(1, errors.New("Db error"))

	body := strings.NewReader(`{"Name": "Dog Shampoo", "Price": 67, "Stock": 23}`)
	req := httptest.NewRequest(http.MethodPost, "/products", body)
	w := httptest.NewRecorder()

	handler := New(slog.Default(), m)
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

	assert.Equal(t, "Failed to create product", resp.Message)
	assert.Equal(t, "Internal server error", resp.Error)
}

// =======================
// Update Product
// =======================

func TestUpdateProduct_Success(t *testing.T) { // (200 OK)
	target := models.Product{Name: "Kitten_Shampoo", Price: 67, Stock: 12}

	m := mocks.NewProducts(t)
	m.On("UpdateProduct", mock.Anything, target).Return(nil)

	body := strings.NewReader(`{"Name": "Kitten_Shampoo", "Price": 67, "Stock": 12}`)
	req := httptest.NewRequest(http.MethodPut, "/products/{id}", body)
	w := httptest.NewRecorder()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "0")
	req = req.WithContext(
		context.WithValue(req.Context(), chi.RouteCtxKey, rctx),
	)

	handler := New(slog.Default(), m)
	handler.UpdateProduct(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var resp struct {
		ID      int            `json:"id"`
		Status  string         `json:"status"`
		Product models.Product `json:"product"`
	}
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode body: %v", err)
	}

	assert.Equal(t, 0, resp.ID)
	assert.Equal(t, "Product updated successfully", resp.Status)
	assert.Equal(t, models.Product{Name: "Kitten_Shampoo", Price: 67, Stock: 12}, resp.Product)
}

func TestUpdateProduct_BadRequest(t *testing.T) { // невалидным JSON (400 Bad Request)
	m := mocks.NewProducts(t)

	body := strings.NewReader(`{"Name": 'Shampoo_cat', "Price': 76, 'Stock': 67}`)
	req := httptest.NewRequest(http.MethodPut, "/products/{id}", body)
	w := httptest.NewRecorder()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	req = req.WithContext(
		context.WithValue(req.Context(), chi.RouteCtxKey, rctx),
	)

	handler := New(slog.Default(), m)
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
	target := models.Product{Name: "Shampoo_for_kittens", Price: 67, Stock: 3}

	m := mocks.NewProducts(t)
	m.On("UpdateProduct", mock.Anything, target).Return(errors.New("DB error"))

	body := strings.NewReader(`{"Name": "Shampoo_for_kittens", "Price": 67, "Stock": 3}`)
	req := httptest.NewRequest(http.MethodPut, "/products/{id}", body)
	w := httptest.NewRecorder()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "0")
	req = req.WithContext(
		context.WithValue(req.Context(), chi.RouteCtxKey, rctx),
	)

	handler := New(slog.Default(), m)
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
	m := mocks.NewProducts(t)
	m.On("DeleteProduct", mock.Anything, 0).Return(nil)

	body := strings.NewReader(`{"Name": "Shampoo_for_kittens", "Price": 67, "Stock": 3}`)
	req := httptest.NewRequest(http.MethodDelete, "/products/{id}", body)
	w := httptest.NewRecorder()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "0")
	req = req.WithContext(
		context.WithValue(req.Context(), chi.RouteCtxKey, rctx),
	)

	handler := New(slog.Default(), m)
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
	assert.Equal(t, 0, resp.ID)
	assert.Equal(t, "Product with ID 0 has been deleted", resp.Message)
}

func TestDeleteProduct_BadRequest(t *testing.T) { //(400 Bad Request) пустой id
	m := mocks.NewProducts(t)

	body := strings.NewReader(`{"Name": "Shampoo_for_kittens", "Price": 67, "Stock": 3}`)
	req := httptest.NewRequest(http.MethodDelete, "/products/{id}", body)
	w := httptest.NewRecorder()

	handler := New(slog.Default(), m)
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
	m := mocks.NewProducts(t)
	m.On("DeleteProduct", mock.Anything, 0).Return(errors.New("DB error"))

	body := strings.NewReader(`{"Name": "Shampoo_for_kittens", "Price": 67, "Stock": 3}`)
	req := httptest.NewRequest(http.MethodDelete, "/products/{id}", body)
	w := httptest.NewRecorder()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "0")
	req = req.WithContext(
		context.WithValue(req.Context(), chi.RouteCtxKey, rctx),
	)

	handler := New(slog.Default(), m)
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
