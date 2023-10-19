package xdb

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/xdblab/xdb-apis/goapi/xdbapi"
	"net/http"
)

type basicClientImpl struct {
	options   ClientOptions
	apiClient *xdbapi.APIClient
}

func (u *basicClientImpl) DescribeCurrentProcessExecution(ctx context.Context, processId string) (*xdbapi.ProcessExecutionDescribeResponse, error) {
	req := u.apiClient.DefaultAPI.ApiV1XdbServiceProcessExecutionDescribePost(ctx)

	reqObj := xdbapi.ProcessExecutionDescribeRequest{
		Namespace: u.options.Namespace,
		ProcessId: processId,
	}

	var resp *xdbapi.ProcessExecutionDescribeResponse
	var httpErr error
	if u.options.EnabledDebugLogging {
		fmt.Println("DescribeCurrentProcessExecution is requested", anyToJson(reqObj))
		defer func() {
			fmt.Println("DescribeCurrentProcessExecution is responded", anyToJson(resp), anyToJson(httpErr))
		}()
	}

	resp, httpResp, httpErr := req.ProcessExecutionDescribeRequest(reqObj).Execute()
	if err := u.processError(httpErr, httpResp); err != nil {
		return nil, err
	}
	return resp, nil
}

func (u *basicClientImpl) StartProcess(
	ctx context.Context, processType string, startStateId, processId string, input interface{}, options *BasicClientProcessOptions,
) (string, error) {
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
	reqObj := xdbapi.ProcessExecutionStartRequest{
		Namespace:          u.options.Namespace,
		ProcessId:          processId,
		ProcessType:        processType,
		WorkerUrl:          u.options.WorkerUrl,
		StartStateId:       startStateIdPtr,
		StartStateInput:    encodedInput,
		StartStateConfig:   startStateConfig,
		ProcessStartConfig: processConfig,
	}

	var resp *xdbapi.ProcessExecutionStartResponse
	var httpErr error
	if u.options.EnabledDebugLogging {
		fmt.Println("ProcessExecutionStartRequest is requested", anyToJson(reqObj))
		defer func() {
			fmt.Println("ProcessExecutionStartRequest is responded", anyToJson(resp), anyToJson(httpErr))
		}()
	}
	resp, httpResp, httpErr := req.ProcessExecutionStartRequest(reqObj).Execute()
	if err := u.processError(err, httpResp); err != nil {
		return "", err
	}
	return resp.GetProcessExecutionId(), nil
}

func (u *basicClientImpl) StopProcess(ctx context.Context, processId string, stopType xdbapi.ProcessExecutionStopType) error {
	req := u.apiClient.DefaultAPI.ApiV1XdbServiceProcessExecutionStopPost(ctx)
	reqObj := xdbapi.ProcessExecutionStopRequest{
		Namespace: u.options.Namespace,
		ProcessId: processId,
		StopType:  stopType.Ptr(),
	}

	var httpErr error
	if u.options.EnabledDebugLogging {
		fmt.Println("ProcessExecutionStopRequest is requested", anyToJson(reqObj))
		defer func() {
			fmt.Println("ProcessExecutionStopRequest is responded", anyToJson(httpErr))
		}()
	}
	httpResp, httpErr := req.ProcessExecutionStopRequest(reqObj).Execute()

	if err := u.processError(httpErr, httpResp); err != nil {
		return err
	}
	return nil
}

func (u *basicClientImpl) processError(err error, httpResp *http.Response) error {
	if httpResp != nil {
		defer httpResp.Body.Close()
	}
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

func anyToJson(req any) string {
	str, err := json.Marshal(req)
	if err != nil {
		fmt.Println("failed to encode to Json", err, req)
		return "failed to encode to json"
	}
	return string(str)
}
