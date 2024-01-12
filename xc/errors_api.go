package xc

import (
	"fmt"
	"github.com/xcherryio/apis/goapi/xcapi"
	"net/http"
)

// InvalidArgumentError represents an invalid input argument
type InvalidArgumentError struct {
	msg string
}

func (w InvalidArgumentError) Error() string {
	return fmt.Sprintf("ProcessDefinitionError: %s", w.msg)
}

func NewInvalidArgumentError(tpl string, arg ...interface{}) error {
	return &ProcessDefinitionError{
		msg: fmt.Sprintf(tpl, arg...),
	}
}

// ProcessDefinitionError represents process code(including its elements like AsyncStates/RPCs) is not valid
type ProcessDefinitionError struct {
	msg string
}

func (w ProcessDefinitionError) Error() string {
	return fmt.Sprintf("ProcessDefinitionError: %s", w.msg)
}

func NewProcessDefinitionError(tpl string, arg ...interface{}) error {
	return &ProcessDefinitionError{
		msg: fmt.Sprintf(tpl, arg...),
	}
}

// InternalSDKError means something wrong within xCherry SDK
type InternalSDKError struct {
	Message string
}

func NewInternalError(format string, args ...interface{}) error {
	return &InternalSDKError{
		Message: fmt.Sprintf(format, args...),
	}
}

func (i InternalSDKError) Error() string {
	return fmt.Sprintf("error in SDK or service: message:%v", i.Message)
}

// ApiError represents error returned from xCherry server
// Could be client side(4xx) or server side(5xx), see below helpers to check details
type ApiError struct {
	StatusCode    int
	OriginalError error
	OpenApiError  *xcapi.GenericOpenAPIError
	HttpResponse  *http.Response
	ErrResponse   *xcapi.ApiErrorResponse
}

func (i *ApiError) Error() string {
	if i.ErrResponse != nil {
		bs, err := i.ErrResponse.MarshalJSON()
		if err != nil {
			return fmt.Errorf("failed to MarshalJSON for ApiErrorResponse: %w", err).Error()
		}
		return fmt.Sprintf("StatusCode: %v , error details: %v", i.StatusCode, string(bs))
	}

	errStr := fmt.Sprintf("StatusCode: %v OriginalError:%v", i.StatusCode, i.OriginalError)

	return errStr
}

func NewApiError(
	originalError error, openApiError *xcapi.GenericOpenAPIError, httpResponse *http.Response,
	errResponse *xcapi.ApiErrorResponse,
) error {
	statusCode := 0
	if httpResponse != nil {
		statusCode = httpResponse.StatusCode
	}
	return &ApiError{
		StatusCode:    statusCode,
		OriginalError: originalError,
		OpenApiError:  openApiError,
		HttpResponse:  httpResponse,
		ErrResponse:   errResponse,
	}
}

func IsClientError(err error) bool {
	apiError, ok := err.(*ApiError)
	if !ok {
		return false
	}
	return apiError.StatusCode >= 400 && apiError.StatusCode < 500
}

func IsProcessAlreadyStartedError(err error) bool {
	apiError, ok := err.(*ApiError)
	if !ok || apiError.ErrResponse == nil {
		return false
	}
	return apiError.StatusCode == http.StatusConflict
}

func IsProcessNotExistsError(err error) bool {
	apiError, ok := err.(*ApiError)
	if !ok || apiError.ErrResponse == nil {
		return false
	}
	return apiError.StatusCode == http.StatusNotFound
}

func IsRPCExecutionError(err error) bool {
	apiError, ok := err.(*ApiError)
	if !ok || apiError.ErrResponse == nil {
		return false
	}
	return apiError.StatusCode == http.StatusFailedDependency
}

func IsRPCLockingFailure(err error) bool {
	apiError, ok := err.(*ApiError)
	if !ok || apiError.ErrResponse == nil {
		return false
	}
	return apiError.StatusCode == http.StatusLocked
}

func IsWaitingExceedingTimeoutError(err error) bool {
	apiError, ok := err.(*ApiError)
	if !ok || apiError.ErrResponse == nil {
		return false
	}
	return apiError.StatusCode == http.StatusRequestTimeout
}

// GetOpenApiErrorBody retrieve the API error body into a string to be human-readable
func GetOpenApiErrorBody(err error) string {
	apiError, ok := err.(*ApiError)
	if !ok {
		return "not an ApiError"
	}
	return string(apiError.OpenApiError.Body())
}

// AsProcessAbnormalExitError will check if it's a ProcessAbnormalExitError and convert it if so
func AsProcessAbnormalExitError(err error) (*ProcessAbnormalExitError, bool) {
	wErr, ok := err.(*ProcessAbnormalExitError)
	return wErr, ok
}

// ProcessAbnormalExitError is returned when process execution doesn't complete successfully when waiting on the completion
type ProcessAbnormalExitError struct {
	ProcessExecutionId string
	// TODO ClosedStatus xcapi.ProcessStatus
	// TODO FailureType    *xcapi.ProcessFailureSubType
	ErrorMessage *string
	// StateResults []xcapi.ProcessCloseOutput
	Encoder ObjectEncoder
}

func (w *ProcessAbnormalExitError) Error() string {
	//errTypeMsg := "<nil>"
	//message := "<nil>"
	//if w.ErrorType != nil {
	//	errTypeMsg = fmt.Sprintf("%v", *w.ErrorType)
	//}
	//if w.ErrorMessage != nil {
	//	message = fmt.Sprintf("%v", *w.ErrorMessage)
	//}
	//return fmt.Sprintf("process is not completed successfully, closedStatus: %v, failedType:%v, error message:%v",
	//	w.ClosedStatus, errTypeMsg, message)
	return "TODO"
}
