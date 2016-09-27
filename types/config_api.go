package types

import (
	"encoding/json"
	"github.com/pkg/errors"
	"strconv"
)

var (
	ConfigServiceEndpoint = "/peer"
	FSMAPIEndpoint        = "/fsm"
)

type ConfigDeletionRequest struct {
	ID uint64 `json:"id"`
}

// Host address should be able to be scraped from the Request on the server-end
// Therefore the request shouldn't be made public unless this changes
type ConfigAdditionRequest struct {
	ID       uint64 `json:"id"`
	RaftPort int    `json:"raft_port"`
	APIPort  int    `json:"api_port"`
	// Host is only for external requests for addition when doing strange things
	Host string `json:"host,omitempty"`
}

type Peer struct {
	IP       string `json:"ip"`
	RaftPort int    `json:"raft_port"`
	APIPort  int    `json:"api_port"`
}

type ConfigServiceResponse struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
	Data    []byte `json:"data,omitempty"`
}

// PeerAdditionAddMe has self-identifying port and id
// With a list of all Peers in the cluster currently
type ConfigAdditionResponseData struct {
	ConfigPeerData
}

type ConfigMembershipResponseData struct {
	ConfigPeerData
}

// This needs to be a different struct because it is important to seperate
// The API/Raft/ID of the node we're pinging from other remote nodes
type ConfigPeerData struct {
	RaftPort    int             `json:"raft_port"`
	APIPort     int             `json:"api_port"`
	ID          uint64          `json:"id"`
	RemotePeers map[uint64]Peer `json:"peers"`
}

func (p *ConfigPeerData) MarshalJSON() ([]byte, error) {
	tmpStruct := &struct {
		RaftPort    int             `json:"raft_port"`
		APIPort     int             `json:"api_port"`
		ID          uint64          `json:"id"`
		RemotePeers map[string]Peer `json:"peers"`
	}{
		RaftPort:    p.RaftPort,
		APIPort:     p.APIPort,
		ID:          p.ID,
		RemotePeers: make(map[string]Peer),
	}

	for key, val := range p.RemotePeers {
		tmpStruct.RemotePeers[strconv.FormatUint(key, 10)] = val
	}

	retJSON, err := json.Marshal(tmpStruct)

	return retJSON, errors.Wrap(err, "Error marshalling JSON for http peer data")
}

func (p *ConfigPeerData) UnmarshalJSON(data []byte) error {
	tmpStruct := &struct {
		RaftPort    int             `json:"raft_port"`
		APIPort     int             `json:"api_port"`
		ID          uint64          `json:"id"`
		RemotePeers map[string]Peer `json:"peers"`
	}{}

	if err := json.Unmarshal(data, tmpStruct); err != nil {
		return errors.Wrap(err, "Error unmarshalling http peer data")
	}

	p.APIPort = tmpStruct.APIPort
	p.RaftPort = tmpStruct.RaftPort
	p.ID = tmpStruct.ID
	p.RemotePeers = make(map[uint64]Peer)

	for key, val := range tmpStruct.RemotePeers {
		convKey, err := strconv.ParseUint(key, 10, 64)
		if err != nil {
			return errors.Wrap(err, "Error parsing peer id from map")
		}
		p.RemotePeers[convKey] = val
	}

	return nil
}
