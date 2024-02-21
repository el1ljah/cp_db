package models

type Item struct {
	ID          int    `valid:"-" json:"id" db:"id"`
	Category    string `valid:"in(ботинки|кроссовки|майка|футболка|куртка|штаны|шорты|ремень|шляпа)" json:"category" db:"category"`
	Size        string `valid:"in(XS|S|M|L|XL|XXL)" json:"size" db:"size"`
	Price       int    `valid:"-" json:"price" db:"price"`
	Sex         string `valid:"in(male|female)" json:"sex" db:"sex"`
	ImageID     int    `valid:"-" json:"image_id" db:"image_id"`
	BrandID     int    `valid:"-" json:"brand_id" db:"brand_id"`
	IsAvailable bool   `valid:"-" json:"is_available" db:"is_available"`
}

const (
	ItemsParamsAny = "any"
	ItemsOrderDesc = "desc"
	ItemsOrderAsc  = "asc"
)

type ItemsPatchPrice struct {
	NewPrice       int    `valid:"-" json:"price" db:"price"`
}

type ItemsParams struct {
	WhereCategory string `valid:"in(ботинки|кроссовки|майка|футболка|куртка|штаны|шорты|ремень|шляпа|any)" json:"WhereCategory" schema:"WhereCategory" example:"ботинки|кроссовки|майка|футболка|куртка|штаны|шорты|ремень|шляпа|any"`
	WhereSex      string `valid:"in(male|female|any)" json:"WhereSex" schema:"WhereSex" example:"male|female|any"`
	WhereBrand    int    `valid:"-" json:"WhereBrand" schema:"WhereBrand" example:1`
	OrderBy    string `valid:"in(asc|desc|any)" json:"OrderBy" schema:"OrderBy" example:"asc|desc|any"`
	Page_size int	`valid:"-" json:"Page_size" schema:"Page_size" example:50`
	Page_num int	`valid:"-" json:"Page_num"  schema:"Page_num" example:1`
}
