package response

// Response structure
type BankResponse struct {
	Data []BankData `json:"data"`
}

type BankData struct {
	CodeBank string `json:"codeBank"`
	BankName string `json:"bankName"`
	Alias    string `json:"alias"`
}
