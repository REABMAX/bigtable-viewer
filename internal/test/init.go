package test

import (
	"cloud.google.com/go/bigtable"
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
)

const (
	itemFamily = "items"
	itemColumn = "item"
	stockFamily = "stock"
)

func InitTestData(ctx context.Context, project string, instance string, tableName string) {
	adminClient, err := bigtable.NewAdminClient(ctx, project, instance)
	if err != nil {
		log.Fatalln(err)
	}
	client, err := bigtable.NewClient(ctx, project, instance)
	defer client.Close()

	err = initSchema(ctx, adminClient, tableName)
	if err != nil {
		log.Fatalln(err)
	}

	err = initData(ctx, client, tableName)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("successfully initialized test database")
}

func initSchema(ctx context.Context, adminClient *bigtable.AdminClient, tableName string) error {
	tables, err := adminClient.Tables(ctx)
	if err != nil {
		return err
	}

	if !sliceContains(tables, tableName) {
		log.Printf("Creating table %s", tableName)
		if err := adminClient.CreateTable(ctx, tableName); err != nil {
			return err
		}
	}

	tblInfo, err := adminClient.TableInfo(ctx, tableName)
	if err != nil {
		return err
	}

	if !sliceContains(tblInfo.Families, itemFamily) {
		if err := adminClient.CreateColumnFamily(ctx, tableName, itemFamily); err != nil {
			return err
		}
	}
	if !sliceContains(tblInfo.Families, stockFamily) {
		if err := adminClient.CreateColumnFamily(ctx, tableName, stockFamily); err != nil {
			return err
		}
	}

	return nil
}

func initData(ctx context.Context, client *bigtable.Client, tableName string) error {
	items := []string{"Item 1", "Item 2", "Item 3", "Item 4"}
	stocks := []int{26, 12, 1, 5}
	table := client.Open(tableName)
	mutations := make([]*bigtable.Mutation, len(items))
	rowKeys := make([]string, len(items))

	for i,item := range items {
		mutations[i] = bigtable.NewMutation()
		mutations[i].Set(itemFamily, "items", bigtable.Now(), []byte(item))
		mutations[i].Set(stockFamily, strconv.Itoa(i), bigtable.Now(), []byte(strconv.Itoa(stocks[i])))
		mutations[i].Set(stockFamily, strconv.Itoa(i + 1), bigtable.Now(), []byte(strconv.Itoa(stocks[i] + 4)))
		rowKeys[i] = fmt.Sprintf("_%d", i)
	}

	rowErrs, err := table.ApplyBulk(ctx, rowKeys, mutations)
	if err != nil {
		return err
	}

	if rowErrs != nil {
		for _, rowErr := range rowErrs {
			log.Println(rowErr)
			return errors.New("see row errors above")
		}
	}

	return nil
}

func sliceContains(list []string, target string) bool {
	for _, s := range list {
		if s == target {
			return true
		}
	}
	return false
}
