package build

import (
	"encoding/json"
	"io"
	"net/http"
)

type HttpClient struct {
}

func (r *HttpClient) Fetch(remote string, platform string, tag string) ([]Build, error) {
	response, err := http.Get(remote)
	if err != nil {
		return nil, err
	}

	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var responseObject Response
	json.Unmarshal(responseData, &responseObject)

	return responseObject.Data, nil
}
