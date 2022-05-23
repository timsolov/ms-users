package postgres

import (
	"encoding/json"
	"errors"
)

type StringInterfaceMap map[string]interface{}

func (m *StringInterfaceMap) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("failed type assertion to []byte")
	}
	return json.Unmarshal(b, &m)
}
