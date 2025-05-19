package ngx

import (
	"errors"
	"om/pkg/db"
	"om/pkg/util"
	"os"
	"strconv"
)

func SaveSite(site *db.Site) error {
	t, err := NewTemplate("srv.tpl", srv_tpl)
	if err != nil {
		return err
	}

	conf, err := t.Parse(site)
	if err != nil {
		return err
	}

	return os.WriteFile(util.NginxDir+"conf/sites/"+strconv.FormatUint(uint64(site.ID), 10)+".conf", conf, 0644)
}

func SaveUpstream(ups *db.Upstream) error {
	t, err := NewTemplate("ups.tpl", ups_tpl)
	if err != nil {
		return err
	}

	conf, err := t.Parse(ups)
	if err != nil {
		return err
	}

	return os.WriteFile(util.NginxDir+"conf/upstreams/"+strconv.FormatUint(uint64(ups.ID), 10)+".conf", conf, 0644)
}

func ReloadOpenresty() error {
	cmd := NewNginxCommand()
	out, err := cmd.Test()
	if err != nil {
		if out != nil {
			return errors.New(string(out))
		}
		return err
	}
	cmd.Reload()

	return nil
}
