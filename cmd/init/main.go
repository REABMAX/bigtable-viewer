package main

import (
	"bigtable-viewer/internal/settings"
	"bigtable-viewer/internal/test"
	"context"
)

const (
	tableName = "testdata"
)

func main() {
	ctx := context.Background()
	project, instance, _ := settings.All()
	test.InitTestData(ctx, project, instance, tableName)
}