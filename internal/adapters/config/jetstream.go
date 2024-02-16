package config

type NATSConnectionConfig struct {
	URL                    string `xconfig:"url"`
	MaxReconnects          int    `xconfig:"maxReconnects=3"`
	ReconnectWaitMilliSecs int    `xconfig:"reconnectWaitMilliSecs=1000"`
	TimeoutMilliSecs       int    `xconfig:"timeoutMilliSecs=3000"`
	CredentialsPath        string `xconfig:"credentialsPath"`
}

type JetStreamStreamConfig struct {
	Name            string
	SourceSubjects  string // comma delimited list of subjects
	DestSubjects    string
	MaxMsgs         int
	MaxAgeMilliSecs int
	MaxBytes        int64
	Replicas        int
	// Retention
}

type JetStreamConsumerConfig struct {
	Name             string
	Subjects         string // comma delimited list of subjects
	IsDurable        bool
	AckWaitMilliSecs int
	MaxDeliver       int
	MaxAckPending    int
}
