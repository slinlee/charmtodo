package main

import (
	"testing"
	"time"
)

func TestGetIndexDate(t *testing.T) {
	now := time.Now()
	expectedX, expectedY := getDateIndex(now)

	actualDate := getIndexDate(expectedX, expectedY)
	actualX, actualY := getDateIndex(actualDate)

	if actualDate.Month() != now.Month() &&
		actualDate.Year() != now.Year() &&
		actualDate.Day() != now.Day() {
		t.Fatalf("Date doesn't match: %v %v", actualDate, now)
	}

	if actualX != expectedX || actualY != expectedY {
		t.Fatalf("Index Date doesn't match")
	}
}

func TestManyDates(t *testing.T) {
	// Generate mock data for debugging

	now := time.Now()

	for i := 0; i < 350; i++ {

		expectedX, expectedY := getDateIndex(now)

		actualDate := getIndexDate(expectedX, expectedY)
		actualX, actualY := getDateIndex(actualDate)

		t.Log("\n------\nactualDate:", actualDate, "\nnow:", now, "\nactualXY:", actualX, actualY, "\nexpectedXY:", expectedX, expectedY)
		if actualDate.Month() != now.Month() ||
			actualDate.Year() != now.Year() ||
			actualDate.Day() != now.Day() {
			t.Fatalf("Date doesn't match: %v %v", actualDate, now)
		}

		if actualX != expectedX || actualY != expectedY {
			t.Fatalf("Index Date doesn't match. now: %v \n actualDate: %v \n actualXY: %v, %v \n expectedXY: %v, %v",
				now,
				actualDate,
				actualX, actualY,
				expectedX, expectedY)
		}
		now = now.AddDate(0, 0, -1)
	}

}

func TestFileDates(t *testing.T) {
	// Get list of dates from file

	readFromFile("./tests/test.json")

	for _, v := range calData {

		expectedX, expectedY := getDateIndex(v.Date)

		actualDate := getIndexDate(expectedX, expectedY)
		actualX, actualY := getDateIndex(actualDate)

		t.Log("\n------\nactualDate:", actualDate, "\ntestdate:", v.Date, "\nactualXY:", actualX, actualY, "\nexpectedXY:", expectedX, expectedY)
		if actualDate.Month() != v.Date.Month() ||
			actualDate.Year() != v.Date.Year() ||
			actualDate.Day() != v.Date.Day() {
			t.Fatalf("Date doesn't match: %v %v", actualDate, v.Date)
		}

		if actualX != expectedX || actualY != expectedY {
			t.Fatalf("Index Date doesn't match. testdate: %v \n actualDate: %v \n actualXY: %v, %v \n expectedXY: %v, %v",
				v.Date,
				actualDate,
				actualX, actualY,
				expectedX, expectedY)
		}
	}

}
