package services

import (
	"context"
	"encoding/json"
	amqpclient "github.com/KBcHMFollower/blog_posts_service/internal/clients/amqp"
	"github.com/KBcHMFollower/blog_posts_service/internal/clients/amqp/messages"
	ctxerrors "github.com/KBcHMFollower/blog_posts_service/internal/domain/errors"
	repositoriestransfer "github.com/KBcHMFollower/blog_posts_service/internal/domain/layers_TOs/repositories"
	servicestransfer "github.com/KBcHMFollower/blog_posts_service/internal/domain/layers_TOs/services"
	"github.com/KBcHMFollower/blog_posts_service/internal/logger"
	servicesdep "github.com/KBcHMFollower/blog_posts_service/internal/services/interfaces/dep"
	"github.com/google/uuid"
)

type PostsStore interface {
	servicesdep.PostCreator
	servicesdep.PostDeleter
	servicesdep.PostGetter
	servicesdep.PostUpdater
}

type RequestStore interface {
	servicesdep.RequestsCreator
	servicesdep.RequestsGetter
}

type EventStore interface {
	servicesdep.EventCreator
}

type PostService struct {
	postRepository     PostsStore
	requestsRepository RequestStore
	eventsRepository   EventStore
	txCreator          servicesdep.TransactionCreator
	log                logger.Logger
}

func NewPostsService(
	postRepository PostsStore,
	requestsRepository RequestStore,
	eventsRepository EventStore,
	txCreator servicesdep.TransactionCreator,
	log logger.Logger,
) *PostService {
	return &PostService{
		postRepository:     postRepository,
		log:                log,
		eventsRepository:   eventsRepository,
		requestsRepository: requestsRepository,
		txCreator:          txCreator,
	}
}

func (g *PostService) GetUserPosts(ctx context.Context, getInfo *servicestransfer.GetUserPostsInfo) (*servicestransfer.GetUserPostsResult, error) {
	posts, err := g.postRepository.Posts(ctx, repositoriestransfer.GetPostsInfo{
		Condition: map[repositoriestransfer.PostConditionTarget]any{
			repositoriestransfer.PostUserIdCondition: getInfo.UserId,
		},
		Page: getInfo.Page,
		Size: getInfo.Size,
	})
	if err != nil {
		return nil, ctxerrors.WrapCtx(ctx, ctxerrors.Wrap("cant get posts from repository", err))
	}
	total, err := g.postRepository.Count(ctx, repositoriestransfer.GetPostsCountInfo{
		Condition: map[repositoriestransfer.PostConditionTarget]any{
			repositoriestransfer.PostUserIdCondition: getInfo.UserId,
		},
	})
	if err != nil {
		return nil, ctxerrors.WrapCtx(ctx, ctxerrors.Wrap("cant get posts count from repository", err))
	}

	return &servicestransfer.GetUserPostsResult{
		Posts:      servicestransfer.ConvertPostsArrayFromModel(posts),
		TotalCount: total,
	}, nil
}

func (g *PostService) GetPost(ctx context.Context, getInfo *servicestransfer.GetPostInfo) (*servicestransfer.GetPostResult, error) {
	post, err := g.postRepository.Post(ctx, repositoriestransfer.GetPostInfo{
		Condition: map[repositoriestransfer.PostConditionTarget]any{
			repositoriestransfer.PostIdCondition: getInfo.PostId,
		},
	})
	if err != nil {
		return nil, ctxerrors.WrapCtx(ctx, ctxerrors.Wrap("cant get post from repository", err))
	}

	return &servicestransfer.GetPostResult{
		Post: servicestransfer.ConvertPostFromModel(post),
	}, nil
}

func (g *PostService) DeletePost(ctx context.Context, deleteInfo *servicestransfer.DeletePostInfo) error {
	logger.UpdateLoggerCtx(ctx, "post-id", deleteInfo.PostId)
	g.log.InfoContext(ctx, "try to delete post", deleteInfo.PostId)

	if err := g.postRepository.Delete(ctx, repositoriestransfer.DeletePostsInfo{
		Condition: map[repositoriestransfer.PostConditionTarget]any{
			repositoriestransfer.PostIdCondition: deleteInfo.PostId,
		},
	}, nil); err != nil {
		return ctxerrors.WrapCtx(ctx, ctxerrors.Wrap("cant delete posts from repository", err))
	}

	g.log.InfoContext(ctx, "post is deleted", deleteInfo.PostId)

	return nil
}

func (g *PostService) CreatePost(ctx context.Context, createInfo *servicestransfer.CreatePostInfo) (*servicestransfer.CreatePostResult, error) {
	logger.UpdateLoggerCtx(ctx, "create-info", createInfo)
	g.log.InfoContext(ctx, "trying create post")

	postId, err := g.postRepository.Create(ctx, repositoriestransfer.CreatePostInfo{
		UserId:        createInfo.UserId,
		Title:         createInfo.Title,
		TextContent:   createInfo.TextContent,
		ImagesContent: createInfo.ImagesContent,
	})
	if err != nil {
		return nil, ctxerrors.WrapCtx(ctx, ctxerrors.Wrap("cant create post in repository", err))
	}
	post, err := g.postRepository.Post(ctx, repositoriestransfer.GetPostInfo{
		Condition: map[repositoriestransfer.PostConditionTarget]any{
			repositoriestransfer.PostIdCondition: postId,
		},
	})
	if err != nil {
		return nil, ctxerrors.WrapCtx(ctx, ctxerrors.Wrap("cant get post from repository", err))
	}

	g.log.InfoContext(ctx, "post is created", postId)

	return &servicestransfer.CreatePostResult{
		PostId: postId,
		Post:   servicestransfer.ConvertPostFromModel(post),
	}, nil
}

