package messenger

import (
	"WSChats/pkg/logger"
	"context"
	"database/sql"
	"errors"
)

type repository struct {
	DB     *sql.DB
	logger *logger.Logger
}

func NewRepository(db *sql.DB, l *logger.Logger) Repository {
	return &repository{
		DB:     db,
		logger: l,
	}
}

type Repository interface {
	NewMessage(ctx context.Context, message *Message) (*Message, error)
	DeleteMessage(ctx context.Context, messageID int) error
	EditTextMessage(ctx context.Context, messageID int, text string) error
	GetChatMessages(ctx context.Context, chatID, limit int, offset int) ([]*Message, error)
	CheckAuthorMessage(ctx context.Context, messageID int, userID string) error

	GetUsernameByUUID(ctx context.Context, uuid string) (string, error)
	NewChat(ctx context.Context, chat *NewChatReq) (*NewChatRes, error)
	NewMembers(ctx context.Context, chat int, users []string) error
	GetChatMembers(ctx context.Context, chatID int) ([]string, error)
	SearchDirectChat(ctx context.Context, member1, member2 string) (int, error)
	CheckChatMember(ctx context.Context, chatID int, userID string) error

	StoreReadStatus(ctx context.Context, msg *ReadMessage) error
	CheckReadStatus(ctx context.Context, msg *ReadMessage) error
	UpdateReadStatus(ctx context.Context, msg *ReadMessage) error
	GetMessageByID(ctx context.Context, id int) (*Message, error)
}

func (r *repository) NewMessage(ctx context.Context, message *Message) (*Message, error) {
	query := "INSERT INTO messages (chat_id ,sender_id, text, time, reply) VALUES ($1, $2, $3, $4, $5) RETURNING id;" //RETURNING

	row := r.DB.QueryRowContext(ctx, query, message.ChatID, message.Sender, message.Text, message.Time, message.Reply)

	var idMes int

	if err := row.Scan(&idMes); err != nil {
		r.logger.Error(err)
		return nil, err
	}
	message.Id = idMes

	return message, nil
}

func (r *repository) GetUsernameByUUID(ctx context.Context, uuid string) (string, error) {
	query := "SELECT username FROM users WHERE uuid = $1;"

	var username string
	err := r.DB.QueryRowContext(ctx, query, uuid).Scan(&username)
	if err != nil {
		r.logger.Error(err)
		return "", err
	}

	return username, nil
}

func (r *repository) NewChat(ctx context.Context, chat *NewChatReq) (*NewChatRes, error) {
	query := "INSERT INTO chats (creator, title, is_direct) values ($1, $2, $3) RETURNING id;"

	res := &NewChatRes{}

	err := r.DB.QueryRowContext(ctx, query, chat.Creator, chat.Title, chat.IsDirect).Scan(&res.Id)
	if err != nil {
		r.logger.Error(err)
		return nil, err
	}

	res.Creator = chat.Creator
	res.Title = chat.Title

	return res, nil
}

func (r *repository) NewMembers(ctx context.Context, chat int, users []string) error {
	query := "BEGIN;"
	_, err := r.DB.ExecContext(ctx, query)
	if err != nil {
		r.logger.Error(err)
		r.DB.ExecContext(ctx, "ROLLBACK;")
		return err
	}

	for i := 0; i < len(users); i++ {
		q := "INSERT INTO chats_members (chat_id, user_id) VALUES($1, " + "'" + users[i] + "'" + ");"
		_, err := r.DB.ExecContext(ctx, q, chat)
		if err != nil {
			r.DB.ExecContext(ctx, "ROLLBACK;")
			return err
		}
	}

	query = "COMMIT;"

	_, err = r.DB.ExecContext(ctx, query)
	if err != nil {
		r.logger.Error(err)
		r.DB.ExecContext(ctx, "ROLLBACK;")
		return err
	}

	return nil
}

func (r *repository) UserPresenceInChat(ctx context.Context, chatID int, userID string) error {
	query := "SELECT * FROM chats_members WHERE chat_id = $1 AND user_id = $2;"

	var ch int
	var u string

	err := r.DB.QueryRowContext(ctx, query, chatID, userID).Scan(&ch, &u)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			r.logger.Error(errors.New("you can't send message to this chat"))
			return errors.New("you can't send message to this chat")
		}
		r.logger.Error(err.Error())
		return err
	}

	return nil
}

func (r *repository) GetChatMembers(ctx context.Context, chatID int) ([]string, error) {
	query := "SElECT (user_id) FROM chats_members WHERE chat_id = $1;"

	res, err := r.DB.QueryContext(ctx, query, chatID)

	if err != nil {
		r.logger.Error(err)
		return nil, err
	}

	defer res.Close()

	var members []string

	for res.Next() {
		var member string
		if err = res.Scan(&member); err != nil {
			r.logger.Error(err)
			return nil, err
		}
		members = append(members, member)
	}

	return members, nil
}

