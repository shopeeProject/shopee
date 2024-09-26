package category

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shopeeProject/shopee/models"
	util "github.com/shopeeProject/shopee/util"
)

const (
	routePrefix      = "/category"
	addCategory      = "/add-category"
	removeCategory   = "/remove-category/:id"
	getAllCategories = "/get-all-categories"
)

type Category struct {
	ID   uint   ` json:"id"`
	Name string `json:"name"`
}

// Add Category Handler
func addCategoryHandler(r *util.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		var category Category
		if err := c.ShouldBindJSON(&category); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := r.DB.Create(&category).Error; err != nil {
			log.Fatal("failed to add category:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to add category"})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"message": "category added successfully", "id": category.ID})
	}
}

// Remove Category Handler
func removeCategoryHandler(r *util.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if err := r.DB.Delete(&models.Category{}, id).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to remove category"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "category removed successfully"})
	}
}

// Get All Categories Handler
func getAllCategoriesHandler(r *util.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		var categories []models.Category
		if err := r.DB.Find(&categories).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch categories"})
			return
		}
		c.JSON(http.StatusOK, categories)
	}
}

func RegisterRoutes(router *gin.Engine, r *util.Repository) {
	v1 := router.Group(routePrefix)
	{
		v1.POST(addCategory, addCategoryHandler(r))
		v1.DELETE(removeCategory, removeCategoryHandler(r))
		v1.GET(getAllCategories, getAllCategoriesHandler(r))
	}
}
