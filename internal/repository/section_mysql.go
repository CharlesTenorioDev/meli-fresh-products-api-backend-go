package repository

import (
	"database/sql"
	"errors"

	"github.com/meli-fresh-products-api-backend-t1/internal"
)

func NewSectionMysql(db *sql.DB) *SectionMysql {
	return &SectionMysql{db}
}

type SectionMysql struct {
	db *sql.DB
}

func (r *SectionMysql) FindAll() ([]internal.Section, error) {
	rows, err := r.db.Query("SELECT `id`, `section_number`, `current_temperature`, `minimum_temperature`, `current_capacity`, `minimum_capacity`, `maximum_capacity`, `warehouse_id`, `product_type_id` FROM sections")
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = internal.ErrSectionNotFound
		}

		return nil, err
	}
	defer rows.Close()

	var sections []internal.Section

	for rows.Next() {
		var s internal.Section

		err := rows.Scan(&s.ID, &s.SectionNumber, &s.CurrentTemperature, &s.MinimumTemperature, &s.CurrentCapacity, &s.MinimumCapacity, &s.MaximumCapacity, &s.WarehouseID, &s.ProductTypeID)
		if err != nil {
			return nil, err
		}

		sections = append(sections, s)
	}

	err = rows.Err()
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = internal.ErrSectionNotFound
		}

		return nil, err
	}

	return sections, nil
}

func (r *SectionMysql) FindByID(id int) (internal.Section, error) {
	query := `
	SELECT 
		id, 
		section_number, 
		current_temperature, 
		minimum_temperature, 
		current_capacity, 
		minimum_capacity, 
		maximum_capacity, 
		warehouse_id, 
		product_type_id 
	FROM 
		sections 
	WHERE 
		id = ?`

	var s internal.Section

	err := r.db.QueryRow(query, id).Scan(&s.ID, &s.SectionNumber, &s.CurrentTemperature, &s.MinimumTemperature, &s.CurrentCapacity, &s.MinimumCapacity, &s.MaximumCapacity, &s.WarehouseID, &s.ProductTypeID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return s, internal.ErrSectionNotFound
		}

		return s, err
	}

	return s, nil
}

func (r *SectionMysql) ReportProducts() ([]internal.ReportProduct, error) {
	query := `
        SELECT 
            s.id AS section_id,
            s.section_number,
            COALESCE(SUM(pb.current_quantity), 0) AS products_count
        FROM 
            sections s
        LEFT JOIN 
            product_batches pb ON s.id = pb.section_id
        GROUP BY 
            s.id, s.section_number;`

	rows, err := r.db.Query(query)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = internal.ErrReportProductNotFound
		}

		return nil, err
	}
	defer rows.Close()

	var report []internal.ReportProduct

	for rows.Next() {
		var rp internal.ReportProduct

		if err := rows.Scan(&rp.SectionID, &rp.SectionNumber, &rp.ProductsCount); err != nil {
			return nil, err
		}

		report = append(report, rp)
	}

	err = rows.Err()
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = internal.ErrReportProductNotFound
		}

		return nil, err
	}

	return report, nil
}

func (r *SectionMysql) ReportProductsByID(sectionID int) (internal.ReportProduct, error) {
	query := `
			SELECT 
				s.id AS section_id,
				s.section_number,
				COALESCE(SUM(pb.current_quantity), 0) AS products_count
			FROM 
				sections s
			LEFT JOIN 
				product_batches pb ON s.id = pb.section_id
			WHERE 
				s.id = ?
			GROUP BY 
				s.id, s.section_number;`

	var rp internal.ReportProduct

	err := r.db.QueryRow(query, sectionID).Scan(&rp.SectionID, &rp.SectionNumber, &rp.ProductsCount)
	if err != nil {
		if err == sql.ErrNoRows {
			rp.ProductsCount = 0
			return rp, nil
		}

		return rp, err
	}

	return rp, nil
}

func (r *SectionMysql) SectionNumberExists(sectionNumber int) (bool, error) {
	query := "SELECT COUNT(*) FROM sections WHERE section_number = ?"

	var count int

	err := r.db.QueryRow(query, sectionNumber).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *SectionMysql) Save(section *internal.Section) error {
	result, err := r.db.Exec(
		"INSERT INTO sections (section_number, current_temperature, minimum_temperature, current_capacity, minimum_capacity, maximum_capacity, warehouse_id, product_type_id) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
		section.SectionNumber,
		section.CurrentTemperature,
		section.MinimumTemperature,
		section.CurrentCapacity,
		section.MinimumCapacity,
		section.MaximumCapacity,
		section.WarehouseID,
		section.ProductTypeID,
	)

	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	section.ID = int(id)

	return nil
}

func (r *SectionMysql) Update(section *internal.Section) error {
	_, err := r.db.Exec(
		"UPDATE sections SET section_number = ?, current_temperature = ?, minimum_temperature = ?, current_capacity = ?, minimum_capacity = ?, maximum_capacity = ?, warehouse_id = ?, product_type_id = ? WHERE id = ?",
		section.SectionNumber,
		section.CurrentTemperature,
		section.MinimumTemperature,
		section.CurrentCapacity,
		section.MinimumCapacity,
		section.MaximumCapacity,
		section.WarehouseID,
		section.ProductTypeID,
		section.ID,
	)

	return err
}

func (r *SectionMysql) Delete(id int) error {
	_, err := r.db.Exec("DELETE FROM sections WHERE id = ?", id)
	return err
}
