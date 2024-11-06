package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/antoniofmoliveira/apis/internal/dto"
	"github.com/antoniofmoliveira/apis/internal/entity"
	"github.com/antoniofmoliveira/apis/internal/infra/database"
	pkgentity "github.com/antoniofmoliveira/apis/pkg/entity"
)

type ProductHandler struct {
	ProductDB database.ProductRepositoryInterface
}

func NewProductHandler(db database.ProductRepositoryInterface) *ProductHandler {
	return &ProductHandler{ProductDB: db}
}

// @Summary      Create a new product
// @Description  Create a new product
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        input  body      dto.CreateProductInput  true  "product request"
// @Success      201
// @Failure      400     {object}  Error
// @Failure      500     {object}  Error
// @Router       /products [post]
// @Security     ApiKeyAuth
func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var product dto.CreateProductInput
	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		json.NewEncoder(w).Encode(Error{Message: err.Error()})
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	p, err := entity.NewProduct(product.Name, product.Price)
	if err != nil {
		json.NewEncoder(w).Encode(Error{Message: err.Error()})
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = h.ProductDB.Create(p)
	if err != nil {
		json.NewEncoder(w).Encode(Error{Message: err.Error()})
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

// @Summary      Get product by ID
// @Description  Get product by ID
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        id  path      string  true  "Product ID"
// @Success      200  {object}  entity.Product
// @Failure      400  {object}  Error
// @Failure      404  {object}  Error
// @Failure      500  {object}  Error
// @Router       /products/{id} [get]
// @Security     ApiKeyAuth
func (h *ProductHandler) GetProduct(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		json.NewEncoder(w).Encode(Error{Message: "Product ID is required"})
		http.Error(w, "Product ID is required", http.StatusBadRequest)
	}
	product, err := h.ProductDB.FindByID(id)
	if err != nil {
		json.NewEncoder(w).Encode(Error{Message: err.Error()})
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(product)
}

// @Summary      Update product by ID
// @Description  Update product by ID
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        id  path      string  true  "Product ID"
// @Param        input  body      dto.UpdateProductInput  true  "product request"
// @Success      200
// @Failure      400  {object}  Error
// @Failure      404  {object}  Error
// @Failure      500  {object}  Error
// @Router       /products/{id} [put]
// @Security     ApiKeyAuth
func (h *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		json.NewEncoder(w).Encode(Error{Message: "Product ID is required"})
		http.Error(w, "Product ID is required", http.StatusBadRequest)
	}
	var product dto.UpdateProductInput
	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		json.NewEncoder(w).Encode(Error{Message: err.Error()})
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	ID, err := pkgentity.ParseId(id)
	if err != nil {
		json.NewEncoder(w).Encode(Error{Message: err.Error()})
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	p := entity.Product{
		ID:    ID,
		Name:  product.Name,
		Price: product.Price,
	}
	rows, err := h.ProductDB.Update(&p)
	if err != nil {
		json.NewEncoder(w).Encode(Error{Message: err.Error()})
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if rows == 0 {
		json.NewEncoder(w).Encode(Error{Message: "Product not found"})
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

// @Summary      Delete product by ID
// @Description  Delete product by ID
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        id  path      string  true  "Product ID"
// @Success      200
// @Failure      400  {object}  Error
// @Failure      404  {object}  Error
// @Failure      500  {object}  Error
// @Router       /products/{id} [delete]
// @Security     ApiKeyAuth
func (h *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		json.NewEncoder(w).Encode(Error{Message: "Product ID is required"})
		http.Error(w, "Product ID is required", http.StatusBadRequest)
	}
	rows, err := h.ProductDB.Delete(id)
	if err != nil {
		json.NewEncoder(w).Encode(Error{Message: err.Error()})
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if rows == 0 {
		json.NewEncoder(w).Encode(Error{Message: "Product not found"})
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

// @Summary      Find all products
// @Description  Find all products
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        page  query     int  false  "Page number"
// @Param        limit  query     int  false  "Number of products per page"
// @Param        sort  query     string  false  "Sort order (asc or desc)"
// @Success      200  {object}  entity.Product
// @Failure      400  {object}  Error
// @Failure      404  {object}  Error
// @Failure      500  {object}  Error
// @Router       /products [get]
// @Security     ApiKeyAuth
func (h *ProductHandler) FindAllProducts(w http.ResponseWriter, r *http.Request) {
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		page = 0
	}
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		limit = 0
	}
	sort := r.URL.Query().Get("sort")
	products, err := h.ProductDB.FindAll(page, limit, sort)
	if err != nil {
		json.NewEncoder(w).Encode(Error{Message: err.Error()})
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}
