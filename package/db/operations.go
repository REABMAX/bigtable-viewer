package db

import (
	"cloud.google.com/go/bigtable"
)

// FetchFamilies returns structure data in shape of the table's column families
func (c *Client) FetchFamilies(table string) ([]bigtable.FamilyInfo, error) {
	t, err := c.admin.TableInfo(c.ctx, table)
	if err != nil {
		return nil, err
	}

	return t.FamilyInfos, nil
}

// FetchTables returns an array of available table names of the bigtable instance given to the client.
func (c *Client) FetchTables() ([]string, error) {
	t, err := c.admin.Tables(c.ctx)
	if err != nil {
		return nil, err
	}
	return t, nil
}

func (c *Client) Close() {
	c.client.Close()
}
