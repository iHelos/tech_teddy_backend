package story

type StoryStorageEngine interface {
	Create(Story) (error)
	Load(string) (Story, error)
	GetAll(category int, order string, page int) ([]Story, error)
	GetMyStories(string) ([]Story, error)
	Search(keyword string) ([]Story, error)
}