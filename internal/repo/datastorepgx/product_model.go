package datastorepgx

import (
	"api/internal/model"

	"github.com/google/uuid"
	"github.com/jackc/pgtype"
)

type Product struct {
	ID                pgtype.UUID       `db:"id"`
	Brand             *string           `db:"brand"`
	Name              *string           `db:"name"`
	ShortDescription  *string           `db:"short_description"`
	Description       *string           `db:"description"`
	Size              *string           `db:"size"`
	Price             *float64          `db:"price"`
	ProductTypeID     *pgtype.UUID      `db:"product_type_id"`
	ProductCategoryID *pgtype.UUID      `db:"product_category_id"`
	SkinTypeID        *pgtype.UUID      `db:"skin_type_id"`
	CountryID         *pgtype.UUID      `db:"country_id"`
	Tags              *pgtype.TextArray `db:"tags"`
	CreatedTime       pgtype.Time       `db:"created_time"`
	DeletedTime       pgtype.Time       `db:"deleted_time"`

	Images []ProductImage `db:"-"`
}

type ProductScan struct {
	ID                pgtype.UUID      `db:"id"`
	Brand             pgtype.Varchar   `db:"brand"`
	Name              pgtype.Varchar   `db:"name"`
	ShortDescription  pgtype.Text      `db:"short_description"`
	Description       pgtype.Text      `db:"description"`
	Size              pgtype.Varchar   `db:"size"`
	Price             pgtype.Float8    `db:"price"`
	ProductTypeID     pgtype.UUID      `db:"product_type_id"`
	ProductCategoryID pgtype.UUID      `db:"product_category_id"`
	SkinTypeID        pgtype.UUID      `db:"skin_type_id"`
	CountryID         pgtype.UUID      `db:"country_id"`
	Tags              pgtype.TextArray `db:"tags"`
	CreatedTime       pgtype.Time      `db:"created_time"`
	DeletedTime       pgtype.Time      `db:"deleted_time"`

	Images []ProductImage `db:"-"`
}

type ProductImage struct {
	ID          pgtype.UUID    `db:"id"`
	ProductID   pgtype.UUID    `db:"product_id"`
	Path        pgtype.Varchar `db:"path"`
	CreatedTime pgtype.Time    `db:"created_time"`
}

func (p *ProductScan) toModel() *model.Product {
	rs := model.Product{
		ID:               p.ID.Bytes,
		Brand:            &p.Brand.String,
		Name:             &p.Name.String,
		ShortDescription: &p.ShortDescription.String,
		Description:      &p.ShortDescription.String,
		Size:             &p.Size.String,
		Price:            &p.Price.Float,
		Tags:             textArrayToSlice(p.Tags),
		ProductTypeID:    (*uuid.UUID)(p.ProductTypeID.Bytes[:]),
	}

	rs.Images = make([]model.ProductImage, 0, len(p.Images))
	for _, image := range p.Images {
		rs.Images = append(rs.Images, *image.toModel())
	}

	return &rs
}

func (p *ProductImage) toModel() *model.ProductImage {
	return &model.ProductImage{
		ID:        p.ID.Bytes,
		ProductID: p.ProductID.Bytes,
		Path:      p.Path.String,
	}
}

func productModelToRow(p model.Product) Product {
	row := Product{
		ID: pgtype.UUID{
			Bytes:  p.ID,
			Status: pgtype.Present,
		},
		Brand:            p.Brand,
		Name:             p.Name,
		ShortDescription: p.ShortDescription,
		Description:      p.Description,
		Size:             p.Size,
		Price:            p.Price,
	}

	if p.ID == uuid.Nil {
		row.ID.Status = pgtype.Null
	}

	if p.ProductTypeID != nil {
		row.ProductTypeID = &pgtype.UUID{
			Bytes:  *p.ProductTypeID,
			Status: pgtype.Present,
		}

		if *p.ProductTypeID == uuid.Nil {
			row.ProductTypeID.Status = pgtype.Null
		}
	}

	if p.ProductCategoryID != nil {
		row.ProductCategoryID = &pgtype.UUID{
			Bytes:  *p.ProductCategoryID,
			Status: pgtype.Present,
		}

		if *p.ProductCategoryID == uuid.Nil {
			row.ProductCategoryID.Status = pgtype.Null
		}
	}

	if p.SkinTypeID != nil {
		row.SkinTypeID = &pgtype.UUID{
			Bytes:  *p.SkinTypeID,
			Status: pgtype.Present,
		}

		if *p.SkinTypeID == uuid.Nil {
			row.SkinTypeID.Status = pgtype.Null
		}
	}

	if p.CountryID != nil {
		row.CountryID = &pgtype.UUID{
			Bytes:  *p.CountryID,
			Status: pgtype.Present,
		}

		if *p.CountryID == uuid.Nil {
			row.CountryID.Status = pgtype.Null
		}
	}

	if p.Tags != nil {
		tags := pgtype.TextArray{}
		if err := tags.Set(p.Tags); err == nil {
			row.Tags = &tags
		}
	}

	return row
}

func productImageModelToRow(p model.ProductImage) ProductImage {
	row := ProductImage{
		ID:        pgtype.UUID{Bytes: p.ID, Status: pgtype.Present},
		ProductID: pgtype.UUID{Bytes: p.ProductID, Status: pgtype.Present},
		Path:      pgtype.Varchar{String: p.Path, Status: pgtype.Present},
	}

	if p.ID == uuid.Nil {
		row.ID.Status = pgtype.Null
	}

	if p.ProductID == uuid.Nil {
		row.ProductID.Status = pgtype.Null
	}

	if p.Path == "" {
		row.Path.Status = pgtype.Null
	}

	return row
}
