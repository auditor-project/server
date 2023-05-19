package db

import (
	"fmt"

	"auditor.z9fr.xyz/server/internal/lib"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Database struct {
	*gorm.DB
}

func NewDatabaseImpl(logger lib.Logger, env *lib.Env) *Database {
	logger.Debug("Init new database impl")
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s  sslmode=%s TimeZone=UTC",
		env.DB_HOST,
		env.DB_HOST,
		env.DB_HOST,
		env.DB_HOST,
		"disable")

	gormConfig := &gorm.Config{}

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true,
	}), gormConfig)

	// err = db.AutoMigrate(&model.User{}, &model.Meta{}, &model.Role{}, &model.Purchase{}, &model.TenetUserMeta{}, &model.Badge{})

	if err != nil {
		logger.Panic(err)
	}

	return &Database{DB: db}
}
