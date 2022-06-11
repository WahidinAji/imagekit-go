package media

import (
	"context"
	"encoding/json"

	"github.com/dhaval070/imagekit-go/api"
	"gopkg.in/validator.v2"
)

// CreateFolderParam represents parameter to create folder api
type CreateFolderParam struct {
	FolderName       string `validate:"nonzero" json:"folderName"`
	ParentFolderPath string `validate:"nonzero" json:"parentFolderPath"`
}

// DeleteFolderParam represents parameter to delete folder api
type DeleteFolderParam struct {
	FolderPath string `validate:"nonzero" json:"folderPath"`
}

// CopyFolderParam represents parameter to copy folder api
type CopyFolderParam struct {
	SourceFolderPath string `validate:"nonzero" json:"sourceFolderPath"`
	destinationPath  string `validate:"nonzero" json:"destinationPath"`
	IncludeVersions  bool   `json:"includeVersions"`
}

// MoveFolderParam represents parameter to move folder api
type MoveFolderParam struct {
	SourceFolderPath string `validate:"nonzero" json:"sourceFolderPath"`
	DestinationPath  string `validate:"nonzero" json:"destinationPath"`
}

// JobIdResponse respresents response struct with JobId for folder operations
type JobIdResponse struct {
	JobId string `json:"jobId"`
}

//MoveFolderResponse respresents struct for response to move folder api.
type MoveFolderResponse struct {
	Data JobIdResponse
	api.Response
}

// CreateFolder creates a new folder in media library
func (m *API) CreateFolder(ctx context.Context, param CreateFolderParam) (*api.Response, error) {
	var err error
	var response = &api.Response{}

	if err = validator.Validate(&param); err != nil {
		return nil, err
	}

	resp, err := m.Post(ctx, "folder", &param)
	defer api.DeferredBodyClose(resp)

	api.SetResponseMeta(resp, response)

	if err != nil {
		return response, err
	}

	if resp.StatusCode != 201 {
		err = response.ParseError()
	}

	return response, err
}

// DeleteFolder removes the folder from media library
func (m *API) DeleteFolder(ctx context.Context, param DeleteFolderParam) (*api.Response, error) {
	var err error
	var response = &api.Response{}

	if err = validator.Validate(&param); err != nil {
		return nil, err
	}

	resp, err := m.Delete(ctx, "folder", &param)

	defer api.DeferredBodyClose(resp)

	api.SetResponseMeta(resp, response)

	if err != nil {
		return response, err
	}
	if resp.StatusCode != 204 {
		err = response.ParseError()
	}
	return response, err
}

// MoveFolder moves given folder to new path in media library
func (m *API) MoveFolder(ctx context.Context, param MoveFolderParam) (*MoveFolderResponse, error) {
	var err error
	var response = &MoveFolderResponse{}

	if err = validator.Validate(&param); err != nil {
		return nil, err
	}

	resp, err := m.Post(ctx, "bulkJobs/moveFolder", &param)
	defer api.DeferredBodyClose(resp)
	api.SetResponseMeta(resp, response)

	if err != nil {
		return response, err
	}

	if resp.StatusCode != 200 {
		err = response.ParseError()
	} else {
		err = json.Unmarshal(response.Body(), &response.Data)
	}
	return response, err
}