func (r *repository) SearchDirectChat(ctx context.Context, member1, member2 string) (int, error) {
	query := `SELECT (chat_id) FROM chats_members 
    		  JOIN chats c on chats_members.chat_id = c.id 
        	  WHERE (creator = $1 AND  user_id = $2 ) OR (creator = $2 AND user_id= $1);`

	var chat int

	err := r.DB.QueryRowContext(ctx, query, member1, member2).Scan(&chat)
	if err != nil {
		r.logger.Error(err)
		return 0, err
	}

	return chat, nil
}

func (r *repository) GetChatMessages(ctx context.Context, chatID, limit int, offset int) ([]*Message, error) {
	query := `
		SELECT * FROM messages
        WHERE chat_id = $1 
        ORDER BY time 
        DESC LIMIT $2 OFFSET $3;
`
	res, err := r.DB.QueryContext(ctx, query, chatID, limit, offset)

	if err != nil {
		r.logger.Error(err)
		return nil, err
	}

	defer res.Close()

	var messages []*Message

	var reply sql.NullInt32

	for res.Next() {
		var message Message
		if err = res.Scan(&message.Id, &message.ChatID, &message.Sender, &message.Text,
			&message.Time, &reply); err != nil {
			r.logger.Error(err)
			return nil, err
		}
		message.Reply = int(reply.Int32)
		messages = append(messages, &message)
	}

	return messages, nil
}

func (r *repository) CheckChatMember(ctx context.Context, chatID int, userID string) error {
	query := "SELECT * FROM chats_members WHERE chat_id = $1 AND user_id = $2;"

	var c sql.NullInt32
	var u sql.NullString

	err := r.DB.QueryRowContext(ctx, query, chatID, userID).Scan(&c, &u)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("user not found in this chat")
		}
		return err
	}

	return nil
}

func (r *repository) StoreReadStatus(ctx context.Context, msg *ReadMessage) error {
	query := `INSERT INTO read_status (user_id, chat_id, last_read_msg) VALUES ($1, $2, $3);`

	_, err := r.DB.ExecContext(ctx, query, msg.UserID, msg.ChatID, msg.LastReadMsg)
	if err != nil {
		r.logger.Error(err)
		return err
	}

	return nil
}

func (r *repository) CheckReadStatus(ctx context.Context, msg *ReadMessage) error {
	query := `
		SELECT (user_id) FROM read_status
		WHERE user_id = $1 AND chat_id = $2;
`
	var u string

	err := r.DB.QueryRowContext(ctx, query, msg.UserID, msg.ChatID).Scan(&u)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("read status not found in this chat")
		}
		return err
	}

	return nil
}

func (r *repository) UpdateReadStatus(ctx context.Context, msg *ReadMessage) error {
	query := `UPDATE read_status SET last_read_msg = $1 WHERE user_id = $2 AND chat_id = $3;`

	_, err := r.DB.ExecContext(ctx, query, msg.LastReadMsg, msg.UserID, msg.ChatID)
	if err != nil {
		r.logger.Error(err)
		return err
	}

	return nil
}

func (r *repository) DeleteMessage(ctx context.Context, messageID int) error {
	query := `DELETE FROM messages WHERE id = $1;`

	_, err := r.DB.ExecContext(ctx, query, messageID)
	if err != nil {
		r.logger.Error(err)
		return err
	}

	return nil
}

func (r *repository) EditTextMessage(ctx context.Context, messageID int, text string) error {
	query := `UPDATE messages SET text = $1 WHERE id = $2;`

	_, err := r.DB.ExecContext(ctx, query, text, messageID)
	if err != nil {
		r.logger.Error(err)
		return err
	}

	return nil
}

func (r *repository) CheckAuthorMessage(ctx context.Context, messageID int, userID string) error {
	query := `SELECT (sender_id) FROM messages WHERE id = $1 AND sender_id = $2;`

	_, err := r.DB.ExecContext(ctx, query, messageID, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("user is not author of message")
		}
		return err
	}

	return nil
}

func (r *repository) GetMessageByID(ctx context.Context, messageID int) (*Message, error) {
	query := "SELECT * FROM messages WHERE id = $1;"

	var reply sql.NullInt32
	msg := &Message{}

	err := r.DB.QueryRowContext(ctx, query, messageID).Scan(
		&msg.Id, &msg.ChatID, &msg.Sender, &msg.Text, &msg.Time, &reply)
	if err != nil {
		r.logger.Error(err)
		return nil, err
	}

	msg.Reply = int(reply.Int32)

	return msg, nil
}
