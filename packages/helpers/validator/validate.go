package validator

import (
	"fmt"
	"reflect"
)

// ValidateFields performs validation on struct fields using reflection.
// It recursively checks if any field in the struct is zero-valued.
// For nested structs, it performs the same validation on their fields.
//
// Parameters:
//   - i: interface{} - The struct to validate. Can be a pointer to struct or struct value.
//
// Returns:
//   - error - Returns nil if all fields are non-zero,
//     otherwise returns an error indicating which field is zero-valued.
//
// Example:
//
//	type Person struct {
//	    Name    string
//	    Age     int
//	    Address struct {
//	        Street string
//	    }
//	}
//
//	err := ValidateFields(person)
//	if err != nil {
//	    fmt.Println(err)
//	}
func ValidateFields(i interface{}) error {
	v := reflect.ValueOf(i)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		if field.Kind() == reflect.Struct {
			if err := ValidateFields(field.Interface()); err != nil {
				return err
			}
		} else if field.IsZero() {
			return fmt.Errorf("field %s is required", t.Field(i).Name)
		}
	}
	return nil
}
