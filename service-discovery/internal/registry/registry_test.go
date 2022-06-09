package registry_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/michalgosek/workout-app-infrastrcutre/service-discovery/internal/registry"
	"github.com/michalgosek/workout-app-infrastrcutre/service-discovery/internal/registry/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRegisterShouldReturnErrorWhenRepoistoryFailureUnit(t *testing.T) {
	assert := assert.New(t)

	// given:
	repo := &mocks.Repository{}
	SUT := registry.NewService(registry.WithRepository(repo))
	component := "service1"
	exepctedInstances := []registry.ServiceInstance{
		{
			Component: component,
			Name:      "node1",
			IP:        "localhost",
			Port:      "8080",
			Endpoint:  "/api/v1/health",
		},
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
		{
			Component: component,
			Name:      "node1",
			IP:        "localhost",
			Port:      "8080",
			Endpoint:  "/api/v1/health",
		},
	}

	// when:
	err := SUT.Register(exepctedInstances...)

	// then:
	assert.Nil(err)

	actualInstances, err := SUT.QueryInstances(component)
	assert.Nil(err)
	assert.Equal(exepctedInstances, actualInstances)
}

func TestRegisterShouldNotReturnErrorForEmptyInstanceUnit(t *testing.T) {
	assert := assert.New(t)

	// given:
	cache := registry.NewCacheRepository()
	SUT := registry.NewService(registry.WithRepository(cache))
	var instances []registry.ServiceInstance

	// when:
	err := SUT.Register(instances...)

	// then:
	assert.Nil(err)
}

func TestRegisterShouldReturnErrorWhenHealthEndpointAddrIsMalformedUnit(t *testing.T) {
	assert := assert.New(t)

	// given:
	cache := registry.NewCacheRepository()
	SUT := registry.NewService(registry.WithRepository(cache))
	instances := []registry.ServiceInstance{
		{
			Component: "service1",
			Name:      "node1",
			IP:        "localhost",
			Port:      "8080",
			Endpoint:  "123123123212",
		},
	}

	// when:
	err := SUT.Register(instances...)

	// then:
	assert.ErrorIs(err, registry.ErrMalformedData)
}

func TestRegisterShouldReturnErrorWhenIPAddrIsMalformedUnit(t *testing.T) {
	assert := assert.New(t)

	// given:
	cache := registry.NewCacheRepository()
	SUT := registry.NewService(registry.WithRepository(cache))
	instances := []registry.ServiceInstance{
		{
			Component: "service1",
			Name:      "node1",
			IP:        "122345",
			Port:      "8080",
			Endpoint:  "/api/v1/health",
		},
	}

	// when:
	err := SUT.Register(instances...)

	// then:
	assert.ErrorIs(err, registry.ErrInvalidDataFormat)
}

func TestRegisterShouldReturnErrorWhenComponentNameIsEmptyUnit(t *testing.T) {
	assert := assert.New(t)

	// given:
	cache := registry.NewCacheRepository()
	SUT := registry.NewService(registry.WithRepository(cache))
	instances := []registry.ServiceInstance{
		{
			Component: "",
			Name:      "node1",
			IP:        "localhost",
			Port:      "8080",
			Endpoint:  "/api/v1/health",
		},
	}

	// when:
	err := SUT.Register(instances...)

	// then:
	assert.ErrorIs(err, registry.ErrMissingData)
}

func TestRegisterShouldReturnErrorWhenInstanceNameIsEmptyUnit(t *testing.T) {
	assert := assert.New(t)

	// given:
	cache := registry.NewCacheRepository()
	SUT := registry.NewService(registry.WithRepository(cache))
	instances := []registry.ServiceInstance{
		{
			Component: "service1",
			Name:      "",
			IP:        "localhost",
			Port:      "8080",
			Endpoint:  "/api/v1/health",
		},
	}

	// when:
	err := SUT.Register(instances...)

	// then:
	assert.ErrorIs(err, registry.ErrMissingData)
}

