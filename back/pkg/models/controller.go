package models

type Model interface {
	// TODO
	// ToModel method recieves any value and tries to map it to
	// model that calls this method. If v == nil, returns nil-value error.
	// ToModel(v any) error

	// ToJson returns raw-data json byte slice, or error, if marshaling fails
	ToJson() ([]byte, error)
}
