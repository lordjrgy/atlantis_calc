package calc

import (
	"fmt"
	"image/color"
	_ "image/png"
	"slices"
	"strings"

	log "github.com/sirupsen/logrus"
)

type MoveQuality int

const (
	BestMove MoveQuality = iota
	GreatMove
	BrilliantMove
)

var (
	bestMoveColor      = color.RGBA{155, 199, 0, 200}
	greatMoveColor     = color.RGBA{0, 121, 211, 200}
	brilliantMoveColor = color.RGBA{48, 162, 197, 200}
)

type Difficulty int

const (
	Easy Difficulty = iota
	Hard
)


type BoostRoom struct {
	Name      string
	Time      float64
	BoostTime float64
	Quality   MoveQuality
}

type Room struct {
	Name          string
	BoostlessTime float64
	Difficulty Difficulty
	BoostStrats   []BoostRoom
}

var RoomMap = map[string]Room{
	// Easy rooms
	"1a": {Name: "1a", BoostlessTime: 13.6, Difficulty: Easy, BoostStrats: []BoostRoom{
		{Name: "cp 1-2", Time: 10.6, BoostTime: 9.6, Quality: BestMove},
	}},
	"1b": {Name: "1b", BoostlessTime: 15.1, Difficulty: Easy, BoostStrats: []BoostRoom{
		{Name: "cp 1-2", Time: 11.1, BoostTime: 8.1, Quality: BestMove},
		{Name: "cp 0-1", Time: 11.4, BoostTime: 3.2, Quality: BestMove},
	}}, 
	"1c": {Name: "1c", BoostlessTime: 11.4, Difficulty: Easy, BoostStrats: []BoostRoom{
		{Name: "cp 1-2", Time: 8.4, BoostTime: 7.0, Quality: BestMove},
		{Name: "cp 0-1", Time: 8.9, BoostTime: 2.1, Quality: GreatMove},
	}},
	"1d": {Name: "1d", BoostlessTime: 17.4, Difficulty: Easy, BoostStrats: []BoostRoom{
		{Name: "cp 0-1", Time: 13.2, BoostTime: 5.5, Quality: BestMove},
	}},
	"1e": {Name: "1e", BoostlessTime: 14.9, Difficulty: Easy, BoostStrats: []BoostRoom{
		{Name: "cp 1-2 (Late)", Time: 12.1, BoostTime: 10.2, Quality: GreatMove},
		{Name: "cp 1-2 (Early)", Time: 12.3, BoostTime: 6.1, Quality: GreatMove},
	}},
	
	"2a": {Name: "2a", BoostlessTime: 13.6, Difficulty: Easy, BoostStrats: []BoostRoom{
		{Name: "cp 1-2", Time: 11.6, BoostTime: 5.6, Quality: GreatMove},
	}},
	"2b": {Name: "2b", BoostlessTime: 16.9, Difficulty: Easy, BoostStrats: []BoostRoom{
		{Name: "cp 1-2 (BIG JRGY)", Time: 12.4, BoostTime: 7.1, Quality: BrilliantMove},
	}},
	"2c": {Name: "2c", BoostlessTime: 19.0, Difficulty: Easy, BoostStrats: []BoostRoom{
		{Name: "cp 1-2", Time: 15.0, BoostTime: 4.2, Quality: BestMove},
		{Name: "cp 2-3", Time: 15.0, BoostTime: 11.9, Quality: BestMove},
	}},
	"2d": {Name: "2d", BoostlessTime: 20.5, Difficulty: Easy, BoostStrats: []BoostRoom{
		{Name: "cp 1-2", Time: 15.8, BoostTime: 6.5, Quality: BestMove},
	}},
	"2e": {Name: "2e", BoostlessTime: 15.9, Difficulty: Easy, BoostStrats: []BoostRoom{
		{Name: "cp 0-1", Time: 12.2, BoostTime: 2.6, Quality: BestMove},
	}},

	"3a": {Name: "3a", BoostlessTime: 15.2, Difficulty: Easy, BoostStrats: []BoostRoom{
		{Name: "cp 0-1", Time: 11.7, BoostTime: 4.8, Quality: BestMove},
	}},
	"3b": {Name: "3b", BoostlessTime: 15.3, Difficulty: Easy, BoostStrats: []BoostRoom{
		{Name: "cp 1-2", Time: 10.3, BoostTime: 9.5, Quality: BestMove},
		{Name: "cp 0-1", Time: 11.8, BoostTime: 3.6, Quality: GreatMove},
	}},
	"3c": {Name: "3c", BoostlessTime: 17.4, Difficulty: Easy, BoostStrats: []BoostRoom{
		{Name: "cp 1-2", Time: 13.2, BoostTime: 7.4, Quality: BestMove},
	}},
	"3d": {Name: "3d", BoostlessTime: 26.9, Difficulty: Easy, BoostStrats: []BoostRoom{
		{Name: "cp 0-1", Time: 20.6, BoostTime: 2.0, Quality: BestMove},
		{Name: "cp 1-2", Time: 21.9, BoostTime: 14.9, Quality: BestMove},
	}},
	"3e": {Name: "3e", BoostlessTime: 15.6, Difficulty: Easy, BoostStrats: []BoostRoom{
		{Name: "cp 1-2", Time: 11.0, BoostTime: 10.0, Quality: BestMove},
	}},

	"4a": {Name: "4a", BoostlessTime: 10.8, Difficulty: Easy, BoostStrats: []BoostRoom{
		{Name: "cp 1-2", Time: 8.0, BoostTime: 5.9, Quality: BestMove},
	}},
	"4b": {Name: "4b", BoostlessTime: 17.8, Difficulty: Easy, BoostStrats: []BoostRoom{
		{Name: "cp 1-2", Time: 13.6, BoostTime: 7.6, Quality: BestMove},
		{Name: "cp 0-1", Time: 14.1, BoostTime: 5.0, Quality: GreatMove},
	}},
	"4c": {Name: "4c", BoostlessTime: 14.7, Difficulty: Easy, BoostStrats: []BoostRoom{
		{Name: "cp 1-2", Time: 11.7, BoostTime: 9.7, Quality: BestMove},
	}},
	"4e": {Name: "4e", BoostlessTime: 18.0, Difficulty: Easy, BoostStrats: []BoostRoom{
		{Name: "cp 0-1", Time: 12.7, BoostTime: 3.0, Quality: BestMove},
	}},

	"5a": {Name: "5a", BoostlessTime: 14.9, Difficulty: Easy, BoostStrats: []BoostRoom{
		{Name: "cp 0-1", Time: 11.6, BoostTime: 6.4, Quality: BestMove},
	}},
	"5b": {Name: "5b", BoostlessTime: 21.1, Difficulty: Easy, BoostStrats: []BoostRoom{
		{Name: "cp 0-1", Time: 13.8, BoostTime: 6.1, Quality: BestMove},
	}},
	"5c": {Name: "5c", BoostlessTime: 20.1, Difficulty: Easy, BoostStrats: []BoostRoom{
		{Name: "cp 2-3", Time: 16.6, BoostTime: 15.6, Quality: BestMove},
	}},
	"5d": {Name: "5d", BoostlessTime: 12.8, Difficulty: Easy, BoostStrats: []BoostRoom{
		{Name: "cp 0-1", Time: 10.1, BoostTime: 3.3, Quality: BestMove},
	}},
	"5e": {Name: "5e", BoostlessTime: 17.6, Difficulty: Easy, BoostStrats: []BoostRoom{
		{Name: "cp 0-1", Time: 13.2, BoostTime: 11.2, Quality: BestMove},
	}},

	// Hard rooms
	"1f": {Name: "1f", BoostlessTime: 27.6, Difficulty: Hard, BoostStrats: []BoostRoom{
		{Name: "cp 2-3", Time: 25.1, BoostTime: 22.2, Quality: BestMove},
	}},
	"1g": {Name: "1g", BoostlessTime: 29.9, Difficulty: Hard, BoostStrats: []BoostRoom{
		{Name: "cp 0-1", Time: 20.4, BoostTime: 0.6, Quality: BestMove},
		{Name: "cp 1-2", Time: 23.2, BoostTime: 15.8, Quality: BestMove},
	}},
	"1h": {Name: "1h", BoostlessTime: 26.2, Difficulty: Hard, BoostStrats: []BoostRoom{
		{Name: "cp 0-1", Time: 20.9, BoostTime: 7.4, Quality: BestMove},
		{Name: "cp 2-3", Time: 22.1, BoostTime: 19.4, Quality: GreatMove},
	}},
	"2f": {Name: "2f", BoostlessTime: 19.2, Difficulty: Hard, BoostStrats: []BoostRoom{
		{Name: "cp 1-2", Time: 16.1, BoostTime: 7.7, Quality: BestMove},
	}},
	"2g": {Name: "2g", BoostlessTime: 21.2, Difficulty: Hard, BoostStrats: []BoostRoom{
		{Name: "cp 1-2-3", Time: 11.9, BoostTime: 7.6, Quality: BestMove},
	}},
	"2h": {Name: "2h", BoostlessTime: 14.6, Difficulty: Hard, BoostStrats: []BoostRoom{
		{Name: "cp 1-2", Time: 10.6, BoostTime: 8.4, Quality: BestMove},
	}},
	"3f": {Name: "3f", BoostlessTime: 26.4, Difficulty: Hard, BoostStrats: []BoostRoom{
		{Name: "cp 1-2 ", Time: 20.0, BoostTime: 8.2, Quality: BestMove},
	}},
	"3g": {Name: "3g", BoostlessTime: 19.4, Difficulty: Hard, BoostStrats: []BoostRoom{
		{Name: "cp 2-3", Time: 15.4, BoostTime: 14.1, Quality: BestMove},
	}},
	"3h": {Name: "3h", BoostlessTime: 30.5, Difficulty: Hard, BoostStrats: []BoostRoom{
		{Name: "cp 2-3", Time: 25.3, BoostTime: 21.5, Quality: BestMove},
	}},
	"4f": {Name: "4f", BoostlessTime: 20.9, Difficulty: Hard, BoostStrats: []BoostRoom{
		{Name: "cp 1-2", Time: 16.2, BoostTime: 7.4, Quality: BestMove},
	}},
	"4g": {Name: "4g", BoostlessTime: 27.7, Difficulty: Hard, BoostStrats: []BoostRoom{
		{Name: "cp 1-2", Time: 23.6, BoostTime: 15.5, Quality: BestMove},
	}},
	"4h": {Name: "4h", BoostlessTime: 21.8, Difficulty: Hard, BoostStrats: []BoostRoom{
		{Name: "cp 0-1", Time: 17.2, BoostTime: 2.0, Quality: BestMove},
		{Name: "cp 2-3", Time: 18.9, BoostTime: 16.4, Quality: GreatMove},	
	}},
	"5f": {Name: "5f", BoostlessTime: 23.8, Difficulty: Hard, BoostStrats: []BoostRoom{
		{Name: "cp 2-3", Time: 20.3, BoostTime: 19.3, Quality: BestMove},
	}},
	"5g": {Name: "5g", BoostlessTime: 22.6, Difficulty: Hard, BoostStrats: []BoostRoom{
		{Name: "cp 1-2", Time: 14.6, BoostTime: 9.8, Quality: BestMove},
	}},
	"5h": {Name: "5h", BoostlessTime: 27.3, Difficulty: Hard, BoostStrats: []BoostRoom{
		{Name: "cp 0-1", Time: 21.0, BoostTime: 2.7, Quality: BestMove},
	}},

	// Finish room
	"finish room": {Name: "finish room", BoostlessTime: 2, Difficulty: Easy, BoostStrats: []BoostRoom{
		{Name: "lol", Time: 0.7, BoostTime: 0.3, Quality: BestMove},
	}},
}

