package category

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shopeeProject/shopee/models"
	util "github.com/shopeeProject/shopee/util"
	"gorm.io/gorm"
)

const (
	routePrefix      = "/category"
	addCategory      = "/add-category"
	removeCategory   = "/remove-category"
	getAllCategories = "/get-all-categories"
)

type Category struct {
	ID   int    ` json:"id"`
	Name string `json:"name"`
}

func ValidateCategory(db *gorm.DB, category int) util.Response {
	categories, err := fetchAllCategories(db)
	if err != nil {
		return util.Response{Message: err.Error()}
	}
	categoryFound := false
	for _, val := range categories {
		if val.Id == category {
			categoryFound = true
			break
		}
	}
	return util.Response{Success: categoryFound, Message: "category found"}
}

// Add Category Handler
func AddCategoryHandler(r *util.Repository) gin.HandlerFunc {

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
func RemoveCategoryHandler(r *util.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if err := r.DB.Delete(&models.Category{}, id).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to remove category"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "category removed successfully"})
	}
}

func fetchAllCategories(db *gorm.DB) ([]models.Category, error) {
	var categories []models.Category
	if err := db.Find(&categories).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch categories:%s", err.Error())
	}
	return categories, nil
}

// Get All Categories Handler
func getAllCategoriesHandler(r *util.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		categories, err := fetchAllCategories(r.DB)
		if err != nil {
			c.JSON(http.StatusBadRequest, "error fetching categories")
			return
		}
		c.JSON(http.StatusOK, categories)
	}
}

func RegisterRoutes(router *gin.Engine, r *util.Repository) {
	v1 := router.Group(routePrefix)
	{
		v1.POST(addCategory, AddCategoryHandler(r))
		v1.DELETE(removeCategory, RemoveCategoryHandler(r))
		v1.GET(getAllCategories, getAllCategoriesHandler(r))
	}
}
