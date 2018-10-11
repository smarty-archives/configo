package configo

import (
	"bytes"
	"errors"
	"strings"
	tt "text/template"
	"text/template/parse"

	"github.com/Masterminds/sprig"
	"github.com/iancoleman/strcase"
)


// Template retrieves values from the provided sources, inserting data
// from the sources into the Template.
// Includes support for sprig template functions (@see http://masterminds.github.io/sprig/)
type Template struct {
	Content    string
	reader     *Reader
	vaultAddr  string
	vaultToken string
}


// NewTemplate initializes a template and reader new reader using the provided sources.
func NewTemplate(content string, sources ...Source) *Template {
	return &Template{
		Content: content,
		reader: NewReader(sources...),
	}
}


// FromDirectory reads the directory path provided.
// If mustExist is true, and the path does not exist, a panic will result.
func FromTemplateJSON(template *Template) *JSONSource {
	if data, err := template.Run(); err == nil {
		return FromJSONContent(data)
	} else {
		panic(err)
	}
}


// Execute the template
func (this *Template) Run() (ret []byte, err error) {
	buf := new(bytes.Buffer)
	tpl := tt.New("self").Funcs(this.funcs())

	if tpl, err = tpl.Parse(this.Content); err != nil {
		return
	}

	//find what keys we need
	data := make(map[string]string, 32)
	this.walkNode(tpl.Root, data)

	//fill referenced keys with data
	this.fillData(data)

	if err = tpl.Execute(buf, data); err != nil {
		return
	}

	ret = buf.Bytes()
	return
}


// Execute the template, return as string
func (this *Template) String() (ret string, err error) {
	data, err := this.Run()
	ret = string(data)
	return
}


// Fill a map of keys with data, removing keys that are not found
func (this *Template) fillData(found map[string]string) {
	for key, _ := range found {
		item, err := this.findDataItem(key)
		if err == nil {
			if len(item) > 0 {
				found[key] = item[0]
			}
		} else {
			delete(found, key)
		}
	}
	return
}


// Attempt to find a key in data sources using various case patterns
func (this *Template) findDataItem(key string) (item []string, err error) {
	list := [8]string{
		key,
		strcase.ToScreamingSnake(key),
		strcase.ToSnake(key),
		strcase.ToKebab(key),
		strcase.ToScreamingKebab(key),
		strcase.ToKebab(key),
		strcase.ToLowerCamel(key),
		strcase.ToCamel(key),
	}

	for _, key = range list {
		if item, err = this.reader.StringsError(key); err != KeyNotFoundError {
			return
		}
	}

	return nil, KeyNotFoundError
}


// Register template functions.
func (this *Template) funcs() tt.FuncMap {
	ret := sprig.TxtFuncMap()
	ret["secret"] = this.funcSecret
	return ret
	//return tt.FuncMap{
	//	"secret": this.funcSecret,
	//}
}


// Template function to read a Vault secret.
func (this *Template) funcSecret(path string) (ret map[string]interface{}, err error) {
	if len(path) == 0 {
		return
	}

	if this.vaultAddr == "" {
		if val, _ := this.findDataItem("vault_addr"); len(val) > 0 {
			this.vaultAddr = val[0]
		}
	}

	if this.vaultToken == "" {
		if val, _ := this.findDataItem("vault_token"); len(val) > 0 {
			this.vaultToken = val[0]
		}
	}

	if src := FromVaultDocument(this.vaultToken, this.vaultAddr, path); src != nil {
		ret = src.values
	} else {
		err = errors.New("No data from Vault: "+ this.vaultAddr + " " + path)
	}

	return
}


// Walk through the parsed template nodes, adding referenced fields to map.
// The walking structure is similar to text/template/exec.go:walk()
func (this *Template) walkNode(node parse.Node, found map[string]string) {
	switch node := node.(type) {
		// ActionNode holds an action (something bounded by delimiters; like FieldNodes)
		case *parse.ActionNode:
			for _, c := range node.Pipe.Cmds {
				for _, a := range c.Args {
					this.walkNode(a, found)
				}
			}

		// this is where our data is, we only need the first item in a chain
		case *parse.FieldNode:
			if items := strings.Split(strings.Trim(node.String(), "."), "."); len(items) > 0 {
				found[items[0]] = ""
			}

		case *parse.IfNode:
			this.walkNodeIf(node.List, node.ElseList, found)

		case *parse.ListNode:
			for _, node := range node.Nodes {
				this.walkNode(node, found)
			}

		case *parse.RangeNode:
			this.walkNodeIf(node.List, node.ElseList, found)

		//case *parse.TemplateNode: //not supported

		case *parse.WithNode:
			this.walkNodeIf(node.List, node.ElseList, found)
	}
}


func (this *Template) walkNodeIf(list, elseList *parse.ListNode,found map[string]string) {
	if list != nil {
		this.walkNode(list, found)
	}

	if elseList != nil {
		this.walkNode(elseList, found)
	}
}
