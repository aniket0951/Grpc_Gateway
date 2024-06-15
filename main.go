package main

import (
	"context"
	"grpcgateway/codec"
	"grpcgateway/config"
	"grpcgateway/gterror"
	"io"
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func main() {
	config.InitAppConfig()
	config.Init()
	defer config.CloseAllConnections()

	http.HandleFunc("/gw", handler)
	http.ListenAndServe(config.GetAppConfig().Port, nil)
}

func handler(w http.ResponseWriter, r *http.Request) {

	reqbody, err := io.ReadAll(r.Body)

	if err != nil {
		log.Error("Body Reader Error : ", err)
	}
	defer r.Body.Close()

	m := strings.Split(gjson.Get(string(reqbody), "method").String(), ".")
	var response []byte

	if len(m) >= 2 {
		service := m[0]
		method := m[1]

		// authorised the request
		err := authorisedRequest(r)
		if err == nil {
			payload := gjson.Get(string(reqbody), "payload").String()
			response = validateServiceAndMethod(service, method, []byte(payload))
		} else {
			response, _ = err.(*gterror.GTError).GTErrorResponse()
		}

	} else {
		response, _ = gterror.SERVICE_NOT_FOUND.GTErrorResponse()
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func validateServiceAndMethod(service, method string, payload []byte) []byte {

	path, ok := config.GrpcConnMap[service]
	if !ok || path.GrpcClientConnection == nil {
		response, _ := gterror.SERVICE_NOT_AVAILABEL.GTErrorResponse()
		return response
	}

	servicePath := path.ServicePath + "/" + method
	var response codec.WrappedBytesPb
	req1 := codec.WrappedBytesPb{Payload: payload}

	md := metadata.Pairs("authorization", "admin")
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	err := path.GrpcClientConnection.Invoke(ctx, servicePath, req1, &response, grpc.CallContentSubtype("gw-codec"))
	if err != nil {
		log.Info("Service not availabel or not register : " + err.Error())
		response, _ := gterror.HandleGrpcError(err)
		return response
	}

	return response.Payload
}

func authorisedRequest(r *http.Request) error {

	if userName, password, ok := r.BasicAuth(); ok {
		if userName == "admin@gmail.com" && password == "admin" {
			return nil
		}
		return gterror.INVALID_CREDENTIALS
	}
	return gterror.UNAUTHORISED_REQUEST
}
