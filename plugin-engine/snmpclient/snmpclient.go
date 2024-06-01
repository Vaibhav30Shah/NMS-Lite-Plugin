package snmpclient

import (
	"NMS-Lite/consts"
	"NMS-Lite/utils"
	"fmt"
	g "github.com/gosnmp/gosnmp"
	"strings"
	"time"
)

type SNMPClient struct {
	GoSNMP *g.GoSNMP
}

func Init(context map[string]interface{}) (*SNMPClient, error) {

	client := &SNMPClient{

		GoSNMP: &g.GoSNMP{

			Target: context[consts.ObjectIp].(string),

			Community: context[consts.SnmpCommunity].(string),

			Port: context[consts.SnmpPort].(uint16),

			Retries: consts.Retries,

			Timeout: consts.Timeout * time.Second,
		},
	}

	switch context[consts.SnmpVersion] {

	case "1":
		client.GoSNMP.Version = g.Version1

	case "3":
		client.GoSNMP.Version = g.Version3

	default:
		client.GoSNMP.Version = g.Version2c
	}

	err := client.GoSNMP.Connect()

	if err != nil {

		return nil, fmt.Errorf("failed to connect to %v: %v", context["ip"], err)
	}

	return client, nil
}

func (c *SNMPClient) Close() error {

	if c.GoSNMP.Conn != nil {

		return c.GoSNMP.Conn.Close()
	}

	return nil
}

func (c *SNMPClient) Get(oids []string) ([]g.SnmpPDU, error) {

	oid, err := c.GoSNMP.Get(oids)

	if err != nil || oid.Variables == nil {

		return nil, err
	}

	return oid.Variables, nil
}

func (c *SNMPClient) Walk(oidMap map[string]string) ([]interface{}, error) {

	Logger := utils.NewLogger("snmp", "Collect")

	interfacesDetails := make([]interface{}, 0)

	results := map[string]map[string]interface{}{}

	for oidName, oid := range oidMap {

		err := c.GoSNMP.BulkWalk(oid, func(dataUnit g.SnmpPDU) error {

			tokens := strings.Split(dataUnit.Name, ".")

			interfaceIndex := tokens[len(tokens)-1]

			if _, ok := results[interfaceIndex]; !ok {

				results[interfaceIndex] = make(map[string]interface{})
			}

			results[interfaceIndex][oidName] = resolveDataType(dataUnit.Value, dataUnit.Type)

			return nil
		})

		if err != nil {

			Logger.Error(err.Error())

			return nil, err
		}

	}

	for _, interfaceData := range results {

		interfacesDetails = append(interfacesDetails, interfaceData)
	}

	Logger.Debug(fmt.Sprintf("%v", interfacesDetails))

	return interfacesDetails, nil
}

func resolveDataType(value interface{}, dataType g.Asn1BER) interface{} {
	switch dataType {

	case g.OctetString:
		return string(value.([]byte))

	case g.Integer:
		return g.ToBigInt(value)

	case g.Counter32:
		return value.(uint)

	case g.Gauge32:
		return value.(uint)

	case g.TimeTicks:
		return value.(uint)

	default:
		return value
	}
}
