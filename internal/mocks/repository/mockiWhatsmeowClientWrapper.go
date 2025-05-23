// Code generated by mockery; DO NOT EDIT.
// github.com/vektra/mockery
// template: testify

package repository

import (
	"context"

	mock "github.com/stretchr/testify/mock"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types"
)

// newMockiWhatsmeowClientWrapper creates a new instance of mockiWhatsmeowClientWrapper. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func newMockiWhatsmeowClientWrapper(t interface {
	mock.TestingT
	Cleanup(func())
}) *mockiWhatsmeowClientWrapper {
	mock := &mockiWhatsmeowClientWrapper{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

// mockiWhatsmeowClientWrapper is an autogenerated mock type for the iWhatsmeowClientWrapper type
type mockiWhatsmeowClientWrapper struct {
	mock.Mock
}

type mockiWhatsmeowClientWrapper_Expecter struct {
	mock *mock.Mock
}

func (_m *mockiWhatsmeowClientWrapper) EXPECT() *mockiWhatsmeowClientWrapper_Expecter {
	return &mockiWhatsmeowClientWrapper_Expecter{mock: &_m.Mock}
}

// Connect provides a mock function for the type mockiWhatsmeowClientWrapper
func (_mock *mockiWhatsmeowClientWrapper) Connect() error {
	ret := _mock.Called()

	if len(ret) == 0 {
		panic("no return value specified for Connect")
	}

	var r0 error
	if returnFunc, ok := ret.Get(0).(func() error); ok {
		r0 = returnFunc()
	} else {
		r0 = ret.Error(0)
	}
	return r0
}

// mockiWhatsmeowClientWrapper_Connect_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Connect'
type mockiWhatsmeowClientWrapper_Connect_Call struct {
	*mock.Call
}

// Connect is a helper method to define mock.On call
func (_e *mockiWhatsmeowClientWrapper_Expecter) Connect() *mockiWhatsmeowClientWrapper_Connect_Call {
	return &mockiWhatsmeowClientWrapper_Connect_Call{Call: _e.mock.On("Connect")}
}

func (_c *mockiWhatsmeowClientWrapper_Connect_Call) Run(run func()) *mockiWhatsmeowClientWrapper_Connect_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *mockiWhatsmeowClientWrapper_Connect_Call) Return(err error) *mockiWhatsmeowClientWrapper_Connect_Call {
	_c.Call.Return(err)
	return _c
}

func (_c *mockiWhatsmeowClientWrapper_Connect_Call) RunAndReturn(run func() error) *mockiWhatsmeowClientWrapper_Connect_Call {
	_c.Call.Return(run)
	return _c
}

// Disconnect provides a mock function for the type mockiWhatsmeowClientWrapper
func (_mock *mockiWhatsmeowClientWrapper) Disconnect() {
	_mock.Called()
	return
}

// mockiWhatsmeowClientWrapper_Disconnect_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Disconnect'
type mockiWhatsmeowClientWrapper_Disconnect_Call struct {
	*mock.Call
}

// Disconnect is a helper method to define mock.On call
func (_e *mockiWhatsmeowClientWrapper_Expecter) Disconnect() *mockiWhatsmeowClientWrapper_Disconnect_Call {
	return &mockiWhatsmeowClientWrapper_Disconnect_Call{Call: _e.mock.On("Disconnect")}
}

func (_c *mockiWhatsmeowClientWrapper_Disconnect_Call) Run(run func()) *mockiWhatsmeowClientWrapper_Disconnect_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *mockiWhatsmeowClientWrapper_Disconnect_Call) Return() *mockiWhatsmeowClientWrapper_Disconnect_Call {
	_c.Call.Return()
	return _c
}

func (_c *mockiWhatsmeowClientWrapper_Disconnect_Call) RunAndReturn(run func()) *mockiWhatsmeowClientWrapper_Disconnect_Call {
	_c.Run(run)
	return _c
}

// GetLoggedInDeviceJID provides a mock function for the type mockiWhatsmeowClientWrapper
func (_mock *mockiWhatsmeowClientWrapper) GetLoggedInDeviceJID() *types.JID {
	ret := _mock.Called()

	if len(ret) == 0 {
		panic("no return value specified for GetLoggedInDeviceJID")
	}

	var r0 *types.JID
	if returnFunc, ok := ret.Get(0).(func() *types.JID); ok {
		r0 = returnFunc()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.JID)
		}
	}
	return r0
}

// mockiWhatsmeowClientWrapper_GetLoggedInDeviceJID_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetLoggedInDeviceJID'
type mockiWhatsmeowClientWrapper_GetLoggedInDeviceJID_Call struct {
	*mock.Call
}

// GetLoggedInDeviceJID is a helper method to define mock.On call
func (_e *mockiWhatsmeowClientWrapper_Expecter) GetLoggedInDeviceJID() *mockiWhatsmeowClientWrapper_GetLoggedInDeviceJID_Call {
	return &mockiWhatsmeowClientWrapper_GetLoggedInDeviceJID_Call{Call: _e.mock.On("GetLoggedInDeviceJID")}
}

