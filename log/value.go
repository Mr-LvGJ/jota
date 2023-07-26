package log

import (
	"go.uber.org/zap"
)

func zapFields(args ...interface{}) []zap.Field {
	if len(args) == 1 {
		fs, ok := args[0].([]interface{})
		if ok {
			args = fs
		}
	}

	fields := make([]zap.Field, 0, len(args)/2+1)
	for i := 0; i < len(args)-1; i += 2 {
		// Make sure this element isn't a dangling key.
		if i == len(args)-1 {
			fields = append(fields, zap.Any("exceeds", args[i]))
			break
		}

		// Consume this value and the next, treating them as a key-value pair. If the
		// key isn't a string, add this pair to the slice of invalid pairs.
		key, val := args[i], args[i+1]
		if keyStr, ok := key.(string); !ok {
			fields = append(fields, zap.Any("invalidKey", val))
		} else {
			fields = append(fields, zap.Any(keyStr, val))
		}
	}
	return fields
}
