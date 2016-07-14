package main

type StubEtcdRegistryClient struct {
	instances []string
}

func (client *StubEtcdRegistryClient) Register() error {
	return nil
}

func (client *StubEtcdRegistryClient) Unregister() error {
	return nil
}

func (client *StubEtcdRegistryClient) ServicesByName(name string) ([]string, error) {
	return client.instances, nil
}