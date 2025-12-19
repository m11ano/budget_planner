package uctypes

type QueryGetListParams struct {
	WithDeleted         bool
	ForShare            bool
	ForUpdate           bool
	ForUpdateSkipLocked bool
	ForUpdateNoWait     bool
	Limit               uint64
	Offset              uint64
}

type QueryGetOneParams struct {
	WithDeleted         bool
	ForShare            bool
	ForUpdate           bool
	ForUpdateSkipLocked bool
	ForUpdateNoWait     bool
}

type CompareType int

const (
	CompareEqual CompareType = iota
	CompareNotEqual
	CompareLess
	CompareMore
	CompareLessOrEqual
	CompareMoreOrEqual
)

type CompareOption[T comparable] struct {
	Value T
	Type  CompareType
}

type SortOption[T comparable] struct {
	Field  T
	IsDesc bool
}
