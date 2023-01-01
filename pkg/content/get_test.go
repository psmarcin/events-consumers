package content

import (
	"reflect"
	"strings"
	"testing"

	"golang.org/x/net/html"
)

func TestGetValue(t *testing.T) {
	bodyStr := "<div><p>1</p><p class=aaa>2</p></div>"
	selector := "div .aaa"
	expectValue := "2"

	body, _ := html.Parse(strings.NewReader(bodyStr))

	value, err := getContent(body, selector)

	if err != nil {
		t.Errorf("can't get value: %s", err)
	}

	if value != expectValue {
		t.Errorf("TestGetValue expect %s, but got %s", expectValue, value)
	}
}

func TestGetRequestAndValue(t *testing.T) {
	url := "https://www.onet.pl/"
	selector := ".MenuIcon_showLabelText__MurLA"
	expectValue := "Sympatia"
	resp, err := runCommand("curl", url)
	if err != nil {
		t.Errorf("TestGetRequest got error %s", err)
	}
	value, err := getContent(resp, selector)

	if err != nil {
		t.Errorf("can't get value: %s", err)
	}

	if value != expectValue {
		t.Errorf("TestGetValue expect %s, but got %s", expectValue, value)
	}
}

func Test_getContent(t *testing.T) {
	type args struct {
		body     *html.Node
		selector string
	}

	simple, _ := html.Parse(strings.NewReader("<p>1</p>"))
	simpleClass, _ := html.Parse(strings.NewReader("<p class=testClass>2</p>"))
	nested, _ := html.Parse(strings.NewReader("<div><p>3</p></div>"))
	nestedClass, _ := html.Parse(strings.NewReader("<div><p class=testClass>4</p></div>"))
	doubleNestedClass, _ := html.Parse(strings.NewReader("<div id=container><p class=testClass>5</p></div>"))

	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "should find one element",
			args: args{
				body:     simple,
				selector: "p",
			},
			want:    "1",
			wantErr: false,
		},
		{
			name: "should find one class element",
			args: args{
				body:     simpleClass,
				selector: ".testClass",
			},
			want:    "2",
			wantErr: false,
		},
		{
			name: "should find one nested element",
			args: args{
				body:     nested,
				selector: "div > p",
			},
			want:    "3",
			wantErr: false,
		},
		{
			name: "should find one nested class element",
			args: args{
				body:     nestedClass,
				selector: "div > .testClass",
			},
			want:    "4",
			wantErr: false,
		},
		{
			name: "should find one double nested class element",
			args: args{
				body:     doubleNestedClass,
				selector: "#container > .testClass",
			},
			want:    "5",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getContent(tt.args.body, tt.args.selector)
			if (err != nil) != tt.wantErr {
				t.Errorf("getContent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("getContent() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_runCommand(t *testing.T) {
	type args struct {
		command   string
		arguments string
	}
	tests := []struct {
		name    string
		args    args
		want    *html.Node
		wantErr bool
	}{
		{
			name: "should return error - command not found",
			args: args{
				command:   "NOT_EXISTING",
				arguments: "--",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "should return no content",
			args: args{
				command:   "echo",
				arguments: "",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := runCommand(tt.args.command, tt.args.arguments)
			if (err != nil) != tt.wantErr {
				t.Errorf("runCommand() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			//fmt.Printf("ttwant: %s\n", tt.want.)
			//fmt.Printf("got: %s\n", got)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("runCommand() got = %v, want %v", got, tt.want)
			}
		})
	}
}
