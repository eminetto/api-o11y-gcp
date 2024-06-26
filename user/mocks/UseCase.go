// Code generated by mockery v2.43.1. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// UseCase is an autogenerated mock type for the UseCase type
type UseCase struct {
	mock.Mock
}

// ValidateUser provides a mock function with given fields: ctx, email, password
func (_m *UseCase) ValidateUser(ctx context.Context, email string, password string) error {
	ret := _m.Called(ctx, email, password)

	if len(ret) == 0 {
		panic("no return value specified for ValidateUser")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) error); ok {
		r0 = rf(ctx, email, password)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewUseCase creates a new instance of UseCase. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewUseCase(t interface {
	mock.TestingT
	Cleanup(func())
}) *UseCase {
	mock := &UseCase{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
