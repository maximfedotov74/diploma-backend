package generator

import (
	"fmt"
	"strconv"
	"strings"
)

type CatalogFilters struct {
	Options          map[string]string
	Slug             string
	Sizes            string
	Brands           string
	SortBy           string
	OnlyWithDiscount string
	Price            string
	Page             string
}

type GeneratedCatalogQuery struct {
	SortStatement string
	Pagination    string
	MainQuery     string
}

func GenerateCatalogQuery(filters CatalogFilters) GeneratedCatalogQuery {

	var isWhereStatement bool = false

	optionsJoins := ""
	optionsWhere := ""
	sizeWhere := ""
	sortStatement := ""
	pagination := ""
	brandsWhere := ""
	priceWhere := ""
	onlyWithDiscountWhere := ""

	if len(filters.Options) > 0 {

		optionsWhereStatement := ""

		optIdx := 1

		for optionSlug, v := range filters.Options {
			join := fmt.Sprintf(`
      inner join product_model_option as pmop%[1]d on pmop%[1]d.product_model_id = pm.product_model_id
      inner join option as op%[1]d on op%[1]d.option_id = pmop%[1]d.option_id
      inner join option_value as v%[1]d on v%[1]d.option_value_id = pmop%[1]d.option_value_id`, optIdx)
			optionsJoins += join
			filterValues := strings.Split(v, ",")
			where := fmt.Sprintf("op%d.slug = '%s' and v%d.option_value_id IN ", optIdx, optionSlug, optIdx)

			if !isWhereStatement {
				where = " WHERE " + where
				isWhereStatement = true
			} else {
				where = " AND " + where
			}
			idsArr := make([]string, 0, len(filterValues))
			for _, optionValue := range filterValues {
				valueId, err := strconv.Atoi(optionValue)
				if err != nil {
					continue
				}
				idStr := fmt.Sprintf("%d", valueId)
				idsArr = append(idsArr, idStr)
			}
			inStatement := fmt.Sprintf("(%s)", strings.Join(idsArr, ","))
			where += inStatement
			optionsWhereStatement += where + " "
			optIdx++
		}
		optionsWhere = optionsWhereStatement
	}

	if filters.Sizes != "" {
		sizesValues := strings.Split(filters.Sizes, ",")
		idsArr := make([]string, 0, len(sizesValues))
		for _, size := range sizesValues {
			idStr := fmt.Sprintf("'%s'", size)
			idsArr = append(idsArr, idStr)
		}
		where := "sz.size_value IN "
		if !isWhereStatement {
			where = " WHERE " + where
			isWhereStatement = true
		} else {
			where = " AND " + where
		}
		inStatement := fmt.Sprintf("(%s)", strings.Join(idsArr, ","))
		where += inStatement
		sizeWhere = where
	}

	if filters.Brands != "" {
		brandsIds := strings.Split(filters.Brands, ",")
		idsArr := make([]string, 0, len(brandsIds))
		for _, brandId := range brandsIds {
			idStr := fmt.Sprintf("%s", brandId)
			idsArr = append(idsArr, idStr)
		}
		where := "b.brand_id IN "
		if !isWhereStatement {
			where = " WHERE " + where
			isWhereStatement = true
		} else {
			where = " AND " + where
		}
		inStatement := fmt.Sprintf("(%s)", strings.Join(idsArr, ","))
		where += inStatement
		brandsWhere = where
	}

	if filters.Price != "" {
		limits := strings.Split(filters.Price, ",")
		minValue, minParseErr := strconv.ParseFloat(limits[0], 32)
		maxValue, maxParseErr := strconv.ParseFloat(limits[1], 32)
		if maxParseErr == nil && minParseErr == nil {

			if !isWhereStatement {
				priceWhere = fmt.Sprintf(" WHERE pm.price BETWEEN %.2f AND %.2f", minValue, maxValue)
				isWhereStatement = true
			} else {
				priceWhere = fmt.Sprintf(" AND pm.price BETWEEN %.2f AND %.2f", minValue, maxValue)
			}
		}
	}

	if filters.OnlyWithDiscount != "" && filters.OnlyWithDiscount == "1" {
		if !isWhereStatement {
			onlyWithDiscountWhere = " WHERE pm.discount != NULL"
			isWhereStatement = true
		} else {
			onlyWithDiscountWhere = " AND pm.discount != NULL"
		}
	}

	if filters.SortBy != "" {
		switch filters.SortBy {
		case "price_asc":
			sortStatement = " ORDER BY pm.price ASC"
			break
		case "price_desc":
			sortStatement = " ORDER BY pm.price DESC"
			break
		case "discount":
			sortStatement = " ORDER BY pm.discount DESC"
			break
		default:
			sortStatement = ""
		}
	}

	limit := 16
	page, err := strconv.Atoi(filters.Page)
	if err != nil {
		page = 1
	}
	offset := page*limit - limit

	pagination = fmt.Sprintf("LIMIT %d OFFSET %d", limit, offset)

	return GeneratedCatalogQuery{
		SortStatement: sortStatement,
		MainQuery:     optionsJoins + optionsWhere + sizeWhere + brandsWhere + priceWhere + onlyWithDiscountWhere,
		Pagination:    pagination,
	}
}
