package converter

import (
	"fmt"
	"strings"

	"github.com/LinkinStars/baileys/internal/parsing"
	"github.com/LinkinStars/baileys/internal/util"
)

var (
	go2pbTypeMapping = map[string]string{
		"float32":       "float",
		"float64":       "double",
		"complex64":     "double",
		"complex128":    "double",
		"int":           "int32",
		"int8":          "int32",
		"int16":         "int32",
		"int32":         "int32",
		"int64":         "int64",
		"uint":          "uint32",
		"uint8":         "uint32",
		"uint16":        "uint32",
		"uint32":        "uint32",
		"uint64":        "uint64",
		"bool":          "bool",
		"string":        "string",
		"[]byte":        "bytes",
		"uintptr":       "bytes",
		"interface":     "bytes",
		"struct":        "bytes",
		"time.Time":     "google.protobuf.Timestamp",
		"time.Duration": "google.protobuf.Duration",
	}
)

// PBFlat pb struct
type PBFlat struct {
	Name        string
	Comment     string
	PBFieldList []*PBField
}

// PBField pb field struct
type PBField struct {
	Repeated bool
	Type     string
	Name     string
	Tag      string
	Comment  string
	Index    int
}

// GoStruct2PB convert golang struct to Protocol Buffers
func GoStruct2PB(structList []*parsing.StructFlat) (pbList []*PBFlat) {
	pbList = make([]*PBFlat, 0)
	for _, s := range structList {
		if s.IsPlural {
			continue
		}
		pbFlat := &PBFlat{
			Name:        s.Name,
			Comment:     s.Comment,
			PBFieldList: make([]*PBField, 0),
		}
		for idx, field := range s.Fields {
			pbField := &PBField{
				Repeated: field.IsPlural,
				Name:     util.ToLowerSnakeCase(field.Name),
				Type:     GoType2PB(field.Type),
				Tag:      trimTagOmit(field.GetJsonTag()),
				Comment:  field.Comment,
				Index:    idx + 1,
			}
			pbFlat.PBFieldList = append(pbFlat.PBFieldList, pbField)
		}
		pbList = append(pbList, pbFlat)
	}
	return pbList
}

// GoType2PB convert golang type to Protocol Buffers
// https://developers.google.com/protocol-buffers/docs/proto3
// .proto    Go
// double -> float64
// float -> float32
// int32 -> int32
// int64 -> int64
// uint32 -> uint32
// uint64 -> uint64
// sint32 -> int32
// sint64 -> int64
// fixed32 -> uint32
// fixed64 -> uint64
// sfixed32 -> int32
// sfixed64 -> int64
// bool -> bool
// string -> string
// bytes -> []byte
func GoType2PB(goType string) (pbType string) {
	pbType = go2pbTypeMapping[goType]
	if len(pbType) > 0 {
		return pbType
	}
	// 处理 slice 的情况
	if strings.HasPrefix(goType, "[]") {
		arrType := strings.TrimLeft(goType, "[]")
		pbType = go2pbTypeMapping[arrType]
		if len(pbType) == 0 {
			pbType = arrType
		}
		return "repeated " + pbType
	}
	// 处理 map 的情况
	if strings.HasPrefix(goType, "map[") {
		goType = strings.TrimLeft(goType, "map[")
		idx := strings.Index(goType, "]")
		key := goType[0:idx]
		val := goType[idx+1:]
		return fmt.Sprintf("map<%s, %s>", key, val)
	}
	// 其他情况可能为嵌套结构，直接返回原类型
	return goType
}

func trimTagOmit(tag string) string {
	if strings.HasSuffix(tag, ",omitempty") {
		return strings.TrimSuffix(tag, ",omitempty")
	}
	return tag
}
