package registry_test

import (
	"errors"
	"net/http"
	"service-discovery/internal/registry"
	"service-discovery/internal/registry/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestShouldNotReturnErrorWhenRepoistoryRegisterSuccessUnit(t *testing.T) {
	assert := assert.New(t)

	// given:
	repo := registry.NewCacheServiceRegistry()
	SUT := registry.New(registry.WithRepository(repo))
	serviceName := "dummy"
	exepctedInstances := []registry.ServiceInstance{
		registry.NewServiceInstance(serviceName, "localhost", 1),
	}

	// when:
	err := SUT.Register(exepctedInstances...)

	// then:
	assert.Nil(err)

	actualInstances, err := SUT.QueryInstances(serviceName)
	assert.Nil(err)
	assert.Equal(exepctedInstances, actualInstances)
}

func TestShouldReturnErrorWhenRepoistoryRegisterFailureUnit(t *testing.T) {
	assert := assert.New(t)

	// given:
	repo := &mocks.ServiceRegistry{}
	SUT := registry.New(registry.WithRepository(repo))
	serviceName := "dummy"
	exepctedInstances := []registry.ServiceInstance{
		registry.NewServiceInstance(serviceName, "localhost", 1),
	}
	repo.EXPECT().Register(mock.Anything).Return(errors.New("repository is down"))

	// when:
	err := SUT.Register(exepctedInstances...)

	// then:
	assert.ErrorIs(err, registry.ErrRepositoryFailure)
	repo.AssertExpectations(t)
}

func TestShouldNotReturnErrorForServiceInstanceWhenRegisterUnit(t *testing.T) {
	assert := assert.New(t)

	// given:
	cache := registry.NewCacheServiceRegistry()
	SUT := registry.New(registry.WithRepository(cache))
	instances := []registry.ServiceInstance{
		registry.NewServiceInstance("dummy", "localhost", 1),
	}

	// when:
	err := SUT.Register(instances...)

	// then:
	assert.Nil(err)
}

func TestShouldReturnErrorForEmptyInstancesWhenRegisterUnit(t *testing.T) {
	assert := assert.New(t)

	// given:
	cache := registry.NewCacheServiceRegistry()
	SUT := registry.New(registry.WithRepository(cache))
	var instances []registry.ServiceInstance

	// when:
	err := SUT.Register(instances...)

	// then:
	assert.ErrorIs(err, registry.ErrEmptyServiceInstances)
}

func TestShouldReturnErrorForMalformedInstanceIPAddrWhenRegisterUnit(t *testing.T) {
	assert := assert.New(t)

	// given:
	cache := registry.NewCacheServiceRegistry()
	SUT := registry.New(registry.WithRepository(cache))
	instances := []registry.ServiceInstance{
		registry.NewServiceInstance("dummy", "12345", 9090),
	}

	// when:
	err := SUT.Register(instances...)

	// then:
	assert.ErrorIs(err, registry.ErrMalformedData)
}

func TestShouldReturnErrorForMalformedInstanceNameWhenRegisterUnit(t *testing.T) {
	assert := assert.New(t)

	// given:
	cache := registry.NewCacheServiceRegistry()
	SUT := registry.New(registry.WithRepository(cache))
	instances := []registry.ServiceInstance{
		registry.NewServiceInstance("", "12345", 9090),
	}

	// when:
	err := SUT.Register(instances...)

	// then:
	assert.ErrorIs(err, registry.ErrMalformedData)
}

func TestShouldReturnErrorForGreaterInstancePortThan65536WhenRegisterUnit(t *testing.T) {
	assert := assert.New(t)

	// given:
	cache := registry.NewCacheServiceRegistry()
	service := registry.New(registry.WithRepository(cache))
	instances := []registry.ServiceInstance{
		registry.NewServiceInstance("dummy", "localhost", 65536),
	}

	// when:
	err := service.Register(instances...)

	// then:
	assert.ErrorIs(err, registry.ErrMalformedData)
}

