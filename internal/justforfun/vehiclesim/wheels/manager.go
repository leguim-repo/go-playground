package wheels

import "fmt"

// WheelManager manages vehicle wheels
type WheelManager struct {
	WheelPair *WheelPair
}

func NewWheelManager(tireSpec string) (*WheelManager, error) {
	wheelPair, err := NewWheelPair(tireSpec)
	if err != nil {
		return nil, fmt.Errorf("error creating wheels: %v", err)
	}

	return &WheelManager{
		WheelPair: wheelPair,
	}, nil
}
