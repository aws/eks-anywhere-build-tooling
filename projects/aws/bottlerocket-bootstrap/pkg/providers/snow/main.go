package main

import (
	"flag"
	"fmt"

	"github.com/eks-anywhere-build-tooling/aws/bottlerocket-bootstrap/pkg/providers/snow/system"
	"github.com/eks-anywhere-build-tooling/aws/bottlerocket-bootstrap/pkg/utils"
)

func main() {
	utils.LogVerbosity = fmt.Sprintf("-v%s", *flag.String("v", "9", "Verbosity of logs"))
	system.NewSnow().Bootstrap()
}
