package operation

import "strconv"

type BidParams struct {
	Limit  int
	Offset int
}

func (res *BidParams) Scan(limit string, offset string) error {
	l, err := strconv.Atoi(limit)
	if err != nil && limit != "" {
		return err
	}
	res.Limit = l
	o, err := strconv.Atoi(offset)
	if err != nil && offset != "" {
		return err
	}
	res.Offset = o
	return nil
}
