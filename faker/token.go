package faker

type TokenResponse struct {
	Token string `json:"token"`
}

func GetTokenResponse() *TokenResponse {
	return &TokenResponse{
		Token: getToken(),
	}
}

func getToken() string {
	return "8A06E39EE694873CD4E4B54F75C6DC15770014B2D92EAE5D09BF15A9233680F3402418B7A3CD30214727BE4B51C81BF5E7CFCE7C7D9061B13280689AD7EEC2C2"
}
