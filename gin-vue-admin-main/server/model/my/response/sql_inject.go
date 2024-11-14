/*
 * @author Crabin
 */

package response

import "github.com/flipped-aurora/gin-vue-admin/server/common"

type Sql struct {
	Urls               []common.VulnerableURL `json:"urls"`
	GeneralMeasurement []common.VulnerableURL `json:"generalMeasurement"`
	TimeBlind          []common.VulnerableURL `json:"timeBlind"`
	SpadeTime          int                    `json:"spadeTime"`
}
