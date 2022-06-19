package types

// Summoner DTO according to v4
type Summoner struct {
	AccountId     string `json:"accountId"`
	ProfileIconId int    `json:"profileIconId"`
	RevisionDate  int64  `json:"revisionDate"`
	Name          string `json:"name"`
	Id            string `json:"id"`
	Puuid         string `json:"puuid"`
	SummonerLevel int64  `json:"summonerLevel"`
}
