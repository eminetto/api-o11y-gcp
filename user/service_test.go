package user_test

import (
	"context"
	"testing"

	tmocks "github.com/eminetto/api-o11y-gcp/internal/telemetry/mocks"
	"github.com/eminetto/api-o11y-gcp/user"
	"github.com/eminetto/api-o11y-gcp/user/mocks"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/trace"
)

func TestValidatePassword(t *testing.T) {
	ctx := context.TODO()
	repo := mocks.NewRepository(t)
	otel := tmocks.NewTelemetry(t)
	otel.On("Start", ctx, "validatePassword").Return(ctx, trace.SpanFromContext(ctx))

	s := user.NewService(repo, otel)
	u := &user.User{
		Email:    "eminetto@email.com",
		Password: "8cb2237d0679ca88db6464eac60da96345513964",
	}
	t.Run("invalid password", func(t *testing.T) {
		err := s.ValidatePassword(ctx, u, "invalid")
		assert.NotNil(t, err)
		assert.Equal(t, "invalid password", err.Error())
	})
	t.Run("valid password", func(t *testing.T) {
		err := s.ValidatePassword(ctx, u, "12345")
		assert.Nil(t, err)
	})
}
