package repository_test

import (
	"testing"

	"github.com/meli-fresh-products-api-backend-t1/internal"
	"github.com/meli-fresh-products-api-backend-t1/internal/repository"
	"github.com/stretchr/testify/require"
)

func TestRepository_NewMapBuyerUnitTest(t *testing.T) {
	t.Run("creates map successfully", func(t *testing.T) {
		mp, e := repository.NewBuyerMap("../../db/buyer.json")

		require.NoError(t, e)
		require.NotNil(t, mp)
	})
	t.Run("fails to create map, invalid path", func(t *testing.T) {
		mp, e := repository.NewBuyerMap("not a valid path")

		require.Nil(t, mp)
		require.Error(t, e)
	})
	t.Run("fails to create map, invalid json structure", func(t *testing.T) {
		mp, e := repository.NewBuyerMap("../../db/buyer_test.json")

		require.Nil(t, mp)
		require.Error(t, e)
	})
}

func TestRepository_MapImplementationsBuyerUnitTest(t *testing.T) {
	mp, e := repository.NewBuyerMap("../../db/buyer.json")

	require.NoError(t, e)
	require.NotNil(t, mp)
	expectedBuyers := map[int]internal.Buyer{
		0: {
			ID:           0,
			CardNumberID: "1234567812345678",
			FirstName:    "John",
			LastName:     "Doe",
		},
		1: {
			ID:           1,
			CardNumberID: "2345678923456789",
			FirstName:    "Jane",
			LastName:     "Smith",
		},
		2: {
			ID:           2,
			CardNumberID: "3456789034567890",
			FirstName:    "Alice",
			LastName:     "Johnson",
		},
		3: {
			ID:           3,
			CardNumberID: "4567890145678901",
			FirstName:    "Bob",
			LastName:     "Williams",
		},
		4: {
			ID:           4,
			CardNumberID: "5678901256789012",
			FirstName:    "Charlie",
			LastName:     "Brown",
		},
	}

	t.Run("getAll", func(t *testing.T) {
		actualBuyers := mp.GetAll()

		require.Equal(t, expectedBuyers, actualBuyers)
	})
	t.Run("add", func(t *testing.T) {
		buyer := internal.Buyer{
			ID:           5,
			FirstName:    "Fabio",
			LastName:     "Nacarelli",
			CardNumberID: "3456",
		}
		mp.Add(&buyer)

		rp := mp.GetAll()
		require.Equal(t, buyer, rp[5])
	})
	t.Run("update", func(t *testing.T) {
		cardNumberID := "1234"
		firstName := "404"
		lastName := "NoName"
		mp.Update(5, internal.BuyerPatch{
			CardNumberID: &cardNumberID,
			FirstName:    &firstName,
			LastName:     &lastName,
		})
		expectedBuyer := internal.Buyer{
			ID:           5,
			CardNumberID: cardNumberID,
			FirstName:    firstName,
			LastName:     lastName,
		}
		rp := mp.GetAll()

		require.Equal(t, expectedBuyer, rp[5])
	})
	t.Run("delete", func(t *testing.T) {
		mp.Delete(5)
		rp := mp.GetAll()

		_, exists := rp[5]
		require.False(t, exists)
	})
	t.Run("reportPurchaseOrders(unimplemented)", func(t *testing.T) {
		mp.ReportPurchaseOrders()
	})
	t.Run("reportPurchaseOrdersByID(unimplemented)", func(t *testing.T) {
		mp.ReportPurchaseOrdersByID(-1)
	})
}
