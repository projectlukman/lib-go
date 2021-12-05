package response

import (
	"encoding/json"
	"net/http"

	"github.com/projectlukman/lib-go/log"
)

type JSONResponse struct {
	Data       interface{} `json:"data,omitempty"`
	Message    string      `json:"message,omitempty"`
	StatusCode int         `json:"-"`
	Result     interface{} `json:"result,omitempty"`
	Latency    float64     `json:"latency"`
}

func NewJSONResponse() *JSONResponse {
	return &JSONResponse{
		StatusCode: http.StatusOK,
	}
}

func (r *JSONResponse) SetData(data interface{}) *JSONResponse {
	r.Data = data
	return r
}

func (r *JSONResponse) SetMessage(msg string) *JSONResponse {
	r.Message = msg
	return r
}

func (r *JSONResponse) SetStatusCode(statusCode int) *JSONResponse {
	r.StatusCode = statusCode
	return r
}

func (r *JSONResponse) SetResult(result interface{}) *JSONResponse {
	r.Result = result
	return r
}

func (r *JSONResponse) SetLatency(latency float64) *JSONResponse {
	r.Latency = latency
	return r
}

func (r *JSONResponse) GetBody() []byte {
	b, _ := json.Marshal(r)
	return b
}

func (r *JSONResponse) Send(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(r.StatusCode)
	err := json.NewEncoder(w).Encode(r)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err.Error(),
		}).Errorln("[JSONResponse] Error encoding response")
	}
}

// APIStatusSuccess for standard request api status success
func (r *JSONResponse) APIStatusSuccess() *JSONResponse {
	r.Message = http.StatusText(http.StatusOK)
	return r
}

// APIStatusCreated
func (r *JSONResponse) APIStatusCreated() *JSONResponse {
	r.StatusCode = http.StatusCreated
	r.Message = http.StatusText(http.StatusCreated)
	return r
}

// APIStatusAccepted
func (r *JSONResponse) APIStatusAccepted() *JSONResponse {
	r.StatusCode = http.StatusAccepted
	r.Message = http.StatusText(http.StatusAccepted)
	return r
}

// APIStatusErrorUnknown
func (r *JSONResponse) APIStatusErrorUnknown() *JSONResponse {
	r.StatusCode = http.StatusBadGateway
	r.Message = http.StatusText(http.StatusBadGateway)
	return r
}

// APIStatusInvalidAuthentication
func (r *JSONResponse) APIStatusInvalidAuthentication() *JSONResponse {
	r.StatusCode = http.StatusProxyAuthRequired
	r.Message = http.StatusText(http.StatusProxyAuthRequired)
	return r
}

// APIStatusUnauthorized
func (r *JSONResponse) APIStatusUnauthorized() *JSONResponse {
	r.StatusCode = http.StatusUnauthorized
	r.Message = http.StatusText(http.StatusUnauthorized)
	return r
}

// APIStatusForbidden
func (r *JSONResponse) APIStatusForbidden() *JSONResponse {
	r.StatusCode = http.StatusForbidden
	r.Message = http.StatusText(http.StatusForbidden)
	return r
}

// APIStatusBadRequest
func (r *JSONResponse) APIStatusBadRequest() *JSONResponse {
	r.StatusCode = http.StatusBadRequest
	r.Message = http.StatusText(http.StatusBadRequest)
	return r
}
