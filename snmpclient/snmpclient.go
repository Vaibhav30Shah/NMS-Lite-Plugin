package snmpclient

import (
	"encoding/hex"
	"fmt"
	g "github.com/gosnmp/gosnmp"
	"strings"
	"time"
)

type SNMPClient struct {
	GoSNMP *g.GoSNMP
}

func Init(ip string, community string, port uint16, version string, timeout int) (*SNMPClient, error) {

	client := &SNMPClient{

		GoSNMP: &g.GoSNMP{

			Target: ip,

			Community: community,

			Port: port,

			Retries: 3,

			Timeout: time.Duration(timeout) * time.Second,
		},
	}

	switch version {

	case "1":
		client.GoSNMP.Version = g.Version1

	case "2c":
		client.GoSNMP.Version = g.Version2c

	case "3":
		client.GoSNMP.Version = g.Version3

	default:
		return nil, fmt.Errorf("unsupported SNMP version: %s", version)
	}

	err := client.GoSNMP.Connect()

	if err != nil {

		return nil, fmt.Errorf("failed to connect to %s: %v", ip, err)
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

	if err != nil {

		return nil, err
	}

	if oid.Error != 0 {

		return nil, fmt.Errorf("SNMP error: %s", oid.Error)
	}

	return oid.Variables, nil
}

//	func (c *SNMPClient) Walk(oid string) ([]g.SnmpPDU, error) {
//		var wg sync.WaitGroup
//		results := make([]g.SnmpPDU, 0)
//
//		wg.Add(1)
//		go func() {
//			defer wg.Done()
//			tempResults, err := c.GoSNMP.WalkAll(oid)
//			if err != nil {
//				return
//			}
//			results = append(results, tempResults...)
//		}()
//
//		wg.Wait()
//		return results, nil
//	}

func (c *SNMPClient) Walk(oidMap map[string]string) (map[string]interface{}, error) {

	//var wg sync.WaitGroup

	filteredResult := make(map[string]interface{})

	resultMap := make(map[string]interface{})

	for oidName, oid := range oidMap {

		oidResult := make(map[string]interface{})

		err := c.GoSNMP.Walk(oid, func(pdu g.SnmpPDU) error {
			switch pdu.Type {
			case g.OctetString:
				oidResult[pdu.Name] = string(pdu.Value.([]byte))
			case g.Integer:
				oidResult[pdu.Name] = g.ToBigInt(pdu.Value)
			case g.TimeTicks:
				oidResult[pdu.Name] = uint32(pdu.Value.(int64))
			case g.Opaque:
				oidResult[pdu.Name] = fmt.Sprintf("0x%X", pdu.Value.([]byte))
			case g.Counter64:
				oidResult[pdu.Name] = g.ToBigInt(pdu.Value)
			case g.UnknownType:
				oidResult[pdu.Name] = hex.Dump(pdu.Value.([]byte))
			default:
				oidResult[pdu.Name] = pdu.Value
			}
			return nil
		})

		if err != nil {
			// Handle error
		}

		for key, value := range resultMap["result"].(map[string]interface{}) {
			if strings.HasPrefix(key, "interface.admin.status") {
				filteredResult[key] = value
			}
		}

		resultMap[oidName] = oidResult

		//wg.Add(1)
		//
		//results, _ := c.GoSNMP.BulkWalk(oid, 0, 10)
		//
		//go func(results []g.SnmpPDU, oidName string) {
		//
		//	defer wg.Done()
		//
		//	for _, result := range results {
		//
		//		switch result.Type {
		//
		//		case g.OctetString:
		//			oidResult[result.Name] = string(result.Value.([]byte))
		//
		//		case g.Integer:
		//			oidResult[result.Name] = g.ToBigInt(result.Value)
		//
		//		case g.TimeTicks:
		//			oidResult[result.Name] = uint32(result.Value.(int64))
		//
		//		case g.Opaque:
		//			oidResult[result.Name] = fmt.Sprintf("0x%X", result.Value.([]byte))
		//
		//		case g.Counter64:
		//			oidResult[result.Name] = g.ToBigInt(result.Value)
		//
		//		case g.UnknownType:
		//			oidResult[result.Name] = hex.Dump(result.Value.([]byte))
		//
		//		default:
		//			oidResult[result.Name] = result.Value
		//		}
		//	}
		//
		//	resultMap[oidName] = oidResult
		//}(results, oidName)
	}

	//wg.Wait()

	fmt.Println(filteredResult)

	return filteredResult, nil
}
