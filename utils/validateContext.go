package utils

import (
	"NMS-Lite/consts"
)

func ValidateContext(context map[string]interface{}) (string, float64, string, float64, string) {
	Logger := NewLogger("snmp", context[consts.PluginName].(string))

	var version, community string

	ip, ok := context[consts.ObjectIp].(string)

	if !ok {
		Logger.Error("IP address not provided in context")

		return "", 0, "", 0, ""
	}

	port, ok := context[consts.SnmpPort].(float64)

	if !ok {
		Logger.Error("Port not provided in context")

		return "", 0, "", 0, ""
	}

	if credentials, ok := context[consts.SnmpCredential].([]interface{}); ok {

		//checking for multiple credential profile
		for _, credential := range credentials {

			credMap, ok := credential.(map[string]interface{})

			if !ok {

				Logger.Error("Invalid credential profile format")

				continue
			}

			tempVersion, ok := credMap["version"].(string)

			if !ok {

				Logger.Error("Version not provided in credential profile")

				continue
			}

			tempCommunity, ok := credMap["community"].(string)

			if !ok {

				Logger.Error("Community string not provided in credential profile")

				continue
			}

			// If both version and community are valid, use them
			version = tempVersion

			community = tempCommunity

			break
		}

		if version == "" || community == "" {

			Logger.Error("Version or community string not found in credential profile")

			return "", 0, "", 0, ""
		}
	}

	timeOut, ok := context[consts.SnmpRetries].(float64)

	if !ok {

		Logger.Error("Time out not provided in credential profile")

		return "", 0, "", 0, ""
	}

	return ip, port, community, timeOut, version
}
