package deprecation

import (
	"context"
	"encoding/json"
	"errors"
	"time"
)

const (
	ContextFieldDeprecationsKey = "fieldDeprecations"
)

type FieldDeprecations struct {
	messages map[string]map[string]string
}

func (f *FieldDeprecations) Add(field, message string, sunsetDate *time.Time) error {
	now := time.Now().UTC()
	if sunsetDate != nil && now.After(*sunsetDate) {
		return errors.New(message)
	}

	f.messages[field] = map[string]string{
		"details":    message,
		"sunsetDate": sunsetDate.Format(time.RFC3339),
	}
	return nil
}

func (f *FieldDeprecations) ToJSON() ([]byte, error) {
	output := make(map[string]string)
	for field, message := range f.messages {
		output[field] = message["details"]
	}
	return json.Marshal(output)
}

func (f *FieldDeprecations) IsEmpty() bool {
	return len(f.messages) == 0
}

func NewFieldDeprecations() FieldDeprecations {
	return FieldDeprecations{messages: make(map[string]map[string]string)}
}

func WithFieldDeprecations(ctx context.Context) context.Context {
	return context.WithValue(ctx, ContextFieldDeprecationsKey, NewFieldDeprecations())
}

func GetFieldDeprecations(ctx context.Context) FieldDeprecations {
	fieldDeprecations, ok := ctx.Value(ContextFieldDeprecationsKey).(FieldDeprecations)
	if !ok {
		return NewFieldDeprecations()
	}
	return fieldDeprecations
}
