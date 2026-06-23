package faker

type LoginResponse struct {
	Success  bool `json:"success"`
	DelaySec int  `json:"delaySec"`
}

func GetLoginResponse() *LoginResponse {
	return &LoginResponse{
		Success:  true,
		DelaySec: 0,
	}
}
