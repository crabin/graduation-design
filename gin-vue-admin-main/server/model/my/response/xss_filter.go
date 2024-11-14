/*
 * @author Crabin
 */

package response

type FilterRes struct {
	SafeHtml     string   `json:"safeHtml"`
	XssSentences []string `json:"xssSentences"`
	HasXSS       int      `json:"hasXSS"`
}
