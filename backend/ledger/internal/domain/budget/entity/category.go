package entity

import (
	"strings"
	"time"

	appErrors "github.com/m11ano/budget_planner/backend/ledger/internal/app/errors"
)

var ErrTitleInvalid = appErrors.ErrBadRequest.Extend("title is invalid")

type Category struct {
	ID    uint64
	Title string

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

func (item *Category) Version() int64 {
	return item.UpdatedAt.UnixMicro()
}

func (item *Category) SetTitle(value string) error {
	value = strings.TrimSpace(value)

	if value == "" {
		return ErrTitleInvalid
	}

	item.Title = value

	return nil
}

func NewCategory(
	title string,
) (*Category, error) {
	timeNow := time.Now().Truncate(time.Microsecond)

	item := &Category{
		CreatedAt: timeNow,
		UpdatedAt: timeNow,
	}

	err := item.SetTitle(title)
	if err != nil {
		return nil, err
	}

	return item, nil
}
