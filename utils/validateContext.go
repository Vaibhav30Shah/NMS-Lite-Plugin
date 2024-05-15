package utils

import "NMS-Lite/consts"

func ValidateContext(context map[string]interface{}) (string, float64, string, float64, string) {
	Logger := NewLogger("snmp", context[consts.PluginName].(string))

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

	credentialProfile, ok := context[consts.SnmpCredential].(map[string]interface{})
	if !ok {
		Logger.Error("Credential profile not provided in context")

		return "", 0, "", 0, ""
	}

	community, ok := credentialProfile[consts.SnmpCommunity].(string)

	if !ok {
		Logger.Error("Community string not provided in credential profile")

		return "", 0, "", 0, ""
	}

	version, ok := credentialProfile[consts.SnmpVersion].(string)

	if !ok {
		Logger.Error("Version not provided in credential profile")

		return "", 0, "", 0, ""
	}

	timeOut, ok := context[consts.SnmpRetries].(float64)

	if !ok {
		Logger.Error("Time out not provided in credential profile")

		return "", 0, "", 0, ""
	}

	return ip, port, community, timeOut, version
}
