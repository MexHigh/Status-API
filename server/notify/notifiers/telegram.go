package notifiers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"status-api/notify"
)

const (
	tgApiSendMessageTmpl = "https://api.telegram.org/bot%s/sendMessage"
)

type telegramConfig struct {
	BotToken  string `json:"bot_token"`
	ChatID    string `json:"chat_id"`
	MuteTimes []struct {
		From string `json:"from"`
		To   string `json:"to"`
	} `json:"mute_times"`
}

type Telegram struct {
	config telegramConfig
}

func (t *Telegram) send(title, msg string) error {
	// don't send anything if within mute time
	shouldMute, err := t.shouldMute()
	if err != nil {
		return err
	}
	if shouldMute {
		return nil // "fail" silently
	}

	// prepare payload
	payload, err := json.Marshal(map[string]string{
		"chat_id":    t.config.ChatID,
		"text":       fmt.Sprintf("*%s*\n\n%s", title, msg),
		"parse_mode": "markdown",
	})
	if err != nil {
		return err
	}

	// do request
	response, err := http.Post(
		fmt.Sprintf(tgApiSendMessageTmpl, t.config.BotToken),
		"application/json",
		bytes.NewBuffer(payload),
	)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	// basic error checking
	if response.StatusCode != 200 {
		return errors.New("received non-200 status code: " + strconv.Itoa(response.StatusCode))
	}

	return nil
}

func (t *Telegram) shouldMute() (bool, error) {
	now := time.Now()
	for _, timeObj := range t.config.MuteTimes {
		from, err := time.ParseInLocation("15:04", timeObj.From, time.Local)
		if err != nil {
			return false, err
		}
		from = from.AddDate(now.Year(), int(now.Month())-1, now.Day()-1)

		to, err := time.ParseInLocation("15:04", timeObj.To, time.Local)
		if err != nil {
			return false, err
		}
		to = to.AddDate(now.Year(), int(now.Month())-1, now.Day()-1)

		if now.After(from) && now.Before(to) {
			return true, nil
		}
	}
	return false, nil
}

func (t *Telegram) NotifyDown(serviceName string, reportedDownAt time.Time, reason string) error {
	title := fmt.Sprintf(downNotificationTitle, serviceName)
	msg := fmt.Sprintf(downNotificationMsg, reportedDownAt.Local().Format(dateTimeFormat), reason)
	if err := t.send(title, msg); err != nil {
		return err
	}
	return nil
}

func (t *Telegram) NotifyUp(serviceName string, reportedDownAt time.Time, wasDownFor time.Duration) error {
	title := fmt.Sprintf(upNotificationTitle, serviceName)
	msg := fmt.Sprintf(upNotificationMsg, reportedDownAt.Local().Format(dateTimeFormat), wasDownFor.String())
	if err := t.send(title, msg); err != nil {
		return err
	}
	return nil
}

func (t *Telegram) UnmarshalConfig(raw json.RawMessage) error {
	var c telegramConfig
	if err := json.Unmarshal(raw, &c); err != nil {
		return err
	}
	if c.ChatID == "" {
		return errors.New("telegram.chat_id is empty but required")
	}
	if c.BotToken == "" {
		return errors.New("telegram.bot_token is empty but required")
	}
	t.config = c
	return nil
}

// Interface guard (CAUTION: This interface guard does
// not detect, if a required function is implemented for
// the reciever type, which is invalid!)
var (
	_ notify.Notifier             = (*Telegram)(nil)
	_ notify.ConfigurableNotifier = (*Telegram)(nil)
)

// Register notifier
func init() {
	notify.Register("telegram", &Telegram{})
}
