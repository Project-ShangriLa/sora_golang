package AnimeAPI

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
	"net/http"
	"strconv"
	"time"
)

// https://cloud.google.com/appengine/docs/go/datastore/reference

var router = mux.NewRouter()

func init() {
	http.Handle("/", router)

	// http://www.gorillatoolkit.org/pkg/mux
	router.HandleFunc("/anime/v1/master/{year_num:[0-9]{4}}/{cours:[1-4]}", animeAPIReadHandler).Methods("GET")

	/*
		Anime API sora-scala play-framework routing
		GET        /anime/v1/master/cours                                 controllers.AnimeV1.masterList
		GET        /anime/v1/master/$year_num<[0-9]{4}+>                  controllers.AnimeV1.year(year_num)
		GET        /anime/v1/master/$year_num<[0-9]{4}+>/$cours<[1-4]+>   controllers.AnimeV1.yearCours(year_num, cours)
	*/

	router.HandleFunc("/anime/v1/master", animeAPIPutHandler).Methods("PUT")

	//router.HandleFunc("/anime/v1/admin/set", setAdminKeyHandler).Methods("GET")
}

// AnimeData Animeマスター情報 MySQLのBasesテーブルに相当
type AnimeData struct {
	BasesID        int       `json:"id"`
	Title          string    `json:"title"`
	TitleShort1    string    `json:"title_short1"`
	TitleShort2    string    `json:"title_short2"`
	TitleShort3    string    `json:"title_short3"`
	PublicURL      string    `json:"public_url"`
	TwitterAccount string    `json:"twitter_account"`
	TwitterHashTag string    `json:"twitter_hash_tag"`
	CoursID        int       `json:"cours_id"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	Sex            int       `json:"sex"`
	Sequel         int       `json:"sequel"`
	CityCode       int       `json:"city_code"`
	CityName       string    `json:"city_name"`
}

// http://localhost:8080/anime/v1/master/2014/2
func animeAPIReadHandler(w http.ResponseWriter, r *http.Request) {

	coursID := year2coursID(r)

	ctx := appengine.NewContext(r)
	q := datastore.NewQuery("bases").Filter("CoursID =", coursID)

	log.Debugf(ctx, "--> %d", coursID)

	var bases []AnimeData
	_, err := q.GetAll(ctx, &bases)
	res, err := json.Marshal(bases)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(res)
}

// TODO 暫定的に手動で計算、本来は管理テーブルからcours_idを算出
func year2coursID(r *http.Request) int {
	vars := mux.Vars(r)

	year, _ := strconv.Atoi(vars["year_num"])
	cours, _ := strconv.Atoi(vars["cours"])
	coursID := (year-2014)*4 + cours

	return coursID
}

// ---- ここまでがsora互換 以降は新規機能でデータマイグレーションのためPUT機能追加

// curl -X PUT -d '{"id": 1}' http://localhost:8080/anime/v1/master
// curl -X PUT -d '{"id": 1, "title": "がりなん", "created_at": "2014-08-25T00:00:00+09:00" }' http://localhost:8080/anime/v1/master
// Time型はRCF3339
// http://qiita.com/taizo/items/2c3a338f1aeea86ce9e2
// curl -H "X-ANIME-API-ADMIN-KEY:zura" -X PUT -d '{"id": 2, "title": "がりなん", "created_at": "2014-08-25T00:00:00+09:00" }' http://localhost:8080/anime/v1/master
func animeAPIPutHandler(w http.ResponseWriter, r *http.Request) {
	//https://gist.github.com/andreagrandi/97263aaf7f9344d3ffe6

	if !validAnimeAPIAdminKey(r) {
		return
	}

	var animeData AnimeData
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&animeData)

	if err != nil {
		panic(err)
	}

	ctx := appengine.NewContext(r)

	// context, entity, stringID, intID, namespace
	key := datastore.NewKey(ctx, "bases", "", int64(animeData.BasesID), nil)

	datastore.Put(ctx, key, &animeData)
}

// Admin 管理者情報
type Admin struct {
	AdminKey string
}

func validAnimeAPIAdminKey(r *http.Request) bool {

	masterAPIKey := r.Header.Get("X-ANIME-API-ADMIN-KEY")

	ctx := appengine.NewContext(r)
	q := datastore.NewQuery("anime_api_admin").Filter("AdminKey =", masterAPIKey)

	log.Debugf(ctx, "master_api_key input=[%s]", masterAPIKey)

	var admin []Admin
	_, err := q.GetAll(ctx, &admin)

	if err != nil {
		panic(err)
	}

	if len(admin) == 0 {
		log.Debugf(ctx, "len(animeAPIAdmin) = %d ", len(admin))
		return false
	}

	return true
}

// 開発環境で直接Entity生成できない場合などに使用
// pwgen 50 1
// http://localhost:8080/anime/v1/admin/set?key=XXXXX
/*
func setAdminKeyHandler(w http.ResponseWriter, r *http.Request) {

	var admin Admin
	admin.AdminKey = r.FormValue("key")
	ctx := appengine.NewContext(r)
	key := datastore.NewKey(ctx, "anime_api_admin", "master_api_key", 0, nil)
	datastore.Put(ctx, key, &admin)
}
*/