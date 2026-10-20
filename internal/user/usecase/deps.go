package usecase

import (
	usermodels "github.com/ardanlabs/service/internal/user/models"
)

// Storer declares the behavior this package needs to persist and retrieve data.
//
// NOTE(AI): alias of usermodels.Storer (migration step (b2)): the repository's
// NewWithTx must return this exact type, so it lives where both layers can see it.
type Storer = usermodels.Storer
