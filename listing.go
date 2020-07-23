package listings

import (
	"fmt"

	"github.com/timbertom-gmbh/go-listings/rpc"
)

const (
	defaultPerPage uint32 = 50
)

type List func(DB, *rpc.ListingOptions) DB

type DB interface {
	Offset(interface{}) DB
	Limit(interface{}) DB
	Where(query string, args ...interface{}) DB
	Order(value string) DB
}

func NewListCreator() List {
	return func(db DB, request *rpc.ListingOptions) DB {
		perPage := request.GetPerPage()
		if perPage <= 0 {
			perPage = defaultPerPage
		}
		page := request.GetPage()
		query := db.Offset(page * perPage).Limit(perPage)

		for _, filter := range request.GetFilters() {
			query = query.Where(fmt.Sprintf("%s = ?", filter.GetField()), filter.GetQuery())
		}

		return query.Order(fmt.Sprintf("%s %s", request.GetSortField(), sortOrderString(request.GetSortOrder())))
	}
}

func sortOrderString(rpcSortOrder rpc.ListingOptions_Order) string {
	switch rpcSortOrder {
	case rpc.ListingOptions_ASC:
		return "ASC"
	case rpc.ListingOptions_DESC:
		return "DESC"
	default:
		panic(fmt.Errorf("unknown listing enum"))
	}
}
