// Package shukujitsu は内閣府が提供している祝日一覧csvファイルを取得・解析します
package shukujitsu

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

// import 省略

type dateint struct {
	Year  int
	Month int
	Day   int
}

// Entry は祝日1日分の情報を保持する構造体です
type Entry struct {
	YMD  string
	Date dateint
	Name string
}

const csvURL = "https://www8.cao.go.jp/chosei/shukujitsu/syukujitsu.csv"

// AllEntriesは内閣府ウェブサイトから祝日 CSV を取得してEntryスライスに変換します。
func AllEntries() ([]Entry, error) {
	resp, err := http.Get(csvURL)
	if err != nil {
		return nil, fmt.Errorf("接続に失敗しました: %w", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("データの取得に失敗しました: %w", err)
	}

	records, err := csv.NewReader(transform.NewReader(bytes.NewReader(body), japanese.ShiftJIS.NewDecoder())).ReadAll()
	if err != nil {
		return nil, fmt.Errorf("データ解析に失敗しました: %w", err)
	}
	var entries []Entry
	for i, row := range records {
		if i == 0 {
			continue // ヘッダーをスキップ
		}
		if len(row) != 2 {
			return nil, fmt.Errorf("想定外のデータに遭遇しました：行 %d = %v", i+1, row)
		}
		var d dateint
		arr := strings.Split(row[0], "/")
		if len(arr) != 3 {
			return nil, fmt.Errorf("想定外の日付のフォーマットです：行 %d = %s", i+1, row[0])
		}
		d.Year, _ = strconv.Atoi(arr[0])
		d.Month, _ = strconv.Atoi(arr[1])
		d.Day, _ = strconv.Atoi(arr[2])
		entries = append(entries, Entry{YMD: row[0], Name: row[1], Date: d})
	}
	return entries, nil
}
