package document

import (
	"fmt"
	"sync"
)

// FieldType represents the type of a field value
type FieldType int

const (
	// StringType represents string field values
	StringType FieldType = iota
	// IntType represents integer field values
	IntType
	// FloatType represents floating-point field values
	FloatType
)

// Field represents a single field in a document
type Field struct {
	Name     string
	Type     FieldType
	Value    interface{}
}

// Document represents a searchable document with multiple fields
type Document struct {
	mu     sync.RWMutex
	fields map[string]Field
}

// NewDocument creates a new Document instance
func NewDocument() *Document {
	return &Document{
		fields: make(map[string]Field),
	}
}

// AddField adds a new field to the document
func (d *Document) AddField(name string, value interface{}) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	fieldType, err := determineFieldType(value)
	if err != nil {
		return fmt.Errorf("failed to add field: %w", err)
	}

	d.fields[name] = Field{
		Name:  name,
		Type:  fieldType,
		Value: value,
	}
	return nil
}

// GetField retrieves a field by name
func (d *Document) GetField(name string) (Field, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	field, exists := d.fields[name]
	if !exists {
		return Field{}, fmt.Errorf("field %s not found", name)
	}
	return field, nil
}

// determineFieldType infers the FieldType from a value
func determineFieldType(value interface{}) (FieldType, error) {
	switch value.(type) {
	case string:
		return StringType, nil
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return IntType, nil
	case float32, float64:
		return FloatType, nil
	default:
		return 0, fmt.Errorf("unsupported field type for value: %v", value)
	}
}
