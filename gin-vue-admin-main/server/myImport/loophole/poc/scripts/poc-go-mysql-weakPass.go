/*
 * @author Crabin
 */

package scripts

import (
	"database/sql"
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/myImport/kunpeng/config"
	"github.com/flipped-aurora/gin-vue-admin/server/myImport/loophole/pkg/util"
	"strconv"
	"strings"
)

func init() {
	ScriptRegister("poc-go-mysql-weakpass", MysqlWeakPass)
}

func MysqlWeakPass(args *ScriptScanArgs) (*util.ScanResult, error) {
	url := args.Host + ":" + strconv.Itoa(int(args.Port))

	if strings.IndexAny(url, "http") == 0 {
		return &util.InVulnerableResult, nil
	}
	userList := []string{
		"root", "www", "bbs", "web", "admin",
	}
	passList, _ := config.LoadDictionary("myImport/kunpeng/files/weakPass.txt")
	for _, user := range userList {
		for _, pass := range passList {
			pass = strings.Replace(pass, "{user}", user, -1)
			connStr := fmt.Sprintf("%s:%s@tcp(%s)/?timeout=%ds", user, pass, url, config.Config.Timeout)
			db, err := sql.Open("mysql", connStr)
			if err != nil {
				break
			}
			err = db.Ping()
			if err == nil {
				db.Close()

				return util.VulnerableTcpOrUdpResult(url, "Mysql weak password.",
					[]string{string(user + ":" + pass)},
					[]string{},
				), nil
				break
			} else if strings.Contains(err.Error(), "Access denied") {
				continue
			} else {
				return &util.InVulnerableResult, nil
			}
		}
	}

	return &util.InVulnerableResult, nil
}
