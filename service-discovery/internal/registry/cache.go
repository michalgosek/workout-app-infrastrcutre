package registry

func (s *ServiceInstance) SetHealth(v bool) {
	s.healthy = v
}

func NewServiceInstance(name, ip string, port uint) ServiceInstance {
	s := ServiceInstance{
		name:    name,
		ip:      ip,
		port:    port,
		healthy: false,
	}
	return s
}

type ServiceClusters map[string]*ServiceCluster

type CacheRepository struct {
	clusters ServiceClusters
}

func (c *CacheRepository) UpdateStatus(s ServiceInstance) error {
	cluster, ok := c.clusters[s.name]
	if !ok {
		return nil
	}
	instances := cluster.Instances
	for i, ins := range instances {
		if ins.ip == s.ip && ins.port == s.port {
			instances[i] = s
		}
	}
	c.clusters[s.name].Instances = instances
	return nil
}

func (c *CacheRepository) Register(ss ...ServiceInstance) error {
	for _, s := range ss {
		v, ok := c.clusters[s.name]
		if !ok {
			c.clusters[s.name] = &ServiceCluster{Name: s.name, Instances: []ServiceInstance{s}}
			continue
		}
		v.Instances = append(v.Instances, s)
	}
	return nil
}

func (c *CacheRepository) QueryInstances(name string) ([]ServiceInstance, error) {
	v, ok := c.clusters[name]
	if !ok {
		return []ServiceInstance{}, nil
	}
	return v.Instances, nil
}

func NewCacheRepository() *CacheRepository {
	c := CacheRepository{
		clusters: make(ServiceClusters),
	}
	return &c
}
