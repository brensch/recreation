package proxy

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/brensch/recreation"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// using globals since i'm using cloud functions and can't pass things like i would regular handler funcs
var (
	o   recreation.Obfuscator
	log *zap.Logger
)

func init() {
	o = recreation.InitAgentRandomiser(context.Background())
	logConfig := zap.NewProductionConfig()
	logConfig.Level.SetLevel(zap.DebugLevel)
	// this ensures google logs pick things up properly
	logConfig.EncoderConfig.MessageKey = "message"
	logConfig.EncoderConfig.LevelKey = "severity"
	logConfig.EncoderConfig.TimeKey = "time"
	// logConfig.Encoding = "console"

	// init logger
	var err error
	log, err = logConfig.Build()
	if err != nil {
		// this indicates a bug or some way that zap can fail i'm not aware of
		panic(err)
	}
}

type GetAvailabilityReq struct {
	CampgroundID string `json:"campground_id,omitempty"`
}

// CloudFunctionGetAvailability is intended to be run as a cloud function in GCP.
// I am spreading these across all the regions of GCP so that they will all have different IP addresses.
// Cloud functions are extremely cheap so the theory is that this will actually lead to less resource usage.
func CloudFunctionGetAvailability(w http.ResponseWriter, r *http.Request) {
	uuid := uuid.New().String()
	log = log.With(zap.String("request_id", uuid))

	var req GetAvailabilityReq
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Error("failed to decode request", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	avail, err := recreation.GetAvailability(context.Background(), log, o, req.CampgroundID, time.Now())
	if err != nil {
		log.Error("failed to get availability", zap.Error(err))
		w.WriteHeader(http.StatusTeapot)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(avail)
	if err != nil {
		log.Error("failed to encode avail into response", zap.Error(err))
		w.WriteHeader(http.StatusTeapot)
		return
	}

}

var (
	sites = []string{
		"273757",
		"10172170",
		"232491",
		"272229",
		"258815",
		"10067346",
		"233359",
		"273819",
		"233754",
		"273878",
		"234015",
		"250000",
		"231958",
		"270421",
		"251550",
		"233683",
		"233779",
		"273869",
		"273874",
		"234073",
		"234072",
		"233116",
		"233118",
		"233431",
		"233664",
		"233537",
		"231959",
		"233775",
		"234761",
		"233701",
		"233414",
		"232262",
		"233509",
		"234654",
		"233440",
		"233839",
		"233842",
		"232451",
		"232452",
		"233772",
		"234756",
		"233981",
		"232254",
		"118995",
		"232293",
		"233879",
		"234111",
		"234172",
		"251615",
		"10077451",
		"232260",
		"233285",
		"232107",
		"232021",
		"232446",
		"232343",
		"234133",
		"274410",
		"10005253",
		"232453",
		"234006",
		"232061",
		"245489",
		"232083",
		"274314",
		"232909",
		"234135",
		"232910",
		"232366",
		"10004152",
		"232878",
		"232881",
		"233708",
		"232125",
		"233768",
		"232121",
		"233162",
		"232879",
		"232802",
		"232912",
		"232880",
		"232045",
		"232801",
		"232888",
		"232882",
		"232349",
		"232450",
		"232449",
		"232887",
		"232911",
		"232447",
		"233437",
		"233532",
		"233129",
		"251578",
		"231954",
		"232049",
		"233104",
		"231953",
		"247867",
		"234600",
		"232070",
		"232348",
		"232302",
		"232263",
		"232321",
		"10040012",
		"232047",
		"10040047",
		"232818",
		"232877",
		"232820",
		"232187",
		"232264",
		"10039993",
		"10039887",
		"234591",
		"233180",
		"232817",
		"232032",
		"234587",
		"232048",
		"10039838",
		"10040022",
		"232448",
		"234739",
		"232769",
		"232821",
		"233183",
		"233117",
		"232876",
		"272246",
		"233521",
		"232185",
		"232906",
		"233314",
		"234592",
		"232810",
		"234589",
		"232261",
		"234457",
		"234548",
		"234547",
		"232875",
		"234549",
		"234542",
		"234538",
		"232422",
		"232768",
		"232119",
		"232811",
		"232814",
		"232813",
		"232874",
		"232079",
		"232812",
		"232090",
		"233102",
		"232815",
		"232058",
		"232117",
		"232367",
		"232053",
		"234543",
		"232136",
		"232884",
		"232808",
		"234330",
		"232755",
		"233568",
		"232398",
		"231956",
		"231957",
		"233235",
		"232883",
		"232757",
		"232268",
		"232756",
		"232759",
		"234534",
		"232269",
		"233161",
		"233439",
		"233130",
		"232805",
		"232767",
		"232766",
		"232806",
		"234541",
		"233404",
		"234329",
		"232804",
		"231955",
		"234290",
		"273846",
		"234544",
		"273872",
		"232803",
		"233860",
		"233830",
		"232270",
		"234117",
		"234210",
		"231977",
		"234535",
		"273870",
		"234536",
		"232271",
		"251445",
		"251446",
		"233729",
		"232267",
		"233728",
		"234539",
		"234537",
		"234540",
		"234114",
		"231963",
		"234545",
		"234546",
		"232907",
		"232908",
		"10114327",
		"232371",
		"232869",
		"234115",
		"232871",
		"10124502",
		"234116",
		"10114366",
		"234752",
		"10124445",
		"251008",
		"234311",
		"232782",
		"234113",
		"232186",
		"232858",
		"232266",
		"232859",
		"233363",
		"232781",
		"10110742",
		"234738",
		"232780",
		"252037",
		"232777",
		"234131",
		"10114392",
		"251363",
		"232238",
		"233907",
		"232239",
		"232396",
		"234663",
		"234458",
		"233101",
		"256932",
		"232784",
		"232783",
		"233692",
		"234672",
		"232785",
		"245528",
		"273836",
		"273833",
		"245526",
		"245525",
		"245541",
		"245494",
		"245539",
		"245538",
		"245552",
		"245564",
		"245555",
		"245563",
		"245493",
		"245565",
		"245566",
		"245567",
		"242337",
		"245561",
		"242333",
		"245554",
		"242338",
		"245558",
		"14806",
		"245568",
		"242340",
		"14807",
		"14808",
		"245560",
		"245551",
		"245559",
		"245556",
		"245570",
		"242349",
		"243903",
	}
)
