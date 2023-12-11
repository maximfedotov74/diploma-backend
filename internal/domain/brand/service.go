package brand

import (
	exception "github.com/maximfedotov74/fiber-psql/internal/shared/error"
	"github.com/maximfedotov74/fiber-psql/internal/shared/utils"
)

type Repository interface {
	CreateBrand(title string, slug string, description *string, imgPath *string) exception.Error
	FindByFeild(field string, value any) (*Brand, exception.Error)
	GetAll() ([]Brand, exception.Error)
	UpdateBrand(dto UpdateBrandDto, newSlug *string, id int) exception.Error
}

type BrandService struct {
	repo Repository
}

func NewBrandService(repo Repository) *BrandService {
	return &BrandService{repo: repo}
}

func (bs *BrandService) GetAll() ([]Brand, exception.Error) {
	brands, err := bs.repo.GetAll()
	if err != nil {
		return nil, err
	}
	return brands, nil
}

func (bs *BrandService) FindByTitle(title string) (*Brand, exception.Error) {
	brandExist, err := bs.repo.FindByFeild("title", title)

	if err != nil {
		return nil, err
	}

	return brandExist, nil
}

func (bs *BrandService) FindById(id int) (*Brand, exception.Error) {
	brandExist, err := bs.repo.FindByFeild("brand_id", id)

	if err != nil {
		return nil, err
	}

	return brandExist, nil
}

func (bs *BrandService) FindBySlug(slug string) (*Brand, exception.Error) {
	brandExist, err := bs.repo.FindByFeild("slug", slug)

	if err != nil {
		return nil, err
	}

	return brandExist, nil
}

func (bs *BrandService) CreateBrand(dto CreateBrandDto) exception.Error {
	brandExist, _ := bs.FindByTitle(dto.Title)

	if brandExist != nil {
		return exception.NewErr(brandExists, exception.STATUS_BAD_REQUEST)
	}
	slug := utils.GenerateSlug(dto.Title)

	ex := bs.repo.CreateBrand(dto.Title, slug, dto.Description, dto.ImgPath)
	if ex != nil {
		return ex
	}
	return nil
}

func (bs *BrandService) UpdateBrand(dto UpdateBrandDto, id int) exception.Error {

	_, err := bs.FindById(id)

	if err != nil {
		return err
	}

	var slug *string

	if dto.Title != nil {

		exist, _ := bs.FindByTitle(*dto.Title)

		if exist != nil {
			return exception.NewErr(brandTitleUnique, 400)
		}

		newSlug := utils.GenerateSlug(*dto.Title)
		slug = &newSlug
	}

	err = bs.repo.UpdateBrand(dto, slug, id)

	if err != nil {
		return err
	}

	return nil

}
