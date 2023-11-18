package brand

import (
	exception "github.com/maximfedotov74/fiber-psql/internal/shared/error"
	"github.com/maximfedotov74/fiber-psql/internal/shared/messages"
	"github.com/maximfedotov74/fiber-psql/internal/shared/utils"
)

type Repository interface {
	CreateBrand(title string, slug string, description *string, imgPath *string) error
	FindByFeild(field string, value any) (*Brand, error)
}

type BrandService struct {
	repo Repository
}

func NewBrandService(repo Repository) *BrandService {
	return &BrandService{repo: repo}
}

func (bs *BrandService) FindByTitle(title string) (*Brand, exception.Error) {
	brandExist, err := bs.repo.FindByFeild("title", title)

	if err != nil {
		return nil, exception.NewErr(err.Error(), 500)
	}

	return brandExist, nil
}

func (bs *BrandService) FindById(id int) (*Brand, exception.Error) {
	brandExist, err := bs.repo.FindByFeild("brand_id", id)

	if err != nil {
		return nil, exception.NewErr(err.Error(), 500)
	}

	return brandExist, nil
}

func (bs *BrandService) FindBySlug(slug string) (*Brand, exception.Error) {
	brandExist, err := bs.repo.FindByFeild("slug", slug)

	if err != nil {
		return nil, exception.NewErr(err.Error(), 500)
	}

	return brandExist, nil
}

func (bs *BrandService) CreateBrand(dto CreateBrandDto) exception.Error {
	brandExist, _ := bs.FindByTitle(dto.Title)

	if brandExist != nil {
		return exception.NewErr(messages.BRAND_EXISTS, 400)
	}
	slug := utils.GenerateSlug(dto.Title)
	ex := bs.repo.CreateBrand(dto.Title, slug, dto.Description, dto.ImgPath)
	if ex != nil {
		return exception.NewErr(messages.BRAND_CREATE_ERROR, 500)
	}
	return nil
}
