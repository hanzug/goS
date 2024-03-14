package analyzer

import (
	logs "github.com/hanzug/goS/pkg/logger"
	"go.uber.org/zap"
	"strings"

	"github.com/hanzug/goS/types"
)

// GseCutForBuildIndex 分词 IK for building index
func GseCutForBuildIndex(docId int64, content string) ([]*types.Tokenization, error) {
	zap.S().Info(logs.RunFuncName())
	content = ignoredChar(content)
	c := GlobalSega.CutSearch(content)
	token := make([]*types.Tokenization, 0)
	for _, v := range c {
		token = append(token, &types.Tokenization{
			Token: v,
			DocId: docId,
		})
	}

	return token, nil
}

func ignoredChar(str string) string {
	zap.S().Info(logs.RunFuncName())
	for _, c := range str {
		switch c {
		case '\f', '\n', '\r', '\t', '\v', '!', '"', '#', '$', '%', '&',
			'\'', '(', ')', '*', '+', ',', '-', '.', '/', ':', ';', '<', '=', '>',
			'?', '@', '[', '\\', '【', '】', ']', '“', '”', '「', '」', '★', '^', '·', '_', '`', '{', '|', '}', '~', '《', '》', '：',
			'（', '）', 0x3000, 0x3001, 0x3002, 0xFF01, 0xFF0C, 0xFF1B, 0xFF1F:
			str = strings.ReplaceAll(str, string(c), "")
		}
	}
	return str
}
