package repository

import (
	"L0/internal/control"
	"L0/internal/models"
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
)

type controlRepo struct {
	db *sqlx.DB
}

func NewControlRepository(db *sqlx.DB) control.Repository {
	return &controlRepo{db: db}
}

func (c *controlRepo) NewOrder(ctx context.Context, params models.Order) error {
	var itemsId []int
	var paymantId int
	var deliveryId int

	err := c.db.GetContext(
		ctx,
		&deliveryId,
		queryNewDelivery,
		params.DeliveryData.Name,
		params.DeliveryData.Phone,
		params.DeliveryData.Zip,
		params.DeliveryData.City,
		params.DeliveryData.Address,
		params.DeliveryData.Region,
		params.DeliveryData.Email,
	)
	// log.Println(params.DeliveryData.Name,
	// 	params.DeliveryData.Phone,
	// 	params.DeliveryData.Zip,
	// 	params.DeliveryData.City,
	// 	params.DeliveryData.Address,
	// 	params.DeliveryData.Region,
	// 	params.DeliveryData.Email)
	if err != nil {
		return errors.Wrapf(err, "controlRepo.NewOrder.ExecContext(Order: %s, Delivery)", params.OrderUid)
	}

	err = c.db.GetContext(
		ctx,
		&paymantId,
		queryNewPaymant,
		params.PaymantData.Transaction,
		params.PaymantData.Request_id,
		params.PaymantData.Currency,
		params.PaymantData.Provider,
		params.PaymantData.Amount,
		params.PaymantData.Payment_dt,
		params.PaymantData.Bank,
		params.PaymantData.Delivery_cost,
		params.PaymantData.Goods_total,
		params.PaymantData.Custom_fee,
	)
	if err != nil {
		return errors.Wrapf(err, "controlRepo.NewOrder.ExecContext(Order: %s, Paymant)", params.OrderUid)
	}

	for _, val := range params.ItemsData {
		var tmp int
		err = c.db.GetContext(
			ctx,
			&tmp,
			queryNewItems,
			val.Chrt_id,
			val.Track_number,
			val.Price,
			val.Rid,
			val.Name,
			val.Sale,
			val.Size,
			val.Total_price,
			val.Nm_id,
			val.Brand,
			val.Status,
		)
		if err != nil {
			return errors.Wrapf(err, "controlRepo.NewOrder.ExecContext(Order: %s, Items)", params.OrderUid)
		}
		itemsId = append(itemsId, tmp)
	}

	_, err = c.db.ExecContext(
		ctx,
		queryNewOrder,
		params.OrderUid,
		params.Track_number,
		params.Entry,
		deliveryId,
		paymantId,
		pq.Array(itemsId),
		params.Locale,
		params.Internal_signature,
		params.Customer_id,
		params.Delivery_service,
		params.Shardkey,
		params.Sm_id,
		params.Date_created,
		params.Oof_shard,
	)
	if err != nil {
		return errors.Wrapf(err, "controlRepo.NewOrder.ExecContext(Order: %s)", params.OrderUid)
	}
	//TODO slise items id

	return nil
}

func (c *controlRepo) DataRecovery(ctx context.Context) (result []byte, err error) {
	err = c.db.GetContext(ctx, &result, queryDataRecovery)
	if err != nil && err.Error() != "sql: no rows in result set" {
		return nil, errors.Wrap(err, "controlRepo.DataRecovery.GetContext()")
	}
	return result, nil
}

func (c *controlRepo) GetOrder(ctx context.Context, uid string) (result []byte, err error) {
	err = c.db.GetContext(ctx, &result, queryGetOrder, uid)
	if err != nil {
		return nil, errors.Wrap(err, "controlRepo.GetOrder.GetContext()")
	}
	return result, nil
}
