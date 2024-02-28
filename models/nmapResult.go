package models

type (
	Result struct {
		Url string `json:"domain"`
		Ip string `json:"ip"`
		Ports []Port `json:"ports"`
	}
	Port struct {
		Port uint16 `json:"port"`
		Protocol string `json:"protocol"`
		State string `json:"state"`
		Service string `json:"service"`
	}
)
