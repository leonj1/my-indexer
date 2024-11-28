package document

import (
	"testing"
)

func TestNewDocument(t *testing.T) {
	doc := NewDocument()
	if doc == nil {
		t.Error("NewDocument() returned nil")
	}
	if doc.fields == nil {
		t.Error("Document fields map not initialized")
	}
}

func TestAddField(t *testing.T) {
	tests := []struct {
		name      string
		fieldName string
		value     interface{}
		wantErr   bool
	}{
		{"string field", "title", "test document", false},
		{"integer field", "count", 42, false},
		{"float field", "score", 3.14, false},
		{"invalid type", "invalid", []string{"test"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc := NewDocument()
			err := doc.AddField(tt.fieldName, tt.value)
			
			if (err != nil) != tt.wantErr {
				t.Errorf("AddField() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				field, err := doc.GetField(tt.fieldName)
				if err != nil {
					t.Errorf("GetField() error = %v", err)
					return
				}

				if field.Name != tt.fieldName {
					t.Errorf("Field name = %v, want %v", field.Name, tt.fieldName)
				}

				if field.Value != tt.value {
					t.Errorf("Field value = %v, want %v", field.Value, tt.value)
				}
			}
		})
	}
}

func TestGetField(t *testing.T) {
	doc := NewDocument()
	fieldName := "test"
	fieldValue := "test value"

	// Test getting non-existent field
	_, err := doc.GetField(fieldName)
	if err == nil {
		t.Error("GetField() should return error for non-existent field")
	}

	// Add field and test retrieval
	err = doc.AddField(fieldName, fieldValue)
	if err != nil {
		t.Errorf("AddField() error = %v", err)
	}

	field, err := doc.GetField(fieldName)
	if err != nil {
		t.Errorf("GetField() error = %v", err)
	}

	if field.Name != fieldName {
		t.Errorf("Field name = %v, want %v", field.Name, fieldName)
	}

	if field.Value != fieldValue {
		t.Errorf("Field value = %v, want %v", field.Value, fieldValue)
	}
}
