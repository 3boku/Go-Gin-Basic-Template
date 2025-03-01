package types

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type BasicModel struct {
	ID       uuid.UUID `gorm:"primarykey"`
	CreateAt time.Time
	UpdateAt time.Time
	DeleteAt gorm.DeletedAt `gorm:"index"`
}
