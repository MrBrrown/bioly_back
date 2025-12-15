package cache

import "bioly/profileservice/internal/types"

type ProfileCache interface {
	GetProfile(username string) (*types.Profile, bool)
	AddProfile(profile types.Profile) error
}
