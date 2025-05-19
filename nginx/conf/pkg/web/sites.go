package web

import (
	"net/http"
	"om/pkg/db"
	"om/pkg/ngx"
	"om/pkg/util"
	"os"
	"strconv"

	"github.com/labstack/echo/v4"
)

func getSites(c echo.Context) error {
	var site db.Site

	sites, err := site.GetAll()
	if err != nil {
		return c.JSON(http.StatusOK, echo.Map{
			"error": err.Error(),
		})
	}

	certOptions, err := getCertOptions()
	if err != nil {
		return c.JSON(http.StatusOK, echo.Map{
			"error": err.Error(),
		})
	}

	upstreamOptions, err := getUpstreamOptions()
	if err != nil {
		return c.JSON(http.StatusOK, echo.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"sites":            sites,
		"cert_options":     certOptions,
		"upstream_options": upstreamOptions,
	})
}

func addSite(c echo.Context) error {
	var site db.Site

	err := c.Bind(&site)
	if err != nil {
		return c.JSON(http.StatusOK, echo.Map{
			"error": err.Error(),
		})
	}

	err = site.Insert()
	if err != nil {
		return c.JSON(http.StatusOK, echo.Map{
			"error": err.Error(),
		})
	}

	err = ngx.SaveSite(&site)
	if err != nil {
		return c.JSON(http.StatusOK, echo.Map{
			"error": err.Error(),
		})
	}

	err = ngx.ReloadOpenresty()
	if err != nil {
		os.Remove(util.NginxDir + "conf/sites/" + strconv.FormatUint(uint64(site.ID), 10) + ".conf")
		site.Delete()
		return c.JSON(http.StatusOK, echo.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, "OK")
}

func setSite(c echo.Context) error {
	var site db.Site

	err := c.Bind(&site)
	if err != nil {
		return c.JSON(http.StatusOK, echo.Map{
			"error": err.Error(),
		})
	}

	oldSite := db.Site{}
	err = oldSite.Get(site.ID)
	if err != nil {
		return c.JSON(http.StatusOK, echo.Map{
			"error": err.Error(),
		})
	}

	err = ngx.SaveSite(&site)
	if err != nil {
		return c.JSON(http.StatusOK, echo.Map{
			"error": err.Error(),
		})
	}

	err = ngx.ReloadOpenresty()
	if err != nil {
		ngx.SaveSite(&oldSite)
		return c.JSON(http.StatusOK, echo.Map{
			"error": err.Error(),
		})
	}

	err = site.Update()
	if err != nil {
		return c.JSON(http.StatusOK, echo.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, "OK")
}

func delSites(c echo.Context) error {
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
	err = site.DeleteAll(req.Keys)
	if err != nil {
		return c.JSON(http.StatusOK, echo.Map{
			"error": err.Error(),
		})
	}

	for _, id := range req.Keys {
		os.Remove(util.NginxDir + "conf/sites/" + strconv.FormatUint(uint64(id), 10) + ".conf")
	}

	err = ngx.ReloadOpenresty()
	if err != nil {
		return c.JSON(http.StatusOK, echo.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, "OK")
}
