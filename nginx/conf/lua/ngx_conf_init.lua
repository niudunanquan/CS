local cjson = require("cjson.safe")

collectgarbage("collect")
local f = io.open("/opt/om/nginx/lua/cfg.json", "r")
local content = f:read("*a")
f:close()
local cfg = cjson.decode(content)

-- init modules
local ok, res
ok, res = pcall(require, "om")
if not ok then
  error("require failed: " .. tostring(res))
else
  res.set_config(cfg)
end
