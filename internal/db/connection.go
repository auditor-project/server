package db

import (
	"auditor.z9fr.xyz/server/internal/lib"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type NewDatabase struct {
	*gorm.DB
}

func NewDatabaseImpl(logger lib.Logger, env *lib.Env) NewDatabase {
	gormConfig := &gorm.Config{}

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  env.DATABASE_DSN,
		PreferSimpleProtocol: true,
	}), gormConfig)

	// err = db.AutoMigrate(&model.User{}, &model.Meta{}, &model.Role{}, &model.Purchase{}, &model.TenetUserMeta{}, &model.Badge{})

	if err != nil {
		logger.Panic(err)
	}

	return NewDatabase{DB: db}
}
