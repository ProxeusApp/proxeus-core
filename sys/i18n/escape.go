package i18n

import (
	"bytes"
	"regexp"
)

/**
This Escape function is making sure that only text formatting html features are preserved.

OK: <span style="font-size:10px">my text</span>
OK: <a href="http...">my link</a> more text..
Not OK: my text <style>...</style>
Not OK: my text <script>...</script>
Not OK: my text <meta...>
Not OK: my text <link...>
...
*/

var linkTagReg = regexp.MustCompile(`(?Ui)<link((\s+[\w\-\_\:]+\=[^"']\S+[^"'])|(\s+[\w\-\_\:]+\=["'][\s\S]*["'])|(\s+[\w\-\_\:]+))*(\s+)?>`)

var styleTagReg = regexp.MustCompile(`(?Ui)<style((\s+[\w\-\_\:]+\=[^"']\S+[^"'])|(\s+[\w\-\_\:]+\=["'][\s\S]*["'])|(\s+[\w\-\_\:]+))*(\s+)?>[\s\S]*</style>`)

var scriptTagReg = regexp.MustCompile(`(?Ui)<script((\s+[\w\-\_\:]+\=[^"']\S+[^"'])|(\s+[\w\-\_\:]+\=["'][\s\S]*["'])|(\s+[\w\-\_\:]+))*(\s+)?>[\s\S]*</script>`)

var metaTagReg = regexp.MustCompile(`(?Ui)<meta((\s+[\w\-\_\:]+\=[^"']\S+[^"'])|(\s+[\w\-\_\:]+\=["'][\s\S]*["'])|(\s+[\w\-\_\:]+))*(\s+)?>`)

var headTagReg = regexp.MustCompile(`(?Ui)<head((\s+[\w\-\_\:]+\=[^"']\S+[^"'])|(\s+[\w\-\_\:]+\=["'][\s\S]*["'])|(\s+[\w\-\_\:]+))*(\s+)?>`)
var headEndTagReg = regexp.MustCompile(`(?Ui)</head>`)

var headerTagReg = regexp.MustCompile(`(?Ui)<header((\s+[\w\-\_\:]+\=[^"']\S+[^"'])|(\s+[\w\-\_\:]+\=["'][\s\S]*["'])|(\s+[\w\-\_\:]+))*(\s+)?>`)
var headerEndTagReg = regexp.MustCompile(`(?Ui)</header>`)

var iframeTagReg = regexp.MustCompile(`(?Ui)<iframe((\s+[\w\-\_\:]+\=[^"']\S+[^"'])|(\s+[\w\-\_\:]+\=["'][\s\S]*["'])|(\s+[\w\-\_\:]+))*(\s+)?>`)
var iframeEndTagReg = regexp.MustCompile(`(?Ui)</iframe>`)

var titleTagReg = regexp.MustCompile(`(?Ui)<title((\s+[\w\-\_\:]+\=[^"']\S+[^"'])|(\s+[\w\-\_\:]+\=["'][\s\S]*["'])|(\s+[\w\-\_\:]+))*(\s+)?>`)
var titleEndTagReg = regexp.MustCompile(`(?Ui)</title>`)

var inputTagReg = regexp.MustCompile(`(?Ui)<input((\s+[\w\-\_\:]+\=[^"']\S+[^"'])|(\s+[\w\-\_\:]+\=["'][\s\S]*["'])|(\s+[\w\-\_\:]+))*(\s+)?>`)
var inputEndTagReg = regexp.MustCompile(`(?Ui)</input>`)

var selectTagReg = regexp.MustCompile(`(?Ui)<select((\s+[\w\-\_\:]+\=[^"']\S+[^"'])|(\s+[\w\-\_\:]+\=["'][\s\S]*["'])|(\s+[\w\-\_\:]+))*(\s+)?>`)
var selectEndTagReg = regexp.MustCompile(`(?Ui)</select>`)

var optionTagReg = regexp.MustCompile(`(?Ui)<option((\s+[\w\-\_\:]+\=[^"']\S+[^"'])|(\s+[\w\-\_\:]+\=["'][\s\S]*["'])|(\s+[\w\-\_\:]+))*(\s+)?>`)
var optionEndTagReg = regexp.MustCompile(`(?Ui)</option>`)

var htmlTagReg = regexp.MustCompile(`(?Ui)<html((\s+[\w\-\_\:]+\=[^"']\S+[^"'])|(\s+[\w\-\_\:]+\=["'][\s\S]*["'])|(\s+[\w\-\_\:]+))*(\s+)?>`)
var htmlEndTagReg = regexp.MustCompile(`(?Ui)</html>`)

var bodyTagReg = regexp.MustCompile(`(?Ui)<body((\s+[\w\-\_\:]+\=[^"']\S+[^"'])|(\s+[\w\-\_\:]+\=["'][\s\S]*["'])|(\s+[\w\-\_\:]+))*(\s+)?>`)
var bodyEndTagReg = regexp.MustCompile(`(?Ui)</body>`)

