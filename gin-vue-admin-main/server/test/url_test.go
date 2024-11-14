/*
 * @author Crabin
 */

package test

import (
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"testing"
)

func TestUrl(t *testing.T) {
	url := "ftp://www.example.com/?asf=2"
	fmt.Println(utils.MatchUrl(url))
}
