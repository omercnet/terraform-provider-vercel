package client

import (
	"context"
	"fmt"
)

// SRV defines the metata required for creating an SRV type DNS Record.
type SRV struct {
	Port     int64  `json:"port"`
	Priority int64  `json:"priority"`
	Target   string `json:"target"`
	Weight   int64  `json:"weight"`
}

// CreateDNSRecordRequest defines the information necessary to create a DNS record within Vercel.
type CreateDNSRecordRequest struct {
	Domain     string `json:"-"`
	MXPriority int64  `json:"mxPriority,omitempty"`
	Name       string `json:"name"`
	SRV        *SRV   `json:"srv,omitempty"`
	TTL        int64  `json:"ttl,omitempty"`
	Type       string `json:"type"`
	Value      string `json:"value,omitempty"`
}

// CreateDNSRecord creates a DNS record for a specified domain name within Vercel.
func (c *Client) CreateDNSRecord(ctx context.Context, teamID string, request CreateDNSRecordRequest) (r DNSRecord, err error) {
	url := fmt.Sprintf("%s/v4/domains/%s/records", c.baseURL, request.Domain)
	if c.teamID(teamID) != "" {
		url = fmt.Sprintf("%s?teamId=%s", url, c.teamID(teamID))
	}

	var response struct {
		RecordID string `json:"uid"`
	}
	err = c.doRequest(clientRequest{
		ctx:    ctx,
		method: "POST",
		url:    url,
		body:   string(mustMarshal(request)),
	}, &response)
	if err != nil {
		return r, err
	}

	return c.GetDNSRecord(ctx, response.RecordID, teamID)
}
