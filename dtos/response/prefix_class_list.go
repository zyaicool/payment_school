package response

type PrefixCLassResponse struct {
	Data []DetailPrefixResponse `json:"data"`
}

type DetailPrefixResponse struct {
	ID         uint   `json:"id"`
	PrefixName string `json:"prefixName"`
}
