package tag

const (
	// TagOsint osint tag
	TagOsint = "osint"
	// TagPrivateExpose private-expose tag
	TagPrivateExpose = "private-expose"
	// TagTakeover takeover tag
	TagTakeover = "takeover"
	// TagCertificateExpiration certificate-expiration tag
	TagCertificateExpiration = "certificate-expiration"
	// Tag certificate-expiration tag
	TagWebTechnology = "web-technology"

	// Osint Resource Type
	// TagOSINTUnknown unknown service tag
	TagOSINTUnknown = "unknown-resource-type"
	// TagDomain domain tag
	TagDomain = "domain"
	// TagWebsite website tag
	TagWebsite = "website"
)

// GetTakeOverList return TakeOverList
func GetTakeOverList() []string {
	return TakeOverList
}

// TakeOverList :takeoverされる危険の高いドメインを列挙
var TakeOverList = []string{
	".herokuapp.com",
	".herokussl.com",
	".azurewebsites.net",
	".cloudapp.net",
	".azure-api.net",
	".trafficmanager.net",
	".azureedge.net",
	".cloudapp.azure.com",
	".s3.amazonaws.com",
	".awsptr.com",
	".elasticbeanstalk.com",
	".uservoice.com",
	".unbouncepages.com",
	".ghs.google.com",
	".ghs.googlehosted.com",
	".ghs-ssl.googlehosted.com",
	".github.io",
	".www.gitbooks.io",
	".sendgrid.net",
	".feedpress.me",
	".fastly.net",
	".webflow.io",
	".proxy.webflow.com",
	".helpscoutdocs.com",
	".readmessl.com",
	".desk.com",
	".zendesk.com",
	".mktoweb.com",
	".wordpress.com",
	".wpengine.com",
	".cloudflare.net",
	".netlify.com",
	".bydiscourse.com",
	".netdna-cdn.com",
	".pageserve.co",
	".pantheonsite.io",
	".arlo.co",
	".apigee.net",
	".pmail5.com",
	".cm-hosting.com",
	".ext-cust.squarespace.com",
	".ext.squarespace.com",
	".www.squarespace6.com",
	".locationinsight.com",
	".helpsite.io",
	".saas.moonami.com",
	".custom.bnc.lt",
	".qualtrics.com",
	".dotcmscloud.net",
	".dotcmscloud.com",
	".knowledgeowl.com",
	".atlashost.eu",
	".headwayapp.co",
	".domain.pixieset.com",
	".cname.bitly.com",
	".awmdm.com",
	".meteor.com",
	".postaffiliatepro.com",
	".na.iso.postaffiliatepro.com",
	".copiny.com",
	".kxcdn.com",
	".phs.getpostman.com",
	".appdirect.com",
	".streamshark.io",
}
