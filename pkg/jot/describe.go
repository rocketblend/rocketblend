package jot

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/rocketblend/rocketblend/pkg/jot/reference"
)

func (d *Driver) Describe(reference reference.Reference, data *interface{}) error {
	url := reference.Url()

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(bodyBytes, &data); err != nil {
		return fmt.Errorf("failed to parse %v", err)
	}

	return nil
}
