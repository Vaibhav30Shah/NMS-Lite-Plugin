package snmp

import "NMS-Lite/consts"

func validateContext(context map[string]interface{}) {

	if _, ok := context[consts.ObjectIp]; !ok {
		context[consts.ObjectIp] = "127.0.0.1"
	}

	if _, ok := context[consts.SnmpCommunity]; !ok {
		context[consts.SnmpCommunity] = "public"
	}

	if _, ok := context[consts.SnmpPort]; !ok {
		context[consts.SnmpPort] = 161
	}

}
