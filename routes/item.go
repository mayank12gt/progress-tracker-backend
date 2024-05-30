package routes

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	database "github.com/mayank12gt/ProgressTracker/db"
	"github.com/mayank12gt/ProgressTracker/model"
	"gorm.io/gorm"
)

type ItemDTO struct {
	ID uint `json:"id"`

	CreatedAt time.Time

	Title string `json:"item_title"`

	Completed bool `json:"completed"`
}

type ItemUpdateReq struct {
	Title string `json:"item_title"`

	Completed bool `json:"completed"`
}

func createItemDTO(item model.Item) ItemDTO {
	return ItemDTO{ID: item.ID, CreatedAt: item.CreatedAt, Title: item.Title, Completed: item.Completed}
}

func CreateItem(c echo.Context) error {
	var item model.Item
	var list model.List

	id := c.Param("id")
	user := c.Get("user").(model.User)
	if err := c.Bind(&item); err != nil {
		return c.String(http.StatusBadRequest, "Bad Request")
	}

	database.Database.Db.Session(&gorm.Session{}).Where("user_id=?", user.ID).First(&list, id)
	if list.ID == 0 {
		return c.JSON(400, "List not found")
	}
	database.Database.Db.Session(&gorm.Session{}).Create(&item)
	setProgress(list)
	return c.String(200, "Item Created")
}
func GetItem(c echo.Context) error {
	var list model.List
	var item model.Item

	list_id := c.Param("id")
	item_id := c.Param("itemId")
	user := c.Get("user").(model.User)

	database.Database.Db.Session(&gorm.Session{}).Where("user_id=?", user.ID).First(&list, list_id)
	if list.ID == 0 {
		return c.JSON(400, "Incorrect List Id")
	}

	database.Database.Db.Session(&gorm.Session{}).Where("list_id=?", list_id).First(&item, item_id)
	if item.ID == 0 {
		return c.JSON(400, "Incorrect Item Id")
	}

	return c.JSON(200, item)
}

func UpdateItem(c echo.Context) error {
	var list model.List
	var item model.Item
	var updateReq ItemUpdateReq

	list_id := c.Param("id")
	item_id := c.Param("itemId")
	user := c.Get("user").(model.User)
	if err := c.Bind(&updateReq); err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	database.Database.Db.Session(&gorm.Session{}).Where("user_id=?", user.ID).First(&list, list_id)
	if list.ID == 0 {
		return c.JSON(400, "Incorrect List Id")
	}

	database.Database.Db.Session(&gorm.Session{}).Where("list_id=?", list_id).First(&item, item_id)
	if item.ID == 0 {
		return c.JSON(400, "Incorrect Item Id")
	}
	item.Title = updateReq.Title
	//item.Completed = updateReq.Completed
	database.Database.Db.Session(&gorm.Session{}).Save(&item)
	setProgress(list)

	return c.JSON(200, "Item Updated")
}

func GetAllItems(c echo.Context) error {
	var list model.List
	var items []model.Item

	list_id := c.Param("id")
	completed := c.QueryParam("completed")
	//item_id := c.Param("itemId")
	user := c.Get("user").(model.User)

	database.Database.Db.Session(&gorm.Session{}).Where("user_id=?", user.ID).First(&list, list_id)
	if list.ID == 0 {
		return c.JSON(400, "Incorrect List Id")
	}

	if completed == "true" {
		database.Database.Db.Session(&gorm.Session{}).Where("list_id=?", list_id).Where("completed=?", true).Find(&items)
	} else if completed == "false" {
		database.Database.Db.Session(&gorm.Session{}).Where("list_id=?", list_id).Where("completed=?", false).Find(&items)
	} else {
		database.Database.Db.Session(&gorm.Session{}).Where("list_id=?", list_id).Find(&items)
	}

	return c.JSON(200, items)
}

func DeleteItem(c echo.Context) error {
	var list model.List
	var item model.Item

	list_id := c.Param("id")
	item_id := c.Param("itemId")
	user := c.Get("user").(model.User)

	database.Database.Db.Session(&gorm.Session{}).Where("user_id=?", user.ID).First(&list, list_id)
	if list.ID == 0 {
		return c.JSON(400, "Incorrect List Id")
	}

	database.Database.Db.Session(&gorm.Session{}).Where("list_id=?", list_id).First(&item, item_id)
	if item.ID == 0 {
		return c.JSON(400, "Incorrect Item Id")
	}
	database.Database.Db.Session(&gorm.Session{}).Delete(&item)
	setProgress(list)

	return c.JSON(200, "Item Deleted")
}

func ToggleItemCompleted(c echo.Context) error {
	var list model.List
	var item model.Item

	list_id := c.Param("id")
	item_id := c.Param("itemId")
	user := c.Get("user").(model.User)

	database.Database.Db.Session(&gorm.Session{}).Where("user_id=?", user.ID).First(&list, list_id)
	if list.ID == 0 {
		return c.JSON(400, "Incorrect List Id")
	}

	database.Database.Db.Session(&gorm.Session{}).Where("list_id=?", list_id).First(&item, item_id)
	if item.ID == 0 {
		return c.JSON(400, "Incorrect Item Id")
	}
	item.Completed = !item.Completed

	database.Database.Db.Session(&gorm.Session{}).Save(&item)
	setProgress(list)
	return c.String(200, "Item marked as completed")
}

func setProgress(list model.List) {
	println("setProgress")
	var completedItemsCount int64

	var totalItemsCount int64

	database.Database.Db.Session(&gorm.Session{}).Model(&model.Item{}).Where("list_id=?", list.ID).Where("completed=?", true).Count(&completedItemsCount)
	println("completed", completedItemsCount)

	database.Database.Db.Session(&gorm.Session{}).Model(&model.Item{}).Where("list_id=?", list.ID).Count(&totalItemsCount)
	println("total", totalItemsCount)
	if totalItemsCount == 0 {
		list.Progress = 0.0
	} else {
		list.Progress = (float32(completedItemsCount) / float32(totalItemsCount)) * 100
	}
	println("progress ", list.Progress)

	database.Database.Db.Session(&gorm.Session{}).Save(&list)
}
