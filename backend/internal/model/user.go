package model

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

// JSONStringSlice 用于存储 JSONB 字符串数组
type JSONStringSlice []string

func (j JSONStringSlice) Value() (driver.Value, error) {
	if j == nil {
		return "[]", nil
	}
	b, err := json.Marshal(j)
	return string(b), err
}

func (j *JSONStringSlice) Scan(value any) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to scan JSONStringSlice")
	}
	return json.Unmarshal(bytes, j)
}

// JSONMap 用于存储 JSONB 任意对象
type JSONMap map[string]any

func (j JSONMap) Value() (driver.Value, error) {
	if j == nil {
		return "{}", nil
	}
	b, err := json.Marshal(j)
	return string(b), err
}

func (j *JSONMap) Scan(value any) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to scan JSONMap")
	}
	return json.Unmarshal(bytes, j)
}

type User struct {
	ID           int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	Phone        string    `json:"phone" gorm:"type:varchar(20);uniqueIndex;not null"`
	PasswordHash string    `json:"-" gorm:"type:varchar(255)"`
	Role         string    `json:"role" gorm:"type:varchar(20);not null;default:consumer"` // admin/merchant/staff/consumer
	Status       string    `json:"status" gorm:"type:varchar(20);not null;default:active"` // active/frozen/pending
	Nickname     string    `json:"nickname" gorm:"type:varchar(100)"`
	AvatarURL    string    `json:"avatar_url" gorm:"type:varchar(500)"`
	Openid       string    `json:"openid,omitempty" gorm:"type:varchar(100);uniqueIndex"`
	Unionid      string    `json:"unionid,omitempty" gorm:"type:varchar(100)"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (User) TableName() string { return "users" }
