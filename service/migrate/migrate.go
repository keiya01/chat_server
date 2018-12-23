package migrate

import (
	"github.com/jinzhu/gorm"
	"github.com/keiya01/chat_room/model"
)

func Set(db *gorm.DB) {
	db.AutoMigrate(&model.Chat{})
}
