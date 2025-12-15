package db

import "github.com/m11ano/budget_planner/backend/auth/pkg/pgclient"

type MasterClient interface {
	pgclient.Client
}
