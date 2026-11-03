package track

import "context"

// Repository defines track storage operations.
// Implementations must return ErrNotFound when an ID does not exist.
type Repository interface {
	Create(ctx context.Context, t Track) (Track, error)
	Get(ctx context.Context, id int64) (Track, error)
	List(ctx context.Context) ([]Track, error)
	Update(ctx context.Context, t Track) (Track, error)
	Delete(ctx context.Context, id int64) error
}
