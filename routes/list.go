package routes

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	database "github.com/mayank12gt/ProgressTracker/db"
	"github.com/mayank12gt/ProgressTracker/model"
	"gorm.io/gorm"
)

type ListDTO struct {
	ID uint `json:"id" `

	CreatedAt time.Time `json:"created_at"`

	Title string `json:"list_title"`

	Progress float32 `json:"progress"`

	Items []model.Item `json:"list_items"`
}

type UpdateListReq struct {
	Title string `json:"list_title"`
}

func createListDTO(list model.List) ListDTO {
	return ListDTO{ID: list.ID, CreatedAt: list.CreatedAt, Title: list.Title, Progress: list.Progress, Items: list.Items}
}

func GetAllLists(c echo.Context) error {
	var lists []model.List
	user := c.Get("user").(model.User)
	database.Database.Db.Session(&gorm.Session{}).Where("user_id=?", user.ID).Find(&lists)

	return c.JSON(200, lists)

}

func GetList(c echo.Context) error {
	var list model.List
	id := c.Param("id")
	user := c.Get("user").(model.User)
	database.Database.Db.Session(&gorm.Session{}).Where("user_id=?", user.ID).Preload("Items").First(&list, id)
	if list.ID == 0 {
		return c.JSON(400, "List not found")
	}
	list.User = user
	return c.JSON(200, list)

}

func CreateList(c echo.Context) error {
	var list model.List

	if err := c.Bind(&list); err != nil {
		return c.String(http.StatusBadRequest, "Bad Request")
	}
	if err := database.Database.Db.Statement.Session(&gorm.Session{}).Create(&list); err != nil {
		setProgress(list)
		return c.String(http.StatusBadRequest, string(err.Name()))
	}

	return c.String(200, "List Created")
}

func UpdateList(c echo.Context) error {
	var list model.List
	var updateReq UpdateListReq

	id := c.Param("id")
	user := c.Get("user").(model.User)
	if err := c.Bind(&updateReq); err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	database.Database.Db.Session(&gorm.Session{}).Where("user_id=?", user.ID).Preload("Items").First(&list, id)
	if list.ID == 0 {
		return c.JSON(400, "List not found")
	}
	list.Title = updateReq.Title
	database.Database.Db.Session(&gorm.Session{}).Save(&list)
	return c.String(200, "List Updated")
}

func DeleteList(c echo.Context) error {
	var list model.List

	id := c.Param("id")
	user := c.Get("user").(model.User)

	database.Database.Db.Session(&gorm.Session{}).Where("user_id=?", user.ID).Preload("Items").First(&list, id)
	if list.ID == 0 {
		return c.JSON(400, "List not found")
	}

	database.Database.Db.Session(&gorm.Session{}).Where("list_id=?", list.ID).Delete(&model.Item{})
	database.Database.Db.Session(&gorm.Session{}).Delete(&list)
	return c.String(200, "List Deleted")
}
