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

func TestShouldReturnErrorWhenRepoistoryFailureUnit(t *testing.T) {
	assert := assert.New(t)

	// given:
	repo := &mocks.ServiceRegistry{}
	SUT := registry.NewService(registry.WithRepository(repo))
	component := "service1"
	exepctedInstances := []registry.ServiceInstance{
		registry.NewServiceInstance(component, "node1", "localhost", "1"),
	}
	repo.EXPECT().Register(mock.Anything).Return(errors.New("repository is down"))

	// when:
	err := SUT.Register(exepctedInstances...)

	// then:
	assert.ErrorIs(err, registry.ErrRepositoryFailure)
	repo.AssertExpectations(t)
}

func TestShouldRegisterSeviceInstanceWithSuccessUnit(t *testing.T) {
	assert := assert.New(t)

	// given:
	repo := registry.NewCacheRepository()
	SUT := registry.NewService(registry.WithRepository(repo))
	component := "dummy"
	exepctedInstances := []registry.ServiceInstance{
		registry.NewServiceInstance(component, "node1", "localhost", "1"),
	}

	// when:
	err := SUT.Register(exepctedInstances...)

	// then:
	assert.Nil(err)

	actualInstances, err := SUT.QueryInstances(component)
	assert.Nil(err)
	assert.Equal(exepctedInstances, actualInstances)
}

func TestShouldReturnErrorForEmptyInstanceUnit(t *testing.T) {
	assert := assert.New(t)

	// given:
	cache := registry.NewCacheRepository()
	SUT := registry.NewService(registry.WithRepository(cache))
	var instances []registry.ServiceInstance

	// when:
	err := SUT.Register(instances...)

	// then:
	assert.ErrorIs(err, registry.ErrEmptyServiceInstances)
}

func TestShouldReturnErrorWhenIPAddrIsMalformedUnit(t *testing.T) {
	assert := assert.New(t)

	// given:
	cache := registry.NewCacheRepository()
	SUT := registry.NewService(registry.WithRepository(cache))
	instances := []registry.ServiceInstance{
		registry.NewServiceInstance("service", "node1", "12345", "9090"),
	}

	// when:
	err := SUT.Register(instances...)

	// then:
	assert.ErrorIs(err, registry.ErrMissingData)
}

func TestShouldReturnErrorWhencomponentNameIsEmptyUnit(t *testing.T) {
	assert := assert.New(t)

	// given:
	cache := registry.NewCacheRepository()
	SUT := registry.NewService(registry.WithRepository(cache))
	instances := []registry.ServiceInstance{
		registry.NewServiceInstance("", "node1", "127.0.0.1", "9090"),
	}

	// when:
	err := SUT.Register(instances...)

	// then:
	assert.ErrorIs(err, registry.ErrMissingData)
}

func TestShouldReturnErrorForMalformedInstanceNameWhenRegisterUnit(t *testing.T) {
	assert := assert.New(t)

	// given:
	cache := registry.NewCacheRepository()
	SUT := registry.NewService(registry.WithRepository(cache))
	instances := []registry.ServiceInstance{
		registry.NewServiceInstance("component", "", "12345", "9090"),
	}

	// when:
	err := SUT.Register(instances...)

	// then:
	assert.ErrorIs(err, registry.ErrMissingData)
}

func TestShouldReturnErrorForGreaterInstancePortThan65536WhenRegisterUnit(t *testing.T) {
	assert := assert.New(t)

	// given:
	cache := registry.NewCacheRepository()
	service := registry.NewService(registry.WithRepository(cache))
	instances := []registry.ServiceInstance{
		registry.NewServiceInstance("component", "node1", "localhost", "65536"),
	}

	// when:
	err := service.Register(instances...)

	// then:
	assert.ErrorIs(err, registry.ErrMissingData)
}

func TestShouldReturnErrorForInstancePortEqualToZeroWhenRegisterUnit(t *testing.T) {
	assert := assert.New(t)

	// given:
	cache := registry.NewCacheRepository()
	service := registry.NewService(registry.WithRepository(cache))
	instances := []registry.ServiceInstance{
		registry.NewServiceInstance("component", "node1", "localhost", "0"),
	}

	// when:
	err := service.Register(instances...)

	// then:
	assert.ErrorIs(err, registry.ErrMissingData)
}

func TestShouldMarkServiceInstanceAsHealthyUnit(t *testing.T) {
	assert := assert.New(t)

	// given:
	cache := registry.NewCacheRepository()
	httpCli := &mocks.HTTPClient{}
	component := "dummy"
	service := registry.NewService(
		registry.WithRepository(cache),
		registry.WithHTTPClient(httpCli),
		registry.WithHealthz(registry.ServiceHealthEndpoints{
			component: "/v1/test",
		}),
	)
	instance := registry.NewServiceInstance(component, "node1", "localhost", "8080")

	expectedInstance := instance
	expectedInstance.SetHealth(true)
	expectedInstances := []registry.ServiceInstance{expectedInstance}

	service.Register(instance)

	instance.SetHealth(true)
	httpCli.EXPECT().Get(mock.Anything).Return(&http.Response{StatusCode: http.StatusOK}, nil)

	// when:
	service.HeartBeat()

	// then:
	actualInstances, err := service.QueryInstances(component)
	assert.Nil(err)

	assert.Equal(expectedInstances, actualInstances)
	httpCli.AssertExpectations(t)

}

func TestShouldMarkSingleServiceInstanceAsNotHealthyUnit(t *testing.T) {
	assert := assert.New(t)

	// given:
	cache := registry.NewCacheRepository()
	httpCli := &mocks.HTTPClient{}
	component := "dummy"
	service := registry.NewService(
		registry.WithRepository(cache),
		registry.WithHTTPClient(httpCli),
		registry.WithHealthz(registry.ServiceHealthEndpoints{
			component: "/v1/test",
		}),
	)

	first := registry.NewServiceInstance(component, "node1", "localhost", "8080")
	second := registry.NewServiceInstance(component, "node2", "localhost", "8090")
	second.SetHealth(true)
	first.SetHealth(true)

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
	actualInstances, err := service.QueryInstances(component)
	assert.Nil(err)

	assert.Equal(expectedInstances, actualInstances)
	httpCli.AssertExpectations(t)
}

func TestShouldLogErrorWhenServiceInstanceUpdateStatusFailure(t *testing.T) {
	// given:
	registryRepository := &mocks.ServiceRegistry{}
	httpCli := &mocks.HTTPClient{}
	logger := &mocks.Logger{}
	component := "dummy"
	service := registry.NewService(
		registry.WithRepository(registryRepository),
		registry.WithHTTPClient(httpCli),
		registry.WithLogger(logger),
		registry.WithHealthz(registry.ServiceHealthEndpoints{
			component: "/v1/test",
		}),
	)
	instance := registry.NewServiceInstance(component, "node1", "localhost", "8080")
	registryRepository.EXPECT().QueryInstances(mock.Anything).Return([]registry.ServiceInstance{instance}, nil)
	registryRepository.EXPECT().UpdateStatus(mock.Anything).Return(errors.New("repository is down"))
	httpCli.EXPECT().Get(mock.Anything).Return(&http.Response{StatusCode: http.StatusOK}, nil)
	logger.EXPECT().Errorf(mock.Anything, mock.Anything).Once()

	// when:
	service.HeartBeat()

	// then:
	mock.AssertExpectationsForObjects(t, registryRepository, logger, httpCli)
}
