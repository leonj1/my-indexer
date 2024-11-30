package document

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"
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
	// TimeType represents time.Time field values
	TimeType
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
	ID     int
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

// GetFields returns a map of all fields in the document
func (d *Document) GetFields() map[string]Field {
	d.mu.RLock()
	defer d.mu.RUnlock()

	// Create a copy of the fields map to prevent concurrent modification
	fields := make(map[string]Field, len(d.fields))
	for k, v := range d.fields {
		fields[k] = v
	}
	return fields
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
	case time.Time:
		return TimeType, nil
	default:
		return 0, fmt.Errorf("unsupported field type for value: %v", value)
	}
}

// MarshalJSON implements json.Marshaler interface
func (d *Document) MarshalJSON() ([]byte, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	// Create a map of field name to field value for JSON serialization
	fields := make(map[string]interface{})
	for name, field := range d.fields {
		fields[name] = field.Value
	}

	return json.Marshal(fields)
}

// UnmarshalJSON implements json.Unmarshaler interface
func (d *Document) UnmarshalJSON(data []byte) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	// Initialize fields map if not already initialized
	if d.fields == nil {
		d.fields = make(map[string]Field)
	}

	// Unmarshal into a temporary map
	var fields map[string]interface{}
	if err := json.Unmarshal(data, &fields); err != nil {
		return err
	}

	// Convert each field into a Document Field
	for name, value := range fields {
		// Determine field type based on the value
		var fieldType FieldType
		switch value.(type) {
		case string:
			fieldType = StringType
		case float64:
			fieldType = FloatType
		case int, int64:
			fieldType = IntType
		default:
			return fmt.Errorf("unsupported field type for field %s", name)
		}

		// Add the field to the document
		d.fields[name] = Field{
			Name:  name,
			Type:  fieldType,
			Value: value,
		}
	}

	return nil
}
