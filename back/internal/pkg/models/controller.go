package models

// TODO
type RawModel interface {
	ToModel() Model
}

type Model interface {
	ToJson() ([]byte, error)
}
