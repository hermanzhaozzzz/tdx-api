package tdx

import (
	"path/filepath"
	"testing"
	"time"

	"xorm.io/core"
	"xorm.io/xorm"
)

func TestCodesNeedUpdateAtNineAMBoundary(t *testing.T) {
	location := time.FixedZone("Asia/Shanghai", 8*60*60)
	parse := func(value string) time.Time {
		result, err := time.ParseInLocation("2006-01-02 15:04:05", value, location)
		if err != nil {
			t.Fatal(err)
		}
		return result
	}
	cases := []struct {
		name, now, updated string
		want               bool
	}{
		{"before nine cached", "2026-09-24 08:59:59", "2026-09-23 09:06:11", false},
		{"before nine at previous node", "2026-09-24 08:59:59", "2026-09-23 09:00:00", false},
		{"before nine stale", "2026-09-24 08:59:59", "2026-09-23 08:59:59", true},
		{"at nine stale", "2026-09-24 09:00:00", "2026-09-23 09:06:11", true},
		{"at nine current", "2026-09-24 09:00:00", "2026-09-24 09:00:00", false},
		{"after nine stale", "2026-09-24 09:00:01", "2026-09-24 08:59:59", true},
		{"after nine current", "2026-09-24 09:00:01", "2026-09-24 09:00:00", false},
	}
	for _, item := range cases {
		t.Run(item.name, func(t *testing.T) {
			got := codesNeedUpdate(parse(item.updated).Unix(), parse(item.now))
			if got != item.want {
				t.Fatalf("got %v, want %v", got, item.want)
			}
		})
	}
}

func TestNewCodesLoadsFreshDatabaseWithoutNetwork(t *testing.T) {
	filename := filepath.Join(t.TempDir(), "codes.db")
	db, err := xorm.NewEngine("sqlite", filename)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	db.SetMapper(core.SameMapper{})
	if err := db.Sync2(new(CodeModel), new(UpdateModel)); err != nil {
		t.Fatal(err)
	}
	if _, err := db.Insert(&CodeModel{Code: "000001", Exchange: "sz", Name: "缓存证券"}); err != nil {
		t.Fatal(err)
	}
	if _, err := db.Insert(&UpdateModel{Key: "codes", Time: time.Now().Unix()}); err != nil {
		t.Fatal(err)
	}
	// There is no network connection: a refresh would fail instead of using the cache.
	codes, err := NewCodes(&Client{}, db)
	if err != nil {
		t.Fatal(err)
	}
	if code := codes.Get("sz000001"); code == nil || code.Name != "缓存证券" {
		t.Fatalf("cached code not loaded: %#v", code)
	}
}
