package data

import (
	"github.com/acs-dl/unverified-svc/resources"
	"gitlab.com/distributed_lab/kit/pgdb"
	"net/http"
	"strconv"
)

const (
	pageParamLimit  = "page[limit]"
	pageParamNumber = "page[number]"
	pageParamOrder  = "page[order]"
	pageParamSort   = "page[order]"
)

type SortParams struct {
	Param string `page:"sort" default:"created_at"`
}

func GetOffsetLinks(r *http.Request, p pgdb.OffsetPageParams, s SortParams) *resources.Links {
	result := resources.Links{
		Next: getOffsetLink(r, p.PageNumber+1, p.Limit, p.Order, s.Param),
		Self: getOffsetLink(r, p.PageNumber, p.Limit, p.Order, s.Param),
	}

	return &result
}

func getOffsetLink(r *http.Request, pageNumber, limit uint64, order, sort string) string {
	u := r.URL
	query := u.Query()
	query.Set(pageParamNumber, strconv.FormatUint(pageNumber, 10))
	query.Set(pageParamLimit, strconv.FormatUint(limit, 10))
	query.Set(pageParamSort, sort)
	query.Set(pageParamOrder, order)
	u.RawQuery = query.Encode()
	return u.String()
}
