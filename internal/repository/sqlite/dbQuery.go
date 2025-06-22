package sqlite

import (
	"cash_piggy_bank/internal/botpkg"
	"cash_piggy_bank/internal/domain"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	_ "github.com/mattn/go-sqlite3" // Драйвер SQLite
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
)

var dbOn int

type Goal struct {
	ID      int
	Name    string
	Current float64
	Target  float64
}

// Рисуем прогресс бар
func progressBar(percent int) string {
	blocks := 10 // Количество блоков в шкале
	if percent >= 100 {
		percent = 100
	}
	filled := int(float64(percent) / 100 * float64(blocks))
	return strings.Repeat("▓", filled) + strings.Repeat("░", blocks-filled) + fmt.Sprintf(" %d%%", percent)
}

// Добавляем новую цель
func addGoal(bot *telego.Bot, ctx *th.Context, message telego.Message, db *sql.DB, userID int64, goalName string, currentAmount, targetAmount float64) error {
	goals, err := getUserGoals(db, message.Chat.ID)
	if err != nil {
		log.Println("Ошибка:", err)
	} else {

		if len(goals) == 10 {
			botpkg.SendText(bot, ctx, message, "<b>🫷 У тебя добавлено максимум целей</b>\n\nБольше 10 целей добавить нельзя.\nВыполни поставленные или удали неактуальные.", "", 0)
			return errors.New("максимум целей")
		}
	}
	// Подготавливаем SQL-запрос (защита от SQL-инъекций)
	stmt, err := db.Prepare("INSERT INTO savings (user_id, goal, amount, target_amount, complete) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Выполняем запрос с параметрами
	_, err = stmt.Exec(userID, goalName, currentAmount, targetAmount, 0)
	return err
}

// Проверяем все записи
func getUserGoals(db *sql.DB, userID int64) ([]Goal, error) {
	// Запрашиваем все записи пользователя
	rows, err := db.Query("SELECT id, goal, amount, target_amount FROM savings WHERE user_id = ?", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var goals []Goal

	// Читаем строки одну за другой
	for rows.Next() {
		var id int
		var goal string
		var amount, target float64

		// Сканируем данные из строки в переменные
		err = rows.Scan(&id, &goal, &amount, &target)
		if err != nil {
			return nil, err
		}

		goals = append(goals, Goal{
			ID:      id,
			Name:    goal,
			Current: amount,
			Target:  target,
		})
	}

	return goals, nil
}

// Обновляем цель
func updateGoalAmount(db *sql.DB, bot *telego.Bot, ctx *th.Context, message telego.Message, userID int64, goalID int, addAmount float64) error {
	var (
		currentAmount float64
		targetAmount  float64
		goalUserID    int64
		complete      int64
	)

	err := db.QueryRow(`
        SELECT amount, target_amount, user_id, complete
        FROM savings 
        WHERE id = ?`,
		goalID,
	).Scan(&currentAmount, &targetAmount, &goalUserID, &complete)

	if err != nil {
		botpkg.SendText(bot, ctx, message, "<b>❌ Такой цели нет</b>\nНачни обновлять сначала и укажи цель корректно.", "", 0)
		return fmt.Errorf("цель не найдена: %v", err)

	}

	// Если user_id цели не совпадает с текущим пользователем
	if goalUserID != userID {
		return fmt.Errorf("это не ваша цель")
	}

	// Обновляем сумму
	_, err = db.Exec(`
        UPDATE savings 
        SET amount = ? 
        WHERE id = ? AND user_id = ?`, // Важно: обновляем ТОЛЬКО если user_id совпадает
		addAmount, goalID, userID,
	)

	if int((currentAmount+addAmount)/targetAmount*100) >= 100 && complete == 0 {
		botpkg.SendText(bot, ctx, message, "<b>🎉 Поздравляем! Ты выполнил цель на 100%</b>\nТы можешь удалить цель или продолжать ее пополнять.", "", 0)
		// меняем значение complete
		_, err = db.Exec(`
        UPDATE savings 
        SET complete = ? 
        WHERE id = ? AND user_id = ?`, // Важно: обновляем ТОЛЬКО если user_id совпадает
			1, goalID, userID,
		)

	}
	return err
}

// Удаляем цель
func deleteGoal(db *sql.DB, userID int64, goalID int) error {
	// Удаляем ТОЛЬКО если цель принадлежит пользователю
	result, err := db.Exec(`
        DELETE FROM savings 
        WHERE id = ? AND user_id = ?`,
		goalID, userID,
	)

	if err != nil {
		return err
	}

	// Проверяем, что удалили хотя бы одну строку
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("цель не найдена или это не ваша цель")
	}

	return nil
}

func getAllUserId(db *sql.DB) []int64 {
	rows, err := db.Query("SELECT DISTINCT user_id FROM savings")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var userIDs []int64
	for rows.Next() {
		var userID int64
		if err := rows.Scan(&userID); err != nil {
			log.Fatal(err)
		}
		userIDs = append(userIDs, userID)
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	return userIDs
}

func DbQuery(command string, bot *telego.Bot, ctx *th.Context, message telego.Message, UserMap map[int64]*domain.User) []int64 {
	// Получаем путь к исполняемому файлу
	execPath, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	// Получаем директорию проекта (на два уровня выше от cmd/bot)
	projectDir := filepath.Dir(filepath.Dir(execPath))

	// Формируем путь к файлу БД
	dbPath := filepath.Join(projectDir, "internal", "repository", "sqlite", "savings.db")
	// Открываем файл БД (если нет — создаётся автоматически)
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		panic(err)
	}
	defer db.Close() // Закрываем соединение при выходе

	// Проверяем подключение
	err = db.Ping()
	if err != nil {
		panic(err)
	}

	if dbOn == 0 {

		// Создаём таблицу, если её нет
		_, err = db.Exec(`
    CREATE TABLE IF NOT EXISTS savings (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        user_id INTEGER NOT NULL,
        goal TEXT NOT NULL,
        amount REAL NOT NULL,
        target_amount REAL NOT NULL,
		complete INTEGER NOT NULL
    );
`)
		if err != nil {
			panic(err)
		}
		dbOn = 1
	}

	if command == "40d9f51c-810c-4023-ae92-c314d90aaf30" {
		return getAllUserId(db)
	}

	if command == "Накопления 🤑" {
		UserMap[message.Chat.ID].Command = "Ждем команду для БД"

		// Отправляем клавиатуру для работы с накоплениями
		botpkg.SendKeyboard(bot, ctx, message, nil, "накопления")
	}
	if command == "Добавить цель" || UserMap[message.Chat.ID].Command == "Добавляем цель" {

		if UserMap[message.Chat.ID].Command == "Добавляем цель" {

			arrString := strings.Split(command, " ")
			if len(arrString) != 2 {
				botpkg.SendText(bot, ctx, message, `❌ Что-то ты ввел не так. Название - 1 слово, пробел, сколько стоит. Например, "Машина 1200000". Попробуй еще раз`, "", 0)
				return nil

			}
			feetFloat, err := strconv.ParseFloat(strings.TrimSpace(arrString[1]), 64)

			if err != nil {
				botpkg.SendText(bot, ctx, message, `❌ Что-то ты ввел не так. Название - 1 слово, пробел, сколько стоит. Например, "Машина 1200000". Попробуй еще раз`, "", 0)
				return nil

			} else {
				err := addGoal(bot, ctx, message, db, message.Chat.ID, arrString[0], 0.0, feetFloat)
				if err != nil {
					log.Println("Ошибка:", err)
				} else {
					botpkg.SendText(bot, ctx, message, fmt.Sprintf(`<b>Я запомнил твою цель %s✅</b>

<b><i>В будущем обновляй свою цель и записывай накопления — я подскажу, сколько осталось накопить до достижения цели.</i></b>`, arrString[0]), "", 0)
				}
			}

			UserMap[message.Chat.ID].Command = "Ждем команду для БД"

		} else {
			UserMap[message.Chat.ID].Command = "Добавляем цель"
			botpkg.SendText(bot, ctx, message, `Отправь мне название цели и сумму сколько хочешь накопить через пробел, например, "Ноутбук 5000"`, "", 0)
		}

	} else if command == "Посмотреть цели" {
		goals, err := getUserGoals(db, message.Chat.ID)
		if err != nil {
			log.Println("Ошибка:", err)
		} else {
			if len(goals) == 0 {
				botpkg.SendText(bot, ctx, message, "Ты не добавил цели. Сначала добавь, а потом смотри 👀", "", 0)
				return nil
			}
			targetSend := "<b>🎯 Твои цели:</b>\n"
			for _, g := range goals {
				targetSend = targetSend + fmt.Sprintf("▫️%s: %.2f / %.2f₽ %v\n", g.Name, g.Current, g.Target, progressBar(int(g.Current/g.Target*100)))
			}
			botpkg.SendText(bot, ctx, message, targetSend, "", 0)
		}
	} else if command == "Обновить цель" || UserMap[message.Chat.ID].Command == "Ждем название цели" || UserMap[message.Chat.ID].Command == "Ждем сумму для цели" {

		if UserMap[message.Chat.ID].Command == "Ждем название цели" {

			UserMap[message.Chat.ID].DbChangeTargetId = UserMap[message.Chat.ID].TargetMap[command]
			botpkg.SendText(bot, ctx, message, fmt.Sprintf(`Давай обновим твою цель 💪🏻

<b>Напиши в ответ сумму, которую ты уже накопил на свою цель %s.</b>

Например, 500`, command), "", 0)

			UserMap[message.Chat.ID].Command = "Ждем сумму для цели"
		} else if UserMap[message.Chat.ID].Command == "Ждем сумму для цели" {
			feetFloat, err := strconv.ParseFloat(command, 64)
			if err != nil {
				botpkg.SendText(bot, ctx, message, `❌ Я не понял что ты мне отправил. Пришли сумму в формате "5000", которую надо добавить к цели`, "", 0)
				return nil
			}
			err = updateGoalAmount(db, bot, ctx, message, message.Chat.ID, UserMap[message.Chat.ID].DbChangeTargetId, feetFloat)
			if err != nil {
				log.Println("Ошибка:", err)
				UserMap[message.Chat.ID].Command = "Ждем название цели"
				return nil
			}
			botpkg.SendText(bot, ctx, message, "Цель обновлена 😉", "", 0)
			UserMap[message.Chat.ID].Command = ""
			botpkg.SendKeyboard(bot, ctx, message, nil, "накопления")
		} else {
			UserMap[message.Chat.ID].Command = "Ждем название цели"

			goals, err := getUserGoals(db, message.Chat.ID)
			if err != nil {
				log.Println("Ошибка:", err)
			} else {
				targetName := make([]string, 0)
				targetSend := "<b>👇 Пришли мне цель, которую хочешь обновить.</b>\nТвои цели:\n"
				for _, g := range goals {
					targetSend = targetSend + fmt.Sprintf("▫️%s: %.2f / %.2f₽\n", g.Name, g.Current, g.Target)
					targetName = append(targetName, g.Name)
					UserMap[message.Chat.ID].TargetMap[g.Name] = g.ID
				}

				if len(targetName) == 0 {
					botpkg.SendText(bot, ctx, message, `Ты не добавил цели. Сначала добавь, а потом обновляй 😉`, "", 0)
					UserMap[message.Chat.ID].Command = ""
					return nil
				}
				botpkg.SendText(bot, ctx, message, targetSend, "", 0)

				botpkg.SendKeyboard(bot, ctx, message, targetName, "цели")
			}

		}
	} else if command == "Удалить цель" || UserMap[message.Chat.ID].Command == "Ждем ID и удаление" {

		if UserMap[message.Chat.ID].Command == "Ждем ID и удаление" {

			err = deleteGoal(db, message.Chat.ID, UserMap[message.Chat.ID].TargetMap[command])
			if err != nil {
				botpkg.SendText(bot, ctx, message, "🤔 Не смогли удалить. Попробуй сначала", "", 0)
				log.Println("Ошибка:", err)
				return nil
			}
			botpkg.SendText(bot, ctx, message, "Цель удалена ❎", "", 0)
			UserMap[message.Chat.ID].Command = ""
			botpkg.SendKeyboard(bot, ctx, message, nil, "накопления")

		} else {

			goals, err := getUserGoals(db, message.Chat.ID)
			if err != nil {
				log.Println("Ошибка:", err)
			} else {
				targetName := make([]string, 0)
				targetSend := "<b>👇 Пришли мне цель, которую хочешь удалить❌</b>\nТвои цели:\n"
				for _, g := range goals {
					targetSend = targetSend + fmt.Sprintf("▫️%s: %.2f / %.2f₽\n", g.Name, g.Current, g.Target)
					targetName = append(targetName, g.Name)
					UserMap[message.Chat.ID].TargetMap[g.Name] = g.ID
				}

				if len(targetName) == 0 {
					botpkg.SendText(bot, ctx, message, `Ты не добавил цели. Сначала добавь, а потом удаляй 😉`, "", 0)
					return nil
				}
				UserMap[message.Chat.ID].Command = "Ждем ID и удаление"
				botpkg.SendText(bot, ctx, message, targetSend, "", 0)

				botpkg.SendKeyboard(bot, ctx, message, targetName, "цели")
			}
		}
	}
	return nil
}
