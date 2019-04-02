package crawler

import (
	"bytes"
	"io"
	"io/ioutil"
	"testing"

	"webcrawler/config"

	"github.com/stretchr/testify/assert"
)

func TestCrawlerInvalidDomain(t *testing.T) {
	cfg, err := config.Parse("../config/config.json")
	assert.Nil(t, err, "Should be nil")
	assert.NotNil(t, cfg, "Shouldn't be nil")

	c := New(cfg)
	assert.NotNil(t, c, "Shouldn't be nil")
	assert.IsType(t, &Crawler{}, c, "Type should be same")

	startURL := "language/en"
	go func() {
		err := c.Crawl(startURL)
		assert.NotNil(t, err, "Shouldn't be nil")
		assert.EqualError(t, err, "Invalid domain", "Error value should be equal")
	}()
	<-c.Quit
}

func TestCrawler(t *testing.T) {
	cfg, err := config.Parse("../config/config.json")
	assert.Nil(t, err, "Should be nil")
	assert.NotNil(t, cfg, "Shouldn't be nil")

	c := New(cfg)
	assert.NotNil(t, c)

	flip := true

	// Fetcher func
	fetcher := func(u string) (io.ReadCloser, error) {
		seed := u
		if flip {
			seed = "./test_data.html"
			flip = false
		}
		b, err := ioutil.ReadFile(seed)
		if err != nil {
			return nil, err
		}
		return ioutil.NopCloser(bytes.NewReader(b)), nil
	}
	c.FetcherFunc = fetcher

	startURL := "http://www.redhat.com/en"
	go c.Crawl(startURL)
	<-c.Quit

	assert.NotNil(t, c.Sitemap, "Should be not nil")
	assert.Equal(t, expectedSitemap, c.Sitemap, "Outcome should be equal")
	assert.IsType(t, map[string][]string{}, c.Sitemap, "Type should be equal")
}

