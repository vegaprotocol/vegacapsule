package state

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
)

func encodeState(state NetworkState) ([]byte, error) {
	stateBytes, err := json.Marshal(state)
	if err != nil {
		return nil, fmt.Errorf("cannot convert network state structure into string: %w", err)
	}

	networkState := make([]byte, hex.EncodedLen(len(stateBytes)))
	hex.Encode(networkState, stateBytes)

	return networkState, nil
}

func decodeState(data []byte) (*NetworkState, error) {
	if data == nil {
		return &NetworkState{}, nil
	}

	stateBytes := make([]byte, hex.DecodedLen(len(data)))
	if _, err := hex.Decode(stateBytes, data); err != nil {
		return nil, fmt.Errorf("cannot decode network state: %w", err)
	}

	networkState := &NetworkState{}
	if err := json.Unmarshal(stateBytes, networkState); err != nil {
		return nil, fmt.Errorf("cannot decode network state from given data: %w", err)
	}

	return networkState, nil
}
