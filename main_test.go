package main

import (
	"testing"
	"time"
)

func makeCount() SadhakaCounts {
	return SadhakaCounts{
		cities:   make(map[string]int),
		sadhakas: make(map[string]time.Time),
	}
}

func TestEmptyCountsIsIncremented(t *testing.T) {
	counts := makeCount()
	updateAndGetSadhakaCount(ParsedRequest{
		City: "Krakow",
	}, &counts)

	if counts.cities["Krakow"] != 1 {
		t.Errorf("For empty counts should add one to the city: got %d, want 1", counts.cities["Krakow"])
	}
}

func TestMultipleRequestsFromTheSameSadhakaOnlyCountOncePerDay(t *testing.T) {
	counts := makeCount()
	updateAndGetSadhakaCount(ParsedRequest{
		City:    "Krk",
		Sadhaka: "uuid",
		Time:    time.Now(),
	}, &counts)
	updateAndGetSadhakaCount(ParsedRequest{
		City:    "Krk",
		Sadhaka: "uuid",
		Time:    time.Now(),
	}, &counts)

	if counts.cities["Krk"] != 1 {
		t.Errorf("Multiple Requests from the same Sadhaka for a given day count only once but was %d", counts.cities["Krk"])
	}
}
