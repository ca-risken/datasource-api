package tag

import "strings"

const (
	// TagGoogle tag
	TagGoogle = "google"
	// TagGCP tag
	TagGCP = "gcp"
	// TagGoogleAssetInventory tag
	TagGoogleAssetInventory = "asset"
	// TagGoogleServiceAccount tag
	TagGoogleServiceAccount = "serviceAccount"
	// TagGoogleCloudSploit tag
	TagGoogleCloudSploit = "cloudsploit"
	// TagGoogleSCC tag
	TagGoogleSCC = "scc"
	// TagGooglePortscan tag
	TagGooglePortscan = "portscan"
)

// UnknownService unknown service name label
const UnknownService string = "unknown"

// GetGCPShortResourceName return short resoruce name from `fullResourceName`. (Resource name format: https://cloud.google.com/asset-inventory/docs/resource-name-format)
func GetGCPShortResourceName(gcpProjectID, fullResourceName string) string {
	service := GetGCPServiceName(fullResourceName)
	array := strings.Split(fullResourceName, "/")
	if len(array) < 2 {
		return getGCPResourceName(gcpProjectID, service, fullResourceName)
	}
	return getGCPResourceName(gcpProjectID, service, array[len(array)-1])
}

// getGCPResourceName return `{gcpProjectID}/{serviceName}/{resourceName}`
func getGCPResourceName(gcpProjectID, serviceName, resourceName string) string {
	return gcpProjectID + "/" + serviceName + "/" + resourceName
}

// GetGCPServiceName return service name from `fullResourceName`. (Resource name format: https://cloud.google.com/asset-inventory/docs/resource-name-format)
func GetGCPServiceName(fullResourceName string) string {
	array := strings.Split(strings.Replace(fullResourceName, "//", "", 1), "/")
	if len(array) < 1 {
		return UnknownService
	}
	svc := array[0]
	if !strings.Contains(svc, ".googleapis.com") {
		return UnknownService
	}
	return strings.ReplaceAll(svc, ".googleapis.com", "")
}
