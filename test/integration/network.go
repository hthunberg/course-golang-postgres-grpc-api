package integration

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	tc "github.com/testcontainers/testcontainers-go"
)

type Network struct {
	tc.Network
	Name string
}

// CreateNetwork creates a network with a random name
func CreateNetwork(ctx context.Context) (*Network, error) {
	networkName := fmt.Sprintf("net_%s", uuid.New())

	n, err := tc.GenericNetwork(ctx, tc.GenericNetworkRequest{
		NetworkRequest: tc.NetworkRequest{
			Name:           networkName,
			CheckDuplicate: true,
		},
	})
	if err != nil {
		return nil, err
	}

	return &Network{n, networkName}, nil
}

// TearDown tears down the network
func (n *Network) TearDown(ctx context.Context) {
	_ = n.Network.Remove(ctx)
}

// ApplyNetworkAlias applies a network alias to a generic container request
func (n *Network) ApplyNetworkAlias(req *tc.GenericContainerRequest, alias string) {
	if req.Networks == nil {
		req.Networks = make([]string, 0)
	}
	req.Networks = append(req.Networks, n.Name)

	if req.NetworkAliases == nil {
		req.NetworkAliases = make(map[string][]string)
	}
	if _, ok := req.NetworkAliases[n.Name]; !ok {
		req.NetworkAliases[n.Name] = make([]string, 0)
	}
	req.NetworkAliases[n.Name] = append(req.NetworkAliases[n.Name], alias)
}
