package client

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/util/wait"
)

// MockIdempotentKeyGen is a mock implementation of IdempotentKeyGen interface
type MockIdempotentKeyGen struct {
	generatedKeys map[string]string
}

func (m *MockIdempotentKeyGen) GenerateKey(argsHash string) string {
	if _, ok := m.generatedKeys[argsHash]; !ok {
		m.generatedKeys[argsHash] = "mockToken"
	}
	return m.generatedKeys[argsHash]
}

func (m *MockIdempotentKeyGen) PutBack(argsHash string, clientToken string) {
	delete(m.generatedKeys, argsHash)
}

// TestCreateNetworkInterfaceOptions_Finish tests the Finish function of CreateNetworkInterfaceOptions
func TestCreateNetworkInterfaceOptions_Finish(t *testing.T) {
	// Prepare the test data
	niOptions := &NetworkInterfaceOptions{
		VSwitchID:        "vsw-xxxxxx",
		SecurityGroupIDs: []string{"sg-xxxxxx"},
		ResourceGroupID:  "rg-xxxxxx",
		Tags:             map[string]string{"key1": "value1", "key2": "value2"},
		Trunk:            true,
		ERDMA:            true,
		IPCount:          2,
		IPv6Count:        1,
	}

	c := &CreateNetworkInterfaceOptions{
		NetworkInterfaceOptions: niOptions,
	}

	// Execute the function to be tested
	req, cleanup, err := c.Finish(&MockIdempotentKeyGen{generatedKeys: map[string]string{}})

	// Verify the result
	assert.NoError(t, err)
	assert.NotNil(t, req)
	assert.NotNil(t, cleanup)

	assert.Equal(t, niOptions.VSwitchID, req.VSwitchId)
	assert.Equal(t, ENITypeTrunk, req.InstanceType)
	assert.Equal(t, ENITrafficModeRDMA, req.NetworkInterfaceTrafficMode)
	assert.Equal(t, 1, len(*req.SecurityGroupIds))
	assert.Equal(t, niOptions.ResourceGroupID, req.ResourceGroupId)
	assert.Equal(t, eniDescription, req.Description)
	assert.Equal(t, "mockToken", req.ClientToken)
	assert.NotNil(t, c.Backoff)

	// Cleanup
	cleanup()
}

func TestCreateNetworkInterfaceOptions_ApplyCreateNetworkInterface(t *testing.T) {
	type fields struct {
		NetworkInterfaceOptions *NetworkInterfaceOptions
		Backoff                 *wait.Backoff
	}
	type args struct {
		options *CreateNetworkInterfaceOptions
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *CreateNetworkInterfaceOptions
	}{
		{
			name: "TestApplyCreateNetworkInterface",
			fields: fields{
				NetworkInterfaceOptions: &NetworkInterfaceOptions{
					VSwitchID:        "vsw-xxxxxx",
					SecurityGroupIDs: []string{"sg-xxxxxx"},
					ResourceGroupID:  "rg-xxxxxx",
					Tags:             map[string]string{"key1": "value1", "key2": "value2"},
					Trunk:            true,
					ERDMA:            true,
					IPCount:          2,
				},
			},
			args: args{
				options: &CreateNetworkInterfaceOptions{
					NetworkInterfaceOptions: &NetworkInterfaceOptions{},
				},
			},
			want: &CreateNetworkInterfaceOptions{
				NetworkInterfaceOptions: &NetworkInterfaceOptions{
					Trunk:                 true,
					ERDMA:                 true,
					VSwitchID:             "vsw-xxxxxx",
					SecurityGroupIDs:      []string{"sg-xxxxxx"},
					ResourceGroupID:       "rg-xxxxxx",
					IPCount:               2,
					IPv6Count:             0,
					Tags:                  map[string]string{"key1": "value1", "key2": "value2"},
					InstanceID:            "",
					InstanceType:          "",
					Status:                "",
					NetworkInterfaceID:    "",
					DeleteENIOnECSRelease: nil,
				},
				Backoff: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CreateNetworkInterfaceOptions{
				NetworkInterfaceOptions: tt.fields.NetworkInterfaceOptions,
				Backoff:                 tt.fields.Backoff,
			}
			c.ApplyCreateNetworkInterface(tt.args.options)
		})
	}
}
