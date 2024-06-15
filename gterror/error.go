package gterror

import (
	"encoding/json"
	"fmt"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GTError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type GTErrorResponse struct {
	Error *GTError `json:"Error,omitempty"`
}

func (gt *GTError) GTErrorResponse() ([]byte, error) {
	response := &GTErrorResponse{gt}
	return json.Marshal(response)
}

func (gt *GTError) Error() string {
	return fmt.Sprintf("code:%s;message%s", gt.Code, gt.Message)
}

func (gt *GTError) UnknowError(message string) ([]byte, error) {
	gt.Code = INTERNAL_SERVER.Code
	gt.Message = message
	return json.Marshal(gt)
}

func HandleGrpcError(err error) ([]byte, error) {
	if s, ok := status.FromError(err); ok {
		switch s.Code() {
		case codes.InvalidArgument:
			return INVALID_ARGUMENTS.GTErrorResponse()
		case codes.Unimplemented:
			return METHOD_NOT_FOUND.GTErrorResponse()
		case codes.Unauthenticated:
			return UNAUTHORISED_REQUEST.GTErrorResponse()
		default:
			logrus.Info("Grpc Error Unhadle Code : ", s.Code(), s.Message())
			return new(GTError).UnknowError(s.Message())
		}
	}

	return INTERNAL_SERVER.GTErrorResponse()
}

var UNAUTHORISED_REQUEST = &GTError{Code: "UNAUTHORISED_REQUEST", Message: "Unauthorised request, please provide a authentication"}
var INVALID_CREDENTIALS = &GTError{Code: "INVALID_CREDENTIALS", Message: "Invalid Credentials found"}
var SERVICE_NOT_AVAILABEL = &GTError{Code: "SERVICE_NOT_AVAILABEL", Message: "Internal Error , Service Not Availabel"}
var SERVICE_NOT_FOUND = &GTError{Code: "SERVICE_NOT_FOUND", Message: "Service not foud !"}
var INVALID_ARGUMENTS = &GTError{Code: "INVALID_ARGUMENTS", Message: "Invalid arguments found !"}
var INTERNAL_SERVER = &GTError{Code: "INTERNAL_SERVER", Message: "Internal Server Error !"}
var METHOD_NOT_FOUND = &GTError{Code: "METHOD_NOT_FOUND", Message: "Invalid method , service does not contain any method !"}
