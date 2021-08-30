package main

import (
	"flag"
	"fmt"

	"github.com/eks-anywhere-build-tooling/aws/bottlerocket-bootstrap/pkg"
	"github.com/eks-anywhere-build-tooling/aws/bottlerocket-bootstrap/pkg/utils"
)

func main() {
	logVerbosity := flag.String("v", "9", "Verbosity of logs")
	utils.LogVerbosity = fmt.Sprintf("-v%s", *logVerbosity)
	pkg.Bootstrap()
}
