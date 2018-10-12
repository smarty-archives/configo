package configo

import (
	"bytes"
	"testing"
	tt "text/template"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
)


func TestTemplateTestFixture(t *testing.T) {
	gunit.Run(new(TemplateTestFixture), t)
}


type TemplateTestFixture struct {
	*gunit.Fixture

	sources []Source
	template  *Template
}


func (this *TemplateTestFixture) Setup() {
	this.sources = []Source{
		&FakeSource{},
		&FakeSource{key: "string", value: []string{"asdf"}},
		&FakeSource{key: "string", value: []string{"qwer"}},
		&FakeSource{key: "string-no-values", value: []string{}},
		&FakeSource{key: "int", value: []string{"42"}},
		&FakeSource{key: "int", value: []string{"-1"}},
		&FakeSource{key: "int-bad", value: []string{"not an integer"}},
		&FakeSource{key: "bool", value: []string{"true"}},
		&FakeSource{key: "bool-bad", value: []string{"not a bool"}},
		&FakeSource{key: "url", value: []string{"http://www.google.com"}},
		&FakeSource{key: "url-bad", value: []string{"%%%%%%"}}, // not a url
		&FakeSource{key: "duration", value: []string{"5s"}},
		&FakeSource{key: "duration-bad", value: []string{"not a duration"}},
		&FakeSource{key: "time", value: []string{"2015-09-15T11:29:00Z"}},
		&FakeSource{key: "time-bad", value: []string{"not a time"}},
		&FakeSource{key: "VAULT_ADDR", value: []string{"http://localhost:29999"}},
		&FakeSource{key: "VAULT_TOKEN", value: []string{"1234-567-89012"}},
	}

	this.template = NewTemplate(this.case1Template(), this.sources...)
}


////////////////////////////////////////////////////////////////


func (this *TemplateTestFixture) TestInitializeSources() {
	for _, source := range this.sources {
		this.So(source.(*FakeSource).initialized, should.Equal, 1)
	}
}


func (this *TemplateTestFixture) TestInitializeContent() {
	this.So(this.template.Content, should.Equal, this.case1Template())
}


func (this *TemplateTestFixture) TestInitializeReader() {
	this.So(this.template.reader, should.NotBeEmpty)
}


////////////////////////////////////////////////////////////////


func (this *TemplateTestFixture) TestSimpleStrings() {
	x, err := this.template.String()
	this.So(err, should.Equal, nil)
	this.So(x, should.Equal, this.getExpected(this.template.Content, this.caseData()))
}


func (this *TemplateTestFixture) TestIf() {
	this.template.Content = this.case2Template()
	x, err := this.template.String()
	this.So(err, should.Equal, nil)
	this.So(x, should.Equal, this.getExpected(this.case2Template(), this.caseData()))
}


func (this *TemplateTestFixture) TestBadTemplate() {
	this.template.Content = `{{ .String`
	_, err := this.template.String()
	this.So(err, should.NotEqual, nil)
}


func (this *TemplateTestFixture) TestBadTemplateFunction() {
	this.template.Content = `{{ .String | no_such_function }}`
	_, err := this.template.String()
	this.So(err, should.NotEqual, nil)
}


func (this *TemplateTestFixture) TestComplexTemplate() {
	this.template.Content = `value: {{- $new_var := (or .no_such_string .url) }}{{ $new_var }}`
	x, err := this.template.String()
	this.So(err, should.BeEmpty)
	this.So(x, should.Equal, "value:http://www.google.com")
}


func (this *TemplateTestFixture) TestSecretInvalid() {
	this.template.Content = `MyEmail: "{{with secret "secret/operations/email"}}{{.string}}{{end}}"`
	this.So(func() { this.template.String() }, should.Panic)
	this.So(this.template.vaultAddr, should.Equal, "http://localhost:29999")
	this.So(this.template.vaultToken, should.Equal, "1234-567-89012")
}


func (this *TemplateTestFixture) TestSecretValid() {
	svr := dummyHTTP(false, nil)
	defer svr.Close()

	this.template.Content = `MyEmail: "{{with secret "secret/operations/email"}}{{.string}}{{end}}"`
	this.template.vaultAddr = svr.URL
	this.template.vaultToken = "1234-567-89012"
	x, err := this.template.String()
	this.So(err, should.BeEmpty)
	this.So(x, should.Equal, `MyEmail: "String"`)
}


func (this *TemplateTestFixture) TestSprig() {
	this.template.Content = `{{ "hello!" | upper | repeat 5 }}`
	x, err := this.template.String()
	this.So(err, should.BeEmpty)
	this.So(x, should.Equal, `HELLO!HELLO!HELLO!HELLO!HELLO!`)
}


func (this *TemplateTestFixture) TestFromTemplateJSON() {
	actual := FromTemplateJSON(this.template)
	x, _ := actual.Strings("string")
	this.So(len(x), should.Equal, 1)
	this.So(x[0], should.Equal, "asdf")
	x, _ = actual.Strings("int")
	this.So(len(x), should.Equal, 1)
	this.So(x[0], should.Equal, "42")
}


////////////////////////////////////////////////////////////////


func (this *TemplateTestFixture) case1Template() string {
	return `{
	"string": "{{ .String }}",
	"string-no-values": "{{ .StringNoValues }}",
	"int": {{ .Int }},
	"int-bad": "{{ .IntBad }}",
	"bool": {{ .Bool }},
	"bool-bad": "{{ .BoolBad }}",
	"url": "{{ .Url }}",
	"url-bad": "{{ .UrlBad }}",
	"duration": "{{ .Duration }}",
	"duration-bad": "{{ .DurationBad }}",
	"inline-string": "my-string",
	"inline-int": 9999
}`
}


func (this *TemplateTestFixture) case2Template() string {
	return `
{{if .Url -}}
		url: "{{ .Url }}",
		url-bad: "{{ .UrlBad }}",
		does-not-exist: "{{ .DoesNotExist }}"
{{- end}}
`
}


func (this *TemplateTestFixture) caseData() map[string]string {
	//make a template-valid data map (e.g., [A-Za-z0-9_])
	return map[string]string{
		"String": this.template.reader.String("string"),
		"StringNoValues": this.template.reader.String("string-no-values"),
		"Int": this.template.reader.String("int"),
		"IntBad": this.template.reader.String("int-bad"),
		"Bool": this.template.reader.String("bool"),
		"BoolBad": this.template.reader.String("bool-bad"),
		"Url": this.template.reader.String("url"),
		"UrlBad": this.template.reader.String("url-bad"),
		"Duration": this.template.reader.String("duration"),
		"DurationBad": this.template.reader.String("duration-bad"),
		"Time": this.template.reader.String("time"),
		"TimeBad": this.template.reader.String("time-bad"),
	}
}


// Use straight text/template with known data to check output
func (this *TemplateTestFixture) getExpected(template string, data map[string]string) string {
	buf := new(bytes.Buffer)
	tpl := tt.New("testing1")
	var err error

	if tpl, err = tpl.Parse(template); err != nil {
		panic("Bad test template: " + err.Error())
	}

	if err = tpl.Execute(buf, data); err != nil {
		panic("Bad test template expected results: " + err.Error())
	}

	return buf.String()
}
