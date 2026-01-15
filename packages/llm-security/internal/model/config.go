package model

import (
	"time"
)

type SysConfig struct {
	ID          int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	ConfigKey   string    `json:"config_key" gorm:"size:100;not null;uniqueIndex"`
	ConfigValue string    `json:"config_value" gorm:"type:text"`
	Description string    `json:"description" gorm:"size:500"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (SysConfig) TableName() string {
	return "sys_config"
}
