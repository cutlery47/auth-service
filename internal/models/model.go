package models

import (
	"github.com/Microsoft/go-winio/pkg/guid"
	"github.com/google/uuid"
)

type InRefresh struct {
	UserId guid.GUID
	Hash   []byte
	Salt   uuid.UUID
	Cost   int
}

type OutRefresh struct {
	Id uuid.UUID
	InRefresh
}
