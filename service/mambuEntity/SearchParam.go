/*
 * @Author: Hugo
 * @Date: 2022-05-12 11:17:37
 * @Last Modified by: Hugo
 * @Last Modified time: 2022-05-20 03:58:11
 */
package mambuEntity

type SearchParam struct {
	SortingCriteria SortingCriteria  `json:"sortingCriteria"`
	FilterCriteria  []FilterCriteria `json:"filterCriteria"`
}
type SortingCriteria struct {
	Field string `json:"field"`
	Order string `json:"order"`
}
type FilterCriteria struct {
	Field       string   `json:"field"`
	Operator    string   `json:"operator"`
	Value       string   `json:"value,omitempty"`
	Values      []string `json:"values,omitempty"`
	SecondValue string   `json:"secondValue,omitempty"`
}
