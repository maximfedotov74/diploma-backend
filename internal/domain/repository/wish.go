package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/maximfedotov74/diploma-backend/internal/domain/model"
	"github.com/maximfedotov74/diploma-backend/internal/domain/msg"
	"github.com/maximfedotov74/diploma-backend/internal/shared/db"
	"github.com/maximfedotov74/diploma-backend/internal/shared/fall"
)

type WishRepository struct {
	db db.PostgresClient
}

func NewWishRepository(db db.PostgresClient) *WishRepository {
	return &WishRepository{db: db}
}

func (r *WishRepository) FindModelInUserCart(ctx context.Context, modelSizeId int, userId int) (*model.CartItemModel, fall.Error) {

	query := `
	SELECT cart.cart_id,cart.user_id,cart.model_size_id,cart.quantity, ms.in_stock,
	pm.product_model_id, pm.price, pm.discount
	FROM cart
	INNER JOIN model_sizes as ms ON cart.model_size_id = ms.model_size_id
	INNER JOIN product_model as pm ON pm.product_model_id = ms.product_model_id
	WHERE cart.user_id = $1 AND cart.model_size_id = $2;
  `
	row := r.db.QueryRow(ctx, query, userId, modelSizeId)

	cartItem := model.CartItemModel{}

	err := row.Scan(&cartItem.CartItemId, &cartItem.UserId, &cartItem.ModelSizeId, &cartItem.Quantity, &cartItem.InStock,
		&cartItem.ModelId, &cartItem.Price, &cartItem.Discount,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fall.NewErr(msg.CartItemNotFound, fall.STATUS_NOT_FOUND)
		}
		return nil, fall.ServerError(err.Error())
	}

	return &cartItem, nil
}

func (r *WishRepository) AddToCart(ctx context.Context, modelSizeId int, userId int) fall.Error {

	query := "INSERT INTO cart (user_id, model_size_id) VALUES ($1, $2);"

	_, err := r.db.Exec(ctx, query, userId, modelSizeId)

	if err != nil {
		return fall.ServerError(err.Error())
	}

	return nil
}

func (r *WishRepository) DeleteFromCart(ctx context.Context, cartItemId int) fall.Error {
	query := "DELETE FROM cart WHERE cart_id = $1;"

	_, err := r.db.Exec(ctx, query, cartItemId)

	if err != nil {
		return fall.ServerError(err.Error())
	}

	return nil
}

func (r *WishRepository) UpdateCartItem(ctx context.Context, cartItemId int, newQuantity int) fall.Error {

	query := `UPDATE cart SET quantity = $1 WHERE cart_id = $2;`

	_, err := r.db.Exec(ctx, query, newQuantity, cartItemId)

	if err != nil {
		return fall.ServerError(err.Error())
	}
	return nil
}

func (r *WishRepository) RemoveSeveralItems(ctx context.Context, tx db.Transaction, cartIds []int) fall.Error {

	query := `DELETE FROM cart WHERE cart_id = ANY ($1);`

	if tx != nil {
		_, err := tx.Exec(ctx, query, cartIds)
		if err != nil {
			return fall.ServerError(err.Error())
		}
		return nil
	}

	_, err := r.db.Exec(ctx, query, cartIds)

	if err != nil {
		return fall.ServerError(err.Error())
	}

	return nil
}

func (r *WishRepository) GetUserCart(ctx context.Context, userId int) ([]model.CartItem, fall.Error) {
	query := `
	SELECT cr.cart_id as cr_id, cr.quantity as cr_q,
	ms.model_size_id as ms_id,
	ms.literal_size as ms_ls, ms.in_stock as ms_in_stock, sz.size_id as sz_id, sz.size_value as sz_value,
	pm.product_model_id as pm_id, pm.price as price, pm.discount as discount, pm.main_image_path as pm_img,
	p.product_id as p_id, p.title as p_title, p.slug as p_slug,
	ct.category_id as ct_id, ct.title as ct_title, ct.short_title as ct_short_title, ct.slug as ct_slug,
	b.brand_id as b_id, b.title as b_title, b.slug as b_slug
	FROM cart as cr
	INNER JOIN model_sizes as ms ON ms.model_size_id = cr.model_size_id
	INNER JOIN sizes as sz ON ms.size_id = sz.size_id
	INNER JOIN product_model as pm ON ms.product_model_id = pm.product_model_id
	INNER JOIN product as p ON pm.product_id = p.product_id
	INNER JOIN category ct ON p.category_id = ct.category_id
	INNER JOIN brand b on p.brand_id = b.brand_id
	WHERE cr.user_id = $1;
	`
	rows, err := r.db.Query(ctx, query, userId)
	if err != nil {
		return nil, fall.ServerError(err.Error())
	}
	defer rows.Close()
	var result []model.CartItem
	for rows.Next() {
		item := model.CartItem{}
		err := rows.Scan(&item.Id, &item.Quantity, &item.ModelSize.Id, &item.ModelSize.LiteralSize, &item.ModelSize.InStock,
			&item.ModelSize.SizeId, &item.ModelSize.SizeValue, &item.ModelSize.ProductModel.Id, &item.ModelSize.ProductModel.Price,
			&item.ModelSize.ProductModel.Discount, &item.ModelSize.ProductModel.ImagePath, &item.ModelSize.ProductModel.ProductId,
			&item.ModelSize.ProductModel.Title, &item.ModelSize.ProductModel.Slug,
			&item.ModelSize.ProductModel.Category.Id, &item.ModelSize.ProductModel.Category.Title,
			&item.ModelSize.ProductModel.Category.ShortTitle, &item.ModelSize.ProductModel.Category.Slug,
			&item.ModelSize.ProductModel.Brand.Id, &item.ModelSize.ProductModel.Brand.Title,
			&item.ModelSize.ProductModel.Brand.Slug,
		)
		if err != nil {
			return nil, fall.ServerError(err.Error())
		}
		result = append(result, item)
	}

	if err := rows.Err(); err != nil {
		return nil, fall.ServerError(err.Error())
	}

	return result, nil

}

