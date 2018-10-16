package pkg

import (
	"io/ioutil"
	"net/http"

	"github.com/samclarke/robotstxt"
)

func retrieveRobotsTxt(client *http.Client, url string) (*robotstxt.RobotsTxt, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	robot, err := robotstxt.Parse(string(contents), url)
	if err != nil {
		return nil, err
	}

	return robot, nil
}
