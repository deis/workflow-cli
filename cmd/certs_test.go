package cmd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/arschles/assert"
	"github.com/deis/controller-sdk-go/api"
	"github.com/deis/workflow-cli/pkg/testutil"
)

func TestCertsList(t *testing.T) {
	t.Parallel()
	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()
	var b bytes.Buffer
	cmdr := DeisCmd{WOut: &b, ConfigFile: cf}

	server.Mux.HandleFunc("/v2/certs/", func(w http.ResponseWriter, r *http.Request) {
		testutil.SetHeaders(w)
		fmt.Fprintf(w, `{
			"count": 2,
			"next": null,
			"previous": null,
			"results": [
				{
					"name": "test-example-com",
					"common_name": "test.example.com",
					"san": [
						"test.com",
						"example.com"
					],
					"domains": [
						"test.com",
						"example.com"
					],
					"created": "2016-06-09T00:00:00UTC",
					"updated": "2016-06-09T00:00:00UTC",
					"expires": "2014-11-10T00:00:00UTC",
					"fingerprint": "12:34:56:78:90"
				},
				{
					"name": "test-deis-com",
					"common_name": "test.deis.com",
					"created": "2016-06-09T00:00:00UTC",
					"updated": "2016-06-09T00:00:00UTC",
					"expires": "2016-08-01T00:00:00UTC",
					"fingerprint": "ab:12:ab:12:ab"
				},
				{
					"name": "test1",
					"common_name": "1.test.deis.com",
					"expires": "2016-06-11T00:00:00UTC"
				},
				{
					"name": "test2",
					"common_name": "2.test.deis.com",
					"expires": "2018-01-01T00:00:00UTC"
				}
			]
		}`)
	})

	err = cmdr.CertsList(-1, time.Date(2016, time.June, 9, 0, 0, 0, 0, time.UTC))
	assert.NoErr(t, err)

	assert.Equal(t, b.String(), `        Name       |   Common Name    |    SubjectAltName    |         Expires          |   Fingerprint   |       Domains        |  Updated   |  Created    
+------------------+------------------+----------------------+--------------------------+-----------------+----------------------+------------+------------+
  test-example-com | test.example.com | test.com,example.com | 10 Nov 2014 (expired)    | 12:34[...]78:90 | test.com,example.com | 9 Jun 2016 | 9 Jun 2016  
  test-deis-com    | test.deis.com    |                      | 1 Aug 2016 (in 2 months) | ab:12[...]12:ab |                      | 9 Jun 2016 | 9 Jun 2016  
  test1            | 1.test.deis.com  |                      | 11 Jun 2016 (in 2 days)  |                 |                      | unknown    | unknown     
  test2            | 2.test.deis.com  |                      | 1 Jan 2018 (in 2 years)  |                 |                      | unknown    | unknown     
`, "output")

	cf, server, err = testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()

	cmdr.ConfigFile = cf
	b.Reset()

	server.Mux.HandleFunc("/v2/certs/", func(w http.ResponseWriter, r *http.Request) {
		testutil.SetHeaders(w)
		fmt.Fprintf(w, `{
			"count": 0,
			"next": null,
			"previous": null,
			"results": []
		}`)
	})

	err = cmdr.CertsList(-1, time.Now())
	assert.NoErr(t, err)

	assert.Equal(t, b.String(), "No certs\n", "output")
}

