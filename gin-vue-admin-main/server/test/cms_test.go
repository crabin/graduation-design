/*
 * @author Crabin
 */

package test

import (
	"crypto/md5"
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/service/myPkg/cmsDiscriminate/fetch"
	"testing"
)

func TestCMS(t *testing.T) {
	content, _, err := fetch.Get("http://gl.sycm.edu.cn//xsweb/images/button/bgbtn2_0.gif")
	if err != nil {
		return
	}
	md5str := fmt.Sprintf("%x", md5.Sum(content))
	if md5str == "061a9376bdb3bfaacfec43986456d455" {
		fmt.Sprintf("Success!")
	}
}
