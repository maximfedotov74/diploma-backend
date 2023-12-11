package option

import exception "github.com/maximfedotov74/fiber-psql/internal/shared/error"

type Repository interface {
	CreateValue(dto CreateOptionValueDto) exception.Error
	GetById(id int) (*Option, exception.Error)
	UpdateOption(dto UpdateOptionDto, id int) exception.Error
	CreateOption(dto CreateOptionDto) exception.Error
	DeleteValue(id int) exception.Error
	DeleteOption(id int) exception.Error
	AddToProductModel(dto AddOptionToProductModelDto) exception.Error
	AddSizeToProductModel(dto AddSizeToProductModelDto) exception.Error
	CreateSize(dto CreateSizeDto) exception.Error
	GetCatalogFilters(categorySlug string) (*CatalogFilters, exception.Error)
	GetAll() ([]Option, exception.Error)
}

type OptionService struct {
	repo Repository
}

func NewOptionService(repo Repository) *OptionService {
	return &OptionService{repo: repo}
}

func (os *OptionService) GetAll() ([]Option, exception.Error) {
	opts, err := os.repo.GetAll()

	if err != nil {
		return nil, err
	}

	return opts, nil
}

func (os *OptionService) CreateOption(dto CreateOptionDto) exception.Error {
	err := os.repo.CreateOption(dto)
	if err != nil {
		return err
	}

	return nil
}

func (os *OptionService) CreateValue(dto CreateOptionValueDto) exception.Error {
	err := os.repo.CreateValue(dto)
	if err != nil {
		return err
	}

	return nil
}

func (os *OptionService) UpdateOption(dto UpdateOptionDto, id int) exception.Error {
	err := os.repo.UpdateOption(dto, id)
	if err != nil {
		return err
	}

	return nil
}

func (os *OptionService) GetById(id int) (*Option, exception.Error) {
	opt, err := os.repo.GetById(id)

	if err != nil {
		return nil, err
	}

	return opt, nil
}

func (os *OptionService) DeleteOption(id int) exception.Error {
	err := os.repo.DeleteOption(id)
	if err != nil {
		return err
	}
	return nil
}

func (os *OptionService) DeleteValue(id int) exception.Error {
	err := os.repo.DeleteValue(id)
	if err != nil {
		return err
	}
	return nil
}

func (os *OptionService) AddOptionToProductModel(dto AddOptionToProductModelDto) exception.Error {
	err := os.repo.AddToProductModel(dto)
	if err != nil {
		return err
	}
	return nil
}

func (os *OptionService) CreateSize(dto CreateSizeDto) exception.Error {
	err := os.repo.CreateSize(dto)
	if err != nil {
		return err
	}

	return nil
}

func (os *OptionService) AddSizeToProductModel(dto AddSizeToProductModelDto) exception.Error {
	err := os.repo.AddSizeToProductModel(dto)
	if err != nil {
		return err
	}
	return nil
}

func (os *OptionService) GetCatalogFilters(categorySlug string) (*CatalogFilters, exception.Error) {
	filters, err := os.repo.GetCatalogFilters(categorySlug)
	if err != nil {
		return nil, err
	}

	return filters, nil
}