func TestShouldReturnErrorForInstancePortEqualToZeroWhenRegisterUnit(t *testing.T) {
	assert := assert.New(t)

	// given:
	cache := registry.NewCacheServiceRegistry()
	service := registry.New(registry.WithRepository(cache))
	instances := []registry.ServiceInstance{
		registry.NewServiceInstance("dummy", "localhost", 0),
	}

	// when:
	err := service.Register(instances...)

	// then:
	assert.ErrorIs(err, registry.ErrMalformedData)
}

func TestShouldMarkServiceInstanceAsHealthyUnit(t *testing.T) {
	assert := assert.New(t)

	// given:
	cache := registry.NewCacheServiceRegistry()
	httpCli := &mocks.HTTPClient{}
	serviceName := "dummy"
	service := registry.New(
		registry.WithRepository(cache),
		registry.WithHTTPClient(httpCli),
		registry.WithHealthz(registry.ServiceHealthEndpoints{
			serviceName: "/v1/test",
		}),
	)
	instance := registry.NewServiceInstance(serviceName, "localhost", 8080)

	expectedInstance := instance
	expectedInstance.SetHealth(true)
	expectedInstances := []registry.ServiceInstance{expectedInstance}

	service.Register(instance)

	instance.SetHealth(true)
	httpCli.EXPECT().Get(mock.Anything).Return(&http.Response{StatusCode: http.StatusOK}, nil)

	// when:
	service.HeartBeat()

	// then:
	actualInstances, err := service.QueryInstances(serviceName)
	assert.Nil(err)

	assert.Equal(expectedInstances, actualInstances)
	httpCli.AssertExpectations(t)

}

func TestShouldMarkSingleServiceInstanceAsNotHealthyUnit(t *testing.T) {
	assert := assert.New(t)

	// given:
	cache := registry.NewCacheServiceRegistry()
	httpCli := &mocks.HTTPClient{}
	serviceName := "dummy"
	service := registry.New(
		registry.WithRepository(cache),
		registry.WithHTTPClient(httpCli),
		registry.WithHealthz(registry.ServiceHealthEndpoints{
			serviceName: "/v1/test",
		}),
	)

	first := registry.NewServiceInstance(serviceName, "localhost", 8080)
	first.SetHealth(true)
	second := registry.NewServiceInstance(serviceName, "localhost", 8090)
	second.SetHealth(true)

	expectedInstance := second
	expectedInstance.SetHealth(false)

	expectedInstances := []registry.ServiceInstance{first, expectedInstance}

	service.Register(first, second)

	second.SetHealth(true)

	httpCli.EXPECT().Get(mock.Anything).Return(&http.Response{StatusCode: http.StatusOK}, nil).Once()
	httpCli.EXPECT().Get(mock.Anything).Return(&http.Response{StatusCode: http.StatusInternalServerError}, nil).Once()

	// when:
	service.HeartBeat()

	// then:
	actualInstances, err := service.QueryInstances(serviceName)
	assert.Nil(err)

	assert.Equal(expectedInstances, actualInstances)
	httpCli.AssertExpectations(t)
}

func TestShouldLogErrorWhenServiceInstanceUpdateStatusFailure(t *testing.T) {
	// given:
	registryRepository := &mocks.ServiceRegistry{}
	httpCli := &mocks.HTTPClient{}
	logger := &mocks.Logger{}
	serviceName := "dummy"
	service := registry.New(
		registry.WithRepository(registryRepository),
		registry.WithHTTPClient(httpCli),
		registry.WithLogger(logger),
		registry.WithHealthz(registry.ServiceHealthEndpoints{
			serviceName: "/v1/test",
		}),
	)
	instance := registry.NewServiceInstance(serviceName, "localhost", 8080)
	registryRepository.EXPECT().QueryInstances(mock.Anything).Return([]registry.ServiceInstance{instance}, nil)
	registryRepository.EXPECT().UpdateStatus(mock.Anything).Return(errors.New("repository is down"))
	httpCli.EXPECT().Get(mock.Anything).Return(&http.Response{StatusCode: http.StatusOK}, nil)
	logger.EXPECT().Errorf(mock.Anything, mock.Anything).Once()

	// when:
	service.HeartBeat()

	// then:
	mock.AssertExpectationsForObjects(t, registryRepository, logger, httpCli)
}