func (_c *mockiWhatsmeowClientWrapper_GetLoggedInDeviceJID_Call) Run(run func()) *mockiWhatsmeowClientWrapper_GetLoggedInDeviceJID_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *mockiWhatsmeowClientWrapper_GetLoggedInDeviceJID_Call) Return(jID *types.JID) *mockiWhatsmeowClientWrapper_GetLoggedInDeviceJID_Call {
	_c.Call.Return(jID)
	return _c
}

func (_c *mockiWhatsmeowClientWrapper_GetLoggedInDeviceJID_Call) RunAndReturn(run func() *types.JID) *mockiWhatsmeowClientWrapper_GetLoggedInDeviceJID_Call {
	_c.Call.Return(run)
	return _c
}

// GetQRChannel provides a mock function for the type mockiWhatsmeowClientWrapper
func (_mock *mockiWhatsmeowClientWrapper) GetQRChannel(ctx context.Context) (<-chan whatsmeow.QRChannelItem, error) {
	ret := _mock.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for GetQRChannel")
	}

	var r0 <-chan whatsmeow.QRChannelItem
	var r1 error
	if returnFunc, ok := ret.Get(0).(func(context.Context) (<-chan whatsmeow.QRChannelItem, error)); ok {
		return returnFunc(ctx)
	}
	if returnFunc, ok := ret.Get(0).(func(context.Context) <-chan whatsmeow.QRChannelItem); ok {
		r0 = returnFunc(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan whatsmeow.QRChannelItem)
		}
	}
	if returnFunc, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = returnFunc(ctx)
	} else {
		r1 = ret.Error(1)
	}
	return r0, r1
}

// mockiWhatsmeowClientWrapper_GetQRChannel_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetQRChannel'
type mockiWhatsmeowClientWrapper_GetQRChannel_Call struct {
	*mock.Call
}

// GetQRChannel is a helper method to define mock.On call
//   - ctx
func (_e *mockiWhatsmeowClientWrapper_Expecter) GetQRChannel(ctx interface{}) *mockiWhatsmeowClientWrapper_GetQRChannel_Call {
	return &mockiWhatsmeowClientWrapper_GetQRChannel_Call{Call: _e.mock.On("GetQRChannel", ctx)}
}

func (_c *mockiWhatsmeowClientWrapper_GetQRChannel_Call) Run(run func(ctx context.Context)) *mockiWhatsmeowClientWrapper_GetQRChannel_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *mockiWhatsmeowClientWrapper_GetQRChannel_Call) Return(qRChannelItemCh <-chan whatsmeow.QRChannelItem, err error) *mockiWhatsmeowClientWrapper_GetQRChannel_Call {
	_c.Call.Return(qRChannelItemCh, err)
	return _c
}

func (_c *mockiWhatsmeowClientWrapper_GetQRChannel_Call) RunAndReturn(run func(ctx context.Context) (<-chan whatsmeow.QRChannelItem, error)) *mockiWhatsmeowClientWrapper_GetQRChannel_Call {
	_c.Call.Return(run)
	return _c
}

// IsLoggedIn provides a mock function for the type mockiWhatsmeowClientWrapper
func (_mock *mockiWhatsmeowClientWrapper) IsLoggedIn() bool {
	ret := _mock.Called()

	if len(ret) == 0 {
		panic("no return value specified for IsLoggedIn")
	}

	var r0 bool
	if returnFunc, ok := ret.Get(0).(func() bool); ok {
		r0 = returnFunc()
	} else {
		r0 = ret.Get(0).(bool)
	}
	return r0
}

// mockiWhatsmeowClientWrapper_IsLoggedIn_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'IsLoggedIn'
type mockiWhatsmeowClientWrapper_IsLoggedIn_Call struct {
	*mock.Call
}

// IsLoggedIn is a helper method to define mock.On call
func (_e *mockiWhatsmeowClientWrapper_Expecter) IsLoggedIn() *mockiWhatsmeowClientWrapper_IsLoggedIn_Call {
	return &mockiWhatsmeowClientWrapper_IsLoggedIn_Call{Call: _e.mock.On("IsLoggedIn")}
}

func (_c *mockiWhatsmeowClientWrapper_IsLoggedIn_Call) Run(run func()) *mockiWhatsmeowClientWrapper_IsLoggedIn_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *mockiWhatsmeowClientWrapper_IsLoggedIn_Call) Return(b bool) *mockiWhatsmeowClientWrapper_IsLoggedIn_Call {
	_c.Call.Return(b)
	return _c
}

func (_c *mockiWhatsmeowClientWrapper_IsLoggedIn_Call) RunAndReturn(run func() bool) *mockiWhatsmeowClientWrapper_IsLoggedIn_Call {
	_c.Call.Return(run)
	return _c
}