//to prevent from script in attribute values
var tagAndAttrReg = regexp.MustCompile(`(?Ui)<[\w\-\_\:]+((\s+[\w\-\_\:]+\=[^"']\S+[^"'])|(\s+[\w\-\_\:]+\=["'][\s\S]*["'])|(\s+[\w\-\_\:]+))*\s*(\/)?>`)
var attrReg = regexp.MustCompile(`(?sU)(\s*([\w\-\_\:]+)\=[\s\S]+["'\s])[\s\/>]`)
var attrEventsMap = map[string]bool{
	"onafterprint":   true,
	"onbeforeprint":  true,
	"onbeforeunload": true,
	"onerror":        true,
	"onhashchange":   true,
	"onload":         true,
	"onmessage":      true,
	"onoffline":      true,
	"ononline":       true,
	"onpagehide":     true,
	"onpageshow":     true,
	"onpopstate":     true,
	"onresize":       true,
	"onstorage":      true,
	"onunload":       true,

	"onblur":        true,
	"onchange":      true,
	"oncontextmenu": true,
	"onfocus":       true,
	"oninput":       true,
	"oninvalid":     true,
	"onreset":       true,
	"onsearch":      true,
	"onselect":      true,
	"onsubmit":      true,

	"onkeydown":  true,
	"onkeypress": true,
	"onkeyup":    true,

	"onclick":      true,
	"ondblclick":   true,
	"onmousedown":  true,
	"onmousemove":  true,
	"onmouseout":   true,
	"onmouseover":  true,
	"onmouseup":    true,
	"onmousewheel": true,
	"onwheel":      true,

	"ondrag":      true,
	"ondragend":   true,
	"ondragenter": true,
	"ondragleave": true,
	"ondragover":  true,
	"ondragstart": true,
	"ondrop":      true,
	"onscroll":    true,

	"oncopy":  true,
	"oncut":   true,
	"onpaste": true,

	"onabort":          true,
	"oncanplay":        true,
	"oncanplaythrough": true,
	"oncuechange":      true,
	"ondurationchange": true,
	"onemptied":        true,
	"onended":          true,
	"onloadeddata":     true,
	"onloadedmetadata": true,
	"onloadstart":      true,
	"onpause":          true,
	"onplay":           true,
	"onplaying":        true,
	"onprogress":       true,
	"onratechange":     true,
	"onseeked":         true,
	"onseeking":        true,
	"onstalled":        true,
	"onsuspend":        true,
	"ontimeupdate":     true,
	"onvolumechange":   true,
	"onwaiting":        true,

	"ontoggle": true,
}

//Escape makes sure javascript and linked CSS is cleaned out. Only embedded CSS is allowed.
func Escape(text string) string {
	btext := []byte(text)
	empty := []byte("")

	//get rid of all html tags we wont need for sure for formatting translations
	btext = styleTagReg.ReplaceAll(btext, empty)
	btext = scriptTagReg.ReplaceAll(btext, empty)
	btext = metaTagReg.ReplaceAll(btext, empty)
	btext = linkTagReg.ReplaceAll(btext, empty)

	btext = iframeTagReg.ReplaceAll(btext, empty)
	btext = iframeEndTagReg.ReplaceAll(btext, empty)

	btext = titleTagReg.ReplaceAll(btext, empty)
	btext = titleEndTagReg.ReplaceAll(btext, empty)

	btext = inputTagReg.ReplaceAll(btext, empty)
	btext = inputEndTagReg.ReplaceAll(btext, empty)

	btext = selectTagReg.ReplaceAll(btext, empty)
	btext = selectEndTagReg.ReplaceAll(btext, empty)

	btext = optionTagReg.ReplaceAll(btext, empty)
	btext = optionEndTagReg.ReplaceAll(btext, empty)

	btext = htmlTagReg.ReplaceAll(btext, empty)
	btext = htmlEndTagReg.ReplaceAll(btext, empty)

	btext = bodyTagReg.ReplaceAll(btext, empty)
	btext = bodyEndTagReg.ReplaceAll(btext, empty)

	btext = headTagReg.ReplaceAll(btext, empty)
	btext = headEndTagReg.ReplaceAll(btext, empty)

	btext = headerTagReg.ReplaceAll(btext, empty)
	btext = headerEndTagReg.ReplaceAll(btext, empty)

	//lookup html attribute event names to prevent from script in attribute values
	for _, b := range tagAndAttrReg.FindAllSubmatch(btext, -1) {
		if len(b) > 0 {
			original := b[0]
			toReplace := original
			for _, attr := range attrReg.FindAllSubmatch(original, -1) {
				if len(attr) == 3 {
					if attrEventsMap[string(attr[2])] {
						toReplace = bytes.Replace(toReplace, attr[1], empty, 1)
					}
				}
			}

			if len(original) != len(toReplace) {
				btext = bytes.Replace(btext, original, toReplace, 1)
			}
		}
	}

	return string(btext)
}
