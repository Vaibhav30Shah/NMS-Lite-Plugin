package snmpclient

import (
	"fmt"
	g "github.com/gosnmp/gosnmp"
	"sync"
	"time"
)

type SNMPClient struct {
	GoSNMP *g.GoSNMP
}

func Init(ip string, community string, port uint16, version string, timeout int) (*SNMPClient, error) {
	client := &SNMPClient{
		GoSNMP: &g.GoSNMP{
			Target:    ip,
			Community: community,
			Port:      port,
			Retries:   3,
			Timeout:   time.Duration(timeout) * time.Second,
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
	packet, err := c.GoSNMP.Get(oids)
	if err != nil {
		return nil, err
	}

	if packet.Error != 0 {
		return nil, fmt.Errorf("SNMP error: %s", packet.Error)
	}

	return packet.Variables, nil
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
	var wg sync.WaitGroup
	resultMap := make(map[string]interface{})

	for oidName, oid := range oidMap {
		oidResult := make(map[string]interface{})
		wg.Add(1)
		results, _ := c.GoSNMP.WalkAll(oid)
		go func(results []g.SnmpPDU, oidName string) {
			defer wg.Done()
			for _, result := range results {
				switch result.Type {
				case g.OctetString:
					oidResult[result.Name] = string(result.Value.([]byte))
				case g.Integer:
					oidResult[result.Name] = g.ToBigInt(result.Value)
				case g.TimeTicks:
					oidResult[result.Name] = uint32(result.Value.(int64))
				case g.Opaque:
					oidResult[result.Name] = fmt.Sprintf("0x%X", result.Value.([]byte))
				case g.Counter64:
					oidResult[result.Name] = g.ToBigInt(result.Value)
				default:
					oidResult[result.Name] = result.Value
				}
			}
			resultMap[oidName] = oidResult
		}(results, oidName)
	}

	wg.Wait()
	return resultMap, nil
}
