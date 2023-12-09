package server

type StatusCode = int32

const (
	StatusOK StatusCode = iota
	StatusInternalServerError
	StatusBadRequest
)
