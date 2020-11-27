package main

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "os/exec"
    "os"
)

var nginx_config_dir = "/etc/nginx/"




type Blacklist struct {
    Data []Item `json:"blacklist"`
}

type Item struct {
    Dis   string `json:"content_disposition"`
    Type   string `json:"content_type"`
}

var lua_file_template = `local function has_value (tab, content_type,content_dis)
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

blacklist = %s

req_header = ngx.req.get_headers()
content_type = req_header['Content-Type']
content_disposition  = req_header['Content-Disposition']

if has_value(blacklist,content_type,content_disposition) then
    ngx.status = 403
    ngx.exit(ngx.HTTP_FORBIDDEN)
end
`

func main() {
//     Get input file
    os_arg :=  os.Args
    blacklist_file := "blacklist.json"
    if len(os_arg) == 2 {
        blacklist_file = os_arg[1]
    }

//     Read blacklist file
    jsonFile, err := os.Open(blacklist_file)
    if err != nil {
        fmt.Println(err)
    }
    fmt.Println("Successfully Opened ",blacklist_file)
    defer jsonFile.Close()

//  Load json file
    byteValue, _ := ioutil.ReadAll(jsonFile)
    var blacklist Blacklist
    json.Unmarshal(byteValue, &blacklist)

//  generate blacklist text
    var blacklist_text = "{"
    for i := 0; i < len(blacklist.Data); i++ {
        blacklist_text +=fmt.Sprintf(`{"%s","%s"},`,blacklist.Data[i].Type,blacklist.Data[i].Dis)
    }
    blacklist_text += "}"

//  generate lua script file
    var lua_file_content  = fmt.Sprintf(lua_file_template, blacklist_text)
    f, err := os.Create(nginx_config_dir+"block_download_file.lua")
    if err != nil {
        fmt.Println("Cant create file")
    }
    defer f.Close()
    _, err2 := f.WriteString(lua_file_content)
    if err2 != nil {
        fmt.Println("Cant write to file")
    }

//  Restart nginx
    cmd := exec.Command("service", "nginx","restart")
    err = cmd.Run()
	if err != nil {
		fmt.Println( err)
	}

    fmt.Println("done")
}