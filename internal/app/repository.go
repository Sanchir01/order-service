package app

import (
	db "github.com/Sanchir01/order-service/pkg/database"
	"log/slog"
)

type Repositories struct {
}

func NewRepositories(databases *db.Database, l *slog.Logger) *Repositories {
	return &Repositories{}
}
