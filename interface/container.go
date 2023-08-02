package container

import (
	_ "fmt"
	"github.com/farisridho/go-skeleton/shared/config"
)

type Container struct {
	Config *config.Config
}

func NewContainer(conf *config.Config) *Container {

}
