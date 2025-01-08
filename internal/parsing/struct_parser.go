package parsing

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"strings"

	"github.com/LinkinStars/baileys/internal/util"
)

const (
	InterfaceTypeDef = "interface"
	StructTypeDef    = "struct"
	TimeTypeDef      = "time.Time"
)

// StructFlat 非嵌套结构体
type StructFlat struct {
	Name       string
	Comment    string
	IsPlural   bool
	PluralName string
	Fields     []*StructField
}

// StructField 结构体字段
type StructField struct {
	Name        string
	ConvertName string
	NamePB      string
	Type        string
	Comment     string
	Tag         string
	IsPlural    bool
}

// GetTag 获取tag
func (s *StructField) GetTag(tagName string) string {
	arr := strings.Split(s.Tag, " ")
	for _, tag := range arr {
		tag = strings.TrimSpace(tag)
		if strings.HasPrefix(tag, tagName) {
			tag = strings.TrimLeft(tag, tagName+":")
			tag = strings.Trim(tag, "\"")
			return tag
		}
	}
	return ""
}

// GetJsonTag 获取json tag
func (s *StructField) GetJsonTag() string {
	tag := s.GetTag("json")
	// ignore json tag is `json:"-"`
	if tag == "-" {
		return ""
	}
	if len(tag) == 0 {
		return s.Name
	}
	return tag
}

// StructParser golang struct 解析器
func StructParser(src string) (structList []*StructFlat, err error) {
	src = addPackageIfNotExist(src)
	fileSet := token.NewFileSet()
	f, err := parser.ParseFile(fileSet, "src.go", src, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	structList = make([]*StructFlat, 0)
	for _, node := range f.Decls {
		switch node.(type) {
		case *ast.GenDecl:
			genDecl := node.(*ast.GenDecl)
			var structComment string
			if genDecl.Doc != nil {
				structComment = strings.TrimSpace(genDecl.Doc.Text())
			}
			for _, spec := range genDecl.Specs {
				switch spec.(type) {
				case *ast.TypeSpec:
					typeSpec := spec.(*ast.TypeSpec)

					// 获取结构体名称
					structFlat := &StructFlat{Name: typeSpec.Name.Name, Comment: structComment}
					structFlat.Fields = make([]*StructField, 0)
					log.Printf("read struct %s %s\n", typeSpec.Name.Name, structComment)

					switch typeSpec.Type.(type) {
					case *ast.StructType:
						structType := typeSpec.Type.(*ast.StructType)
						for _, reField := range structType.Fields.List {
							structField := &StructField{}
							if reField.Tag != nil {
								structField.Tag = strings.Trim(reField.Tag.Value, "`")
							}
							switch v := reField.Type.(type) {
							case *ast.Ident:
								iDent := reField.Type.(*ast.Ident)
								structField.Type = iDent.Name
							case *ast.InterfaceType:
								structField.Type = InterfaceTypeDef
							case *ast.MapType:
								iDent := reField.Type.(*ast.MapType)
								structField.Type = fmt.Sprintf("map[%s]%s", iDent.Key, iDent.Value)
							case *ast.ArrayType:
								iDent := reField.Type.(*ast.ArrayType)
								iDentElem := util.ReflectAccess(iDent.Elt)
								structField.Type = fmt.Sprintf("[]%s", iDentElem)
							case *ast.StructType:
								structField.Type = StructTypeDef
							case *ast.StarExpr:
								structField.Type = "*" + fmt.Sprintf("%+v", v.X)
							case *ast.SelectorExpr:
								iDent := reField.Type.(*ast.SelectorExpr)
								if iDent.Sel.Name == "Time" {
									structField.Type = TimeTypeDef
								} else {
									log.Printf("undefined reField names %+v, type %+v", reField.Names, reField.Type)
								}
							default:
								log.Printf("undefined reField type %+v", reField.Type)
							}
							var plural bool
							for _, name := range reField.Names {
								if !name.IsExported() {
									continue
								}
								structField.NamePB = util.ToUpperCamelCase(name.Name)
								structField.Name = name.Name
								structField.Comment = fixMultiLineComment(strings.TrimSpace(reField.Doc.Text()))
								structField.Type, plural = util.ArrayToPlural(structField.Type)
								structField.ConvertName = structField.Type
								if structField.ConvertName == "" {
									structField.ConvertName = name.Name
								}
								if plural {
									structList = append(structList, &StructFlat{
										IsPlural:   true,
										PluralName: util.ToPlural(structField.ConvertName),
										Name:       structField.ConvertName,
										Comment:    structField.Comment,
									})
									structField.IsPlural = true
									structField.ConvertName = util.ToPlural(structField.ConvertName)
								}
								structFlat.Fields = append(structFlat.Fields, structField)
								log.Printf("name=%s type=%s comment=%s tag=%s\n", name.Name, structField.Type, structField.Comment, structField.Tag)
							}
						}
					}
					structList = append(structList, structFlat)
				}
			}
		}
	}
	return structList, nil
}

func addPackageIfNotExist(src string) string {
	if strings.HasPrefix(src, "package") {
		return src
	}
	return "package mypackage\n" + src
}

func fixMultiLineComment(comment string) string {
	return strings.ReplaceAll(comment, "\n", "\n\t// ")
}
