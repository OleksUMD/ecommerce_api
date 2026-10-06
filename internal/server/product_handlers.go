package server

import (
	"strconv"

	"github.com/OleksUMD/ecommerce_api/internal/dto"
	"github.com/OleksUMD/ecommerce_api/internal/services"
	"github.com/OleksUMD/ecommerce_api/internal/utils"
	"github.com/gin-gonic/gin"
)

func (s *Server) createCategory(c *gin.Context) {
	var req dto.CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "Invalid request data", err)
		return
	}
	prodService := services.NewProductService(s.db, s.config)
	response, err := prodService.CreateCategory(&req)
	if err != nil {
		utils.BadRequestResponse(c, "Create category failed", err)
		return
	}
	utils.CreatedResponse(c, "Category created successfully", response)
}

func (s *Server) getCategories(c *gin.Context) {
	prodService := services.NewProductService(s.db, s.config)
	response, err := prodService.GetCategories()
	if err != nil {
		utils.BadRequestResponse(c, "Get categories failed", err)
		return
	}
	utils.SuccessResponse(c, "Categories retrieved successfully", response)
}

func (s *Server) updateCategory(c *gin.Context) {
	var req dto.UpdateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "Invalid request data", err)
		return
	}
	categoryID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "Update category failed", err)
		return
	}
	prodService := services.NewProductService(s.db, s.config)
	response, err := prodService.UpdateCategory(uint(categoryID), &req)
	if err != nil {
		utils.BadRequestResponse(c, "Update category failed", err)
		return
	}
	utils.SuccessResponse(c, "Category updated successfully", response)
}

func (s *Server) deleteCategory(c *gin.Context) {
	categoryID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "Delete category failed", err)
		return
	}
	prodService := services.NewProductService(s.db, s.config)
	err = prodService.DeleteCategory(uint(categoryID))
	if err != nil {
		utils.BadRequestResponse(c, "Delete category failed", err)
		return
	}
	utils.SuccessNoContentResponse(c, "Category deleted successfully", nil)
}

func (s *Server) createProduct(c *gin.Context) {
	var req dto.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "Invalid request data", err)
		return
	}
	prodService := services.NewProductService(s.db, s.config)
	response, err := prodService.CreateProduct(&req)
	if err != nil {
		utils.BadRequestResponse(c, "Create product failed", err)
		return
	}
	utils.CreatedResponse(c, "Product created successfully", response)
}

func (s *Server) getProducts(c *gin.Context) {
	page, err := strconv.Atoi(c.Param("page"))
	if err != nil {
		utils.BadRequestResponse(c, "Get products failed", err)
		return
	}
	limit, err := strconv.Atoi(c.Param("limit"))
	if err != nil {
		utils.BadRequestResponse(c, "Get products failed", err)
		return
	}
	prodService := services.NewProductService(s.db, s.config)
	response, meta, err := prodService.GetProducts(page, limit)
	if err != nil {
		utils.BadRequestResponse(c, "Get products failed", err)
		return
	}
	utils.PaginatedSuccessResponse(c, "Products retrieved successfully", response, *meta)
}

func (s *Server) getProduct(c *gin.Context) {
	productID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "Get product failed", err)
		return
	}
	prodService := services.NewProductService(s.db, s.config)
	response, err := prodService.GetProduct(uint(productID))
	if err != nil {
		utils.BadRequestResponse(c, "Get product failed", err)
		return
	}
	utils.SuccessResponse(c, "Product updated successfully", response)
}

func (s *Server) updateProduct(c *gin.Context) {
	var req dto.UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "Invalid request data", err)
		return
	}
	productID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "Update product failed", err)
		return
	}
	prodService := services.NewProductService(s.db, s.config)
	response, err := prodService.UpdateProduct(uint(productID), &req)
	if err != nil {
		utils.BadRequestResponse(c, "Update product failed", err)
		return
	}
	utils.SuccessResponse(c, "Product updated successfully", response)
}

func (s *Server) deleteProduct(c *gin.Context) {
	productID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "Delete product failed", err)
		return
	}
	prodService := services.NewProductService(s.db, s.config)
	err = prodService.DeleteProduct(uint(productID))
	if err != nil {
		utils.BadRequestResponse(c, "Delete product failed", err)
		return
	}
	utils.SuccessNoContentResponse(c, "Product deleted successfully", nil)
}
