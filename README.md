# Block File Downloading 
Automatic tool for generating a lua script file. This lua script file will be imported into nginx configuration file to block all requests that having specific pairs of Content-Type and Content-Disposition.

This tool was programed by golang. It take a json file as an input parameter. The output is the new lua script file, which placed at the nginx configuration 's folder. 
## 1. Prepare Input file (*blacklist.json*):
```
{
  "blacklist": [
    {
      "content_disposition": "b",
      "content_type": "a"
    },
    {
      "content_disposition": "*",
      "content_type": "a"
    },
    {
      "content_disposition": "b",
      "content_type": "*"
    },
    {
      "content_disposition": "ee",
      "content_type": "abcd"
    }
  ]
}
```

## 2. Config virtualhost 
add `rewrite_by_lua_file /etc/nginx/block_download_file.lua;`  into server directive.
```
server {
    listen       80;
    server_name  localhost;
    rewrite_by_lua_file /etc/nginx/block_download_file.lua; 
```

## 3. Run
```
go run block_download_file.go blacklist.json
```

## 4.  Lua file output:
The tool actually just convert input file `blacklist.json` to string and replace the variable `blacklist` 
in lua template file.
```
local function has_value (tab, content_type,content_dis)
    for index, value in ipairs(tab) do
        if (value[1] == content_type and value[2] == content_dis ) or
          (value[1] == "*" and value[2] == "*" ) or
          (value[1] == "*" and value[2] == content_dis )  or
          (value[1] == content_type and value[2] == "*" ) then
            return true
        end
    end

    return false
end

blacklist = {{"a","b"},{"a","*"},{"*","b"},{"abcd","ee"},}

req_header = ngx.req.get_headers()
content_type = req_header['Content-Type']
content_disposition  = req_header['Content-Disposition']

if has_value(blacklist,content_type,content_disposition) then
    ngx.status = 403
    ngx.exit(ngx.HTTP_FORBIDDEN)
end
```