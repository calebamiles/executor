// This file was generated by counterfeiter
package fakes

import (
	. "github.com/cloudfoundry-incubator/executor/executor"

	"github.com/cloudfoundry-incubator/runtime-schema/models"

	"sync"
)

type FakeClient struct {
	RunContainerStub        func(guid string, actions []models.ExecutorAction, completeURL string) error
	runContainerMutex       sync.RWMutex
	runContainerArgsForCall []struct {
		arg1 string
		arg2 []models.ExecutorAction
		arg3 string
	}
	runContainerReturns struct {
		result1 error
	}
}

func (fake *FakeClient) RunContainer(arg1 string, arg2 []models.ExecutorAction, arg3 string) error {
	fake.runContainerMutex.Lock()
	defer fake.runContainerMutex.Unlock()
	fake.runContainerArgsForCall = append(fake.runContainerArgsForCall, struct {
		arg1 string
		arg2 []models.ExecutorAction
		arg3 string
	}{arg1, arg2, arg3})
	if fake.RunContainerStub != nil {
		return fake.RunContainerStub(arg1, arg2, arg3)
	} else {
		return fake.runContainerReturns.result1
	}
}

func (fake *FakeClient) RunContainerCallCount() int {
	fake.runContainerMutex.RLock()
	defer fake.runContainerMutex.RUnlock()
	return len(fake.runContainerArgsForCall)
}

func (fake *FakeClient) RunContainerArgsForCall(i int) (string, []models.ExecutorAction, string) {
	fake.runContainerMutex.RLock()
	defer fake.runContainerMutex.RUnlock()
	return fake.runContainerArgsForCall[i].arg1, fake.runContainerArgsForCall[i].arg2, fake.runContainerArgsForCall[i].arg3
}

func (fake *FakeClient) RunContainerReturns(result1 error) {
	fake.runContainerReturns = struct {
		result1 error
	}{result1}
}

var _ Client = new(FakeClient)
