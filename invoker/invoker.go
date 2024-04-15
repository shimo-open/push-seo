package invoker

import (
	"fmt"
	"os"

	"github.com/ego-component/eos"
	"github.com/gotomicro/ego/client/ehttp"
)

var (
	S3Cli   *eos.Component
	HttpCli *ehttp.Component
)

func Init() error {
	akID := os.Getenv("AK_ID")
	if akID == "" {
		return fmt.Errorf("environment AK_ID can't be empty")
	}
	akSecret := os.Getenv("AK_SECRET")
	if akSecret == "" {
		return fmt.Errorf("environment AK_SECRET can't be empty")
	}
	S3Cli = eos.Load("eos").Build(eos.WithAccessKeyID(akID), eos.WithAccessKeySecret(akSecret))
	HttpCli = ehttp.Load("http.baidu").Build()
	return nil
}
