package repository

import (
	"context"

	"github.com/linximanager/backend/internal/model"
	"gorm.io/gorm"
)

type ProductRepo struct {
	db *gorm.DB
}

func NewProductRepo(db *gorm.DB) *ProductRepo {
	return &ProductRepo{db: db}
}

type ListProductsFilter struct {
	MerchantID int64
	CategoryID *int64
	Status     string
	Keyword    string
	Offset     int
	Limit      int
}

func (r *ProductRepo) List(ctx context.Context, f ListProductsFilter) ([]*model.Product, int64, error) {
	q := r.db.WithContext(ctx).Model(&model.Product{}).
		Where("merchant_id = ?", f.MerchantID)
	if f.Status != "" {
		q = q.Where("status = ?", f.Status)
	} else {
		q = q.Where("status != 'deleted'")
	}
	if f.CategoryID != nil {
		q = q.Where("category_id = ?", *f.CategoryID)
	}
	if f.Keyword != "" {
		q = q.Where("name ILIKE ?", "%"+f.Keyword+"%")
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var products []*model.Product
	err := q.Offset(f.Offset).Limit(f.Limit).Order("sort_order ASC, id ASC").Find(&products).Error
	return products, total, err
}

func (r *ProductRepo) FindByID(ctx context.Context, id, merchantID int64) (*model.Product, error) {
	var p model.Product
	err := r.db.WithContext(ctx).
		Where("id = ? AND merchant_id = ? AND status != 'deleted'", id, merchantID).
		First(&p).Error
	return &p, err
}

func (r *ProductRepo) FindPublicByID(ctx context.Context, id int64) (*model.Product, error) {
	var p model.Product
	err := r.db.WithContext(ctx).
		Where("id = ? AND status = 'active'", id).
		First(&p).Error
	return &p, err
}

func (r *ProductRepo) Create(ctx context.Context, p *model.Product) error {
	return r.db.WithContext(ctx).Create(p).Error
}

func (r *ProductRepo) Update(ctx context.Context, p *model.Product) error {
	return r.db.WithContext(ctx).Save(p).Error
}

func (r *ProductRepo) SoftDelete(ctx context.Context, id, merchantID int64) error {
	return r.db.WithContext(ctx).Model(&model.Product{}).
		Where("id = ? AND merchant_id = ?", id, merchantID).
		Update("status", "deleted").Error
}

func (r *ProductRepo) GetBOM(ctx context.Context, productID int64) ([]*model.ProductBOM, error) {
	var bom []*model.ProductBOM
	err := r.db.WithContext(ctx).
		Preload("Material").
		Where("product_id = ?", productID).
		Find(&bom).Error
	return bom, err
}

func (r *ProductRepo) ReplaceBOM(ctx context.Context, tx *gorm.DB, productID int64, bom []*model.ProductBOM) error {
	if err := tx.Where("product_id = ?", productID).Delete(&model.ProductBOM{}).Error; err != nil {
		return err
	}
	if len(bom) == 0 {
		return nil
	}
	return tx.Create(&bom).Error
}

func (r *ProductRepo) ListCategories(ctx context.Context, merchantID int64) ([]*model.ProductCategory, error) {
	var cats []*model.ProductCategory
	err := r.db.WithContext(ctx).
		Where("merchant_id = ?", merchantID).
		Order("sort_order ASC, id ASC").
		Find(&cats).Error
	return cats, err
}

func (r *ProductRepo) DB() *gorm.DB {
	return r.db
}
