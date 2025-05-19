package web

import (
	"fmt"
	"net/http"
	"om/pkg/acme"
	"om/pkg/db"
	"om/pkg/ngx"
	"om/pkg/util"
	"os"
	"regexp"
	"strconv"

	"github.com/tidwall/gjson"

	"github.com/labstack/echo/v4"
)

var drgx = regexp.MustCompile(`^(\*\.)?[0-9a-z\-\.]{2,}$`)

type CertOption struct {
	Label string `json:"label"`
	Value uint   `json:"value"`
}

func getCertOptions() ([]CertOption, error) {
	var cert db.Cert

	rows, err := cert.GetFields([]string{"id", "name"})
	if err != nil {
		return nil, err
	}

	certOptions := []CertOption{}
	for _, row := range rows {
		certOptions = append(certOptions, CertOption{Label: row.Name, Value: row.ID})
	}

	return certOptions, nil
}

func getCerts(c echo.Context) error {
	var cert db.Cert

	certs, err := cert.GetAll()
	if err != nil {
		return c.JSON(http.StatusOK, echo.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, certs)
}

func ApplyCert(cert *db.Cert) error {
	domains := []string{}
	jds := gjson.Parse(cert.Domains).Array()
	for _, jd := range jds {
		jdStr := jd.String()
		if !drgx.MatchString(jdStr) {
			return fmt.Errorf("illegal domain: %s", jdStr)
		}
		domains = append(domains, jdStr)
	}

	if cert.Type == 0 {
		if cert.DnsChallenge {
			crt, key, err := acme.DnsObtain(cert.Email, cert.DnsProvider, cert.DnsCredential, domains)
			if err != nil {
				return err
			}
			cert.Crt = string(crt)
			cert.Key = string(key)
		} else {
			crt, key, err := acme.HttpObtain(cert.Email, domains)
			if err != nil {
				return err
			}
			cert.Crt = string(crt)
			cert.Key = string(key)
		}
	}

	dnsNames, notAfter, err := util.GetCertInfo([]byte(cert.Crt))
	if err != nil {
		return err
	}
	if cert.Type == 1 {
		jsonData, err := util.Json.Marshal(dnsNames)
		if err != nil {
			return err
		}
		cert.Domains = string(jsonData)
	}
	cert.Expires = *notAfter

	return nil
}

func addCert(c echo.Context) error {
	var cert db.Cert

	err := c.Bind(&cert)
	if err != nil {
		return c.JSON(http.StatusOK, echo.Map{
			"error": err.Error(),
		})
	}

	err = ApplyCert(&cert)
	if err != nil {
		return c.JSON(http.StatusOK, echo.Map{
			"error": err.Error(),
		})
	}

	err = cert.Insert()
	if err != nil {
		return c.JSON(http.StatusOK, echo.Map{
			"error": err.Error(),
		})
	}

	err = saveCert(&cert)
	if err != nil {
		return c.JSON(http.StatusOK, echo.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, "OK")
}

func setCert(c echo.Context) error {
	var cert db.Cert

	err := c.Bind(&cert)
	if err != nil {
		return c.JSON(http.StatusOK, echo.Map{
			"error": err.Error(),
		})
	}

	err = ApplyCert(&cert)
	if err != nil {
		return c.JSON(http.StatusOK, echo.Map{
			"error": err.Error(),
		})
	}

	err = cert.Update()
	if err != nil {
		return c.JSON(http.StatusOK, echo.Map{
			"error": err.Error(),
		})
	}

	err = saveCert(&cert)
	if err != nil {
		return c.JSON(http.StatusOK, echo.Map{
			"error": err.Error(),
		})
	}

	err = ngx.ReloadOpenresty()
	if err != nil {
		return c.JSON(http.StatusOK, echo.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, "OK")
}

func saveCert(cert *db.Cert) error {
	err := os.WriteFile(util.NginxDir+"conf/certs/"+strconv.FormatUint(uint64(cert.ID), 10)+".crt", []byte(cert.Crt), 0644)
	if err != nil {
		return err
	}
	return os.WriteFile(util.NginxDir+"conf/certs/"+strconv.FormatUint(uint64(cert.ID), 10)+".key", []byte(cert.Key), 0644)
}

func delCerts(c echo.Context) error {
	var req struct {
		Keys []uint `json:"keys"`
	}
	err := c.Bind(&req)
	if err != nil {
		return c.JSON(http.StatusOK, echo.Map{
			"error": err.Error(),
		})
	}

	site := db.Site{}
	count, err := site.CertCount(req.Keys)
	if err != nil {
		return c.JSON(http.StatusOK, echo.Map{
			"error": err.Error(),
		})
	}
	if count > 0 {
		msg := "There is one site using this certificate"
		if count > 1 {
			msg = fmt.Sprintf("There are %d sites using this certificate", count)
		}
		return c.JSON(http.StatusOK, echo.Map{
			"error": msg,
		})
	}

	cert := db.Cert{}
	err = cert.DeleteAll(req.Keys)
	if err != nil {
		return c.JSON(http.StatusOK, echo.Map{
			"error": err.Error(),
		})
	}

	for _, id := range req.Keys {
		os.Remove(util.NginxDir + "conf/certs/" + strconv.FormatUint(uint64(id), 10) + ".crt")
		os.Remove(util.NginxDir + "conf/certs/" + strconv.FormatUint(uint64(id), 10) + ".key")
	}

	return c.JSON(http.StatusOK, "OK")
}
