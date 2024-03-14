package input_data

import (
	"fmt"
	logs "github.com/hanzug/goS/pkg/logger"
	"go.uber.org/zap"

	"github.com/hanzug/goS/pkg/fileutils"
)

const InputDataPath = "./input_data"

func WukongData2MapReduce() {
	zap.S().Info(logs.RunFuncName())
	res := fileutils.GetFiles(InputDataPath)
	fmt.Println(res)
}
