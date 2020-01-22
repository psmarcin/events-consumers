package get_content

import (
	"strings"
	"testing"

	"golang.org/x/net/html"
)

func TestGetRequest(t *testing.T) {
	url := "https://www.onet.pl/"

	_, err := requestCommand(url)
	if err != nil {
		t.Errorf("TestGetRequest got error %s", err)
	}
}

func TestGetValue(t *testing.T){
	bodyStr := "<div><p>1</p><p class=aaa>2</p></div>"
	selector := "div .aaa"
	expectValue := "2"

	body, _ := html.Parse(strings.NewReader(bodyStr))

	value, err := getContent(body, selector)

	if err != nil{
		t.Errorf("can't get value: %s", err)
	}

	if value != expectValue{
		t.Errorf("TestGetValue expect %s, but got %s", expectValue, value)
	}
}

func TestGetRequestAndValue(t *testing.T){
	url := "https://www.onet.pl/"
	selector:=".serviceName"
	expectValue := "Sympatia"
	resp, err := requestCommand(url)
	if err != nil {
		t.Errorf("TestGetRequest got error %s", err)
	}
	value, err := getContent(resp, selector)

	if err != nil{
		t.Errorf("can't get value: %s", err)
	}

	if value != expectValue{
		t.Errorf("TestGetValue expect %s, but got %s", expectValue, value)
	}
}