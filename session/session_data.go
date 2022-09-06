package session

type SessionData interface {
	GetData() map[string]interface{}
	Set(key string, value interface{})
	Get(key string) (interface{}, bool)
	GetString(key, def string) string
	Delete(key string)
}
