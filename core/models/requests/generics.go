package requests

type GenericIdRequest struct {
	Id string `json:"id" validate:"required"`
}

type GenericRandomizedIdRequest struct {
	Id string `json:"id" validate:"required,len=7"`
}

type GenericRandomizedIdsRequest struct {
	Ids []string `json:"ids" validate:"required"`
}
