package registry_test

import (
	"testing"

	"github.com/michalgosek/workout-app-infrastrcutre/service-discovery/internal/registry"

	"github.com/stretchr/testify/assert"
)

func TestShouldReturnAllRegisterdClustesUnit(t *testing.T) {
	assert := assert.New(t)

	// given:
	SUT := registry.NewCacheRepository()
	component := "service1"
	instance1 := registry.ServiceInstance{
		Component: component,
		Name:      "node1",
		IP:        "localhost",
		Port:      "8080",
		Endpoint:  "/v1/health",
	}
	instance2 := registry.ServiceInstance{
		Component: component,
		Name:      "node2",
		IP:        "localhost",
		Port:      "8090",
		Endpoint:  "/v1/health",
	}
	expectedInstances := []registry.ServiceInstance{
		instance1,
		instance2,
	}
	expectedClusters := []registry.ServiceCluster{
		{
			Name: component,
			Instances: map[string]registry.ServiceInstance{
				instance1.Name: instance1,
				instance2.Name: instance2,
			},
		},
	}

	SUT.Register(expectedInstances...)

	// when:
	actualClusters := SUT.ListClusters()

	// then:
	assert.EqualValues(expectedClusters, actualClusters)
}

func TestShouldNotRegisterDuplicateNodeInstancesUnit(t *testing.T) {
	assert := assert.New(t)

	// given:
	SUT := registry.NewCacheRepository()
	component := "service1"
	instance := registry.ServiceInstance{
		Component: component,
		Name:      "node1",
		IP:        "localhost",
		Port:      "8080",
		Endpoint:  "/v1/health",
	}

	duplicateInstances := []registry.ServiceInstance{instance, instance}
	expectedInstance := instance
	expectedInstances := []registry.ServiceInstance{expectedInstance}

	// when:
	err := SUT.Register(duplicateInstances...)

	// then:
	assert.Nil(err)

	actualInstances, err := SUT.QueryInstances(component)
	assert.Nil(err)
	assert.Equal(expectedInstances, actualInstances)
}

func TestShouldNotUpdateNonExistingServiceInstanceUnit(t *testing.T) {
	assert := assert.New(t)

	// given:
	SUT := registry.NewCacheRepository()
	component := "service1"
	instance := registry.ServiceInstance{
		Component: component,
		Name:      "node1",
		IP:        "localhost",
		Port:      "8080",
		Endpoint:  "/v1/health",
	}

	// when:
	err := SUT.UpdateStatus(instance)

	// then:
	assert.Nil(err)

	actualInstances, err := SUT.QueryInstances(component)
	assert.Nil(err)
	assert.Empty(actualInstances)
}

func TestShouldUpdateServiceInstanceStatusToTrueUnit(t *testing.T) {
	assert := assert.New(t)

	// given:
	SUT := registry.NewCacheRepository()
	component := "service1"
	instance := registry.ServiceInstance{
		Component: component,
		Name:      "node1",
		IP:        "localhost",
		Port:      "8080",
		Endpoint:  "/v1/health",
	}
	instances := []registry.ServiceInstance{instance}
	expectedInstance := instance
	expectedInstance.SetHealth(true)
	expectedInstances := []registry.ServiceInstance{expectedInstance}

	SUT.Register(instances...)
	instance.SetHealth(true)

	// when:
	err := SUT.UpdateStatus(instance)

	// then:
	assert.Nil(err)

	actualInstances, err := SUT.QueryInstances(component)
	assert.Nil(err)
	assert.Equal(expectedInstances, actualInstances)
}

func TestShouldUpdateServiceSingleInstanceStatusToTrueUnit(t *testing.T) {
	assert := assert.New(t)

	// given:
	SUT := registry.NewCacheRepository()
	component := "dummy"
	first := registry.ServiceInstance{
		Component: component,
		Name:      "node1",
		IP:        "localhost",
		Port:      "8080",
		Endpoint:  "/v1/health",
	}
	second := registry.ServiceInstance{
		Component: component,
		Name:      "node2",
		IP:        "localhost",
		Port:      "8090",
		Endpoint:  "/v1/health",
	}
	third := registry.ServiceInstance{
		Component: component,
		Name:      "node3",
		IP:        "localhost",
		Port:      "8070",
		Endpoint:  "/v1/health",
	}
	instances := []registry.ServiceInstance{first, second, third}

	SUT.Register(instances...)
	second.SetHealth(true)

	expectedInstances := []registry.ServiceInstance{first, second, third}

	// when:
	err := SUT.UpdateStatus(second)

	// then:
	assert.Nil(err)

	actualInstances, err := SUT.QueryInstances(component)
	assert.Nil(err)
	assert.Equal(expectedInstances, actualInstances)
}

func TestShouldReturnEmptyServiceInstancesAfterQueryForNonExistingServiceUnit(t *testing.T) {
	assert := assert.New(t)

	// given:
	SUT := registry.NewCacheRepository()

	// when:
	actualInstances, err := SUT.QueryInstances("dummy")

	// then:
	assert.Nil(err)
	assert.Empty(actualInstances)
}

func TestShouldRegisterOneServiceInstanceUnit(t *testing.T) {
	assert := assert.New(t)

	// given:
	component := "service1"
	expectedInstance := registry.ServiceInstance{
		Component: component,
		Name:      "node1",
		IP:        "localhost",
		Port:      "8080",
		Endpoint:  "/v1/health",
	}
	expectedInstances := []registry.ServiceInstance{expectedInstance}

	SUT := registry.NewCacheRepository()

	// when:
	err := SUT.Register(expectedInstance)

	// then:
	assert.Nil(err)

	actualInstances, err := SUT.QueryInstances(component)
	assert.Nil(err)
	assert.Equal(expectedInstances, actualInstances)
}

func TestShouldRegisterServiceIsntancesUnit(t *testing.T) {
	assert := assert.New(t)

	// given:
	component := "dummy"
	first := registry.ServiceInstance{
		Component: component,
		Name:      "node1",
		IP:        "localhost",
		Port:      "8080",
		Endpoint:  "/v1/health",
	}
	second := registry.ServiceInstance{
		Component: component,
		Name:      "node2",
		IP:        "localhost",
		Port:      "8090",
		Endpoint:  "/v1/health",
	}
	expectedInstances := []registry.ServiceInstance{first, second}
	SUT := registry.NewCacheRepository()

	// when:
	err := SUT.Register(first, second)

	// then:
	assert.Nil(err)

	actualInstances, err := SUT.QueryInstances(component)
	assert.Nil(err)
	assert.Equal(expectedInstances, actualInstances)
}
