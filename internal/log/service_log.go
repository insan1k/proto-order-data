package log

import apexlog "github.com/apex/log"

// ServiceLog holds the information for this particular service log
type ServiceLog struct {
	ServiceName string
}

// LoadServiceLog creates the log used by a service or subsystem... it's a child of the global log
func LoadServiceLog(serviceName string) (entry *apexlog.Entry) {
	if l.Handler == nil {
		Load(apexlog.DebugLevel)
	}
	return l.WithFields(apexlog.Fields{
		"name": serviceName,
	})
}
