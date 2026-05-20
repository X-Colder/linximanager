package service

import (
	"context"

	"github.com/linximanager/backend/internal/model"
	"github.com/linximanager/backend/internal/repository"
)

type ProductService struct {
	productRepo *repository.ProductRepo
}

func NewProductService(productRepo *repository.ProductRepo) *ProductService {
	return &ProductService{productRepo: productRepo}
}

func (s *ProductService) List(ctx context.Context, f repository.ListProductsFilter) ([]*model.Product, int64, error) {
	return s.productRepo.List(ctx, f)
}

func (s *ProductService) GetByID(ctx context.Context, id, merchantID int64) (*model.Product, error) {
	return s.productRepo.FindByID(ctx, id, merchantID)
}

func (s *ProductService) GetPublicByID(ctx context.Context, id int64) (*model.Product, error) {
	return s.productRepo.FindPublicByID(ctx, id)
}

func (s *ProductService) Create(ctx context.Context, p *model.Product) error {
	return s.productRepo.Create(ctx, p)
}

func (s *ProductService) Update(ctx context.Context, p *model.Product) error {
	return s.productRepo.Update(ctx, p)
}

func (s *ProductService) Delete(ctx context.Context, id, merchantID int64) error {
	return s.productRepo.SoftDelete(ctx, id, merchantID)
}

func (s *ProductService) GetBOM(ctx context.Context, productID int64) ([]*model.ProductBOM, error) {
	return s.productRepo.GetBOM(ctx, productID)
}

func (s *ProductService) ReplaceBOM(ctx context.Context, productID int64, bom []*model.ProductBOM) error {
	db := s.productRepo.DB()
	tx := db.Begin()
	if err := s.productRepo.ReplaceBOM(ctx, tx, productID, bom); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

func (s *ProductService) ListCategories(ctx context.Context, merchantID int64) ([]*model.ProductCategory, error) {
	return s.productRepo.ListCategories(ctx, merchantID)
}
