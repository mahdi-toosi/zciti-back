package schema

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"time"
)

type Base struct {
	CreatedAt time.Time      `gorm:"autoCreateTime" json:",omitempty"`
	UpdatedAt time.Time      `json:",omitempty"`
	DeletedAt gorm.DeletedAt `json:",omitempty" faker:"-"`
}

type JSON json.RawMessage

// Scan value into Jsonb, implements sql.Scanner interface
func (j *JSON) Scan(value any) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}
	result := json.RawMessage{}
	err := json.Unmarshal(bytes, &result)
	*j = JSON(result)
	return err
}

// Value return json value, implement driver.Valuer interface
func (j JSON) Value() (driver.Value, error) {
	if len(j) == 0 {
		return nil, nil
	}
	return json.RawMessage(j).MarshalJSON()
}
