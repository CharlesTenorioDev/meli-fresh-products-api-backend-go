package dto

type SellersPostRequest struct {
	CID         int    `json:"cid"`
	CompanyName string `json:"company_name"`
	Address     string `json:"address"`
	Telephone   string `json:"telephone"`
}

type SellersUpdateRequest struct {
	CID         *int    `json:"cid"`
	CompanyName *string `json:"company_name"`
	Address     *string `json:"address"`
	Telephone   *string `json:"telephone"`
}

type SellersGetResponse struct {
	Id          int    `json:"id"`
	Cid         int    `json:"cid"`
	CompanyName string `json:"company_name"`
	Address     string `json:"address"`
	Telephone   string `json:"telephone"`
}

type DataResponse struct {
	Data interface{} `json:"data"`
}

type SellerPostResponse struct {
	Id int `json:"id"`
}
