package transfer

import (
	"encoding/json"
	"net/http"
)

type UnsignedRequest struct {
	Data interface{}
}

type UnsignedResponse struct {
	Success bool
	Data    interface{}
}

/***********/
/* Request */
/***********/
func ParseUnsigned(r *http.Request, data interface{}) (UnsignedRequest, error) {
	var req UnsignedRequest

	req.Data = data
	err := json.NewDecoder(r.Body).Decode(&req)

	return req, err
}

/************/
/* Response */
/************/
func RespondUnsigned(w http.ResponseWriter, data interface{}, success bool) error {
	var resp UnsignedResponse

	resp.Success = success
	resp.Data = data

	return json.NewEncoder(w).Encode(&resp)
}
