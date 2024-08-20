package runware

type ConnectRequestData struct {
	TaskType string `json:"taskType"`
	ApiKey   string `json:"apiKey"`
}

func NewConnectRequestData(apikey string) *ConnectRequestData {
	return &ConnectRequestData{
		ApiKey:   apikey,
		TaskType: "authentication",
	}
}
