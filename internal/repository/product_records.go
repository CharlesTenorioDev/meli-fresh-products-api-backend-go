package repository

import (
	"database/sql"
	"errors"

	"github.com/go-sql-driver/mysql"
	"github.com/meli-fresh-products-api-backend-t1/internal"
	"github.com/meli-fresh-products-api-backend-t1/utils/rest_err"
)

type ProductRecordsSQL struct {
	db *sql.DB
}

// Construtor para criar uma nova instância do repositório
func NewProductRecordsSQL(db *sql.DB) *ProductRecordsSQL {
	return &ProductRecordsSQL{db}
}

// Implementação de FindAll
func (psql *ProductRecordsSQL) FindAll() ([]internal.ProductRecords, error) {
	rows, err := psql.db.Query("SELECT `id`, `last_update_date`, `purchase_price`, `sale_price`, `product_id` FROM `product_records`")
	if err != nil {
		return nil, rest_err.NewInternalServerError("Erro ao buscar todos os registros de produtos")
	}
	defer rows.Close()

	var productRecords []internal.ProductRecords

	for rows.Next() {
		var productRecord internal.ProductRecords
		err := rows.Scan(&productRecord.Id, &productRecord.LastUpdateDate, &productRecord.PurchasePrice, &productRecord.SalePrice, &productRecord.ProductID)
		if err != nil {
			return nil, rest_err.NewInternalServerError("Erro ao processar registros de produtos")
		}
		productRecords = append(productRecords, productRecord)
	}

	if err = rows.Err(); err != nil {
		return nil, rest_err.NewInternalServerError("Erro ao iterar pelos registros de produtos")
	}

	return productRecords, nil
}

// Implementação de FindByID
func (psql *ProductRecordsSQL) FindByID(id int) (internal.ProductRecords, error) {
	var productRecord internal.ProductRecords

	row := psql.db.QueryRow("SELECT `id`, `last_update_date`, `purchase_price`, `sale_price`, `product_id` FROM `product_records` WHERE `id` = ?", id)
	err := row.Scan(&productRecord.Id, &productRecord.LastUpdateDate, &productRecord.PurchasePrice, &productRecord.SalePrice, &productRecord.ProductID)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return productRecord, rest_err.NewNotFoundError("Registro de produto não encontrado")
		}
		return productRecord, rest_err.NewInternalServerError("Erro ao buscar registro de produto por ID")
	}

	return productRecord, nil
}

// Implementação de Save
func (psql *ProductRecordsSQL) Save(productRec internal.ProductRecords) (internal.ProductRecords, error) {
	_, err := psql.db.Exec(
		"INSERT INTO `product_records` (`last_update_date`, `purchase_price`, `sale_price`, `product_id`) VALUES (?, ?, ?, ?)",
		productRec.LastUpdateDate, productRec.PurchasePrice, productRec.SalePrice, productRec.ProductID,
	)

	if err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) {
			switch mysqlErr.Number {
			case 1062:
				return productRec, rest_err.NewBadRequestError("Registro de produto já existe")
			default:
				return productRec, rest_err.NewInternalServerError("Erro ao salvar registro de produto")
			}
		}
	}

	return productRec, nil
}
