package xc

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"reflect"

	"github.com/google/uuid"
	"github.com/xcherryio/apis/goapi/xcapi"
)

type basicClientImpl struct {
	options   ClientOptions
	apiClient *xcapi.APIClient
}

func (u *basicClientImpl) DescribeCurrentProcessExecution(
	ctx context.Context, processId string,
) (*xcapi.ProcessExecutionDescribeResponse, error) {
	req := u.apiClient.DefaultAPI.ApiV1XcherryServiceProcessExecutionDescribePost(ctx)

	reqObj := xcapi.ProcessExecutionDescribeRequest{
		Namespace: u.options.Namespace,
		ProcessId: processId,
	}

	var resp *xcapi.ProcessExecutionDescribeResponse
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
	ctx context.Context, processType string, startStateId, processId string, input interface{},
	options *BasicClientProcessOptions,
) (string, error) {
	var encodedInput *xcapi.EncodedObject
	if input != nil {
		var err error
		encodedInput, err = u.options.ObjectEncoder.Encode(input)
		if err != nil {
			return "", err
		}
	}

	var startStateIdPtr *string
	if startStateId != "" {
		startStateIdPtr = &startStateId
	}
	var startStateConfig *xcapi.AsyncStateConfig
	var processConfig *xcapi.ProcessStartConfig
	if options != nil {
		startStateConfig = options.StartStateOptions
		processConfig = &xcapi.ProcessStartConfig{
			IdReusePolicy:        options.ProcessIdReusePolicy,
			TimeoutSeconds:       &options.TimeoutSeconds,
			LocalAttributeConfig: options.LocalAttributeConfig,
		}
	}

	if u.options.DefaultProcessTimeoutSecondsOverride > 0 {
		if processConfig == nil {
			processConfig = &xcapi.ProcessStartConfig{}
		}
		if processConfig.TimeoutSeconds == nil || *processConfig.TimeoutSeconds == 0 {
			processConfig.TimeoutSeconds = &u.options.DefaultProcessTimeoutSecondsOverride
		}
	}

	req := u.apiClient.DefaultAPI.ApiV1XcherryServiceProcessExecutionStartPost(ctx)
	reqObj := xcapi.ProcessExecutionStartRequest{
		Namespace:          u.options.Namespace,
		ProcessId:          processId,
		ProcessType:        processType,
		WorkerUrl:          u.options.WorkerUrl,
		StartStateId:       startStateIdPtr,
		StartStateInput:    encodedInput,
		StartStateConfig:   startStateConfig,
		ProcessStartConfig: processConfig,
	}

	var resp *xcapi.ProcessExecutionStartResponse
	var httpErr error
	if u.options.EnabledDebugLogging {
		fmt.Println("ProcessExecutionStartRequest is requested", anyToJson(reqObj))
		defer func() {
			fmt.Println("ProcessExecutionStartRequest is responded", anyToJson(resp), anyToJson(httpErr))
		}()
	}
	resp, httpResp, httpErr := req.ProcessExecutionStartRequest(reqObj).Execute()
	if err := u.processError(httpErr, httpResp); err != nil {
		return "", err
	}
	return resp.GetProcessExecutionId(), nil
}

func (u *basicClientImpl) StopProcess(
	ctx context.Context, processId string, stopType xcapi.ProcessExecutionStopType,
) error {
	req := u.apiClient.DefaultAPI.ApiV1XcherryServiceProcessExecutionStopPost(ctx)
	reqObj := xcapi.ProcessExecutionStopRequest{
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

func (u *basicClientImpl) PublishToLocalQueue(
	ctx context.Context, processId string, messages []xcapi.LocalQueueMessage,
) error {
	for _, m := range messages {
		if m.DedupId != nil {
			_, err := uuid.Parse(*m.DedupId)
			if err != nil {
				return fmt.Errorf("invalid dedupUUId %v , err: %w", *m.DedupId, err)
			}
		}

	}

	req := u.apiClient.DefaultAPI.ApiV1XcherryServiceProcessExecutionPublishToLocalQueuePost(ctx)

	reqObj := xcapi.PublishToLocalQueueRequest{
		Namespace: u.options.Namespace,
		ProcessId: processId,
		Messages:  messages,
	}

	var httpErr error
	if u.options.EnabledDebugLogging {
		fmt.Println("PublishToLocalQueue is requested", anyToJson(reqObj))
		defer func() {
			fmt.Println("PublishToLocalQueue is responded", anyToJson(httpErr))
		}()
	}

	httpResp, httpErr := req.PublishToLocalQueueRequest(reqObj).Execute()
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
	if u.options.EnabledDebugLogging {
		if err != nil {
			uerr, ok := err.(*url.Error)
			if ok {
				fmt.Println("encounter url.Error", uerr.Err, uerr.Err.Error())
				uet := reflect.TypeOf(uerr.Err)
				fmt.Println("url.Error.Err type", uet.String(), uet.Name(), uet.Kind())
			}
		}
	}
	var resp *xcapi.ApiErrorResponse
	oerr, ok := err.(*xcapi.GenericOpenAPIError)
	if ok {
		rsp, ok := oerr.Model().(xcapi.ApiErrorResponse)
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
