package main

type Target struct {
	Name     string `json:"name" validate:"required,max=200,min=1"`
	Country  string `json:"country" validate:"required,max=200,min=1"`
	Notes    string `json:"notes" validate:"required,max=255,min=1"`
	Complete bool   `json:"complete,default=false"`
}