var expectedSitemap = map[string][]string{
	"http://www.redhat.com/en": {
		"https://www.redhat.com/wapps/ugc/register.html",
		"https://www.redhat.com/wapps/ugc/protected/account.html",
		"http://www.redhat.com/en/technologies",
		"http://www.redhat.com/en/challenges",
		"http://www.redhat.com/en/services",
		"http://www.redhat.com/en/about",
		"http://www.redhat.com/en/partners",
		"http://www.redhat.com/en/about/open-source",
		"http://www.redhat.com/en/store",
		"http://www.redhat.com/en/search",
		"http://www.redhat.com/en/about/around-the-world",
		"http://www.redhat.com/en/about/our-culture?intcmp=701f2000000u65hAAA",
		"http://www.redhat.com/en/about?intcmp=701f2000001D6x7AAC",
		"http://www.redhat.com/en/topics/cloud-computing/what-is-hybrid-cloud?intcmp=701f2000000u65XAAQ",
		"http://www.redhat.com/en/command-line-heroes?intcmp=701f2000001D982AAC",
		"http://www.redhat.com/en/summit/2019?intcmp=701f2000001D97dAAC",
		"http://www.redhat.com/en/technologies/cloud-computing/openshift/application-runtimes?intcmp=701f2000001OEGrAAO",
		"http://www.redhat.com/en/technologies/cloud-computing/openshift?intcmp=701f2000001OEGrAAO",
		"http://www.redhat.com/en/technologies/linux-platforms/openstack-platform?intcmp=701f2000001OEGrAAO",
		"http://www.redhat.com/en/technologies/storage/ceph?intcmp=701f2000001OEGrAAO",
		"http://www.redhat.com/en/technologies/cloud-computing/cloud-suite?intcmp=701f2000001OEGrAAO",
		"http://www.redhat.com/en/technologies/cloud-computing/cloud-infrastructure?intcmp=701f2000001OEGrAAO",
		"http://www.redhat.com/en/technologies/cloud-computing/directory-server?intcmp=701f2000001OEGrAAO",
		"http://www.redhat.com/en/technologies/cloud-computing/old-certificate-system?intcmp=701f2000001OEGrAAO",
		"http://www.redhat.com/en/technologies/management/cloudforms?intcmp=701f2000001OEGrAAO",
		"http://www.redhat.com/en/technologies/management/insights?intcmp=701f2000001OEGrAAO",
		"http://www.redhat.com/en/technologies/linux-platforms/enterprise-linux?intcmp=701f2000001OEGhAAO",
		"http://www.redhat.com/en/technologies/linux-platforms/openstack-platform?intcmp=701f2000001OEGhAAO",
		"http://www.redhat.com/en/technologies/management/ansible-old?intcmp=701f2000001OEGhAAO",
		"http://www.redhat.com/en/technologies/management/satellite?intcmp=701f2000001OEGhAAO",
		"http://www.redhat.com/en/technologies/virtualization/enterprise-virtualization?intcmp=701f2000001OEGhAAO",
		"http://www.redhat.com/en/technologies/storage/old-gluster?intcmp=701f2000001OEGhAAO",
		"http://www.redhat.com/en/technologies/management/insights?intcmp=701f2000001OEGhAAO",
		"http://www.redhat.com/en/technologies/jboss-middleware/fuse?intcmp=701f2000001OEH1AAO",
		"http://www.redhat.com/en/technologies/jboss-middleware/fuse-online?intcmp=701f2000001OEH1AAO",
		"http://www.redhat.com/en/technologies/jboss-middleware/amq?intcmp=701f2000001OEH1AAO",
		"http://www.redhat.com/en/technologies/jboss-middleware/3scale?intcmp=701f2000001OEH1AAO",
		"http://www.redhat.com/en/technologies/jboss-middleware/old-data-virtualization?intcmp=701f2000001OEH1AAO",
		"http://www.redhat.com/en/technologies/cloud-computing/openshift/application-runtimes?intcmp=701f2000001OEGwAAO",
		"http://www.redhat.com/en/technologies/cloud-computing/openshift?intcmp=701f2000001OEGwAAO",
		"http://www.redhat.com/en/technologies/jboss-middleware/developer-studio?intcmp=701f2000001OEGwAAO",
		"http://www.redhat.com/en/technologies/jboss-middleware/decision-manager?intcmp=701f2000001OEGwAAO",
		"http://www.redhat.com/en/technologies/jboss-middleware/old-data-grid?intcmp=701f2000001OEGwAAO",
		"http://www.redhat.com/en/technologies/jboss-middleware/process-automation-manager?intcmp=701f2000001OEGwAAO",
		"http://www.redhat.com/en/technologies/jboss-middleware/application-platform-old?intcmp=701f2000001OEGwAAO",
		"http://www.redhat.com/en/technologies/jboss-middleware/old-web-server?intcmp=701f2000001OEGwAAO",
		"http://www.redhat.com/en/technologies/management/ansible-old?intcmp=701f2000001OEGmAAO",
		"http://www.redhat.com/en/technologies/management/cloudforms?intcmp=701f2000001OEGmAAO",
		"http://www.redhat.com/en/technologies/management/old-satellite?intcmp=701f2000001OEGmAAO",
		"http://www.redhat.com/en/technologies/management/insights?intcmp=701f2000001OEGmAAO",
		"http://www.redhat.com/en/technologies/cloud-computing/directory-server",
		"http://www.redhat.com/en/technologies/management/old-satellite?intcmp=701f2000001OEGhAAO",
		"http://www.redhat.com/en/technologies?intcmp=701f2000001D6xCAAS",
		"http://www.redhat.com/en/about/open-source?intcmp=701f2000001OEHBAA4",
		"http://www.redhat.com/en/about/company",
		"http://www.redhat.com/en/about/office-locations",
		"http://www.redhat.com/en/blog",
		"http://www.redhat.com/en/about/newsroom",
		"http://www.redhat.com/en/about/development-model",
		"http://www.redhat.com/en/events",
		"http://www.redhat.com/en/jobs",
		"http://www.redhat.com/en/technologies/linux-platforms/enterprise-linux",
		"http://www.redhat.com/en/technologies/storage/old-gluster",
		"http://www.redhat.com/en/technologies/management/old-satellite",
		"http://www.redhat.com/en/technologies/cloud-computing/openshift",
		"http://www.redhat.com/en/technologies/linux-platforms/openstack-platform",
		"https://www.redhat.com/wapps/sso/login.html",
		"https://www.redhat.com/en/partners",
		"http://www.redhat.com/en/resources",
		"http://www.redhat.com/en/about/japan-buy",
		"http://www.redhat.com/en/contact",
		"http://www.redhat.com/en/services/training/contact",
		"http://www.redhat.com/en/services/consulting",
		"http://www.redhat.com/en/about/feedback",
		"http://www.redhat.com/en/about/social",
		"http://www.redhat.com/en/all-blogs",
		"http://www.redhat.com/en/about/videos",
		"http://www.redhat.com/en/topics",
		"http://www.redhat.com/en/about/privacy-policy",
		"http://www.redhat.com/en/about/terms-use",
		"http://www.redhat.com/en/about/all-policies-guidelines",
		"http://www.redhat.com/summit/",
	},
}
