package xdb

import (
	"context"
	"github.com/xdblab/xdb-apis/goapi/xdbapi"
	"net/http"
)

type basicClientImpl struct {
	options   *ClientOptions
	apiClient *xdbapi.APIClient
}

func (u *basicClientImpl) DescribeCurrentProcessExecution(ctx context.Context, processId string) (*xdbapi.ProcessExecutionDescribeResponse, error) {
	req := u.apiClient.DefaultAPI.ApiV1XdbServiceProcessExecutionDescribePost(ctx)
	resp, httpResp, err := req.ProcessExecutionDescribeRequest(xdbapi.ProcessExecutionDescribeRequest{
		Namespace: u.options.Namespace,
		ProcessId: processId,
	}).Execute()
	if err := u.processError(err, httpResp); err != nil {
		return nil, err
	}
	return resp, nil
}

func (u *basicClientImpl) StartProcess(ctx context.Context, processType string, startStateId, processId string, input interface{}, options *BasicClientProcessOptions) (string, error) {
	var encodedInput *xdbapi.EncodedObject
	var err error
	if input != nil {
		encodedInput, err = u.options.ObjectEncoder.Encode(input)
		if err != nil {
			return "", err
		}
	}

	var startStateIdPtr *string
	if startStateId != "" {
		startStateIdPtr = &startStateId
	}
	var startStateConfig *xdbapi.AsyncStateConfig
	var processConfig *xdbapi.ProcessStartConfig
	if options != nil {
		startStateConfig = options.StartStateOptions
		processConfig = &xdbapi.ProcessStartConfig{
			IdReusePolicy:  options.ProcessIdReusePolicy,
			TimeoutSeconds: &options.TimeoutSeconds,
		}
	}

	req := u.apiClient.DefaultAPI.ApiV1XdbServiceProcessExecutionStartPost(ctx)
	resp, httpResp, err := req.ProcessExecutionStartRequest(xdbapi.ProcessExecutionStartRequest{
		Namespace:          u.options.Namespace,
		ProcessId:          processId,
		ProcessType:        processType,
		WorkerUrl:          u.options.WorkerUrl,
		StartStateId:       startStateIdPtr,
		StartStateInput:    encodedInput,
		StartStateConfig:   startStateConfig,
		ProcessStartConfig: processConfig,
	}).Execute()
	if err := u.processError(err, httpResp); err != nil {
		return "", err
	}
	return resp.GetProcessExecutionId(), nil
}

func (u *basicClientImpl) processError(err error, httpResp *http.Response) error {
	if err == nil && httpResp != nil && httpResp.StatusCode == http.StatusOK {
		return nil
	}
	var resp *xdbapi.ApiErrorResponse
	oerr, ok := err.(*xdbapi.GenericOpenAPIError)
	if ok {
		rsp, ok := oerr.Model().(xdbapi.ApiErrorResponse)
		if ok {
			resp = &rsp
		}
	}
	return NewApiError(err, oerr, httpResp, resp)
}
