package story

type StoryStorageEngine interface {
	Create(Story) (int, error)
	Load(string) (Story, error)
	GetAll(order string, order_type string, page int) ([]Story, error)
	GetAllByCategory(order string, order_type string, page int, category int) ([]Story, error)
	GetMyStories(string) ([]Story, error)
	Search(keyword string) ([]Story, error)
}