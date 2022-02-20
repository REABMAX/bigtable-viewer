package db

import (
	"cloud.google.com/go/bigtable"
	"context"
	"encoding/json"
	"log"
)


type Client struct {
	ctx context.Context
	client *bigtable.Client
	admin *bigtable.AdminClient
}

func NewClient(ctx context.Context, project string, instance string) (*Client, error) {
	client, err := bigtable.NewClient(ctx, project, instance)
	if err != nil {
		return nil, err
	}
	admin, err := bigtable.NewAdminClient(ctx, project, instance)
	if err != nil {
		return nil, err
	}
	return &Client{
		ctx: ctx,
		client: client,
		admin: admin,
	}, nil
}

type Row struct {
	ID string
	Families map[string]*Family
}

type Family struct {
	Name string
	Columns map[string]*Column
}

type Column struct {
	Name string
	Cells []*Cell
}

type Cell struct {
	Time bigtable.Timestamp
	Value string
}

// TODO: Following methods: This should not be done that way

func (c *Cell) IsJSON() bool {
	var js json.RawMessage
	return json.Unmarshal([]byte(c.Value), &js) == nil
}

func (c *Cell) PrettyPrintJSON() string {
	var js json.RawMessage
	err := json.Unmarshal([]byte(c.Value), &js)
	if err != nil {
		log.Println(err)
		return c.Value
	}

	json, err := json.MarshalIndent(js, "", "    ")
	if err != nil {
		log.Println(err)
		return c.Value
	}
	return string(json)
}