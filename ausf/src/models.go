package ausf

type AuthRequest struct {
	UEID string `json:"ue_id"`
	SUPI string `json:"supi"`
}

type AuthResponse struct {
	RAND  string `json:"rand"`
	AUTN  string `json:"autn"`
	XRES  string `json:"xres"`
	KASME string `json:"kasme"`
}

type VerifyRequest struct {
	UEID       string `json:"ue_id"`
	Response   string `json:"response"`
}

type VerifyResponse struct {
	Result string `json:"result"`
}

