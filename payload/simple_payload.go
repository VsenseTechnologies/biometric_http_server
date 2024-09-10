package payload

type SimpleSuccessPayload struct{
	Message string `json:"message"`
}

type SimpleFailedPayload struct{
	ErrorMessage string `json:"message"`
}

type SuccessPayloadWithData struct{
	Message string `json:"message"`
	Data any `json:"data"`
}

type JWTTokenPayload struct{
	Token string `json:"token"`
}