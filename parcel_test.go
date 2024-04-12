package main

import (
	"database/sql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"math/rand"
	"testing"
	"time"
)

var (
	// randSource is the source of pseudorandom numbers.
	// To increase uniqueness, the current time in Unix format (as a number) is used as the seed.
	randSource = rand.NewSource(time.Now().UnixNano())
	// randRange uses randSource to generate random numbers.
	randRange = rand.New(randSource)
)

// getTestParcel returns a test parcel.
func getTestParcel() Parcel {
	return Parcel{
		Client:    1000,
		Status:    ParcelStatusRegistered,
		Address:   "test",
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
	}
}

// TestAddGetDelete tests adding, retrieving, and deleting a parcel.
func TestAddGetDelete(t *testing.T) {
	// prepare
	db, err := sql.Open("sqlite", "tracker.db")
	require.NoError(t, err)
	defer db.Close()
	store := NewParcelStore(db)
	parcel := getTestParcel()

	id, err := store.Add(parcel)
	assert.NoError(t, err)
	assert.Equal(t, id, parcel.Number)

	storedParcel, err := store.Get(id)
	assert.NoError(t, err)
	assert.Equal(t, parcel, storedParcel)

	err = store.Delete(id)
	assert.NoError(t, err)
}

// TestSetAddress tests updating the address.
func TestSetAddress(t *testing.T) {
	// prepare
	db, err := sql.Open("sqlite", "tracker.db")
	require.NoError(t, err)
	defer db.Close()
	store := NewParcelStore(db)
	parcel := getTestParcel()

	id, err := store.Add(parcel)
	assert.NoError(t, err)
	assert.Equal(t, id, parcel.Number)

	newAddress := "new test address"
	err = store.SetAddress(id, newAddress)
	assert.NoError(t, err)

	storedParcel, err := store.Get(id)
	assert.NoError(t, err)
	assert.Equal(t, newAddress, storedParcel.Address)

	err = store.Delete(id)
	assert.NoError(t, err)
}

// TestSetStatus tests updating the status.
func TestSetStatus(t *testing.T) {
	// prepare
	db, err := sql.Open("sqlite", "tracker.db")
	require.NoError(t, err)
	defer db.Close()
	store := NewParcelStore(db)
	parcel := getTestParcel()

	id, err := store.Add(parcel)
	assert.NoError(t, err)
	assert.Equal(t, id, parcel.Number)

	err = store.SetStatus(id, ParcelStatusDelivered)
	assert.NoError(t, err)

	storedParcel, err := store.Get(id)
	assert.NoError(t, err)
	assert.Equal(t, ParcelStatusDelivered, storedParcel.Status)

	err = store.SetStatus(id, ParcelStatusRegistered)
	assert.NoError(t, err)

	err = store.Delete(id)
	assert.NoError(t, err)
}

// TestGetByClient tests retrieving parcels by client identifier.
func TestGetByClient(t *testing.T) {
	// prepare
	db, err := sql.Open("sqlite", "tracker.db")
	require.NoError(t, err)
	defer db.Close()
	store := NewParcelStore(db)

	parcels := []Parcel{
		getTestParcel(),
		getTestParcel(),
		getTestParcel(),
	}
	parcelMap := map[int]Parcel{}

	client := randRange.Intn(10_000_000)
	parcels[0].Client = client
	parcels[1].Client = client
	parcels[2].Client = client

	parcels[1].Number = 1
	parcels[2].Number = 2

	// add
	for i := 0; i < len(parcels); i++ {
		id, err := store.Add(parcels[i])
		assert.NoError(t, err)
		assert.Equal(t, id, parcels[i].Number)

		// Update the identifier of the added parcel.
		parcels[i].Number = id

		// Save the added parcel in a map structure, so it can be easily retrieved by parcel identifier.
		parcelMap[id] = parcels[i]
	}

	storedParcels, err := store.GetByClient(client)
	assert.NoError(t, err)
	assert.Len(t, storedParcels, len(parcels))
	for _, parcel := range storedParcels {
		// In parcelMap lie the added parcels, the key - parcel identifier, the value - the parcel itself.
		require.Contains(t, parcelMap, parcel.Number)
	}

	for _, parcel := range storedParcels {
		err = store.Delete(parcel.Number)
		assert.NoError(t, err)
	}
}
