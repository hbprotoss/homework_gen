package main

import "html/template"

const Tpl = `
<html>
<HEAD>
    <META http-equiv=Content-Type content='text/html; charset=utf-8'>
</head>
<body>
    <center><br><br>
        <h2>{{.grade}}{{if .isQuestion}}口算题{{else}}答案{{end}} 第{{index}}份<br></h2>
        <h3>{{.work}}</h3>{{if .isQuestion}}班级__________ 学号________ 姓名__________<br>{{end}}<br>
        <table border=0 width=640 cellspacing=10>
            {{range $idx,$v := .result}}
                {{if isRowBegin $idx}}
                <tr>
                {{end}}
                <td>{{.Question}} = {{if not $.isQuestion}}{{.Answer}}{{end}}</td>
                {{if isRowEnd $idx (len $.result)}}
                </tr>
                {{end}}
            {{end}}
        </table>
    </center>
</body>
</html>
`

func initTemplate() (tpl *template.Template) {
	tpl, err := template.New("tpl").Funcs(template.FuncMap{
		"isRowBegin": func(idx int) bool {
			return idx%4 == 0
		},
		"isRowEnd": func(idx, length int) bool {
			return idx%4 == 3 || idx == length-1
		},
	}).Parse(Tpl)
	if err != nil {
		panic(err)
	}
	return
}
