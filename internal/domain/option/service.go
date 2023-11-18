package option

import exception "github.com/maximfedotov74/fiber-psql/internal/shared/error"

type Repository interface {
	CreateValue(dto CreateOptionValueDto) error
	GetById(id int) (*Option, error)
	UpdateOption(dto UpdateOptionDto, id int) error
	CreateOption(dto CreateOptionDto) error
	DeleteValue(id int) error
	DeleteOption(id int) error
	AddToProductModel(dto AddOptionToProductModelDto) error
}

type OptionService struct {
	repo Repository
}

func NewOptionService(repo Repository) *OptionService {
	return &OptionService{repo: repo}
}

func (os *OptionService) CreateOption(dto CreateOptionDto) exception.Error {
	err := os.repo.CreateOption(dto)
	if err != nil {
		return exception.NewErr(err.Error(), 500)
	}

	return nil
}

func (os *OptionService) CreateValue(dto CreateOptionValueDto) exception.Error {
	err := os.repo.CreateValue(dto)
	if err != nil {
		return exception.NewErr(err.Error(), 500)
	}

	return nil
}

func (os *OptionService) UpdateOption(dto UpdateOptionDto, id int) exception.Error {
	err := os.repo.UpdateOption(dto, id)
	if err != nil {
		return exception.NewErr(err.Error(), 500)
	}

	return nil
}

func (os *OptionService) GetById(id int) (*Option, exception.Error) {
	opt, err := os.repo.GetById(id)

	if err != nil {
		return nil, exception.NewErr(err.Error(), 404)
	}

	return opt, nil
}

func (os *OptionService) DeleteOption(id int) exception.Error {
	err := os.repo.DeleteOption(id)
	if err != nil {
		return exception.NewErr(err.Error(), 500)
	}
	return nil
}

func (os *OptionService) DeleteValue(id int) exception.Error {
	err := os.repo.DeleteValue(id)
	if err != nil {
		return exception.NewErr(err.Error(), 500)
	}
	return nil
}

func (os *OptionService) AddOptionToProductModel(dto AddOptionToProductModelDto) exception.Error {
	err := os.repo.AddToProductModel(dto)
	if err != nil {
		return exception.NewErr(err.Error(), 500)
	}
	return nil
}
