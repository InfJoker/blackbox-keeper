package configuration

type ConfigXml struct {
	Apps struct {
		Command     string `yaml:"command" xml:"command"`
		HealthCheck struct {
			Http struct {
				Host                string `yaml:"host" xml:"host"`
				Port                int64  `yaml:"port" xml:"port"`
				Path                string `yaml:"path" xml:"path"`
				WaitAfterStartMilli int64  `yaml:"wait-after-start" xml:"wait-after-start"`
				RepeatAfterMilli    int64  `yaml:"repeat-after" xml:"repeat-after"`
				TimeoutMilli        int64  `yaml:"timeout" xml:"timeout"`
			} `yaml:"http" xml:"http"`
			StopAction struct {
				Signal struct {
					SignalType   string `yaml:"signal-type" xml:"signal-type"`
					TimeoutMilli int64  `yaml:"timeout" xml:"timeout"`
				} `yaml:"signal" xml:"signal"`
			} `yaml:"stop-action" xml:"stop-action"`
		} `yaml:"health-check" xml:"health-check"`
	} `yaml:"apps" xml:"apps"`
}
