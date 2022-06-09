package registry_test

import (
	"service-discovery/internal/registry"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShouldNotRegisterDuplicateNodeInstancesUnit(t *testing.T) {
	assert := assert.New(t)

	// given:
	SUT := registry.NewCacheRepository()
	component := "service1"
	instance := registry.NewServiceInstance(component, "dummy", "localhost", "8080")
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
	instance := registry.NewServiceInstance(component, "dummy", "localhost", "8080")

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
	instance := registry.NewServiceInstance(component, "dummy", "localhost", "8080")
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
	first := registry.NewServiceInstance(component, "node1", "localhost", "8090")
	second := registry.NewServiceInstance(component, "node2", "localhost", "8080")
	third := registry.NewServiceInstance(component, "node3", "localhost", "8030")
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
	expectedInstance := registry.NewServiceInstance(component, "dummy", "localhost", "8080")
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
	first := registry.NewServiceInstance(component, "node1", "localhost", "8080")
	second := registry.NewServiceInstance(component, "node2", "localhost", "8090")
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
