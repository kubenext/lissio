package configuration

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kubenext/lissio/internal/controllers"
	"github.com/kubenext/lissio/internal/log"
	"github.com/kubenext/lissio/internal/testutil"
	"github.com/kubenext/lissio/pkg/action"
	actionFake "github.com/kubenext/lissio/pkg/action/fake"
	"github.com/kubenext/lissio/pkg/store"
	storeFake "github.com/kubenext/lissio/pkg/store/fake"
)

func TestObjectDeleter_ActionName(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	objectStore := storeFake.NewMockStore(controller)

	logger := log.NopLogger()

	d := NewObjectDeleter(logger, objectStore)
	require.Equal(t, controllers.ActionDeleteObject, d.ActionName())
}

func TestObjectDeleter_Handle(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	objectStore := storeFake.NewMockStore(controller)
	alerter := actionFake.NewMockAlerter(controller)

	pod := testutil.CreatePod("pod")
	key, err := store.KeyFromObject(pod)
	require.NoError(t, err)

	objectStore.EXPECT().
		Delete(gomock.Any(), key).
		Return(nil)

	alerter.EXPECT().
		SendAlert(gomock.Any()).
		DoAndReturn(func(alert action.Alert) {
			assert.Equal(t, action.AlertTypeInfo, alert.Type)
			assert.Equal(t, `Deleted Pod "pod"`, alert.Message)
			assert.NotNil(t, alert.Expiration)
		})

	logger := log.NopLogger()

	d := NewObjectDeleter(logger, objectStore)

	ctx := context.Background()

	err = d.Handle(ctx, alerter, key.ToActionPayload())
	require.NoError(t, err)
}
