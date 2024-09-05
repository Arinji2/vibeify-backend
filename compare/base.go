package compare

import (
	"fmt"
	"sync"

	spotify_helpers "github.com/Arinji2/vibeify-backend/tasks/helpers/spotify"
)

func ComparePlaylists() {
	const playlist1 = "https://open.spotify.com/playlist/0qWtwYs2NAUWM3AxAIplC6?si=4819f344a55f4e48" // Travelling Vibes
	const playlist2 = "https://open.spotify.com/playlist/6kuIGnkym1n4RRO1xTFAqc?si=544719742ea64195" // Hindi Playlist

	fields := "id,name,tracks.items(track(id,name)),tracks.total,tracks.offset,tracks.limit"

	Playlist1, _, _ := spotify_helpers.GetSpotifyPlaylist(playlist1, nil, fields)
	Playlist2, _, _ := spotify_helpers.GetSpotifyPlaylist(playlist2, nil, fields)

	commonMap := make(map[string]struct{})
	missingIn1Map := make(map[string]struct{})
	missingIn2Map := make(map[string]struct{})

	set1 := make(map[string]struct{})
	set2 := make(map[string]struct{})

	for _, trackItem := range Playlist1 {
		set1[trackItem.Track.ID] = struct{}{}
	}

	for _, trackItem := range Playlist2 {
		set2[trackItem.Track.ID] = struct{}{}
	}

	var mu sync.Mutex
	wg := sync.WaitGroup{}
	chunkSize := 50

	processChunk := func(start, end int) {
		defer wg.Done()
		for i := start; i < end; i++ {
			id1 := Playlist1[i].Track.ID

			mu.Lock()
			if _, foundIn2 := set2[id1]; foundIn2 {
				commonMap[id1] = struct{}{}
			} else {
				missingIn2Map[id1] = struct{}{}
			}
			mu.Unlock()
		}
	}

	for i := 0; i < len(Playlist1); i += chunkSize {
		end := i + chunkSize
		if end > len(Playlist1) {
			end = len(Playlist1)
		}
		wg.Add(1)
		go processChunk(i, end)
	}

	processChunk2 := func(start, end int) {
		defer wg.Done()
		for i := start; i < end; i++ {
			id2 := Playlist2[i].Track.ID

			mu.Lock()
			if _, foundIn1 := set1[id2]; !foundIn1 {
				missingIn1Map[id2] = struct{}{}
			}
			mu.Unlock()
		}
	}

	for i := 0; i < len(Playlist2); i += chunkSize {
		end := i + chunkSize
		if end > len(Playlist2) {
			end = len(Playlist2)
		}
		wg.Add(1)
		go processChunk2(i, end)
	}

	wg.Wait()

	common := []string{}
	missingIn1 := []string{}
	missingIn2 := []string{}

	for id := range commonMap {
		common = append(common, id)
	}
	for id := range missingIn1Map {
		missingIn1 = append(missingIn1, id)
	}
	for id := range missingIn2Map {
		missingIn2 = append(missingIn2, id)
	}

	fmt.Println("Common tracks:", len(common), common)
	fmt.Println("Tracks in Playlist1 missing in Playlist2:", len(missingIn1))
	fmt.Println("Tracks in Playlist2 missing in Playlist1:", len(missingIn2))
}