func (g *PostService) UpdatePost(ctx context.Context, updateInfo *servicestransfer.UpdatePostInfo) (*servicestransfer.UpdatePostResult, error) {
	logger.UpdateLoggerCtx(ctx, "post-id", updateInfo.PostId)
	logger.UpdateLoggerCtx(ctx, "update-fields", updateInfo.Fields)

	g.log.InfoContext(ctx, "trying to update post")

	if err := g.postRepository.Update(ctx, repositoriestransfer.UpdatePostInfo{
		Condition: map[repositoriestransfer.PostConditionTarget]any{
			repositoriestransfer.PostIdCondition: updateInfo.PostId,
		},
		UpdateData: updateInfo.Fields,
	}); err != nil {
		return nil, ctxerrors.WrapCtx(ctx, ctxerrors.Wrap("cant update posts from repository", err))
	}

	post, err := g.postRepository.Post(ctx, repositoriestransfer.GetPostInfo{
		Condition: map[repositoriestransfer.PostConditionTarget]any{
			repositoriestransfer.PostIdCondition: updateInfo.PostId,
		},
	})
	if err != nil {
		return nil, ctxerrors.WrapCtx(ctx, ctxerrors.Wrap("cant get post from repository", err))
	}

	g.log.InfoContext(ctx, "post is updated", updateInfo.PostId)

	return &servicestransfer.UpdatePostResult{
		PostId: post.Id,
		Post:   servicestransfer.ConvertPostFromModel(post),
	}, nil
}

// todo: имеет смысл вынести в отдельный сервис
func (g *PostService) DeleteUserPosts(ctx context.Context, deleteInfo servicestransfer.DeleteUserPostInfo) (resErr error) { //TODO: СДЕЛАТЬ DEFER
	logger.UpdateLoggerCtx(ctx, "user-id", deleteInfo.UserId)

	g.log.InfoContext(ctx, "trying to delete user-posts")

	//todo: нужна какая-то amqp-мидлвара для этого
	if err := g.requestsRepository.Create(ctx, repositoriestransfer.CreateRequestInfo{
		Key: deleteInfo.EventId,
	}, nil); err != nil {
		return ctxerrors.WrapCtx(ctx, ctxerrors.Wrap("cant create request from repository", err))
	}

	tx, err := g.txCreator.BeginTxCtx(ctx, nil)
	if err != nil {
		return ctxerrors.WrapCtx(ctx, ctxerrors.Wrap("cant begin transaction", err))
	}
	defer func() {
		if resErr != nil {
			if err := tx.Rollback(); err != nil {
				resErr = ctxerrors.WrapCtx(ctx, ctxerrors.Wrap("cant rollback transaction", err))
			}

			payload, err := json.Marshal(messages.PostsDeleted{
				Status:  messages.Failed,
				EventId: deleteInfo.EventId,
			})
			if err != nil {
				resErr = ctxerrors.WrapCtx(ctx, ctxerrors.Wrap("cant marshal post event", err))
				return
			}

			if err := g.eventsRepository.Create(ctx, repositoriestransfer.CreateEventInfo{
				EventId:   uuid.New(),
				EventType: amqpclient.PostsDeletedEventKey,
				Payload:   payload,
			}, nil); err != nil {
				resErr = ctxerrors.WrapCtx(ctx, ctxerrors.Wrap("cant create post event", err))
			}
		}
	}()

	if err := g.postRepository.Delete(ctx, repositoriestransfer.DeletePostsInfo{
		Condition: map[repositoriestransfer.PostConditionTarget]any{
			repositoriestransfer.PostUserIdCondition: deleteInfo.UserId,
		},
	}, tx); err != nil {
		return ctxerrors.WrapCtx(ctx, ctxerrors.Wrap("cant delete user from repository", err))
	}

	payload, err := json.Marshal(messages.PostsDeleted{
		Status:  messages.Success,
		EventId: deleteInfo.EventId,
	})
	if err != nil {
		return ctxerrors.WrapCtx(ctx, ctxerrors.Wrap("cant marshal post event", err))
	}

	if err := g.eventsRepository.Create(ctx, repositoriestransfer.CreateEventInfo{
		EventId:   uuid.New(),
		EventType: amqpclient.PostsDeletedEventKey,
		Payload:   payload,
	}, tx); err != nil {
		return ctxerrors.WrapCtx(ctx, ctxerrors.Wrap("cant create post event", err))
	}

	if err := tx.Commit(); err != nil {
		return ctxerrors.WrapCtx(ctx, ctxerrors.Wrap("cant commit transaction", err))
	}

	g.log.InfoContext(ctx, "user-posts is deleted")

	return nil
}
