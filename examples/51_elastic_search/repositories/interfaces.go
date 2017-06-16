package repositories

import "golang-learn/examples/51_elastic_search/domain"

type IRepository interface {
	CreateIndex(index string, mapping map[string]interface{}) error
	Insert(data []domain.Something) error
}