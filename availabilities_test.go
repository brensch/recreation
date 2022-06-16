package main

// func TestCheckForAvailabilityChange(t *testing.T) {

// 	ctx := context.Background()
// 	// TODO: use a local firestore instance for this
// 	fs, err := InitFirestore(ctx)
// 	if err != nil {
// 		t.Log(err)
// 		t.FailNow()
// 	}
// 	defer fs.Close()

// 	log, _ := zap.NewProduction()

// 	now := time.Date(2022, 6, 12, 12, 0, 0, 0, time.UTC)
// 	targetTime := time.Date(2022, 7, 12, 12, 0, 0, 0, time.UTC)

// 	// TODO: set up testing server to test against
// 	newAvailabilities, deltas, err := CheckForAvailabilityChange(ctx, log, "http://www.recreation.gov", fs, targetTime, now, "232784")
// 	if err != nil {
// 		t.Log(err)
// 		t.FailNow()
// 	}

// 	// TODO: since this is an integration test with firebase need to figure out what to expect
// 	t.Log(newAvailabilities, deltas)
// }

// func TestDoAvailabilitiesSync(t *testing.T) {

// 	ctx := context.Background()
// 	// TODO: use a local firestore instance for this
// 	fs, err := InitFirestore(ctx)
// 	if err != nil {
// 		t.Log(err)
// 		t.FailNow()
// 	}
// 	defer fs.Close()

// 	log, _ := zap.NewDevelopment()
// 	now := time.Date(2022, 6, 12, 12, 0, 0, 0, time.UTC)
// 	targetTime := time.Date(2022, 7, 12, 12, 0, 0, 0, time.UTC)

// 	err = AvailabilitiesSync(ctx, log, fs, targetTime, now)
// 	if err != nil {
// 		t.Log(err)
// 		t.FailNow()
// 	}

// 	// TODO: since this is an integration test with firebase need to figure out what to expect

// }

// func TestChunkGroundsUp(t *testing.T) {

// 	proxies := []string{
// 		"send",
// 	}
// 	for i := 0; i < 100; i++ {
// 		chunks := ChunkGroundsUp(proxies, campgroundIDs)

// 		// check chunk count matches proxy count
// 		if len(chunks) > len(proxies) {
// 			t.Errorf("got %d chunks, expected %d", len(chunks), len(proxies))
// 			for _, chunk := range chunks {
// 				t.Log(len(chunk))
// 			}
// 			t.FailNow()
// 		}

// 		// check total number of campgrounds is correct
// 		totalGrounds := 0
// 		for _, chunk := range chunks {
// 			totalGrounds += len(chunk)
// 		}
// 		if len(campgroundIDs) != totalGrounds {
// 			t.Errorf("got %d campgrounds, expected %d", totalGrounds, len(campgroundIDs))
// 		}

// 		proxies = append(proxies, "sendo")
// 	}

// }
