package responses

// S3PostPolicyForm lists the fields the client must send in multipart/form-data.
type S3PostPolicyForm struct {
	Key               string `json:"key"`
	Policy            string `json:"policy"`
	Algorithm         string `json:"x-amz-algorithm"`
	Credential        string `json:"x-amz-credential"`
	Date              string `json:"x-amz-date"`
	Signature         string `json:"x-amz-signature"`
	ContentType       string `json:"Content-Type"`
	MetaSHA256        string `json:"x-amz-meta-sha256,omitempty"`
	SuccessActionCode string `json:"success_action_status"`
}

// CreateUploadResponsePostPolicy is the payload returned by POST /api/uploads.
type CreateUploadResponsePostPolicy struct {
	UploadURL   string           `json:"uploadUrl"`
	ResourceURL string           `json:"resourceUrl"`
	ExpiresAt   string           `json:"expiresAt"`
	Form        S3PostPolicyForm `json:"form"`
}
