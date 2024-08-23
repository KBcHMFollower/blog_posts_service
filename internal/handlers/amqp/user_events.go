package amqp

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/KBcHMFollower/blog_posts_service/internal/clients/amqp/messages"
	services_transfer "github.com/KBcHMFollower/blog_posts_service/internal/domain/layers_TOs/services"
	"github.com/KBcHMFollower/blog_posts_service/internal/services"
)

type UserEventsHandler struct {
	postsService   *services.PostService //TODO: ТУТ ДОЛЖНЫ БЫТЬ ИНТЕРФЕЙСЫ
	requestService *services.RequestsService
}

func (uh *UserEventsHandler) HandleUserDeletedEvent(message []byte) error {

	var userMessage messages.UserDeletedMessage
	if err := json.Unmarshal(message, &userMessage); err != nil {
		return fmt.Errorf("cant`t pars message: %w", err)
	}

	res, err := uh.requestService.CheckExists(context.TODO(), services_transfer.RequestsCheckExistsInfo{
		Key: userMessage.EventId,
	})
	if err != nil {
		return fmt.Errorf("check exists error: %w", err)
	}
	if res {
		return fmt.Errorf("request is exists: %w", err)
	}

	if err := uh.postsService.DeleteUserPosts(context.TODO(), services_transfer.DeleteUserPostInfo{
		UserId: userMessage.UserId,
	}); err != nil {
		return fmt.Errorf("delete user posts error: %w", err)
	}

	return nil
}
