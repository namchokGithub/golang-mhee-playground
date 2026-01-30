package di

import "go.uber.org/zap"

type Deps struct {
	Log *zap.Logger
	// DB *gorm.DB
	// Redis *redis.Client
	// Config Config
}

func (d *Deps) Logger() *zap.Logger { return d.Log }
