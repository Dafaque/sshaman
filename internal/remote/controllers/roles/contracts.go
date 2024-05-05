package roles

import "context"

type repository interface {
	Create(ctx context.Context, role *Role) error
	Update(ctx context.Context, role *Role) error
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context) ([]Role, error)
	Get(ctx context.Context, ids ...int64) ([]Role, error)
}

type permissions interface {
	// Read(space string) bool
	// Write(space string) bool
	// Delete(space string) bool
	// Overwrite(space string) bool
	SU() bool
}
