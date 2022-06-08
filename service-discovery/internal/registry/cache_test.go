package registry_test

import (
	"service-discovery/internal/registry"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShouldUpdateServiceInstanceStatusToTrueUnit(t *testing.T) {
	assert := assert.New(t)

	// given:
	SUT := registry.NewCacheRepository()
	name := "dummy"
	instance := registry.NewServiceInstance(name, "localhost", 8080)
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

	actualInstances, err := SUT.QueryInstances(name)
	assert.Nil(err)
	assert.Equal(expectedInstances, actualInstances)
}

func TestShouldUpdateServiceSingleInstanceStatusToTrueUnit(t *testing.T) {
	assert := assert.New(t)

	// given:
	SUT := registry.NewCacheRepository()
	name := "dummy"
	first := registry.NewServiceInstance(name, "localhost", 8090)
	second := registry.NewServiceInstance(name, "localhost", 8080)
	third := registry.NewServiceInstance(name, "localhost", 8030)
	instances := []registry.ServiceInstance{first, second, third}

	SUT.Register(instances...)
	second.SetHealth(true)

	expectedInstances := []registry.ServiceInstance{first, second, third}

	// when:
	err := SUT.UpdateStatus(second)

	// then:
	assert.Nil(err)

	actualInstances, err := SUT.QueryInstances(name)
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
	serviceName := "dummy-v1"
	expectedInstance := registry.NewServiceInstance(serviceName, "localhost", 8080)
	expectedInstances := []registry.ServiceInstance{expectedInstance}

	SUT := registry.NewCacheRepository()

	// when:
	err := SUT.Register(expectedInstance)

	// then:
	assert.Nil(err)

	actualInstances, err := SUT.QueryInstances(serviceName)
	assert.Nil(err)
	assert.Equal(expectedInstances, actualInstances)
}

func TestShouldRegisterServiceIsntancesUnit(t *testing.T) {
	assert := assert.New(t)

	// given:
	serviceName := "dummy"
	first := registry.NewServiceInstance(serviceName, "localhost", 8080)
	second := registry.NewServiceInstance(serviceName, "localhost", 8090)
	expectedInstances := []registry.ServiceInstance{first, second}
	SUT := registry.NewCacheRepository()

	// when:
	err := SUT.Register(first, second)

	// then:
	assert.Nil(err)

	actualInstances, err := SUT.QueryInstances(serviceName)
	assert.Nil(err)
	assert.Equal(expectedInstances, actualInstances)
}
