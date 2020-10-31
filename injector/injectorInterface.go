package injector

import "github.com/newity/crawler/parser"

type Injector interface {
	Inject(data *parser.Data) error
}
