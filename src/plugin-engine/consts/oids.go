package consts

var ScalerOids = map[string]string{
	"system.name":        ".1.3.6.1.2.1.1.5.0",
	"system.description": ".1.3.6.1.2.1.1.1.0",
	"system.location":    ".1.3.6.1.2.1.1.6.0",
	"system.objectId":    ".1.3.6.1.2.1.1.2.0",
	"system.uptime":      ".1.3.6.1.2.1.1.3.0",
	"system.interfaces":  ".1.3.6.1.2.1.2.1.0",
}

var TabularOids = map[string]string{
	"interface.index":                 ".1.3.6.1.2.1.2.2.1.1",
	"interface.alias":                 ".1.3.6.1.2.1.31.1.1.1.18",
	"interface.name":                  ".1.3.6.1.2.1.31.1.1.1.1",
	"interface.operational.status":    ".1.3.6.1.2.1.2.2.1.8",
	"interface.admin.status":          ".1.3.6.1.2.1.2.2.1.7",
	"interface.description":           ".1.3.6.1.2.1.2.2.1.2",
	"interface.sent.error.packet":     ".1.3.6.1.2.1.2.2.1.20",
	"interface.received.error.packet": ".1.3.6.1.2.1.2.2.1.14",
	"interface.sent.octets":           ".1.3.6.1.2.1.2.2.1.16",
	"interface.received.octets":       ".1.3.6.1.2.1.2.2.1.10",
	"interface.speed":                 ".1.3.6.1.2.1.2.2.1.5",
}
