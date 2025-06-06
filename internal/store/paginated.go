package store

import (
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
)

type Paginated struct {
	LIMIT  int64  `json:"limit"`
	OFFSET int64  `json:"offset"`
	SORT   string `json:"sort" validate:"oneof=des asc"`
}

var Validate *validator.Validate

func init() {
	Validate = validator.New(
		validator.WithRequiredStructEnabled(),
	)
}

func (p *Paginated) Parse(r *http.Request) error {
	p.LIMIT = 20
	p.OFFSET = 0
	p.SORT = "asc"

	values := r.URL.Query()
	limit, err := strconv.ParseInt(values.Get("limit"), 10, 16)

	if err == nil {
		p.LIMIT = limit
	}
	offset, err := strconv.ParseInt(values.Get("offset"), 10, 16)
	if err == nil {
		p.OFFSET = offset
	}
	sort := values.Get("sort")
	println(sort, p.SORT)
	if len(sort) > 0 {
		p.SORT = sort
	}
	err = Validate.Struct(p)
	if err != nil {
		return err
	}
	return nil
}
