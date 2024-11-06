package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

var (
	baseURL           = "https://api.airtable.com/v0"
	listRecordsPath   = "/%s/%s" // baseID, tableID/tableName
	listRecordsFields = "createdBy,count"
)

type (
	AirtableClient struct {
		APIToken string
	}

	Record struct{}

	OptionalParams struct {
		PageSize int
		Offset   string
	}

	OptionalResponse struct {
		Offset string
	}

	ListRecordsResponse struct {
		Records []Record `json:"records"`
		Offset  string   `json:"offset"`
	}
)

func NewAirtableClient(token string) *AirtableClient {
	return &AirtableClient{
		APIToken: token,
	}
}

func (c *AirtableClient) ListRecords(baseID, tableID string, reqParams OptionalParams) ([]Record, OptionalResponse, error) {
	request, err := http.NewRequest("GET", fmt.Sprintf("%s%s", baseURL, listRecordsPath, baseID, tableID), nil)
	if err != nil {
		return nil, OptionalResponse{}, err
	}

	query := request.URL.Query()
	if reqParams.PageSize > 0 {
		query.Add("pageSize", fmt.Sprintf("%d", reqParams.PageSize))
	}
	if reqParams.Offset != "" {
		query.Add("offset", reqParams.Offset)
	}
	query.Add("fields", listRecordsFields)
	request.URL.RawQuery = query.Encode()
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.APIToken))

	httpClient := &http.Client{}
	httpResp, err := httpClient.Do(request)
	if err != nil {
		return nil, OptionalResponse{}, err
	}

	defer httpResp.Body.Close()

	body, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, OptionalResponse{}, err
	}

	var response ListRecordsResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, OptionalResponse{}, err
	}

	return response.Records, OptionalResponse{
		Offset: response.Offset,
	}, nil
}
