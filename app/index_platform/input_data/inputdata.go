package input_data

import (
	logs "github.com/hanzug/goS/pkg/logger"
	"go.uber.org/zap"
	"strings"

	"github.com/spf13/cast"

	"github.com/hanzug/goS/pkg/util/stringutils"
	"github.com/hanzug/goS/types"
)

// Doc2Struct 从csv转换到struct
func Doc2Struct(docStr string) (*types.Document, error) {
	zap.S().Info(logs.RunFuncName())
	docStr = strings.Replace(docStr, "\"", "", -1)
	d := strings.Split(docStr, ",")
	something2Str := make([]string, 0)

	for i := 2; i < 5; i++ {
		if len(d) > i && d[i] != "" {
			something2Str = append(something2Str, d[i])
		}
	}

	doc := &types.Document{
		DocId: cast.ToInt64(d[0]),
		Title: d[1],
		Body:  stringutils.StrConcat(something2Str),
	}

	return doc, nil
}
