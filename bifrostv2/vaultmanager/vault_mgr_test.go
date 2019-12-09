package vaultmanager

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	. "gopkg.in/check.v1"

	"gitlab.com/thorchain/thornode/bifrostv2/config"
	"gitlab.com/thorchain/thornode/bifrostv2/metrics"
)

func Test(t *testing.T) {
	TestingT(t)
}

type VaultsMgrSuite struct {
	server *httptest.Server
}

var _ = Suite(&VaultsMgrSuite{})

func (s *VaultsMgrSuite) SetUpSuite(c *C) {
	s.server = httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		switch req.RequestURI {
		case "/thorchain/vaults/pubkeys":
			vaultsHandle(c, rw)
		}
	}))
}

func vaultsHandle(c *C, rw http.ResponseWriter) {
	content, err := ioutil.ReadFile("../../test/fixtures/endpoints/vaults/pubKeys.json")
	if err != nil {
		c.Fatal(err)
	}

	rw.Header().Set("Content-Type", "application/json")
	if _, err := rw.Write(content); err != nil {
		c.Fatal(err)
	}
}

func getMetricForTest(c *C) *metrics.Metrics {
	m, err := metrics.NewMetrics(config.MetricConfiguration{
		Enabled:      false,
		ListenPort:   9000,
		ReadTimeout:  time.Second,
		WriteTimeout: time.Second,
	})
	c.Assert(m, NotNil)
	c.Assert(err, IsNil)
	return m
}

func (s *VaultsMgrSuite) TestGetVaults(c *C) {
	vaultMgr, err := NewVaultManager(s.server.URL, getMetricForTest(c))
	c.Assert(err, IsNil)
	c.Assert(vaultMgr, NotNil)

	vaults, err := vaultMgr.getVaults()
	c.Assert(err, IsNil)
	c.Assert(vaults, NotNil)
	c.Assert(vaults.Asgard[0].String(), Equals, "thorpub1addwnpepqflvfv08t6qt95lmttd6wpf3ss8wx63e9vf6fvyuj2yy6nnyna5763e2kck")
}