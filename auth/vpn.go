package auth

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
)

const (
	apiUrl = "http://check.getipintel.net/check.php?ip=%s&contact=ipinfo@xvh.lol" // The URL for the IP check API.
)

// VPNCheckerInstance is a global instance of VPNChecker.
var VPNCheckerInstance = NewVPNChecker()

// VPNChecker is a struct that checks if an IP address is associated with a VPN.
type VPNChecker struct {
	apiUrl string
}

// NewVPNChecker creates a new instance of VPNChecker.
func NewVPNChecker() *VPNChecker {
	return &VPNChecker{
		apiUrl: apiUrl,
	}
}

// Check checks if the given IP address is associated with a VPN.
// It sends a request to the IP intelligence API and compares the response with a threshold value.
// If the IP intelligence score is greater than 0.99, it considers the IP as a VPN and returns true.
// Otherwise, it returns false.
func (v *VPNChecker) Check(ip string) (bool, error) {
	resp, err := http.Get(fmt.Sprintf(v.apiUrl, ip))
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}

	ipintel, err := strconv.ParseFloat(string(body), 64)
	if err != nil {
		return false, err
	}

	if ipintel > 0.99 {
		return true, nil
	}

	return false, nil
}
