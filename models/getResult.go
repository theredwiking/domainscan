package models

type (
	Headers struct {
		Id int64`json:"id"`
		Protocol string `json:"protocol"`
		ContentType string `json:"contentType"`
		Server string `json:"server"`
	}
)
