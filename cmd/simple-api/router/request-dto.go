package router

type RequestDto struct {
	Version       int           `json:"version" binding:"required"`
	TransactionId string        `json:"operationId" binding:"required"`
	TenantId      string        `json:"tenantId" binding:"required"`
	SessionId     string        `json:"sessionId" binding:"required"`
	Source        string        `json:"source" binding:"required"`
	Family        string        `json:"family" binding:"required"`
	Events        []interface{} `json:"events" binding:"required"`
}
