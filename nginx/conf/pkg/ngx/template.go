package ngx

import (
	"om/pkg/db"
	"strconv"
	"strings"
	"text/template"

	"github.com/tidwall/gjson"
)

const (
	defBufferSize = 65535
	srv_tpl       = `{{ $site := . }}
server {
{{ buildListeners $site.Listeners $site.Ipv6 }}
{{ if $site.Http2 }}
    http2  on;
{{ end }}
    server_name {{ buildDomains $site.Domains }};
{{ if (gt $site.CertId 0) }}
    ssl_certificate      certs/{{ itoa $site.CertId }}.crt;
    ssl_certificate_key  certs/{{ itoa $site.CertId }}.key;
{{ end }}
    include log.conf;
    include acme_challenge.conf;
{{ if $site.ForceSsl }}
    include force_ssl.conf;
{{ end }}
{{buildLocations $site }}
}
`
	ups_tpl = `{{ $ups := . }}
upstream {{ itoa $ups.ID }} {
{{ $ups.Config }}
}
`
)

type Template struct {
	tmpl *template.Template

	bp *BufferPool
}

func NewTemplate(name, data string) (*Template, error) {
	tmpl, err := template.New(name).Funcs(funcMap).Parse(data)
	if err != nil {
		return nil, err
	}

	return &Template{
		tmpl: tmpl,
		bp:   NewBufferPool(defBufferSize),
	}, nil
}

func (t *Template) Parse(conf any) ([]byte, error) {
	tmplBuf := t.bp.Get()
	defer t.bp.Put(tmplBuf)

	err := t.tmpl.Execute(tmplBuf, conf)
	if err != nil {
		return nil, err
	}

	// make a copy to ensure that we are no longer modifying the content of the buffer
	out := tmplBuf.Bytes()
	res := make([]byte, len(out))
	copy(res, out)

	return res, nil
}

var funcMap = template.FuncMap{
	"empty": func(input interface{}) bool {
		check, ok := input.(string)
		if ok {
			return check == ""
		}
		return true
	},
	"buildListeners": buildListeners,
	"buildDomains":   buildDomains,
	"buildLocations": buildLocations,
	"itoa":           itoa,
}

func buildListeners(input string, ipv6 bool) string {
	listeners := gjson.Parse(input).Array()
	sb := strings.Builder{}
	for _, listener := range listeners {
		port := listener.Get("port").Uint()
		if port != 0 {
			ssl := listener.Get("ssl").Bool()
			if ipv6 {
				sb.WriteString("    listen [::]:" + strconv.FormatUint(uint64(port), 10))
			} else {
				sb.WriteString("    listen " + strconv.FormatUint(uint64(port), 10))
			}

			if ssl {
				sb.WriteString(" ssl")
			}
			sb.WriteString(";\n")
		}
	}
	return sb.String()
}

func buildDomains(input string) string {
	domains := gjson.Parse(input).Array()
	sb := strings.Builder{}
	for _, domain := range domains {
		if domain.String() != "" {
			sb.WriteString(" '" + strings.ReplaceAll(domain.String(), "'", "\\'") + "'")
		}
	}
	return sb.String()
}

func buildLocations(site *db.Site) string {
	locations := gjson.Parse(site.Locations).Array()
	sb := strings.Builder{}
	for _, location := range locations {
		path := location.Get("path").String()
		if path != "" && !strings.ContainsAny(path, "{;\n") {
			sb.WriteString("    location " + path + " {\n")
			sb.WriteString("        set $proxy_scheme  '" + strings.ReplaceAll(location.Get("protocol").String(), "'", "\\'") + "';\n")
			sb.WriteString("        set $proxy_backend '" + strings.ReplaceAll(location.Get("upstream_id").String(), "'", "\\'") + "';\n")
			sb.WriteString("        set $proxy_uri     '" + strings.ReplaceAll(location.Get("upstream_path").String(), "'", "\\'") + "';\n")
			sb.WriteString("        include websocket.conf;\n")
			sb.WriteString("        include proxy_ip.conf;\n")
			if !site.Gzip {
				sb.WriteString("        gzip    off;\n")
			}
			if site.Cache {
				sb.WriteString("        include cache.conf;\n")
			}
			if site.Hsts {
				sb.WriteString("        include hsts.conf;\n")
			}
			sb.WriteString("        include proxy.conf;\n")
			sb.WriteString("        include error_page.conf;\n")
			sb.WriteString("    }\n")
		}
	}
	return sb.String()
}

func itoa(input uint) string {
	return strconv.FormatUint(uint64(input), 10)
}
