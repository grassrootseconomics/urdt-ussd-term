package testutil

import (
	"context"
	"errors"

	"git.defalsify.org/vise.git/persist"
	"git.defalsify.org/vise.git/resource"
	"git.defalsify.org/vise.git/db"
)

// TestStorageService wraps db for nats subscription.
type TestStorageService struct {
	Store  db.Db
}

// GetUserdataDb implements urdt/ussd/common.StorageServices.
func(ss *TestStorageService) GetUserdataDb(ctx context.Context) (db.Db, error) {
	return ss.Store, nil
}

// GetPersister implements urdt/ussd/common.StorageServices.
func(ts *TestStorageService) GetPersister(ctx context.Context) (*persist.Persister, error) {
	return nil, errors.New("not implemented")
}

// GetResource implements urdt/ussd/common.StorageServices.
func(ts *TestStorageService) GetResource(ctx context.Context) (resource.Resource, error) {
	return nil, errors.New("not implemented")
}

// EnsureDbDir implements urdt/ussd/common.StorageServices.
func(ss *TestStorageService) EnsureDbDir() error {
	return nil
}
