package services

import (
	"mime/multipart"

	"github.com/OleksUMD/ecommerce_api/internal/dto"
	"github.com/OleksUMD/ecommerce_api/internal/models"
	"github.com/OleksUMD/ecommerce_api/internal/utils"
	"gorm.io/gorm"
)

type IAuthService interface {
	Register(req *dto.RegisterRequest) (*dto.UserResponse, error)
	Login(req *dto.LoginRequest) (*dto.TmpAuthResponse, error)
	RefreshToken(req *dto.RefreshTokenRequest) (*dto.AuthResponse, error)
	Logout(userID uint) error
	VerifyEmail(token string) error
	VerifyTOTP(req *dto.TotpRequest) (*dto.AuthResponse, error)
	ResubmitVerificationLink(token string) (*dto.UserResponse, error)
	AddMFA(user *models.User) (*dto.AddMFAResponse, error)
	UseBackup(req *dto.BackupRequest) (*dto.AddMFAResponse, error)
	GetUnverifiedUser(token string) (*models.User, error)
	ResetPassword(req *dto.ResetPasswordRequest) error
	generateAuthResponse(user *models.User) (*dto.AuthResponse, error)
	generateAddMFAUserResponse(user *models.User) (*dto.AddMFAResponse, error)
	checkBackupCode(backupCode, inputCode string) error
}

type IUserService interface {
	GetProfile(userID uint) (*dto.UserResponse, error)
	UpdateProfile(userID uint, req *dto.UpdateProfileRequest) (*dto.UserResponse, error)
	UserIsVerified(userID uint) bool
}

type IProductService interface {
	CreateCategory(req *dto.CreateCategoryRequest) (*dto.CategoryResponse, error)
	GetCategories() ([]dto.CategoryResponse, error)
	UpdateCategory(id uint, req *dto.UpdateCategoryRequest) (*dto.CategoryResponse, error)
	DeleteCategory(id uint) error
	CreateProduct(req *dto.CreateProductRequest) (*dto.ProductResponse, error)
	GetProduct(id uint) (*dto.ProductResponse, error)
	GetProducts(page, limit int) ([]dto.ProductResponse, *utils.PaginationMeta, error)
	UpdateProduct(id uint, req *dto.UpdateProductRequest) (*dto.ProductResponse, error)
	DeleteProduct(id uint) error
	AddProductImage(productID uint, url, altText string) error
	convertToProductResponse(product *models.Product) dto.ProductResponse
	SearchProducts(req *dto.SearchProductRequest) ([]dto.ProductSearchResult, *utils.PaginationMeta, error)
}

type ICartService interface {
	GetCart(userID uint) (*dto.CartResponse, error)
	AddToCart(userID uint, req *dto.AddToCartRequest) (*dto.CartResponse, error)
	UpdateToCart(userID, itemID uint, req *dto.UpdateCartItemRequest) (*dto.CartResponse, error)
	RemoveFromCart(userID, itemID uint) error
	convertToCartResponse(cart *models.Cart) *dto.CartResponse
}

type IOrderService interface {
	CreateOrder(userID uint) (*dto.OrderResponse, error)
	GetOrders(userID uint, page, limit int) ([]dto.OrderResponse, *utils.PaginationMeta, error)
	GetOrder(userID, orderID uint) (*dto.OrderResponse, error)
	getOrderResponse(tx *gorm.DB, id uint) (*dto.OrderResponse, error)
	convertToOrderResponse(order *models.Order) dto.OrderResponse
}

type IUploadService interface {
	UploadProductImage(productID uint, file *multipart.FileHeader) (string, error)
}
