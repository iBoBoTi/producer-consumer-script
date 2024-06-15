package models

var (
	CreateAction = "create"
	DeleteAction = "delete"
)


type Payload struct {
	Action string `json:"action"`
	Data struct{
		Name string `json:"name"`
		Age int `json:"age"`
	} `json:"data"`
}