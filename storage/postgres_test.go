package storage_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"www.github.com/techbrolakes/go-fiber-postgres/storage"
)

func TestNewConnectionSuccess(t *testing.T) {
	config := &storage.Config{
		Host:     "localhost",
		Port:     "5432",
		Password: "1001",
		User:     "postgres",
		SSLMode:  "disable",
		DBName:   "fiber_demo",
	}

	dbCon, err := storage.NewConnection(config)
	assert.Nil(t, err)
	assert.NotNil(t, dbCon)
}

func TestNewConnectionFail(t *testing.T) {
	config := &storage.Config{
		Host:     "localhost",
		Port:     "5432",
		Password: "1001",
		User:     "postgres",
		SSLMode:  "disable",
		DBName:   "fiber_demo1",
	}

	dbCon, err := storage.NewConnection(config)
	assert.NotNil(t, err)
	assert.NotNil(t, dbCon)
}
