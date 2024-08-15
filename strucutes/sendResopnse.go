package strucutes

import "go/types"

type SendResponse struct {
	Id          string       `json:"id"`
	Description string       `json:"description"`
	CreatedAt   string       `json:"createdAt"`
	CreatedBy   string       `json:"createdBy"`
	ModifiedAt  string       `json:"modifiedAt"`
	Done        bool         `json:"done"`
	Metadata    types.Object `json:"metadata"`
	Error       struct {
		Code    float64        `json:"code"`
		Message string         `json:"message"`
		Details []types.Object `json:"details"`
	} `json:"error"`
	Response types.Object `json:"response"`
}