func GetRooms() []string {
    res := []string{} // start with empty slice
    for _, v := range RoomMap {
        if strings.ToLower(v.Name) == "finish room" {
            continue
        }
        res = append(res, strings.ToLower(v.Name))
    }
    return res
}


func calcBoostless(roomList []string, splits map[string]Room) float64 {
	time := 0.0
	for _, room := range roomList {
		time += splits[room].BoostlessTime
	}

	// timesave := calcTimesave(roomList, nil)

	return time
}

type CalcResultBoost struct {
	Ind      int
	StratInd int
	Pacelock float64
}

type calcResult struct {
	time       float64
	boostRooms []CalcResultBoost
}

func calcTwoBoost(roomList []string, splits map[string]Room) ([]calcResult, error) {
	if strings.ToLower(roomList[len(roomList)-1]) != "finish room" {
		err := fmt.Errorf("last room is supposed to be finish room. this is a programming error")
		log.Warn(err)
		return nil, err
	}

	boostlessTime := calcBoostless(roomList, splits)
	results := make([]calcResult, 0, 81)

	for i := 0; i < 9; i++ {
		for j := i + 1; j < 9; j++ {
			firstBoostRoom := splits[roomList[i]]
			secondBoostRoom := splits[roomList[j]]

			timeBetweenBoosts := 0.0
			for k := i + 1; k < j; k++ {
				timeBetweenBoosts += splits[roomList[k]].BoostlessTime
			}

			for firstBoostStrat := 0; firstBoostStrat < len(firstBoostRoom.BoostStrats); firstBoostStrat++ {
				for secondBoostStrat := 0; secondBoostStrat < len(secondBoostRoom.BoostStrats); secondBoostStrat++ {
					pacelock := max(0, 60-(timeBetweenBoosts+firstBoostRoom.BoostStrats[firstBoostStrat].Time-firstBoostRoom.BoostStrats[firstBoostStrat].BoostTime+secondBoostRoom.BoostStrats[secondBoostStrat].BoostTime))

					boostTime := boostlessTime - (firstBoostRoom.BoostlessTime - firstBoostRoom.BoostStrats[firstBoostStrat].Time) - (secondBoostRoom.BoostlessTime - secondBoostRoom.BoostStrats[secondBoostStrat].Time) + pacelock

					boostStrat := []CalcResultBoost{
						{
							Ind:      i,
							StratInd: firstBoostStrat,
						},
						{
							Ind:      j,
							StratInd: secondBoostStrat,
							Pacelock: pacelock,
						},
					}

					results = append(results, calcResult{
						time:       boostTime,
						boostRooms: boostStrat,
					})
				}
			}
		}
	}

	slices.SortFunc(results, func(a, b calcResult) int {
		if a.time < b.time {
			return -1
		}

		return 1
	})

	return results, nil
}

