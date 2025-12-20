package pgclient

import "context"

type mockImpl struct{}

func NewMock() *mockImpl {
	return &mockImpl{}
}

var _ Client = (*mockImpl)(nil)

func (m *mockImpl) ServerID() string {
	return "sid"
}

func (m *mockImpl) Pool() Pool {
	return nil
}

func (m *mockImpl) GetConn(ctx context.Context) Conn {
	return nil
}

func (m *mockImpl) Do(ctx context.Context, fn func(context.Context) error) error {
	fn(ctx)
	return nil
}

func (m *mockImpl) DoWithIsoLvl(ctx context.Context, isoLvl TxIsoLevel, fn func(context.Context) error) error {
	fn(ctx)
	return nil
}

func (m *mockImpl) Close() {}
