package response

type Ports struct {
	Ports []int `json:"ports"`
	Time  int   `json:"time"`
}

type DirbDns struct {
	Total int `json:"total"`
	Time  int `json:"time"`
}

type Params struct {
	Ip      string `json:"ip"`
	Port    string `json:"port"`
	Process int    `json:"process"`
	Timeout int    `json:"timeout"`
	Debug   int    `json:"debug"`
}
