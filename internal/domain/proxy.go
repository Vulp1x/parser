package domain

import (
	"fmt"
	"net"
	"strconv"
	"strings"

	datasetsservice "github.com/inst-api/parser/gen/datasets_service"
	"github.com/inst-api/parser/internal/pb/instaproxy"
)

type Proxy instaproxy.Proxy
type Proxies []*Proxy

const proxyUploadErrorType = 2

func ParseProxies(proxyRecords []*datasetsservice.ProxyRecord, uploadErrs []*datasetsservice.UploadError) Proxies {
	domainProxies := make([]*Proxy, 0, len(proxyRecords))
	var err error
	for _, proxyRecord := range proxyRecords {
		proxy := Proxy{}
		err = proxy.parse(proxyRecord.Record)
		if err != nil {
			uploadErrs = append(uploadErrs, &datasetsservice.UploadError{
				Type:   proxyUploadErrorType,
				Line:   proxyRecord.LineNumber,
				Input:  strings.Join(proxyRecord.Record, ":"),
				Reason: err.Error(),
			})

			continue
		}

		domainProxies = append(domainProxies, &proxy)
	}

	return domainProxies
}

func (p *Proxy) parse(proxyRecord []string) error {
	ip := net.ParseIP(proxyRecord[0])
	if ip == nil {
		return fmt.Errorf("failed to parse ip")
	}

	port, err := strconv.ParseInt(proxyRecord[1], 10, 32)
	if err != nil {
		return fmt.Errorf("failed to parse port: %v", err)
	}

	p.Host = ip.String()
	p.Port = int32(port)
	p.Login = proxyRecord[2]
	p.Pass = proxyRecord[3]

	return nil
}
