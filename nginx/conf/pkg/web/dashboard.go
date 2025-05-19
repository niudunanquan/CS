package web

import (
	"net/http"
	"om/pkg/db"
	"om/pkg/util"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/mem"
	"github.com/tidwall/gjson"
)

type Stats struct {
	Sites     int64  `json:"sites"`
	Certs     int64  `json:"certs"`
	Upstreams int64  `json:"upstreams"`
	Users     int64  `json:"users"`
	Req       string `json:"req"`
	Geo       string `json:"geo"`
}

type Rts struct {
	Cpu    float64 `json:"cpu"`
	Memory float64 `json:"mem"`
	Disk   float64 `json:"disk"`
	Req    uint64  `json:"req"`
}

func getStats(c echo.Context) error {
	var stats Stats
	var err error

	site := db.Site{}
	stats.Sites, err = site.Count()
	if err != nil {
		return c.JSON(http.StatusOK, echo.Map{
			"error": err.Error(),
		})
	}

	cert := db.Cert{}
	stats.Certs, err = cert.Count()
	if err != nil {
		return c.JSON(http.StatusOK, echo.Map{
			"error": err.Error(),
		})
	}

	upstream := db.Upstream{}
	stats.Upstreams, err = upstream.Count()
	if err != nil {
		return c.JSON(http.StatusOK, echo.Map{
			"error": err.Error(),
		})
	}

	user := db.User{}
	stats.Users, err = user.Count()
	if err != nil {
		return c.JSON(http.StatusOK, echo.Map{
			"error": err.Error(),
		})
	}

	data, err := util.HttpGet("http://127.0.0.1/om_stats?t=s")
	if err != nil {
		return c.JSON(http.StatusOK, echo.Map{
			"error": err.Error(),
		})
	}
	stats.Req = gjson.GetBytes(data, "req").Raw
	stats.Geo = gjson.GetBytes(data, "geo").Raw

	return c.JSON(http.StatusOK, stats)
}

func getRts(c echo.Context) error {
	var rts Rts
	var err error

	data, err := util.HttpGet("http://127.0.0.1/om_stats")
	if err != nil {
		return c.JSON(http.StatusOK, echo.Map{
			"error": err.Error(),
		})
	}

	rts.Req, err = strconv.ParseUint(string(data), 10, 64)
	if err != nil {
		return c.JSON(http.StatusOK, echo.Map{
			"error": err.Error(),
		})
	}

	cpuPercent, err := cpu.Percent(time.Second, false)
	if err != nil {
		return c.JSON(http.StatusOK, echo.Map{
			"error": err.Error(),
		})
	}
	rts.Cpu = cpuPercent[0]

	memInfo, err := mem.VirtualMemory()
	if err != nil {
		return c.JSON(http.StatusOK, echo.Map{
			"error": err.Error(),
		})
	}
	rts.Memory = memInfo.UsedPercent

	diskUsage, err := disk.Usage("/")
	if err != nil {
		return c.JSON(http.StatusOK, echo.Map{
			"error": err.Error(),
		})
	}
	rts.Disk = diskUsage.UsedPercent

	return c.JSON(http.StatusOK, rts)
}
