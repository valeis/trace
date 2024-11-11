package chat

type IChatRepository interface {
	GetContactList(username string) (ChatServiceResponse, error)
	GetChatHistory(participants ChatParticipants) (ChatHistory, error)
}

type ChatService struct {
	repo IChatRepository
}

func NewChatService(repo IChatRepository) *ChatService {
	return &ChatService{repo: repo}
}

func (svc *ChatService) GetContactList(username string) (ChatServiceResponse, error) {
	return svc.repo.GetContactList(username)
}

func (svc *ChatService) GetChatHistory(participants ChatParticipants) (ChatHistory, error) {
	return svc.repo.GetChatHistory(participants)
}
