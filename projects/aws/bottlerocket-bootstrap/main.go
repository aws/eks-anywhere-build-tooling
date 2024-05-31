package main

import (
	"flag"
	"fmt"

	"github.com/aws/eks-anywhere-build-tooling/bottlerocket-bootstrap/pkg"
	"github.com/aws/eks-anywhere-build-tooling/bottlerocket-bootstrap/pkg/utils"
)

func main() {
	logVerbosity := flag.String("v", "9", "Verbosity of logs")
	utils.LogVerbosity = fmt.Sprintf("-v%s", *logVerbosity)
	pkg.Bootstrap()
}
