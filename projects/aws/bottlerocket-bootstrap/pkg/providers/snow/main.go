package main

import (
	"flag"
	"fmt"

	"github.com/aws/eks-anywhere-build-tooling/bottlerocket-bootstrap/pkg/providers/snow/system"
	"github.com/aws/eks-anywhere-build-tooling/bottlerocket-bootstrap/pkg/utils"
)

func main() {
	utils.LogVerbosity = fmt.Sprintf("-v%s", *flag.String("v", "9", "Verbosity of logs"))
	system.NewSnow().Bootstrap()
}
