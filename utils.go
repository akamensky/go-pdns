package pdns

import "fmt"

func isFQDN(hostname string) bool {
	if len(hostname) == 0 {
		return false
	} else {
		return hostname[len(hostname)-1] == '.'
	}
}

func FQDN(hostname string) string {
	if isFQDN(hostname) {
		return hostname
	} else {
		return fmt.Sprintf("%s.", hostname)
	}
}
