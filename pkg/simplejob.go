package pkg

type SimpleJob struct {
	Name       string      `json: "name"`
	Containers []Container `json: "containers"`
	Scheduled  bool        `json: "scheduled"`
	MaxRetries int         `json: "maxRetries"`
	Cron       string      `json: "cron,omitempty"`
}

type Container struct {
	Name    string   `json: "name"`
	Image   string   `json: "image"`
	Command []string `json: "command"`
}
