// Code generated by mockery v2.43.1. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	trace "go.opentelemetry.io/otel/trace"
)

// Telemetry is an autogenerated mock type for the Telemetry type
type Telemetry struct {
	mock.Mock
}

// Shutdown provides a mock function with given fields: ctx
func (_m *Telemetry) Shutdown(ctx context.Context) {
	_m.Called(ctx)
}

// Start provides a mock function with given fields: ctx, name, opts
func (_m *Telemetry) Start(ctx context.Context, name string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, name)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for Start")
	}

	var r0 context.Context
	var r1 trace.Span
	if rf, ok := ret.Get(0).(func(context.Context, string, ...trace.SpanStartOption) (context.Context, trace.Span)); ok {
		return rf(ctx, name, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, ...trace.SpanStartOption) context.Context); ok {
		r0 = rf(ctx, name, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(context.Context)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, ...trace.SpanStartOption) trace.Span); ok {
		r1 = rf(ctx, name, opts...)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(trace.Span)
		}
	}

	return r0, r1
}

// NewTelemetry creates a new instance of Telemetry. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewTelemetry(t interface {
	mock.TestingT
	Cleanup(func())
}) *Telemetry {
	mock := &Telemetry{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
