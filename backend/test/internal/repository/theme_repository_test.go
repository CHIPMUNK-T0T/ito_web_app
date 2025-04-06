package repository_test

import (
	"CHIPMUNK-T0T/ito_web_app/test/mock"
	"reflect"
	"testing"

	"golang.org/x/exp/slices"
)

var themes = []string{
	"朝食の値段",
	"理想の気温",
	"好きな数字",
	"自分の身長",
	"一番高い買い物",
	"一番安い買い物",
	"理想の年収",
	"一番長く歩いた距離",
}

func TestThemeRepository(t *testing.T) {
	t.Run("GetRandom", func(t *testing.T) {
		repo := mock.NewThemeRepository()

		theme, err := repo.GetRandom()
		if err != nil {
			t.Errorf("予期しないエラー: %v", err)
		}
		if theme == "" {
			t.Error("テーマが空文字列です")
		}

		// 取得したテーマがthemesに含まれているか
		if !slices.Contains(themes, theme) {
			t.Errorf("テーマが含まれていません: %s", theme)
		}
	})

	t.Run("FindAll", func(t *testing.T) {
		repo := mock.NewThemeRepository()

		themes, err := repo.FindAll()
		if err != nil {
			t.Errorf("予期しないエラー: %v", err)
		}
		if len(themes) == 0 {
			t.Error("テーマリストが空です")
		}

		// 取得したテーマがthemesと同等の内容であるか
		if !reflect.DeepEqual(themes, themes) {
			t.Errorf("テーマリストが同等ではありません: %v", themes)
		}
	})
}
