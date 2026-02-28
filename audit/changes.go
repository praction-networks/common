package audit

import (
	"fmt"
	"reflect"
)

// DiffChanges compares two structs (before/after) and returns a list of Changes
// for fields that differ. Uses json struct tags for field names if available,
// otherwise uses the Go field name.
//
// Usage:
//
//	changes := audit.DiffChanges(oldSubscriber, newSubscriber)
//	// → [{Field:"mobile", OldValue:"0501234567", NewValue:"0509876543"}, ...]
func DiffChanges(before, after any) []Change {
	if before == nil || after == nil {
		return nil
	}

	beforeVal := reflect.ValueOf(before)
	afterVal := reflect.ValueOf(after)

	// Dereference pointers
	if beforeVal.Kind() == reflect.Ptr {
		beforeVal = beforeVal.Elem()
	}
	if afterVal.Kind() == reflect.Ptr {
		afterVal = afterVal.Elem()
	}

	// Must be same type
	if beforeVal.Type() != afterVal.Type() {
		return nil
	}

	if beforeVal.Kind() != reflect.Struct {
		return nil
	}

	var changes []Change
	t := beforeVal.Type()

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		// Skip unexported fields
		if !field.IsExported() {
			continue
		}

		// Skip fields with `audit:"-"` tag
		if tag := field.Tag.Get("audit"); tag == "-" {
			continue
		}

		beforeField := beforeVal.Field(i)
		afterField := afterVal.Field(i)

		// Compare values
		beforeStr := formatValue(beforeField)
		afterStr := formatValue(afterField)

		if beforeStr != afterStr {
			// Use json tag for field name if available
			fieldName := field.Name
			if jsonTag := field.Tag.Get("json"); jsonTag != "" && jsonTag != "-" {
				// Parse json tag (take first part before comma)
				for j := 0; j < len(jsonTag); j++ {
					if jsonTag[j] == ',' {
						jsonTag = jsonTag[:j]
						break
					}
				}
				if jsonTag != "" {
					fieldName = jsonTag
				}
			}

			changes = append(changes, Change{
				Field:    fieldName,
				OldValue: beforeStr,
				NewValue: afterStr,
			})
		}
	}

	return changes
}

// formatValue converts a reflect.Value to its string representation
func formatValue(v reflect.Value) string {
	if !v.IsValid() {
		return ""
	}

	// Handle nil pointers/interfaces
	if v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface {
		if v.IsNil() {
			return ""
		}
		v = v.Elem()
	}

	return fmt.Sprintf("%v", v.Interface())
}
