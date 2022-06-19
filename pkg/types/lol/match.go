package types

// Match reflects the MatchDTO object according to riot api documentation
type Match struct {
	MetaData MetaData `bson:"metaData" json:"metadata"`
	Info     Info     `json:"info"`
}

type MetaData struct {
	Dataversion  string   `json:"dataversion"`
	MatchID      string   `json:"matchId"`
	Participants []string `json:"participants"`
}

type Info struct {
	GameCreation       int64         `bson:"gameCreation" json:"gameCreation"`
	GameDuration       int64         `bson:"gameDuration" json:"gameDuration"`
	GameEndTimestamp   int64         `bson:"gameEndTimestamp" json:"gameEndTimestamp"`
	GameID             int64         `bson:"gameId" json:"gameId"`
	GameMode           string        `bson:"gameMode" json:"gameMode"`
	GameName           string        `bson:"gameName" json:"gameName"`
	GameStartTimestamp int64         `bson:"gameStartTimestamp" json:"gameStartTimestamp"`
	GameType           string        `bson:"gameType" json:"gameType"`
	GameVersion        string        `bson:"gameVersion" json:"gameVersion"`
	MapID              int           `bson:"mapId" json:"mapId"`
	Participants       []Participant `bson:"participants" json:"participants"`
	PlatformID         string        `bson:"platformId" json:"platformId"`
	QueueID            int           `bson:"queueId" json:"queueId"`
	Teams              []Team        `bson:"teams" json:"teams"`
	TournamentCode     string        `bson:"tournamentCode" json:"tournamentCode"`
}

