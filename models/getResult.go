package models

type (
	Headers struct {
		Id int `json:"id"`
		Protocol string `json:"protocol"`
		ContentType string `json:"contentType"`
		Server string `json:"server"`
	}
)
