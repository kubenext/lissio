package loading

import (
	"context"

	"github.com/kubenext/lissio/pkg/store"
)

// IsObjectLoading returns true if objects described by a key are loading.
func IsObjectLoading(ctx context.Context, namespace string, key store.Key, objectStore store.Store) bool {
	key.Namespace = namespace
	return objectStore.IsLoading(ctx, key)
}
