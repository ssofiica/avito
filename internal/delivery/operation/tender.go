package operation

import "strconv"

type TenderListParams struct {
	Limit       uint32
	Offset      uint32
	ServiceType []string
}

func (res TenderListParams) Scan(limit string, offset string, service string) error {
	l, err := strconv.Atoi(limit)
	if err != nil && limit != "" {
		return err
	}
	res.Limit = uint32(l)
	o, err := strconv.Atoi(offset)
	if err != nil && offset != "" {
		return err
	}
	res.Offset = uint32(o)
	return nil
}
