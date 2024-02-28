package models

type (
	Headers struct {
		Protocol string `json:"protocol"`
		ContentType []string `json:"contentType"`
		Server []string `json:"server"`
	}
)
