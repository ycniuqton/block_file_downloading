# Block File Downloading 
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