type Participant struct {
	Assists                        int    `json:"assists"`
	Baronkills                     int    `json:"baronKills"`
	Bountylevel                    int    `json:"bountyLevel"`
	Champexperience                int    `json:"champExperience"`
	Champlevel                     int    `json:"champLevel"`
	Championid                     int    `json:"championId"` // Prior to patch 11.4, on Feb 18th, 2021, this field returned invalid championIds. We recommend determining the champion based on the championName field for matches played prior to patch 11.4.
	Championname                   string `json:"championName"`
	Championtransform              int    `json:"championTransform"` // This field is currently only utilized for Kayn's transformations. (Legal values: 0 - None, 1 - Slayer, 2 - Assassin)
	Consumablespurchased           int    `json:"consumablesPurchased"`
	Damagedealttobuildings         int    `json:"damageDealtToBuildings"`
	Damagedealttoobjectives        int    `json:"damageDealtToObjectives"`
	Damagedealttoturrets           int    `json:"damageDealtToTurrets"`
	Damageselfmitigated            int    `json:"damageSelfMitigated"`
	Deaths                         int    `json:"deaths"`
	Detectorwardsplaced            int    `json:"detectorWardsPlaced"`
	Doublekills                    int    `json:"doubleKills"`
	Dragonkills                    int    `json:"dragonKills"`
	Firstbloodassist               bool   `json:"firstBloodAssist"`
	Firstbloodkill                 bool   `json:"firstBloodKill"`
	Firsttowerassist               bool   `json:"firstTowerAssist"`
	Firsttowerkill                 bool   `json:"firstTowerKill"`
	Gameendedinearlysurrender      bool   `json:"gameEndedInEarlySurrender"`
	Gameendedinsurrender           bool   `json:"gameEndedInSurrender"`
	Goldearned                     int    `json:"goldEarned"`
	Goldspent                      int    `json:"goldSpent"`
	Individualposition             string `json:"individualPosition"` // Both individualPosition and teamPosition are computed by the game server and are different versions of the most likely position played by a player. The individualPosition is the best guess for which position the player actually played in isolation of anything else. The teamPosition is the best guess for which position the player actually played if we add the constraint that each team must have one top player, one jungle, one middle, etc. Generally the recommendation is to use the teamPosition field over the individualPosition field.
	Inhibitorkills                 int    `json:"inhibitorKills"`
	Inhibitortakedowns             int    `json:"inhibitorTakedowns"`
	Inhibitorslost                 int    `json:"inhibitorsLost"`
	Item0                          int    `json:"item0"`
	Item1                          int    `json:"item1"`
	Item2                          int    `json:"item2"`
	Item3                          int    `json:"item3"`
	Item4                          int    `json:"item4"`
	Item5                          int    `json:"item5"`
	Item6                          int    `json:"item6"`
	Itemspurchased                 int    `json:"itemsPurchased"`
	Killingsprees                  int    `json:"killingSprees"`
	Kills                          int    `json:"kills"`
	Lane                           string `json:"lane"`
	Largestcriticalstrike          int    `json:"largestCriticalStrike"`
	Largestkillingspree            int    `json:"largestKillingSpree"`
	Largestmultikill               int    `json:"largestMultiKill"`
	Longesttimespentliving         int    `json:"longestTimeSpentLiving"`
	Magicdamagedealt               int    `json:"magicDamageDealt"`
	Magicdamagedealttochampions    int    `json:"magicDamageDealtToChampions"`
	Magicdamagetaken               int    `json:"magicDamageTaken"`
	Neutralminionskilled           int    `json:"neutralMinionsKilled"`
	Nexuskills                     int    `json:"nexusKills"`
	Nexustakedowns                 int    `json:"nexusTakedowns"`
	Nexuslost                      int    `json:"nexuslost"`
	Objectivesstolen               int    `json:"objectivesStolen"`
	Objectivesstolenassists        int    `json:"objectivesStolenAssists"`
	Participantid                  int    `json:"participantId"`
	Pentakills                     int    `json:"pentaKills"`
	Perks                          Perk   `json:"perks"`
	Physicaldamagedealt            int    `json:"physicalDamageDealt"`
	Physicaldamagedealttochampions int    `json:"physicalDamageDealtToChampions"`
	Physicaldamagetaken            int    `json:"physicalDamageTaken"`
	Profileicon                    int    `json:"profileIcon"`
	Puuid                          string `json:"puuid"`
	Quadrakills                    int    `json:"quadraKills"`
	Riotidname                     string `json:"riotIdName"`
	Riotidtagline                  string `json:"riotIdTagline"`
	Role                           string `json:"role"`
	Sightwardsboughtingame         int    `json:"sightWardsBoughtInGame"`
	Spell1casts                    int    `json:"spell1Casts"`
	Spell2casts                    int    `json:"spell2Casts"`
	Spell3casts                    int    `json:"spell3Casts"`
	Spell4casts                    int    `json:"spell4Casts"`
	Summoner1casts                 int    `json:"summoner1Casts"`
	Summoner1id                    int    `json:"summoner1Id"`
	Summoner2casts                 int    `json:"summoner2Casts"`
	Summoner2id                    int    `json:"summoner2Id"`
	Summonerid                     string `json:"summonerId"`
	Summonerlevel                  int    `json:"summonerLevel"`
	Summonername                   string `json:"summonerName"`
	Teamearlysurrendered           bool   `json:"teamEarlySurrendered"`
	Teamid                         int    `json:"teamId"`
	Teamposition                   string `json:"teamPosition"` // 	Both individualPosition and teamPosition are computed by the game server and are different versions of the most likely position played by a player. The individualPosition is the best guess for which position the player actually played in isolation of anything else. The teamPosition is the best guess for which position the player actually played if we add the constraint that each team must have one top player, one jungle, one middle, etc. Generally the recommendation is to use the teamPosition field over the individualPosition field.
	Timeccingothers                int    `json:"timeCCingOthers"`
	Timeplayed                     int    `json:"timePlayed"`
	Totaldamagedealt               int    `json:"totalDamageDealt"`
	Totaldamagedealttochampions    int    `json:"totalDamageDealtToChampions"`
	Totaldamageshieldedonteammates int    `json:"totalDamageShieldedOnTeammates"`
	Totaldamagetaken               int    `json:"totalDamageTaken"`
	Totalheal                      int    `json:"totalHeal"`
	Totalhealsonteammates          int    `json:"totalHealsOnTeammates"`
	Totalminionskilled             int    `json:"totalMinionsKilled"`
	Totaltimeccdealt               int    `json:"totalTimeCCDealt"`
	Totaltimespentdead             int    `json:"totalTimeSpentDead"`
	Totalunitshealed               int    `json:"totalUnitsHealed"`
	Triplekills                    int    `json:"tripleKills"`
	Truedamagedealt                int    `json:"trueDamageDealt"`
	Truedamagedealttochampions     int    `json:"trueDamageDealtToChampions"`
	Truedamagetaken                int    `json:"trueDamageTaken"`
	Turretkills                    int    `json:"turretKills"`
	Turrettakedowns                int    `json:"turretTakedowns"`
	Turretslost                    int    `json:"turretsLost"`
	Unrealkills                    int    `json:"unrealKills"`
	Visionscore                    int    `json:"visionScore"`
	Visionwardsboughtingame        int    `json:"visionWardsBoughtInGame"`
	Wardskilled                    int    `json:"wardsKilled"`
	Wardsplaced                    int    `json:"wardsPlaced"`
	Win                            bool   `json:"win"`
}

type Perk struct {
	StatPerks PerkStats   `json:"statPerks"`
	Styles    []PerkStyle `json:"styles"`
}

type PerkStats struct {
	Defense int `json:"defense"`
	Flex    int `json:"flex"`
	Offense int `json:"offense"`
}

type PerkStyle struct {
	Description string               `json:"description"`
	Selections  []PerkStyleSelection `json:"selections"`
	Style       int                  `json:"style"`
}

type PerkStyleSelection struct {
	Perk int `json:"perk"`
	Var1 int `json:"var1"`
	Var2 int `json:"var2"`
	Var3 int `json:"var3"`
}

type Team struct {
	Bans       []Ban     `json:"bans"`
	Objectives Objective `json:"objectives"`
	Teamid     int       `json:"teamId"`
	Win        bool      `json:"win"`
}

type Ban struct {
	ChampionId int `json:"championId"`
	PickTurn   int `json:"pickTurn"`
}

type Objectives struct {
	Baron      Objective `json:"baron"`
	Champion   Objective `json:"champion"`
	Dragon     Objective `json:"dragon"`
	Inhibitor  Objective `json:"inhibitor"`
	Riftherald Objective `json:"riftHerald"`
	Tower      Objective `json:"tower"`
}

type Objective struct {
	First bool `json:"first"`
	Kills int  `json:"kills"`
}
