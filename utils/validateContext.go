package utils

func ValidateContext(context map[string]interface{}) (string, float64, string, float64, string) {
	Logger := NewLogger("snmp", context["plugin_type"].(string))

	ip, ok := context["ip"].(string)

	if !ok {
		Logger.Error("IP address not provided in context")

		return "", 0, "", 0, ""
	}

	port, ok := context["port"].(float64)

	if !ok {
		Logger.Error("Port not provided in context")

		return "", 0, "", 0, ""
	}

	credentialProfile, ok := context["credential_profile"].(map[string]interface{})
	if !ok {
		Logger.Error("Credential profile not provided in context")

		return "", 0, "", 0, ""
	}

	community, ok := credentialProfile["community"].(string)

	if !ok {
		Logger.Error("Community string not provided in credential profile")

		return "", 0, "", 0, ""
	}

	version, ok := credentialProfile["version"].(string)

	if !ok {
		Logger.Error("Version not provided in credential profile")

		return "", 0, "", 0, ""
	}

	timeOut, ok := context["retry_count"].(float64)

	if !ok {
		Logger.Error("Time out not provided in credential profile")

		return "", 0, "", 0, ""
	}

	return ip, port, community, timeOut, version
}
