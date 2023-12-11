package wish

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	exception "github.com/maximfedotov74/fiber-psql/internal/shared/error"
)

type WishRepository struct {
	db *pgxpool.Pool
}

func NewWishRepository(db *pgxpool.Pool) *WishRepository {
	return &WishRepository{db: db}
}

func (wr *WishRepository) FindModelInUserCart(modelSizeId int, userId int) (*CartItemModel, exception.Error) {

	query := `
	SELECT cart.cart_id,cart.user_id,cart.model_size_id,cart.quantity, ms.in_stock
	FROM cart
	INNER JOIN model_sizes as ms ON cart.model_size_id = ms.model_size_id
	WHERE cart.user_id = $1 AND cart.model_size_id = $2;
  `
	row := wr.db.QueryRow(context.Background(), query, userId, modelSizeId)

	cartItem := CartItemModel{}

	err := row.Scan(&cartItem.CartItemId, &cartItem.UserId, &cartItem.ModelSizeId, &cartItem.Quantity, &cartItem.InStock)

	if err != nil {
		return nil, exception.NewErr(CART_ITEM_NOT_FOUND, exception.STATUS_NOT_FOUND)
	}

	return &cartItem, nil
}

func (wr *WishRepository) AddToCart(modelSizeId int, userId int) exception.Error {

	query := "INSERT INTO cart (user_id, model_size_id) VALUES ($1, $2);"

	_, err := wr.db.Exec(context.Background(), query, userId, modelSizeId)

	if err != nil {
		return exception.ServerError(err.Error())
	}

	return nil
}

func (wr *WishRepository) DeleteFromCart(cartItemId int) exception.Error {
	query := "DELETE FROM cart WHERE cart_id = $1;"

	_, err := wr.db.Exec(context.Background(), query, cartItemId)

	if err != nil {
		return exception.ServerError(err.Error())
	}

	return nil
}

func (wr *WishRepository) UpdateCartItem(cartItemId int, newQuantity int) exception.Error {

	query := `UPDATE cart SET quantity = $1 WHERE cart_id = $2;`

	_, err := wr.db.Exec(context.Background(), query, newQuantity, cartItemId)

	if err != nil {
		exception.ServerError(err.Error())
	}
	return nil
}

func (wr *WishRepository) RemoveSeveralItems(cartIds string) exception.Error {
	query := `DELETE FROM cart WHERE cart_id IN ($1);`

	_, err := wr.db.Exec(context.Background(), query, cartIds)

	if err != nil {
		return exception.ServerError(err.Error())
	}

	return nil
}

func (wr *WishRepository) GetUserCart(userId int) ([]CartItem, exception.Error) {
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
	rows, err := wr.db.Query(context.Background(), query, userId)
	if err != nil {
		return nil, exception.ServerError(err.Error())
	}
	defer rows.Close()
	var result []CartItem
	for rows.Next() {
		item := CartItem{}
		err := rows.Scan(&item.Id, &item.Quantity, &item.ModelSize.Id, &item.ModelSize.LiteralSize, &item.ModelSize.InStock,
			&item.ModelSize.SizeId, &item.ModelSize.SizeValue, &item.ModelSize.ProductModel.Id, &item.ModelSize.ProductModel.Price,
			&item.ModelSize.ProductModel.Discount, &item.ModelSize.ProductModel.ImagePath, &item.ModelSize.ProductModel.Product.Id,
			&item.ModelSize.ProductModel.Product.Title, &item.ModelSize.ProductModel.Product.Slug,
			&item.ModelSize.ProductModel.Product.Category.Id, &item.ModelSize.ProductModel.Product.Category.Title,
			&item.ModelSize.ProductModel.Product.Category.ShortTitle, &item.ModelSize.ProductModel.Product.Category.Slug,
			&item.ModelSize.ProductModel.Product.Brand.Id, &item.ModelSize.ProductModel.Product.Brand.Title,
			&item.ModelSize.ProductModel.Product.Brand.Slug,
		)
		if err != nil {
			return nil, exception.ServerError(err.Error())
		}
		result = append(result, item)
	}

	if rows.Err() != nil {
		return nil, exception.ServerError(rows.Err().Error())
	}

	return result, nil

}

func (wr *WishRepository) AddToWish(modelId int, userId int) exception.Error {

	query := "INSERT INTO wish (user_id, product_model_id) VALUES ($1, $2);"

	_, err := wr.db.Exec(context.Background(), query, userId, modelId)

	if err != nil {
		return exception.ServerError(err.Error())
	}

	return nil
}

func (wr *WishRepository) FindWishItem(modelId int, userId int) (*int, exception.Error) {
	query := `
	SELECT wish_id
	FROM wish
	WHERE user_id = $1 AND product_model_id = $2;
  `
	row := wr.db.QueryRow(context.Background(), query, userId, modelId)

	var id int

	err := row.Scan(&id)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, exception.NewErr(WISH_ITEM_NOT_FOUND, exception.STATUS_NOT_FOUND)
		}
		return nil, exception.NewErr(CART_ITEM_NOT_FOUND, exception.STATUS_NOT_FOUND)
	}

	return &id, nil
}

func (wr *WishRepository) DeleteFromWish(wishId int) exception.Error {

	query := "DELETE FROM wish WHERE wish_id = $1;"

	_, err := wr.db.Exec(context.Background(), query, wishId)

	if err != nil {
		return exception.ServerError(err.Error())
	}

	return nil
}
