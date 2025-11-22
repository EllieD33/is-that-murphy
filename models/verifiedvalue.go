package models

import "errors"

type VerifiedValue struct {
	Value string `json:"value"`
	Type string `json:"type"`
}

func (v *VerifiedValue) Validate() error {
    if v.Value == "" {
        return errors.New("value is required")
    }
    if v.Type == "" {
        return errors.New("type is required")
    }

    return nil
}