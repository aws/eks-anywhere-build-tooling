package arn

import (
	"fmt"
	"strings"
)

func ForImageBuilderObject(account, region, kind, name string) string {
	return fmt.Sprintf("arn:aws:imagebuilder:%s:%s:%s/%s", region, account, kind, NameForARN(name))
}

func ForVersionImageBuilderObject(account, region, kind, name, version string) string {
	return fmt.Sprintf("%s/%s", ForImageBuilderObject(account, region, kind, name), version)
}

func ForLastVersionImageBuilderObject(account, region, kind, name string) string {
	return ForVersionImageBuilderObject(account, region, kind, name, "x.x.x")
}

func NameForARN(name string) string {
	return strings.ToLower(strings.ReplaceAll(name, " ", "-"))
}

func ForRegion(arn, region string) string {
	return strings.Replace(arn, "{region}", region, 1)
}
