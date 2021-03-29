package log

import apexlog "github.com/apex/log"

type ServiceLog struct {
	ServiceName string
}

func LoadServiceLog(serviceName string) (entry *apexlog.Entry) {
	if l.Handler == nil{
		Load(apexlog.DebugLevel)
	}
	return l.WithFields(apexlog.Fields{
		"name": serviceName,
	})
}