func calcThreeBoost(roomList []string, splits map[string]Room) ([]calcResult, error) {
	if strings.ToLower(roomList[len(roomList)-1]) != "finish room" {
		err := fmt.Errorf("last room is supposed to be finish room. this is a programming error")
		log.Warn(err)
		return nil, err
	}

	boostlessTime := calcBoostless(roomList, splits)
	results := make([]calcResult, 0, 729)

	for i := 0; i < 9; i++ {
		for j := i + 1; j < 9; j++ {
			for k := j + 1; k < 9; k++ {
				firstBoostRoom := splits[roomList[i]]
				secondBoostRoom := splits[roomList[j]]
				thirdBoostRoom := splits[roomList[k]]

				timeBetweenBoosts12 := 0.0
				for m := i + 1; m < j; m++ {
					timeBetweenBoosts12 += splits[roomList[m]].BoostlessTime
				}
				timeBetweenBoosts23 := 0.0
				for m := j + 1; m < k; m++ {
					timeBetweenBoosts23 += splits[roomList[m]].BoostlessTime
				}

				for firstBoostStrat := 0; firstBoostStrat < len(firstBoostRoom.BoostStrats); firstBoostStrat++ {
					for secondBoostStrat := 0; secondBoostStrat < len(secondBoostRoom.BoostStrats); secondBoostStrat++ {
						for thirdBoostStrat := 0; thirdBoostStrat < len(thirdBoostRoom.BoostStrats); thirdBoostStrat++ {
							pacelock1 := max(0, 60-(timeBetweenBoosts12+firstBoostRoom.BoostStrats[firstBoostStrat].Time-firstBoostRoom.BoostStrats[firstBoostStrat].BoostTime+secondBoostRoom.BoostStrats[secondBoostStrat].BoostTime))

							pacelock2 := max(0, 60-(timeBetweenBoosts23+secondBoostRoom.BoostStrats[secondBoostStrat].Time-secondBoostRoom.BoostStrats[secondBoostStrat].BoostTime+thirdBoostRoom.BoostStrats[thirdBoostStrat].BoostTime))
							boostTime := boostlessTime - (firstBoostRoom.BoostlessTime - firstBoostRoom.BoostStrats[firstBoostStrat].Time) - (secondBoostRoom.BoostlessTime - secondBoostRoom.BoostStrats[secondBoostStrat].Time) - (thirdBoostRoom.BoostlessTime - thirdBoostRoom.BoostStrats[thirdBoostStrat].Time) + pacelock1 + pacelock2

							boostStrat := []CalcResultBoost{
								{
									Ind:      i,
									StratInd: firstBoostStrat,
								},
								{
									Ind:      j,
									StratInd: secondBoostStrat,
									Pacelock: pacelock1,
								},
								{
									Ind:      k,
									StratInd: thirdBoostStrat,
									Pacelock: pacelock2,
								},
							}

							results = append(results, calcResult{
								time:       boostTime,
								boostRooms: boostStrat,
							})
						}
					}
				}
			}
		}
	}

	slices.SortFunc(results, func(a, b calcResult) int {
		if a.time < b.time {
			return -1
		}

		return 1
	})

	return results, nil
}

