package db

import (
	"cloud.google.com/go/bigtable"
	"sort"
	"strings"
	"sync"
)

// FetchRows searches the given table for rows and returns an array of Row, limited by the limit parameter.
// You can provide a start string telling it to only search rows which key's are greater than start. This is useful
// the search string.
func (c *Client) FetchRows(table string, start string, limit int, search string) ([]Row, error) {
	t := c.client.Open(table)

	sampleRowKeys, err := t.SampleRowKeys(c.ctx)
	if err != nil {
		return nil, err
	}

	channel := make(chan []Row, len(sampleRowKeys))
	var wg sync.WaitGroup
	wg.Add(len(sampleRowKeys))

	go func() {
		wg.Wait()
		close(channel)
	}()

	for i, k := range sampleRowKeys {
		go func(i int, rowKey string) {
			defer wg.Done()

			var rows []Row
			counter := 0

			t.ReadRows(c.ctx, createRowRange(sampleRowKeys, rowKey, i), func(row bigtable.Row) bool {
				switch {
				case search != "":
					if strings.Contains(row.Key(), search) {
						rows = append(rows, mapRow(row))
						counter++
					}
				case start != "":
					if row.Key() > start {
						rows = append(rows, mapRow(row))
						counter++
					}
				default:
					rows = append(rows, mapRow(row))
					counter++
				}

				return true
			})

			channel <- rows
		}(i, k)
	}

	var items []Row
	count := 0
	for res := range channel {
		for _, r := range res {
			if count >= limit {
				break
			}

			items = append(items, r)
			count++
		}
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].ID < items[j].ID
	})

	return items, nil
}

func createRowRange(sampleRowKeys []string, rowKey string, iteration int) bigtable.RowRangeList {
	var list bigtable.RowRangeList

	switch {
	case len(sampleRowKeys) == 1:
		return append(list, bigtable.InfiniteRange(""))
	case iteration == 0:
		return append(list, bigtable.NewRange("", rowKey))
	case len(sampleRowKeys)-1 == iteration:
		return append(list, bigtable.NewRange(rowKey, ""))
	default:
		return append(list, bigtable.NewRange(rowKey, sampleRowKeys[iteration+1]))
	}
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
