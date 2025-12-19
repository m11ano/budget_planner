package pg

import (
	"time"

	"github.com/m11ano/budget_planner/backend/ledger/internal/domain/budget/entity"
	"github.com/m11ano/budget_planner/backend/ledger/pkg/dbhelper"
)

const (
	CategoryTable = "category"
)

var CategoryTableFields = []string{}

func init() {
	CategoryTableFields = dbhelper.ExtractDBFields(&CategoryDBModel{})
}

type CategoryDBModel struct {
	ID    uint64 `db:"id"`
	Title string `db:"title"`

	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt time.Time  `db:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at"`
}

func (db *CategoryDBModel) ToEntity() *entity.Category {
	return &entity.Category{
		ID:    db.ID,
		Title: db.Title,

		CreatedAt: db.CreatedAt,
		UpdatedAt: db.UpdatedAt,
		DeletedAt: db.DeletedAt,
	}
}

func MapCategoryEntityToDBModel(entity *entity.Category) *CategoryDBModel {
	return &CategoryDBModel{
		ID:    entity.ID,
		Title: entity.Title,

		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
		DeletedAt: entity.DeletedAt,
	}
}
