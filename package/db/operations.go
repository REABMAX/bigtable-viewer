package db

import (
	"cloud.google.com/go/bigtable"
	"strings"
)

func (c *Client) FetchFamilies(table string) ([]bigtable.FamilyInfo, error) {
	t, err := c.admin.TableInfo(c.ctx, table)
	if err != nil {
		return nil, err
	}

	return t.FamilyInfos, nil
}

func (c *Client) FetchRows(table string, start string, limit int, search string) ([]Row, error) {
	t := c.client.Open(table)
	var items []Row
	count := 0

	appendItem := func(items []Row, item bigtable.Row, counter *int) []Row {
		items = append(items, mapRow(item))
		count++
		return items
	}

	err := t.ReadRows(c.ctx, bigtable.InfiniteRange(start), func(row bigtable.Row) bool {
		if search != "" {
			if strings.Contains(row.Key(), search) {
				items = appendItem(items, row, &count)
			}
		} else {
			items = appendItem(items, row, &count)
		}

		if count >= limit {
			return false
		}
		return true
	})
	if err != nil {
		return nil, err
	}

	if len(items) > 1 {
		return items[0:], nil
	}
	return items, nil
}

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

func mapRow(orig bigtable.Row) Row {
	id := orig.Key()
	families := make(map[string]*Family)
	for familyName, cells := range orig {
		columns := make(map[string]*Column)
		for _, cell := range cells {
			if _, ok := columns[cell.Column]; !ok {
				columnName := strings.TrimPrefix(cell.Column, familyName+":")
				columns[cell.Column] = &Column{
					Name:  columnName,
					Cells: []*Cell{},
				}
			}

			columns[cell.Column].Cells = append(columns[cell.Column].Cells, &Cell{
				Value: string(cell.Value),
				Time:  cell.Timestamp,
			})
		}

		families[familyName] = &Family{
			Name:    familyName,
			Columns: columns,
		}
	}

	return Row{
		ID:       id,
		Families: families,
	}
}
