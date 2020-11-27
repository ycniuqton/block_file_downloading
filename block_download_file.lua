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
