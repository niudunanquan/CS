package task

import (
	"om/pkg/cron"
	"om/pkg/db"
	"om/pkg/util"
	"om/pkg/web"
	"os/exec"
	"time"

	log "github.com/sirupsen/logrus"
)

var Task *cron.Cron

func InitTask() error {
	c := cron.New()
	c.SetTimezone(time.Local)
	err := c.Add("certs", "0 1 * * *", checkCerts)
	if err != nil {
		return err
	}

	err = c.Add("certs", "0 0 * * *", logrotate)
	if err != nil {
		return err
	}
	c.Start()
	Task = c
	return nil
}

func Stop() {
	if Task != nil {
		Task.Stop()
	}
}

func checkCerts() {
	var cert db.Cert

	certs, err := cert.GetAll()
	if err != nil {
		log.Errorf("failed to get certs: %s", err.Error())
		return
	}

	for _, cert := range certs {
		if cert.Type == 0 && cert.Expires.Sub(time.Now().AddDate(0, 1, 0)) < 0 {
			err = web.ApplyCert(&cert)
			if err != nil {
				log.Errorf("failed to apply cert %d: %s", cert.ID, err.Error())
			}
		}
	}
}

func logrotate() {
	out, err := exec.Command("/usr/bin/logrotate", "-f", util.NginxDir+"logrotate.d/nginx").CombinedOutput()
	if err != nil {
		if out != nil {
			log.Errorf("failed to logrotate nginx: %s", string(out))
			return
		}
		log.Errorf("failed to logrotate nginx: %s", err.Error())
	}
}