func TestCertsInfo(t *testing.T) {
	t.Parallel()
	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()
	var b bytes.Buffer
	cmdr := DeisCmd{WOut: &b, ConfigFile: cf}

	server.Mux.HandleFunc("/v2/certs/test-example-com", func(w http.ResponseWriter, r *http.Request) {
		testutil.SetHeaders(w)
		fmt.Fprintf(w, `{
			"name": "test-example-com",
			"owner": "admin",
			"issuer": "testca",
			"subject": "testing",
			"common_name": "test.deis.com",
			"created": "2016-06-09T00:00:00UTC",
			"updated": "2016-06-09T00:00:00UTC",
			"expires": "2016-06-09T00:00:00UTC",
			"starts": "2016-06-09T00:00:00UTC",
			"fingerprint": "ab:12:ab:12:ab",
			"san": [
				"test.com",
				"example.com"
			],
			"domains": [
				"test.com",
				"example.com"
			]
		}`)
	})

	err = cmdr.CertInfo("test-example-com")
	assert.NoErr(t, err)
	assert.Equal(t, b.String(), `=== test-example-com Certificate
Common Name(s):     test.deis.com
Expires At:         9 Jun 2016
Starts At:          9 Jun 2016
Fingerprint:        ab:12:ab:12:ab
Subject Alt Name:   test.com,example.com
Issuer:             testca
Subject:            testing

Connected Domains:  test.com,example.com
Owner:              admin
Created:            9 Jun 2016
Updated:            9 Jun 2016
`, "output")

	server.Mux.HandleFunc("/v2/certs/test-deis-com", func(w http.ResponseWriter, r *http.Request) {
		testutil.SetHeaders(w)
		fmt.Fprintf(w, `{
			"name": "test-deis-com"
		}`)
	})
	b.Reset()

	err = cmdr.CertInfo("test-deis-com")
	assert.NoErr(t, err)
	assert.Equal(t, b.String(), `=== test-deis-com Certificate
Common Name(s):     
Expires At:         unknown
Starts At:          unknown
Fingerprint:        
Subject Alt Name:   N/A
Issuer:             
Subject:            

Connected Domains:  No connected domains
Owner:              
Created:            unknown
Updated:            unknown
`, "output")
}

func TestCertsRemove(t *testing.T) {
	t.Parallel()
	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()
	var b bytes.Buffer
	cmdr := DeisCmd{WOut: &b, ConfigFile: cf}

	server.Mux.HandleFunc("/v2/certs/test-example-com", func(w http.ResponseWriter, r *http.Request) {
		testutil.SetHeaders(w)
		w.WriteHeader(http.StatusNoContent)
	})

	err = cmdr.CertRemove("test-example-com")
	assert.NoErr(t, err)

	assert.Equal(t, testutil.StripProgress(b.String()), "Removing test-example-com... done\n", "output")
}

func TestCertsAttach(t *testing.T) {
	t.Parallel()
	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()
	var b bytes.Buffer
	cmdr := DeisCmd{WOut: &b, ConfigFile: cf}

	server.Mux.HandleFunc("/v2/certs/test-example-com/domain/", func(w http.ResponseWriter, r *http.Request) {
		testutil.SetHeaders(w)
		testutil.AssertBody(t, api.CertAttachRequest{Domain: "deis.com"}, r)
		w.WriteHeader(http.StatusCreated)
	})

	err = cmdr.CertAttach("test-example-com", "deis.com")
	assert.NoErr(t, err)

	assert.Equal(t, testutil.StripProgress(b.String()), "Attaching certificate test-example-com to domain deis.com... done\n", "output")
}

func TestCertsDetach(t *testing.T) {
	t.Parallel()
	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()
	var b bytes.Buffer
	cmdr := DeisCmd{WOut: &b, ConfigFile: cf}

	server.Mux.HandleFunc("/v2/certs/test-example-com/domain/deis.com", func(w http.ResponseWriter, r *http.Request) {
		testutil.SetHeaders(w)
		w.WriteHeader(http.StatusNoContent)
	})

	err = cmdr.CertDetach("test-example-com", "deis.com")
	assert.NoErr(t, err)

	assert.Equal(t, testutil.StripProgress(b.String()), "Detaching certificate test-example-com from domain deis.com... done\n", "output")
}

func TestCertsAdd(t *testing.T) {
	t.Parallel()
	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()
	var b bytes.Buffer
	cmdr := DeisCmd{WOut: &b, ConfigFile: cf}

	server.Mux.HandleFunc("/v2/certs/", func(w http.ResponseWriter, r *http.Request) {
		testutil.SetHeaders(w)
		testutil.AssertBody(t, api.CertCreateRequest{Certificate: "cert", Key: "key", Name: "testcert"}, r)
		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, "{}")
	})

	keyFile, err := ioutil.TempFile("", "deis-cli-unit-test-key")
	assert.NoErr(t, err)
	_, err = keyFile.Write([]byte("key"))
	assert.NoErr(t, err)
	keyFile.Close()

	certFile, err := ioutil.TempFile("", "deis-cli-unit-test-cert")
	assert.NoErr(t, err)
	_, err = certFile.Write([]byte("cert"))
	assert.NoErr(t, err)
	certFile.Close()

	err = cmdr.CertAdd(certFile.Name(), keyFile.Name(), "testcert")
	assert.NoErr(t, err)

	assert.Equal(t, testutil.StripProgress(b.String()), "Adding SSL endpoint... done\n", "output")
}
