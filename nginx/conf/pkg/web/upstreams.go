package web

import (
	"net/http"
	"om/pkg/db"
	"om/pkg/ngx"
	"om/pkg/util"
	"os"
	"slices"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/tidwall/gjson"
)

type UpstreamOption struct {
	Label string `json:"label"`
	Value uint   `json:"value"`
}

func getUpstreamOptions() ([]UpstreamOption, error) {
	var upstream db.Upstream

	rows, err := upstream.GetFields([]string{"id", "name"})
	if err != nil {
		return nil, err
	}

	upstreamOptions := []UpstreamOption{}
	for _, row := range rows {
		upstreamOptions = append(upstreamOptions, UpstreamOption{Label: row.Name, Value: row.ID})
	}

	return upstreamOptions, nil
}

func getUpstreams(c echo.Context) error {
	var upstream db.Upstream

	upstreams, err := upstream.GetAll()
	if err != nil {
		return c.JSON(http.StatusOK, echo.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, upstreams)
}

func addUpstream(c echo.Context) error {
	var upstream db.Upstream

	err := c.Bind(&upstream)
	if err != nil {
		return c.JSON(http.StatusOK, echo.Map{
			"error": err.Error(),
		})
	}

	if strings.Contains(upstream.Config, "}") {
		return c.JSON(http.StatusOK, echo.Map{
			"error": "illegal config: " + upstream.Config,
		})
	}

	err = upstream.Insert()
	if err != nil {
		return c.JSON(http.StatusOK, echo.Map{
			"error": err.Error(),
		})
	}

	err = ngx.SaveUpstream(&upstream)
	if err != nil {
		return c.JSON(http.StatusOK, echo.Map{
			"error": err.Error(),
		})
	}

	err = ngx.ReloadOpenresty()
	if err != nil {
		os.Remove(util.NginxDir + "conf/upstreams/" + strconv.FormatUint(uint64(upstream.ID), 10) + ".conf")
		upstream.Delete()
		return c.JSON(http.StatusOK, echo.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, "OK")
}

func setUpstream(c echo.Context) error {
	var upstream db.Upstream

	err := c.Bind(&upstream)
	if err != nil {
		return c.JSON(http.StatusOK, echo.Map{
			"error": err.Error(),
		})
	}

	if strings.Contains(upstream.Config, "}") {
		return c.JSON(http.StatusOK, echo.Map{
			"error": "illegal config: " + upstream.Config,
		})
	}

	oldUpstream := db.Upstream{}
	err = oldUpstream.Get(upstream.ID)
	if err != nil {
		return c.JSON(http.StatusOK, echo.Map{
			"error": err.Error(),
		})
	}

	err = ngx.SaveUpstream(&upstream)
	if err != nil {
		return c.JSON(http.StatusOK, echo.Map{
			"error": err.Error(),
		})
	}

	err = ngx.ReloadOpenresty()
	if err != nil {
		ngx.SaveUpstream(&oldUpstream)
		return c.JSON(http.StatusOK, echo.Map{
			"error": err.Error(),
		})
	}

	err = upstream.Update()
	if err != nil {
		return c.JSON(http.StatusOK, echo.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, "OK")
}

func delUpstreams(c echo.Context) error {
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
	sites, err := site.GetLocations()
	if err != nil {
		return c.JSON(http.StatusOK, echo.Map{
			"error": err.Error(),
		})
	}

	for _, row := range sites {
		result := gjson.Get(row.Locations, "#.upstream_id")
		for _, value := range result.Array() {
			if slices.Contains(req.Keys, uint(value.Uint())) {
				return c.JSON(http.StatusOK, echo.Map{
					"error": "There are sites using these upstreams",
				})
			}
		}
	}

	upstream := db.Upstream{}
	err = upstream.DeleteAll(req.Keys)
	if err != nil {
		return c.JSON(http.StatusOK, echo.Map{
			"error": err.Error(),
		})
	}

	for _, id := range req.Keys {
		os.Remove(util.NginxDir + "conf/upstreams/" + strconv.FormatUint(uint64(id), 10) + ".conf")
	}

	return c.JSON(http.StatusOK, "OK")
}
