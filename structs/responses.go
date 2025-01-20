package structs

type Response struct {
	Error bool                   `json:"error"`
	Code  int32                  `json:"code"`
	Data  map[string]interface{} `json:"data"`
}
