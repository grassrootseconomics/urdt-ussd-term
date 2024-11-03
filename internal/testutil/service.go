package testutil

import (
	"context"
	"errors"

	"git.defalsify.org/vise.git/persist"
	"git.defalsify.org/vise.git/resource"
	"git.defalsify.org/vise.git/db"
)

type TestStorageService struct {
	Store  db.Db
}
	
func(ss *TestStorageService) GetUserdataDb(ctx context.Context) (db.Db, error) {
	return ss.Store, nil
}

func(ts *TestStorageService) GetPersister(ctx context.Context) (*persist.Persister, error) {
	return nil, errors.New("not implemented")
}

func(ts *TestStorageService) GetResource(ctx context.Context) (resource.Resource, error) {
	return nil, errors.New("not implemented")
}

func(ss *TestStorageService) EnsureDbDir() error {
	return nil
}
