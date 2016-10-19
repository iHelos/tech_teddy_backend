package story

type StoryStorageEngine interface {
	Create(Story) (error)
	Load(string) (Story, error)
	CheckLogin(string, string) (error)
	CheckIsLogged(string) (error)
}