func (r *WishRepository) AddToWish(ctx context.Context, modelId int, userId int) fall.Error {

	query := "INSERT INTO wish (user_id, product_model_id) VALUES ($1, $2);"

	_, err := r.db.Exec(ctx, query, userId, modelId)

	if err != nil {
		return fall.ServerError(err.Error())
	}

	return nil
}

func (r *WishRepository) FindWishItem(ctx context.Context, modelId int, userId int) (*int, fall.Error) {
	query := `
	SELECT wish_id
	FROM wish
	WHERE user_id = $1 AND product_model_id = $2;
  `
	row := r.db.QueryRow(ctx, query, userId, modelId)

	var id int

	err := row.Scan(&id)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fall.NewErr(msg.WishItemNotFound, fall.STATUS_NOT_FOUND)
		}
		return nil, fall.ServerError(err.Error())
	}

	return &id, nil
}

func (r *WishRepository) DeleteFromWish(ctx context.Context, wishId int) fall.Error {

	query := "DELETE FROM wish WHERE wish_id = $1;"

	_, err := r.db.Exec(ctx, query, wishId)

	if err != nil {
		return fall.ServerError(err.Error())
	}

	return nil
}

func (r *WishRepository) GetUserWish(ctx context.Context, userId int) ([]*model.CatalogProductModel, fall.Error) {
	query := `
	SELECT p.product_id as p_id, p.title as p_title,
	b.brand_id as b_id, b.title as b_title, b.slug as b_slug, ct.category_id as ct_id, ct.title as ct_title, ct.slug as ct_slug,
	pm.product_model_id as model_id,pm.slug as m_slug, pm_artice as m_artice, pm.price as model_price, pm.discount as model_discount,
	pm.main_image_path as pm_main_img,
	pimg.product_img_id as pimg_id, pimg.product_model_id as pimg_model_id, pimg.img_path as pimg_img_path,
	sz.size_id as size_id, sz.size_value as size_value, ms.literal_size as literal_size,
	ms.product_model_id as ms_pm_id, ms.in_stock as ms_in_stock,
	ms.model_size_id as ms_m_sz_id
	FROM wish w
	INNER JOIN product_model pm ON pm.product_model_id = w.product_model_id
	INNER JOIN product p ON pm.product_id = p.product_id
	INNER JOIN category ct ON p.category_id = ct.category_id
	INNER JOIN brand b on p.brand_id = b.brand_id
	inner join model_sizes ms on ms.product_model_id = pm.product_model_id
	inner join sizes sz on ms.size_id = sz.size_id
	inner join product_model_img as pimg on pimg.product_model_id = pm.product_model_id

	WHERE w.user_id = $1;
	`

	rows, err := r.db.Query(ctx, query, userId)

	if err != nil {
		return nil, fall.ServerError(err.Error())
	}
	defer rows.Close()

	imagesMap := make(map[int]*model.ProductModelImg)
	sizesMap := make(map[int]*model.ProductModelSize)
	modelsMap := make(map[int]*model.CatalogProductModel)
	var modelOrder []int
	var imgOrder []int
	var sizeOrder []int

	for rows.Next() {
		sz := model.ProductModelSize{}
		img := model.ProductModelImg{}
		m := model.CatalogProductModel{}

		err := rows.Scan(&m.ProductId, &m.Title, &m.Brand.Id, &m.Brand.Title, &m.Brand.Slug,
			&m.Category.Id, &m.Category.Title, &m.Category.Slug, &m.ModelId, &m.Slug, &m.Article, &m.Price, &m.Discount,
			&m.MainImagePath, &img.Id, &img.ProductModelId, &img.ImgPath, &sz.SizeId, &sz.Value, &sz.Literal, &sz.ModelId, &sz.InStock, &sz.SizeModelId,
		)
		if err != nil {
			return nil, fall.ServerError(err.Error())
		}
		_, ok := modelsMap[m.ModelId]
		if !ok {
			modelsMap[m.ModelId] = &m
			modelOrder = append(modelOrder, m.ModelId)
		}
		_, ok = imagesMap[img.Id]
		if !ok {
			imagesMap[img.Id] = &img
			imgOrder = append(imgOrder, img.Id)
		}
		_, ok = sizesMap[sz.SizeModelId]
		if !ok {
			sizesMap[sz.SizeModelId] = &sz
			sizeOrder = append(sizeOrder, sz.SizeModelId)
		}
	}

	if err := rows.Err(); err != nil {
		return nil, fall.ServerError(err.Error())
	}

	for _, v := range imgOrder {
		img := imagesMap[v]
		m := modelsMap[img.ProductModelId]
		m.Images = append(m.Images, img)
	}

	for _, v := range sizeOrder {
		sz := sizesMap[v]
		m := modelsMap[sz.ModelId]
		m.Sizes = append(m.Sizes, sz)
	}

	result := make([]*model.CatalogProductModel, 0, len(modelsMap))

	for _, id := range modelOrder {
		m := modelsMap[id]
		result = append(result, m)
	}

	return result, nil
}
