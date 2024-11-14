/*
 * @author Crabin
 */

package test

import (
	"fmt"
	wappalyzer "github.com/projectdiscovery/wappalyzergo"
	"io/ioutil"
	"log"
	"net/http"
	"testing"
)

func TestGetAllApi(t *testing.T) {
	resp, err := http.DefaultClient.Get("https://www.iwate-u.ac.jp/iuic/foreigner/type/research-student.html")
	if err != nil {
		log.Fatal(err)
	}
	data, _ := ioutil.ReadAll(resp.Body) // Ignoring error for example

	wappalyzerClient, err := wappalyzer.New()

	fingerprints := wappalyzerClient.Fingerprint(resp.Header, data)
	fmt.Printf("%v\n", fingerprints)

	fmt.Printf("%v\n", fingerprints["MariaDB"])

	// Output: map[Acquia Cloud Platform:{} Amazon EC2:{} Apache:{} Cloudflare:{} Drupal:{} PHP:{} Percona:{} React:{} Varnish:{}]
}
