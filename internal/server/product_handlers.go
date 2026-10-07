package server

import (
	"strconv"

	"github.com/OleksUMD/ecommerce_api/internal/dto"
	"github.com/OleksUMD/ecommerce_api/internal/utils"
	"github.com/gin-gonic/gin"
)

func (s *Server) createCategory(c *gin.Context) {
	var req dto.CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "Invalid request data", err)
		return
	}
	response, err := s.productService.CreateCategory(&req)
	if err != nil {
		utils.BadRequestResponse(c, "Create category failed", err)
		return
	}
	utils.CreatedResponse(c, "Category created successfully", response)
}

func (s *Server) getCategories(c *gin.Context) {
	response, err := s.productService.GetCategories()
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
	response, err := s.productService.UpdateCategory(uint(categoryID), &req)
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
	err = s.productService.DeleteCategory(uint(categoryID))
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
	response, err := s.productService.CreateProduct(&req)
	if err != nil {
		utils.BadRequestResponse(c, "Create product failed", err)
		return
	}
	utils.CreatedResponse(c, "Product created successfully", response)
}

func (s *Server) getProducts(c *gin.Context) {
	page, err := strconv.Atoi(c.Query("page"))
	if err != nil {
		utils.BadRequestResponse(c, "Get products failed", err)
		return
	}
	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil {
		utils.BadRequestResponse(c, "Get products failed", err)
		return
	}
	response, meta, err := s.productService.GetProducts(page, limit)
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
	response, err := s.productService.GetProduct(uint(productID))
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
	response, err := s.productService.UpdateProduct(uint(productID), &req)
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
	err = s.productService.DeleteProduct(uint(productID))
	if err != nil {
		utils.BadRequestResponse(c, "Delete product failed", err)
		return
	}
	utils.SuccessNoContentResponse(c, "Product deleted successfully", nil)
}

func (s *Server) uploadProductImage(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid product ID", err)
		return
	}
	file, err := c.FormFile("image")
	if err != nil {
		utils.BadRequestResponse(c, "No file uploaded", err)
		return
	}

	url, err := s.uploadService.UploadProductImage(uint(id), file)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to upload image", err)
		return
	}
	if err := s.productService.AddProductImage(uint(id), url, file.Filename); err != nil {
		utils.InternalServerErrorResponse(c, "Failed to save image record", err)
		return
	}

	utils.SuccessResponse(c, "Image uploaded successfully", map[string]string{"url": url})
}
