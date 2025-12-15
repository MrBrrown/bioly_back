package usecases

import (
	"bioly/common/asynclogger"
	"bioly/profileservice/internal/cache"
	"bioly/profileservice/internal/repositories"
	"bioly/profileservice/internal/types"
	"context"
	"fmt"
)

type ProfileService interface {
	GetProfile(ctx context.Context, username string) (types.Profile, error)
	GetProfileCached(ctx context.Context, username string) (types.Profile, error)
}

type profileImpl struct {
	profileRepo repositories.Profile
	cache       cache.ProfileCache
}

func NewProfile(profileRepo repositories.Profile, cache cache.ProfileCache) ProfileService {
	return &profileImpl{profileRepo: profileRepo, cache: cache}
}

func (p *profileImpl) GetProfile(ctx context.Context, username string) (types.Profile, error) {
	id, err := p.profileRepo.GetUserId(ctx, username)
	if err != nil {
		return types.Profile{}, err
	}

	profile, err := p.profileRepo.GetProfile(ctx, id)
	if err != nil {
		asynclogger.Error("failed to get profile %s/%d: %v", username, id, err)
		return types.Profile{}, fmt.Errorf("failed to get profile")
	}

	profile.Username = username
	return profile, nil
}

func (p *profileImpl) GetProfileCached(ctx context.Context, username string) (types.Profile, error) {
	if p.cache != nil {
		if cachedProfile, found := p.cache.GetProfile(username); found {
			return *cachedProfile, nil
		}
	}

	id, err := p.profileRepo.GetUserId(ctx, username)
	if err != nil {
		return types.Profile{}, err
	}

	profile, err := p.profileRepo.GetProfile(ctx, id)
	if err != nil {
		asynclogger.Error("failed to get profile %s/%d: %v", username, id, err)
		return types.Profile{}, fmt.Errorf("failed to get profile")
	}

	profile.Username = username

	if p.cache != nil {
		if err := p.cache.AddProfile(profile); err != nil {
			asynclogger.Error("failed to cache profile %s/%d: %v", username, id, err)
		}
	}

	return profile, nil
}
