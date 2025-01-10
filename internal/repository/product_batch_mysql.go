package repository

import (
	"database/sql"
	"errors"

	"github.com/go-sql-driver/mysql"
	"github.com/meli-fresh-products-api-backend-t1/internal"
)

var (
	ProductBatchAlreadyExists = errors.New("product-batch already exists")
	ProductBatchNotFound      = errors.New("product-batch not found")
)

func NewProductBatchMysql(db *sql.DB) *ProductBatchDB {
	return &ProductBatchDB{db}
}

type ProductBatchDB struct {
	db *sql.DB
}

func (r *ProductBatchDB) FindByID(id int) (internal.ProductBatch, error) {
	query := `
	SELECT 
		pb.id,
		pb.batch_number,
		pb.current_quantity,
		pb.current_temperature,
		pb.due_date,
		pb.initial_quantity,
		pb.manufacturing_date,
		pb.manufacturing_hour,
		pb.minumum_temperature,           
		pb.product_id,           
		pb.section_id           
	FROM 
		product_batches pb
	WHERE 
		id = ?`

	var pb internal.ProductBatch
	err := r.db.QueryRow(query, id).Scan(
		&pb.ID,
		&pb.BatchNumber,
		&pb.CurrentQuantity,
		&pb.CurrentTemperature,
		&pb.DueDate,
		&pb.InitialQuantity,
		&pb.ManufacturingDate,
		&pb.ManufacturingHour,
		&pb.MinumumTemperature,
		&pb.ProductId,
		&pb.SectionId,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return pb, ProductBatchNotFound
		}
		return pb, err
	}

	return pb, nil
}

func (r *ProductBatchDB) Save(prodBatch *internal.ProductBatch) error {
	result, err := r.db.Exec(
		"INSERT INTO product_batches (batch_number, current_quantity, current_temperature, due_date, initial_quantity, manufacturing_date, manufacturing_hour, minumum_temperature, product_id, section_id) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		prodBatch.BatchNumber,
		prodBatch.CurrentQuantity,
		prodBatch.CurrentTemperature,
		prodBatch.DueDate,
		prodBatch.InitialQuantity,
		prodBatch.ManufacturingDate,
		prodBatch.ManufacturingHour,
		prodBatch.MinumumTemperature,
		prodBatch.ProductId,
		prodBatch.SectionId,
	)

	if err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) {
			switch mysqlErr.Number {
			case 1062:
				return ProductBatchAlreadyExists
			}
		}
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	prodBatch.ID = int(id)

	return nil
}

func (r *ProductBatchDB) ProductBatchNumberExists(batchNumber int) (bool, error) {
	query := "SELECT COUNT(*) FROM product_batches WHERE batch_number = ?"

	var count int
	err := r.db.QueryRow(query, batchNumber).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *ProductBatchDB) ReportProducts() (prodBatches []internal.ProductBatch, err error) {
	query := `
	SELECT 
		pb.batch_number,
		pb.current_quantity,
		pb.current_temperature,
		pb.due_date,
		pb.initial_quantity,
		pb.manufacturing_date,
		pb.manufacturing_hour,
		pb.minumum_temperature,
		p.product_code,        
		s.section_number       
	FROM 
		product_batches pb
	JOIN 
		products p ON pb.product_id = p.id
	JOIN 
		sections s ON pb.section_id = s.id
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, ProductBatchNotFound
	}
	defer rows.Close()

	for rows.Next() {
		var pb internal.ProductBatch
		if err := rows.Scan(
			&pb.BatchNumber,
			&pb.CurrentQuantity,
			&pb.CurrentTemperature,
			&pb.DueDate,
			&pb.InitialQuantity,
			&pb.ManufacturingDate,
			&pb.ManufacturingHour,
			&pb.MinumumTemperature,
			&pb.ProductId,
			&pb.SectionId,
		); err != nil {
			return nil, ProductBatchNotFound
		}
		prodBatches = append(prodBatches, pb)
	}

	if err := rows.Err(); err != nil {
		return nil, ProductBatchNotFound
	}

	return prodBatches, nil
}

func (r *ProductBatchDB) ReportProductsByID(id int) (prodBatches []internal.ProductBatch, err error) {
	query := `
	SELECT 
		pb.batch_number,
		pb.current_quantity,
		pb.current_temperature,
		pb.due_date,
		pb.initial_quantity,
		pb.manufacturing_date,
		pb.manufacturing_hour,
		pb.minumum_temperature,
		p.product_code,          -- Supondo que você queira o código do produto
		s.section_number         -- Supondo que você queira o número da seção
	FROM 
		product_batches pb
	JOIN 
		products p ON pb.product_id = p.id
	JOIN 
		sections s ON pb.section_id = s.id
	WHERE 
		pb.id = ?
	`

	row := r.db.QueryRow(query, id)

	var pb internal.ProductBatch
	if err := row.Scan(
		&pb.BatchNumber,
		&pb.CurrentQuantity,
		&pb.CurrentTemperature,
		&pb.DueDate,
		&pb.InitialQuantity,
		&pb.ManufacturingDate,
		&pb.ManufacturingHour,
		&pb.MinumumTemperature,
		&pb.ProductId, // Se também precisar do ProductId
		&pb.SectionId, // Se também precisar do SectionId
		// Adicione aqui mais campos que você precisa retornar
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, internal.ErrLocalityNotFound // Ajuste conforme necessário
		}
		return nil, ProductBatchNotFound
	}

	prodBatches = append(prodBatches, pb)

	return prodBatches, nil
}
