/*
 * @Author: Hugo
 * @Date: 2022-05-12 11:17:37
 * @Last Modified by: Hugo
 * @Last Modified time: 2022-05-12 11:23:14
 */
package mambuEntity

type SearchParam struct {
	Sortingcriteria Sortingcriteria  `json:"sortingCriteria"`
	Filtercriteria  []Filtercriteria `json:"filterCriteria"`
}
type Sortingcriteria struct {
	Field string `json:"field"`
	Order string `json:"order"`
}
type Filtercriteria struct {
	Field       string   `json:"field"`
	Operator    string   `json:"operator"`
	Value       string   `json:"value,omitempty"`
	Values      []string `json:"values,omitempty"`
	Secondvalue string   `json:"secondValue,omitempty"`
}
