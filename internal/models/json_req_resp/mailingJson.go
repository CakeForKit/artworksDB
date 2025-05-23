package jsonreqresp

type MailingResponse struct {
	MsgText string   `json:"msg_text"`
	UserIDs []string `json:"user_ids"`
}

type ChangeSubscribeToMailingRequest struct {
	Subscribe bool `json:"subscribe" binding:"boolean" example:"true"`
}
