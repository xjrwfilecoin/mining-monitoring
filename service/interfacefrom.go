package service

import (
	"fmt"
)

type MinerInfoForm struct {
	MinerId string `json:"minerId"`
}

func (m *MinerInfoForm) Valid() error {
	if len(m.MinerId) == 0 {
		return fmt.Errorf("minerId is empty \n")
	}
	return nil
}
