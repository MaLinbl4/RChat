package main

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/glebarez/go-sqlite" // Импорт драйвера SQLite
)

// App struct — это наш главный объект бэкенда
type App struct {
	ctx context.Context
	db  *sql.DB // Сюда мы сохраним подключение к базе
}

// NewApp создает экземпляр приложения
func NewApp() *App {
	return &App{}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	// Используем полный путь для проверки, если относительный не работает
	db, err := sql.Open("sqlite", "chat_history.db")
	if err != nil {
		fmt.Println("КРИТИЧЕСКАЯ ОШИБКА БАЗЫ:", err)
		return
	}
	a.db = db

	// Пытаемся реально что-то записать, чтобы файл создался физически
	_, err = a.db.Exec(`
		CREATE TABLE IF NOT EXISTS messages (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			sender TEXT,
			encrypted_text TEXT,
			timestamp DATETIME DEFAULT CURRENT_TIMESTAMP
		);
	`)

	if err != nil {
		fmt.Println("ОШИБКА СОЗДАНИЯ ТАБЛИЦЫ:", err)
	} else {
		fmt.Println("БАЗА ДАННЫХ ПОДКЛЮЧЕНА УСПЕШНО")
	}
}

func (a *App) SaveMessage(sender string, encryptedText string) string {
	// 1. Сохраняем само сообщение
	_, err := a.db.Exec("INSERT INTO messages (sender, encrypted_text) VALUES (?, ?)", sender, encryptedText)
	if err != nil {
		return "Ошибка сохранения сообщения"
	}

	// 2. БЭКЕНД-ЛОГИКА: Обновляем или добавляем контакт
	// Используем INSERT OR REPLACE, чтобы обновить последнее сообщение
	query := `
	INSERT INTO contacts (name, last_message) 
	VALUES (?, ?) 
	ON CONFLICT(name) DO UPDATE SET last_message = EXCLUDED.last_message;`

	_, err = a.db.Exec(query, sender, "Зашифровано...")
	if err != nil {
		fmt.Println("Ошибка обновления контактов:", err)
	}

	return "OK"
}

// GetContacts возвращает список всех собеседников
func (a *App) GetContacts() []map[string]string {
	rows, err := a.db.Query("SELECT name, last_message FROM contacts ORDER BY id DESC")
	if err != nil {
		return nil
	}
	defer rows.Close()

	var contacts []map[string]string
	for rows.Next() {
		var name, lastMsg string
		rows.Scan(&name, &lastMsg)
		contacts = append(contacts, map[string]string{
			"name":         name,
			"last_message": lastMsg,
		})
	}
	return contacts
}

// GetMessagesByContact возвращает историю переписки именно с этим человеком
func (a *App) GetMessagesByContact(contactName string) []map[string]string {
	rows, err := a.db.Query("SELECT sender, encrypted_text, type FROM messages WHERE contact_name = ? ORDER BY id ASC", contactName)
	if err != nil {
		return nil
	}
	defer rows.Close()

	var messages []map[string]string
	for rows.Next() {
		var s, t, mt string
		rows.Scan(&s, &t, &mt)
		messages = append(messages, map[string]string{
			"sender": s,
			"text":   t,
			"type":   mt,
		})
	}
	return messages
}
