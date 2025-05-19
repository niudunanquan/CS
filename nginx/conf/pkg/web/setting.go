package web

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"om/pkg"
	"om/pkg/ngx"
	"om/pkg/selfupdate"
	"om/pkg/util"

	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
)

type OMConfig struct {
	Resolver  string `json:"resolver"`
	RealIp    string `json:"real_ip"`
	Ssl       string `json:"ssl"`
	Gzip      string `json:"gzip"`
	Cache     string `json:"cache"`
	Hsts      string `json:"hsts"`
	ForceSsl  string `json:"force_ssl"`
	Proxy     string `json:"proxy"`
	ProxyIp   string `json:"proxy_ip"`
	Websocket string `json:"websocket"`
	ErrorPage string `json:"error_page"`
	Log       string `json:"log"`
}

func getOmConfig(c echo.Context) error {
	if !isAdmin(c) {
		return c.JSON(http.StatusOK, echo.Map{
			"error": "Not allowed",
		})
	}

	omc, err := readConfig()
	if err != nil {
		return c.JSON(http.StatusOK, echo.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, omc)
}

func setOmConfig(c echo.Context) error {
	if !isAdmin(c) {
		return c.JSON(http.StatusOK, echo.Map{
			"error": "Not allowed",
		})
	}

	var omc OMConfig

	err := c.Bind(&omc)
	if err != nil {
		return c.JSON(http.StatusOK, echo.Map{
			"error": err.Error(),
		})
	}

	oldOmc, err := readConfig()
	if err != nil {
		return c.JSON(http.StatusOK, echo.Map{
			"error": err.Error(),
		})
	}

	err = saveConfig(&omc)
	if err != nil {
		return c.JSON(http.StatusOK, echo.Map{
			"error": err.Error(),
		})
	}

	err = ngx.ReloadOpenresty()
	if err != nil {
		saveConfig(oldOmc)
		return c.JSON(http.StatusOK, echo.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, omc)
}

func readConfig() (*OMConfig, error) {
	omc := OMConfig{}

	data, err := os.ReadFile(util.NginxDir + "conf/resolver.conf")
	if err != nil {
		return nil, err
	}
	omc.Resolver = string(data)

	data, err = os.ReadFile(util.NginxDir + "conf/real_ip.conf")
	if err != nil {
		return nil, err
	}
	omc.RealIp = string(data)

	data, err = os.ReadFile(util.NginxDir + "conf/ssl.conf")
	if err != nil {
		return nil, err
	}
	omc.Ssl = string(data)

	data, err = os.ReadFile(util.NginxDir + "conf/gzip.conf")
	if err != nil {
		return nil, err
	}
	omc.Gzip = string(data)

	data, err = os.ReadFile(util.NginxDir + "conf/cache.conf")
	if err != nil {
		return nil, err
	}
	omc.Cache = string(data)

	data, err = os.ReadFile(util.NginxDir + "conf/hsts.conf")
	if err != nil {
		return nil, err
	}
	omc.Hsts = string(data)

	data, err = os.ReadFile(util.NginxDir + "conf/force_ssl.conf")
	if err != nil {
		return nil, err
	}
	omc.ForceSsl = string(data)

	data, err = os.ReadFile(util.NginxDir + "conf/proxy.conf")
	if err != nil {
		return nil, err
	}
	omc.Proxy = string(data)

	data, err = os.ReadFile(util.NginxDir + "conf/proxy_ip.conf")
	if err != nil {
		return nil, err
	}
	omc.ProxyIp = string(data)

	data, err = os.ReadFile(util.NginxDir + "conf/websocket.conf")
	if err != nil {
		return nil, err
	}
	omc.Websocket = string(data)

	data, err = os.ReadFile(util.NginxDir + "conf/error_page.conf")
	if err != nil {
		return nil, err
	}
	omc.ErrorPage = string(data)

	data, err = os.ReadFile(util.NginxDir + "conf/log.conf")
	if err != nil {
		return nil, err
	}
	omc.Log = string(data)

	return &omc, nil
}

func saveConfig(omc *OMConfig) error {

	if err := os.WriteFile(util.NginxDir+"conf/resolver.conf", []byte(omc.Resolver), 0644); err != nil {
		return err
	}

	if err := os.WriteFile(util.NginxDir+"conf/real_ip.conf", []byte(omc.RealIp), 0644); err != nil {
		return err
	}

	if err := os.WriteFile(util.NginxDir+"conf/ssl.conf", []byte(omc.Ssl), 0644); err != nil {
		return err
	}

	if err := os.WriteFile(util.NginxDir+"conf/gzip.conf", []byte(omc.Gzip), 0644); err != nil {
		return err
	}

	if err := os.WriteFile(util.NginxDir+"conf/cache.conf", []byte(omc.Cache), 0644); err != nil {
		return err
	}

	if err := os.WriteFile(util.NginxDir+"conf/hsts.conf", []byte(omc.Hsts), 0644); err != nil {
		return err
	}

	if err := os.WriteFile(util.NginxDir+"conf/force_ssl.conf", []byte(omc.ForceSsl), 0644); err != nil {
		return err
	}

	if err := os.WriteFile(util.NginxDir+"conf/proxy.conf", []byte(omc.Proxy), 0644); err != nil {
		return err
	}

	if err := os.WriteFile(util.NginxDir+"conf/proxy_ip.conf", []byte(omc.ProxyIp), 0644); err != nil {
		return err
	}

	if err := os.WriteFile(util.NginxDir+"conf/websocket.conf", []byte(omc.Websocket), 0644); err != nil {
		return err
	}

	if err := os.WriteFile(util.NginxDir+"conf/error_page.conf", []byte(omc.ErrorPage), 0644); err != nil {
		return err
	}

	if err := os.WriteFile(util.NginxDir+"conf/log.conf", []byte(omc.Log), 0644); err != nil {
		return err
	}

	return nil
}

func getVersion(c echo.Context) error {
	url := "https://om.uusec.com/version"
	resp, err := http.Get(url)
	if err != nil {
		return c.JSON(http.StatusOK, echo.Map{
			"error": err.Error(),
		})
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return c.JSON(http.StatusOK, echo.Map{
			"error": fmt.Sprintf("bad http status from %s: %v", url, resp.Status),
		})
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return c.JSON(http.StatusOK, echo.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"latest_version": strings.TrimSpace(string(data)),
		"version":        pkg.Version,
	})
}

func update(c echo.Context) error {
	if !isAdmin(c) {
		return c.JSON(http.StatusOK, echo.Map{
			"error": "Not allowed",
		})
	}

	if _, err := os.Stat("/.dockerenv"); err == nil {
		return c.JSON(http.StatusOK, echo.Map{
			"error": "You can update the docker version by om.sh",
		})
	}

	updater := &selfupdate.Updater{
		CurrentVersion: pkg.Version,
		ApiURL:         "https://om.uusec.com/download/", //$CmdName/$GOOS-$ARCH.json { "Version": "2", "Sha256": "..." }
		BinURL:         "https://om.uusec.com/download/", //$CmdName/$NewVersion/$GOOS-$ARCH.gz
		DiffURL:        "",                               //$CmdName/$CurrentVersion/$NewVersion/$GOOS-$ARCH
		Dir:            "update/",
		CmdName:        "oms",
		ForceCheck:     true,
	}

	updater.OnSuccessfulUpdate = func() {
		out, err := exec.Command(util.RootDir+"oms", "-s", "restart").CombinedOutput()
		if err != nil {
			if out != nil {
				log.Errorf("failed to restart OpenResty Manager: %s", string(out))
				return
			}
			log.Errorf("failed to restart OpenResty Manager: %s", err.Error())
		}
	}

	err := updater.BackgroundRun()
	if err != nil {
		return c.JSON(http.StatusOK, echo.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, "OK")
}
