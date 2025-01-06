package generator

import (
	"bytes"
	"log"
	"strings"
	"text/template"

	"github.com/LinkinStars/baileys/internal/parsing"
)

const (
	struct2PBFuncTpl = `

{{- if .IsPlural -}}
	type (
		{{.PluralName}} = []*ent.{{.Name}}
		{{.PluralName}}PB = []*pb.{{.Name}}
	)
	{{if and (ne .Comment "")}}// Convert{{.PluralName}}2PB {{ .Comment}}
	{{end -}}func Convert{{.PluralName}}2PB(gosModel {{.PluralName}}) (pbsModel {{.PluralName}}PB) {
		for _, model := range gosModel {
			pbsModel = append(pbsModel, Convert{{.Name}}2PB(model))
		}
		return pbsModel
	}

	{{if and (ne .Comment "")}}// Convert{{.PluralName}}PB2Object {{ .Comment}}
	{{end -}}func Convert{{.PluralName}}PB2Object(pbsModel {{.PluralName}}) (gosModel {{.PluralName}}) {
		for _, model := range pbsModel {
			gosModel = append(gosModel, Convert{{.Name}}2PB(model))
		}
		return gosModel
	}

{{- else}}
	type (
		{{.Name}} = ent.{{.Name}}
		{{.Name}}PB = pb.{{.Name}}
	)

	{{if and (ne .Comment "")}}// Convert{{.Name}}2PB {{ .Comment}}
	{{end -}}func Convert{{.Name}}2PB(goModel *{{.Name}}) (pbModel *{{.Name}}PB) {
    pbModel = &{{.Name}}PB{} 
	if goModel == nil {
		return pbModel
	}
	{{range .Fields}}
    {{if eq .Type "bool" -}}
    pbModel.{{.NamePB}} = goModel.{{.Name}}
    {{- else if eq .Type "int" -}}
    pbModel.{{.NamePB}} = int32(goModel.{{.Name}})
    {{- else if eq .Type "int8" -}}
    pbModel.{{.NamePB}} = int32(goModel.{{.Name}})
    {{- else if eq .Type "int16" -}}
    pbModel.{{.NamePB}} = int32(goModel.{{.Name}})
    {{- else if eq .Type "int32" -}}
    pbModel.{{.NamePB}} = int32(goModel.{{.Name}})
    {{- else if eq .Type "int64" -}}
    pbModel.{{.NamePB}} = int64(goModel.{{.Name}})
    {{- else if eq .Type "uint" -}}
    pbModel.{{.NamePB}} = int32(goModel.{{.Name}})
    {{- else if eq .Type "uint8" -}}
    pbModel.{{.NamePB}} = int32(goModel.{{.Name}})
    {{- else if eq .Type "uint16" -}}
    pbModel.{{.NamePB}} = int32(goModel.{{.Name}})
    {{- else if eq .Type "uint32" -}}
    pbModel.{{.NamePB}} = int32(goModel.{{.Name}})
    {{- else if eq .Type "uint64" -}}
    pbModel.{{.NamePB}} = int64(goModel.{{.Name}})
    {{- else if eq .Type "uintptr" -}}
    pbModel.{{.NamePB}} = goModel.{{.Name}}
    {{- else if eq .Type "float32" -}}
    pbModel.{{.NamePB}} = goModel.{{.Name}}
    {{- else if eq .Type "float64" -}}
    pbModel.{{.NamePB}} = goModel.{{.Name}}
    {{- else if eq .Type "complex64" -}}
    pbModel.{{.NamePB}} = goModel.{{.Name}}
    {{- else if eq .Type "complex128" -}}
    pbModel.{{.NamePB}} = goModel.{{.Name}}
    {{- else if eq .Type "interface{}" -}}
    pbModel.{{.NamePB}} = goModel.{{.Name}}
    {{- else if eq .Type "map[string]string" -}}
    pbModel.{{.NamePB}} = goModel.{{.Name}}
    {{- else if eq .Type "string" -}}
    pbModel.{{.NamePB}} = goModel.{{.Name}}
    {{- else if eq .Type "[]string" -}}
    pbModel.{{.NamePB}} = goModel.{{.Name}}
    {{- else if eq .Type "struct{}" -}}
    pbModel.{{.NamePB}} = goModel.{{.Name}}
    {{- else if eq .Type "time.Time" -}}
    pbModel.{{.NamePB}} = timestamppb.New(goModel.{{.Name}})
    {{- else -}}
    pbModel.{{.NamePB}} = Convert{{.ConvertName}}2PB(goModel.{{.Name}})
    {{- end}}
	{{- end}}
    return pbModel
}

	{{if and (ne .Comment "")}}// Convert{{.Name}}PB2Object {{ .Comment}}
	{{end -}}func Convert{{.Name}}PB2Object(pbModel *{{.Name}}PB) (goModel *{{.Name}}) {
    goModel = &{{.Name}}{}
	if pbModel == nil {
		return goModel
	}
	{{range .Fields}}
    {{if eq .Type "bool" -}}
    goModel.{{.Name}} = pbModel.{{.NamePB}}
    {{- else if eq .Type "int" -}}
    goModel.{{.Name}} = int32(pbModel.{{.NamePB}})
    {{- else if eq .Type "int8" -}}
    goModel.{{.Name}} = int32(pbModel.{{.NamePB}})
    {{- else if eq .Type "int16" -}}
    goModel.{{.Name}} = int32(pbModel.{{.NamePB}})
    {{- else if eq .Type "int32" -}}
    goModel.{{.Name}} = int32(pbModel.{{.NamePB}})
    {{- else if eq .Type "int64" -}}
    goModel.{{.Name}} = int64(pbModel.{{.NamePB}})
    {{- else if eq .Type "uint" -}}
    goModel.{{.Name}} = int32(pbModel.{{.NamePB}})
    {{- else if eq .Type "uint8" -}}
    goModel.{{.Name}} = int32(pbModel.{{.NamePB}})
    {{- else if eq .Type "uint16" -}}
    goModel.{{.Name}} = int32(pbModel.{{.NamePB}})
    {{- else if eq .Type "uint32" -}}
    goModel.{{.Name}} = int32(pbModel.{{.NamePB}})
    {{- else if eq .Type "uint64" -}}
    goModel.{{.Name}} = int64(pbModel.{{.NamePB}})
    {{- else if eq .Type "uintptr" -}}
    goModel.{{.Name}} = pbModel.{{.NamePB}}
    {{- else if eq .Type "float32" -}}
    goModel.{{.Name}} = pbModel.{{.NamePB}}
    {{- else if eq .Type "float64" -}}
    goModel.{{.Name}} = pbModel.{{.NamePB}}
    {{- else if eq .Type "complex64" -}}
    goModel.{{.Name}} = pbModel.{{.NamePB}}
    {{- else if eq .Type "complex128" -}}
    goModel.{{.Name}} = pbModel.{{.NamePB}}
    {{- else if eq .Type "interface{}" -}}
    goModel.{{.Name}} = pbModel.{{.NamePB}}
    {{- else if eq .Type "map[string]string" -}}
    goModel.{{.Name}} = pbModel.{{.NamePB}}
    {{- else if eq .Type "string" -}}
    goModel.{{.Name}} = pbModel.{{.NamePB}}
    {{- else if eq .Type "[]string" -}}
    goModel.{{.Name}} = pbModel.{{.NamePB}}
    {{- else if eq .Type "struct{}" -}}
    goModel.{{.Name}} = pbModel.{{.NamePB}}
    {{- else if eq .Type "time.Time" -}}
    goModel.{{.Name}} = pbModel.{{.NamePB}}.AsTime()
    {{- else -}}
    goModel.{{.Name}} = Convert{{.ConvertName}}PB2Object(pbModel.{{.NamePB}})
    {{- end}}
	{{- end}}
    return goModel
	}

{{- end}}
`
)

// GenerateStruct2PBFunc 生成 golang struct 转换为 protobuf 的方法
func GenerateStruct2PBFunc(structList []*parsing.StructFlat) (res string, err error) {
	funcs := map[string]interface{}{
		"contains": strings.Contains,
	}
	t, err := template.New("struct2PBFuncTpl.tpl").Funcs(funcs).Parse(struct2PBFuncTpl)
	if err != nil {
		log.Printf("could not parse template: %s\n", err.Error())
		return "", err
	}

	for _, s := range structList {
		resBytes := bytes.NewBufferString("")
		err := t.Execute(resBytes, s)
		if err != nil {
			log.Printf("could not generate %s", err.Error())
		}
		resBytes.WriteString("\n")
		res += resBytes.String()
	}
	return
}
