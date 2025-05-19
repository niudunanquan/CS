package acme

import (
	"errors"
	"om/pkg/util"
	"os"
	"strings"

	"golang.org/x/net/idna"
)

func sanitizedDomain(domain string) string {
	safe, err := idna.ToASCII(strings.NewReplacer(":", "-", "*", "_").Replace(domain))
	if err != nil {
		return ""
	}
	return safe
}

func HttpObtain(email string, domains []string) ([]byte, []byte, error) {
	lc := NewLegoCommand()
	args := []string{}
	args = append(args, "--email", email, "--http", "--http.webroot", util.NginxDir+"html")
	for _, domain := range domains {
		args = append(args, "--domains", domain)
	}
	args = append(args, "run")

	out, err := lc.ExecCommand(args...).CombinedOutput()
	if err != nil {
		if out != nil {
			return nil, nil, errors.New(string(out))
		}
		return nil, nil, err
	}

	path := util.RootDir + "acme/certificates/" + sanitizedDomain(domains[0])
	crt, err := os.ReadFile(path + ".crt")
	if err != nil {
		return nil, nil, err
	}
	key, err := os.ReadFile(path + ".key")
	if err != nil {
		return nil, nil, err
	}

	return crt, key, nil
}

func DnsObtain(email, dnsProvider, dnsCredential string, domains []string) ([]byte, []byte, error) {
	lc := NewLegoCommand()
	args := []string{}
	args = append(args, "--email", email, "--dns", dnsProvider)
	for _, domain := range domains {
		args = append(args, "--domains", domain)
	}
	args = append(args, "run")
	cmd := lc.ExecCommand(args...)

	dnsEnv := strings.Split(dnsCredential, "\n")
	for i, s := range dnsEnv {
		dnsEnv[i] = strings.TrimSpace(s)
	}
	cmd.Env = append(os.Environ(), dnsEnv...)

	out, err := cmd.CombinedOutput()
	if err != nil {
		if out != nil {
			return nil, nil, errors.New(string(out))
		}
		return nil, nil, err
	}

	path := util.RootDir + "acme/certificates/" + sanitizedDomain(domains[0])
	crt, err := os.ReadFile(path + ".crt")
	if err != nil {
		return nil, nil, err
	}
	key, err := os.ReadFile(path + ".key")
	if err != nil {
		return nil, nil, err
	}

	return crt, key, nil
}