type CalcSeedResult struct {
	BoostlessTime float64
	BoostTime     float64
	BoostRooms    []CalcResultBoost
}

func mergeSortedResults(a, b []calcResult) []calcResult {
	merged := make([]calcResult, 0, len(a)+len(b))
	i, j := 0, 0

	for i < len(a) && j < len(b) {
		if a[i].time <= b[j].time {
			merged = append(merged, a[i])
			i++
		} else {
			merged = append(merged, b[j])
			j++
		}
	}

	// Append remaining elements from either array
	merged = append(merged, a[i:]...)
	merged = append(merged, b[j:]...)

	return merged
}

func calcSeedInternal(roomList []string, splits map[string]Room) ([]CalcSeedResult, error) {
	boostlessTime := calcBoostless(roomList, splits)

	res := make([]CalcSeedResult, 0, 5)

	twoBoost, err := calcTwoBoost(roomList, splits)
	if err != nil {
		log.Warn(err)
		return nil, err
	}

	if len(twoBoost) == 0 {
		err := fmt.Errorf("two boost calculation returned an empty array")
		log.Warn(err)
		return nil, err
	}

	threeBoost, err := calcThreeBoost(roomList, splits)
	if err != nil {
		log.Warn(err)
		return nil, err
	}

	if len(threeBoost) == 0 {
		err := fmt.Errorf("three boost calculation returned an empty array")
		log.Warn(err)
		return nil, err
	}

	for _, r := range mergeSortedResults(twoBoost, threeBoost) {
		res = append(res, CalcSeedResult{
			BoostlessTime: boostlessTime,
			BoostTime:     r.time,
			BoostRooms:    r.boostRooms,
		})
	}

	return res, nil
}


func CalcSeed(roomList []string) ([]CalcSeedResult, error) {
	if roomList[len(roomList)-1] != "finish room" {
		roomList = append(roomList, "finish room")
	}
	return calcSeedInternal(roomList, RoomMap)
}

func CalcSeedCustom(roomList []string, splits map[string]Room) ([]CalcSeedResult, error) {
	if roomList[len(roomList)-1] != "finish room" {
		roomList = append(roomList, "finish room")
	}

	for _, r := range roomList {
		log.Debugf("%+v", splits[r])
	}

	return calcSeedInternal(roomList, splits)
}
