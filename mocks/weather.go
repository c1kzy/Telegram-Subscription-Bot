// Code generated by MockGen. DO NOT EDIT.
// Source: subscriptionbot/weather (interfaces: WeatherService)

// Package mocks is a generated GoMock package.
package mocks

import (
	url "net/url"
	reflect "reflect"
	db "subscriptionbot/db"

	gomock "github.com/golang/mock/gomock"
)

// WeatherService is a mock of WeatherService interface.
type WeatherService struct {
	ctrl     *gomock.Controller
	recorder *WeatherServiceMockRecorder
}

// WeatherServiceMockRecorder is the mock recorder for WeatherService.
type WeatherServiceMockRecorder struct {
	mock *WeatherService
}

// NewWeatherService creates a new mock instance.
func NewWeatherService(ctrl *gomock.Controller) *WeatherService {
	mock := &WeatherService{ctrl: ctrl}
	mock.recorder = &WeatherServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *WeatherService) EXPECT() *WeatherServiceMockRecorder {
	return m.recorder
}

// WeatherRequest mocks base method.
func (m *WeatherService) WeatherRequest(arg0 db.User) (url.Values, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WeatherRequest", arg0)
	ret0, _ := ret[0].(url.Values)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// WeatherRequest indicates an expected call of WeatherRequest.
func (mr *WeatherServiceMockRecorder) WeatherRequest(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WeatherRequest", reflect.TypeOf((*WeatherService)(nil).WeatherRequest), arg0)
}
