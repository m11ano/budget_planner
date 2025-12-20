package category

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/m11ano/budget_planner/backend/ledger/internal/common/uctypes"
	"github.com/m11ano/budget_planner/backend/ledger/internal/domain/budget/entity"
	"github.com/m11ano/budget_planner/backend/ledger/internal/domain/budget/usecase"
)

func TestUsecase_Category_FindOneByID_Table(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx         context.Context
		id          uint64
		queryParams *uctypes.QueryGetOneParams
	}

	tests := []struct {
		name string
		args args

		setup func(t *testing.T, d *dependencies, a args)

		wantErr bool
	}{
		{
			name: "Positive_ok",
			args: args{
				ctx:         context.Background(),
				id:          10,
				queryParams: &uctypes.QueryGetOneParams{},
			},
			setup: func(t *testing.T, d *dependencies, a args) {
				t.Helper()

				expectedEntity := &entity.Category{}

				d.categoryRepo.
					FindOneByIDMock.
					Expect(a.ctx, a.id, a.queryParams).
					Inspect(func(ctx context.Context, id uint64, qp *uctypes.QueryGetOneParams) {
						require.Same(t, a.queryParams, qp)
					}).
					Return(expectedEntity, nil)
			},
			wantErr: false,
		},
		{
			name: "Negative_repo_error",
			args: args{
				ctx:         context.Background(),
				id:          11,
				queryParams: &uctypes.QueryGetOneParams{},
			},
			setup: func(t *testing.T, d *dependencies, a args) {
				t.Helper()

				d.categoryRepo.
					FindOneByIDMock.
					Expect(a.ctx, a.id, a.queryParams).
					Return(nil, errors.New("repo failed"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			d := newDependencies(t)
			t.Cleanup(func() { finishDependencies(t, d) })

			if tt.setup != nil {
				tt.setup(t, d, tt.args)
			}

			out, err := d.uc.FindOneByID(tt.args.ctx, tt.args.id, tt.args.queryParams)

			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, out)
			require.NotNil(t, out.Category)
		})
	}
}

func TestUsecase_Category_FindList_Table(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx         context.Context
		listOptions *usecase.CategoryListOptions
		queryParams *uctypes.QueryGetListParams
	}

	tests := []struct {
		name string
		args args

		setup func(t *testing.T, d *dependencies, a args)

		wantLen int
		wantErr bool
	}{
		{
			name: "Positive_ok_two_items",
			args: args{
				ctx: context.Background(),
				listOptions: &usecase.CategoryListOptions{
					FilterIDs: nil,
					Sort:      nil,
				},
				queryParams: &uctypes.QueryGetListParams{},
			},
			setup: func(t *testing.T, d *dependencies, a args) {
				t.Helper()

				items := []*entity.Category{
					{},
					{},
				}

				d.categoryRepo.
					FindListMock.
					Expect(a.ctx, a.listOptions, a.queryParams).
					Inspect(func(ctx context.Context, lo *usecase.CategoryListOptions, qp *uctypes.QueryGetListParams) {
						require.Same(t, a.listOptions, lo)
						require.Same(t, a.queryParams, qp)
					}).
					Return(items, nil)
			},
			wantLen: 2,
			wantErr: false,
		},
		{
			name: "Positive_ok_empty",
			args: args{
				ctx:         context.Background(),
				listOptions: &usecase.CategoryListOptions{},
				queryParams: &uctypes.QueryGetListParams{},
			},
			setup: func(t *testing.T, d *dependencies, a args) {
				t.Helper()

				d.categoryRepo.
					FindListMock.
					Expect(a.ctx, a.listOptions, a.queryParams).
					Return([]*entity.Category{}, nil)
			},
			wantLen: 0,
			wantErr: false,
		},
		{
			name: "Negative_repo_error",
			args: args{
				ctx:         context.Background(),
				listOptions: &usecase.CategoryListOptions{},
				queryParams: &uctypes.QueryGetListParams{},
			},
			setup: func(t *testing.T, d *dependencies, a args) {
				t.Helper()

				d.categoryRepo.
					FindListMock.
					Expect(a.ctx, a.listOptions, a.queryParams).
					Return(nil, errors.New("repo failed"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			d := newDependencies(t)
			t.Cleanup(func() { finishDependencies(t, d) })

			if tt.setup != nil {
				tt.setup(t, d, tt.args)
			}

			out, err := d.uc.FindList(tt.args.ctx, tt.args.listOptions, tt.args.queryParams)

			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Len(t, out, tt.wantLen)
			for _, dto := range out {
				require.NotNil(t, dto)
				require.NotNil(t, dto.Category)
			}
		})
	}
}
