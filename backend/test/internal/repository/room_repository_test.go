package repository_test

import (
	"CHIPMUNK-T0T/ito_web_app/internal/entity/domain"
	"CHIPMUNK-T0T/ito_web_app/internal/repository"
	"CHIPMUNK-T0T/ito_web_app/test/mock"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestRoomRepository(t *testing.T) {
	t.Run("Create", func(t *testing.T) {
		db, mock, err := mock.NewMockDB()
		if err != nil {
			t.Fatalf("モックデータベースの作成に失敗しました: %v", err)
		}

		repo := repository.NewRoomRepository(db)
		room, _ := domain.NewRoom("testroom", "password", 4, 1, "description", false)

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `rooms` (`created_at`,`updated_at`,`deleted_at`,`name`,`password`,`max_players`,`creator_id`,`description`,`is_private`) VALUES (?,?,?,?,?,?,?,?,?)")).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), room.Name(), room.Password(), room.MaxPlayers(), room.CreatorID(), room.Description(), room.IsPrivate()).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err = repo.Create(&room)
		if err != nil {
			t.Errorf("予期しないエラー: %v", err)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("未実行のクエリが存在します: %v", err)
		}
	})

	t.Run("FindByID", func(t *testing.T) {
		db, mock, err := mock.NewMockDB()
		if err != nil {
			t.Fatalf("モックデータベースの作成に失敗しました: %v", err)
		}

		repo := repository.NewRoomRepository(db)

		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `rooms` WHERE `rooms`.`id` = ? AND `rooms`.`deleted_at` IS NULL ORDER BY `rooms`.`id` LIMIT ?")).
			WithArgs(1, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "password", "max_players", "creator_id", "description", "is_private", "created_at", "updated_at", "deleted_at"}).
				AddRow(1, "testroom", "password", 4, 1, "description", false, time.Now(), time.Now(), nil))

		room, err := repo.FindByID(1)
		if err != nil {
			t.Errorf("予期しないエラー: %v", err)
		}
		if room.Name() != "testroom" {
			t.Errorf("期待値: testroom, 実際: %s", room.Name())
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("未実行のクエリが存在します: %v", err)
		}
	})

	t.Run("FindByRoomNameAndPassword", func(t *testing.T) {
		db, mock, err := mock.NewMockDB()
		if err != nil {
			t.Fatalf("モックデータベースの作成に失敗しました: %v", err)
		}

		repo := repository.NewRoomRepository(db)

		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `rooms` WHERE (name = ? AND password = ?) AND `rooms`.`deleted_at` IS NULL ORDER BY `rooms`.`id` LIMIT ?")).
			WithArgs("testroom", "password", 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "password", "max_players", "creator_id", "description", "is_private", "created_at", "updated_at", "deleted_at"}).
				AddRow(1, "testroom", "password", 4, 1, "description", false, time.Now(), time.Now(), nil))

		room, err := repo.FindByRoomNameAndPassword("testroom", "password")
		if err != nil {
			t.Errorf("予期しないエラー: %v", err)
		}
		if room.Name() != "testroom" {
			t.Errorf("期待値: testroom, 実際: %s", room.Name())
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("未実行のクエリが存在します: %v", err)
		}
	})

	t.Run("FindAll", func(t *testing.T) {
		db, mock, err := mock.NewMockDB()
		if err != nil {
			t.Fatalf("モックデータベースの作成に失敗しました: %v", err)
		}

		repo := repository.NewRoomRepository(db)

		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `rooms` WHERE `rooms`.`deleted_at` IS NULL")).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "password", "max_players", "creator_id", "description", "is_private", "created_at", "updated_at", "deleted_at"}).
				AddRow(1, "room1", "pass1", 4, 1, "desc1", false, time.Now(), time.Now(), nil).
				AddRow(2, "room2", "pass2", 6, 2, "desc2", true, time.Now(), time.Now(), nil))

		rooms, err := repo.FindAll()
		if err != nil {
			t.Errorf("予期しないエラー: %v", err)
		}
		if len(rooms) != 2 {
			t.Errorf("期待値: 2, 実際: %d", len(rooms))
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("未実行のクエリが存在します: %v", err)
		}
	})

	t.Run("Update", func(t *testing.T) {
		db, mock, err := mock.NewMockDB()
		if err != nil {
			t.Fatalf("モックデータベースの作成に失敗しました: %v", err)
		}

		repo := repository.NewRoomRepository(db)
		room := domain.NewRoomWithID(1, "updated", "newpassword", 6, 1, "new description", true)

		// 更新前のルーム存在確認のSELECTクエリ
		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `rooms` WHERE `rooms`.`id` = ? AND `rooms`.`deleted_at` IS NULL ORDER BY `rooms`.`id` LIMIT ?")).
			WithArgs(1, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "password", "max_players", "creator_id", "description", "is_private", "created_at", "updated_at", "deleted_at"}).
				AddRow(1, "testroom", "oldpassword", 4, 1, "old description", false, time.Now(), time.Now(), nil))

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta("UPDATE `rooms` SET `created_at`=?,`updated_at`=?,`deleted_at`=?,`name`=?,`password`=?,`max_players`=?,`creator_id`=?,`description`=?,`is_private`=? WHERE `rooms`.`deleted_at` IS NULL AND `id` = ?")).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), room.Name(), room.Password(), room.MaxPlayers(), room.CreatorID(), room.Description(), room.IsPrivate(), room.ID()).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err = repo.Update(room)
		if err != nil {
			t.Errorf("予期しないエラー: %v", err)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("未実行のクエリが存在します: %v", err)
		}
	})

	// t.Run("Delete", func(t *testing.T) {
	// 	db, mock, err := mock.NewMockDB()
	// 	if err != nil {
	// 		t.Fatalf("モックデータベースの作成に失敗しました: %v", err)
	// 	}

	// 	repo := repository.NewRoomRepository(db)

	// 	// 削除前のルーム存在確認のSELECTクエリ
	// 	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `rooms` WHERE `rooms`.`id` = ? AND `rooms`.`deleted_at` IS NULL ORDER BY `rooms`.`id` LIMIT ?")).
	// 		WithArgs(1, 1).
	// 		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "password", "max_players", "creator_id", "description", "is_private", "created_at", "updated_at", "deleted_at"}).
	// 			AddRow(1, "testroom", "password", 4, 1, "description", false, time.Now(), time.Now(), nil))

	// 	mock.ExpectBegin()
	// 	mock.ExpectExec(regexp.QuoteMeta("UPDATE `rooms` SET `deleted_at`=? WHERE `rooms`.`id` = ? AND `rooms`.`deleted_at` IS NULL")).
	// 		WithArgs(sqlmock.AnyArg(), 1).
	// 		WillReturnResult(sqlmock.NewResult(1, 1))
	// 	mock.ExpectCommit()

	// 	err = repo.Delete(1)
	// 	if err != nil {
	// 		t.Errorf("予期しないエラー: %v", err)
	// 	}

	// 	// 削除確認用のクエリ
	// 	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `rooms` WHERE `rooms`.`id` = ? AND `rooms`.`deleted_at` IS NULL ORDER BY `rooms`.`id` LIMIT ?")).
	// 		WithArgs(1, 1).
	// 		WillReturnError(sql.ErrNoRows)

	// 	_, err = repo.FindByID(1)
	// 	if err == nil {
	// 		t.Error("削除されたルームが見つかりました")
	// 	}

	// 	if err := mock.ExpectationsWereMet(); err != nil {
	// 		t.Errorf("未実行のクエリが存在します: %v", err)
	// 	}
	// })
}