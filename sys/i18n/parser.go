package i18n

import (
	"bytes"
	"regexp"

	"github.com/robertkrimen/otto"
)

//UIParses takes key(first arg) and value(second arg) out of structures like:
// $t('verification.blockchain.hint.valid', 'The file {filename} is valid.', {filename: file.name})
type UIParser struct {
	trans  map[string]string
	vm     *otto.Otto
	first  *regexp.Regexp
	second *regexp.Regexp
}

func NewUIParser() *UIParser {
	me := &UIParser{
		first:  regexp.MustCompile(`(?U)\$t\(('|")([\S\s]+)('|")([\S\s]*,[\S\s]*('|")([\S\s]+)('|"))?([\S\s]*,[^\)]+)?\)`),
		second: regexp.MustCompile(`(?U)('|")[\S\s]*(,[^'"]+\))`),
		vm:     otto.New(),
		trans:  map[string]string{},
	}
	_ = me.vm.Set("$t", func(call otto.FunctionCall) otto.Value {
		var k string
		var v string
		args := call.ArgumentList
		if len(args) > 0 {
			k = args[0].String()
			if len(k) > 0 {
				if len(args) > 1 {
					v = args[1].String()
				} else {
					v = k
				}
				if a, ok := me.trans[k]; !ok || len(a) == 0 {
					me.trans[k] = v
				}
			}
		}
		return otto.Value{}
	})
	return me
}

func (me *UIParser) Parse(content []byte) {
	firstRun := make([][]byte, 0)
	for _, match := range me.first.FindAllSubmatch(content, -1) {
		if len(match) > 0 {
			firstRun = append(firstRun, match[0])
		}
	}
	for i, b := range firstRun {
		for _, match := range me.second.FindAllSubmatch(b, -1) {
			if len(match) > 1 && len(match[2]) > 0 {
				firstRun[i] = bytes.Replace(b, match[2], []byte(")"), 1)
			}
		}
	}
	for _, b := range firstRun {
		_, _ = me.vm.Run(string(b))
	}
}

func (me *UIParser) Translations() map[string]string {
	return me.trans
}
