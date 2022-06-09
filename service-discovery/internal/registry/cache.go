package registry

import "sort"

func NewServiceInstance(Component, Name, IP string, Port string) ServiceInstance {
	s := ServiceInstance{
		Component: Component,
		Name:      Name,
		IP:        IP,
		Port:      Port,
		healthy:   false,
	}
	return s
}

type ServiceClusters map[string]*ServiceCluster

type CacheRepository struct {
	clusters ServiceClusters
}

func (c *CacheRepository) UpdateStatus(s ServiceInstance) error {
	cluster, ok := c.clusters[s.Component]
	if !ok {
		return nil
	}
	_, ok = cluster.Instances[s.Name]
	if !ok {
		return nil
	}
	c.clusters[s.Component].Instances[s.Name] = s
	return nil
}

func (c *CacheRepository) Register(ss ...ServiceInstance) error {
	for _, s := range ss {
		v, ok := c.clusters[s.Component]
		if !ok {
			c.clusters[s.Component] = &ServiceCluster{
				Name: s.Component,
				Instances: map[string]ServiceInstance{
					s.Name: s,
				},
			}
			continue
		}
		v.Instances[s.Name] = s
	}
	return nil
}

func (c *CacheRepository) QueryInstances(componentName string) ([]ServiceInstance, error) {
	cluster, ok := c.clusters[componentName]
	if !ok {
		return []ServiceInstance{}, nil
	}

	var keys []string
	for k := range cluster.Instances {
		keys = append(keys, k)
	}
	sort.SliceStable(keys, func(i, j int) bool { return keys[i] < keys[j] })

	var instances []ServiceInstance
	for _, k := range keys {
		instances = append(instances, cluster.Instances[k])
	}
	return instances, nil
}

func NewCacheRepository() *CacheRepository {
	c := CacheRepository{
		clusters: make(ServiceClusters),
	}
	return &c
}
