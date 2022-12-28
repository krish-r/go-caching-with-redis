package main

type TitleBasics struct {
	Tconst         string `json:"tconst"`
	TitleType      string `json:"title_type"`
	PrimaryTitle   string `json:"primary_title"`
	OriginalTitle  string `json:"original_title"`
	IsAdult        string `json:"is_adult"`
	StartYear      string `json:"start_year"`
	EndYear        string `json:"end_year"`
	RuntimeMinutes string `json:"runtime_minutes"`
	Genres         string `json:"genres"`
}

type CacheClient interface {
	Close()
	Get(key string) (*TitleBasics, error)
	Set(key string, value *TitleBasics) error
}

type DbClient interface {
	Close()
	Get(key string) (*TitleBasics, error)
}
