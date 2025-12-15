package types

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type JSONB map[string]any

func (j *JSONB) Scan(src any) error {
	if src == nil {
		*j = nil
		return nil
	}

	var data []byte
	switch v := src.(type) {
	case []byte:
		data = v
	case string:
		data = []byte(v)
	default:
		return fmt.Errorf("types.JSONB: unsupported type %T", src)
	}

	if len(data) == 0 {
		*j = JSONB{}
		return nil
	}

	var result map[string]any
	if err := json.Unmarshal(data, &result); err != nil {
		return fmt.Errorf("types.JSONB: %w", err)
	}

	*j = JSONB(result)
	return nil
}

func (j JSONB) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}

	data, err := json.Marshal(map[string]any(j))
	if err != nil {
		return nil, fmt.Errorf("types.JSONB: %w", err)
	}

	return data, nil
}
