package di

import "go.uber.org/zap"

type HasLogger interface {
	Logger() *zap.Logger
}
