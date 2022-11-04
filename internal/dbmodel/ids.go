package dbmodel

import "github.com/google/uuid"

func (b Bot) GetID() uuid.UUID {
	return b.ID
}

func (t Target) GetID() uuid.UUID {
	return t.ID
}