func TestRegisterShouldReturnErrorForGreaterInstancePortThan65536Unit(t *testing.T) {
	assert := assert.New(t)

	// given:
	cache := registry.NewCacheRepository()
	service := registry.NewService(registry.WithRepository(cache))
	instances := []registry.ServiceInstance{
		{
			Component: "component",
			Name:      "node1",
			IP:        "localhost",
			Port:      "65537",
			Endpoint:  "/api/v1/health",
		},
	}

	// when:
	err := service.Register(instances...)

	// then:
	assert.ErrorIs(err, registry.ErrMissingData)
}

func TestRegisterShouldReturnErrorForInstancePortEqualToZeroUnit(t *testing.T) {
	assert := assert.New(t)

	// given:
	cache := registry.NewCacheRepository()
	service := registry.NewService(registry.WithRepository(cache))
	instances := []registry.ServiceInstance{
		{
			Component: "component",
			Name:      "node1",
			IP:        "localhost",
			Port:      "0",
			Endpoint:  "/api/v1/health",
		},
	}

	// when:
	err := service.Register(instances...)

	// then:
	assert.ErrorIs(err, registry.ErrMissingData)
}

func TestProcessClusterShouldMarkServiceInstanceAsHealthyUnit(t *testing.T) {
	assert := assert.New(t)

	// given:
	cache := registry.NewCacheRepository()
	httpCli := &mocks.HTTPClient{}
	component := "component"
	service := registry.NewService(
		registry.WithRepository(cache),
		registry.WithHTTPClient(httpCli),
		registry.WithHealthz(registry.ServiceHealthEndpoints{
			component: "/v1/test",
		}),
	)
	instance := registry.ServiceInstance{
		Component: component,
		Name:      "node1",
		IP:        "localhost",
		Port:      "8080",
		Endpoint:  "/v1/health",
	}

	expectedInstance := instance
	expectedInstance.SetHealth(true)
	expectedInstances := []registry.ServiceInstance{expectedInstance}

	clusters := []registry.ServiceCluster{
		{
			Name: instance.Component,
			Instances: map[string]registry.ServiceInstance{
				instance.Name: instance,
			},
		},
	}

	service.Register(instance)
	instance.SetHealth(true)

	httpCli.EXPECT().Get(mock.Anything).Return(&http.Response{StatusCode: http.StatusOK}, nil)

	// when:
	service.ProcessClusters(clusters...)

	// then:
	actualInstances, err := service.QueryInstances(component)
	assert.Nil(err)

	assert.Equal(expectedInstances, actualInstances)
	httpCli.AssertExpectations(t)
}

func TestProcessClusterShouldMarkServiceInstanceAsUnhealthyUnit(t *testing.T) {
	assert := assert.New(t)

	// given:
	cache := registry.NewCacheRepository()
	httpCli := &mocks.HTTPClient{}
	component := "component"
	service := registry.NewService(
		registry.WithRepository(cache),
		registry.WithHTTPClient(httpCli),
		registry.WithHealthz(registry.ServiceHealthEndpoints{
			component: "/v1/test",
		}),
	)
	instance := registry.ServiceInstance{
		Component: component,
		Name:      "node1",
		IP:        "localhost",
		Port:      "8080",
		Endpoint:  "/v1/health",
	}
	instance.SetHealth(true)

	expectedInstance := instance
	expectedInstance.SetHealth(false)
	expectedInstances := []registry.ServiceInstance{expectedInstance}

	clusters := []registry.ServiceCluster{
		{
			Name: instance.Component,
			Instances: map[string]registry.ServiceInstance{
				instance.Name: instance,
			},
		},
	}

	service.Register(instance)

	httpCli.EXPECT().Get(mock.Anything).Return(&http.Response{StatusCode: http.StatusServiceUnavailable}, nil)

	// when:
	service.ProcessClusters(clusters...)

	// then:
	actualInstances, err := service.QueryInstances(component)
	assert.Nil(err)

	assert.Equal(expectedInstances, actualInstances)
	httpCli.AssertExpectations(t)
}
