local ngx = ngx
local ngx_log = ngx.log
local ngx_err = ngx.ERR
local ngx_var = ngx.var
local echo = ngx.print
local str_find = string.find
local str_format = string.format
local str_sub = string.sub
local util = require("util")
local cjson = require("cjson.safe")

local _M = {}

local config, geo_ip

local function init_mmdb()
    local err
    local mmdb = require("mmdb")
    geo_ip, err = mmdb.load_database("/opt/om/nginx/geoip.mmdb")
    if not geo_ip then
        ngx_log(ngx_err, err)
    end
end

function _M.set_config(new_config)
    config = new_config
end

function _M.init_worker()
    init_mmdb()
end

function _M.stats()
    local db = ngx.shared.kvdb
    local t = ngx_var.arg_t

    if t and t == "s" then
        local keys = db:get_keys()
        local req = {}
        local geo = {}
        for _, key in ipairs(keys) do
            if str_find(key, "req:", 1, true) == 1 then
                req[str_sub(key, 5)] = db:get(key)
            elseif str_find(key, "geo:", 1, true) == 1 then
                geo[str_sub(key, 5)] = db:get(key)
            end
        end
        local data = cjson.encode({ req = req, geo = geo })
        echo(data)
    else
        local c = db:get("req")
        c = c or 0
        echo(tostring(c))
    end
end

function _M.log()
    if ngx_var.remote_addr == "127.0.0.1" then
        return
    end
    
    local db = ngx.shared.kvdb
    db:incr("req", 1, 0)
    local m, err = ngx.re.match(ngx.localtime(), [[^\d{4}-\d{2}-\d{2}\s+(\d{2}):\d{2}:\d{2}$]], "jo")
    if err then
        ngx_log(ngx_err, str_format("could not match localtime: %s", err))
    end
    db:incr("req:" .. m[1], 1, 0, 43200)

    local iso_code = util.ip2loc(geo_ip, ngx_var.remote_addr)
    if iso_code and iso_code ~= "" then
        db:incr("geo:" .. iso_code, 1, 0)
    end
end

return _M
