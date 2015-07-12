# beeper-api
REST API wrapper for Beeper messaging server
## PACKAGE DOCUMENTATION
```
import "github.com/arvinkulagin/beeperapi/"

type Client struct {
    // contains filtered or unexported fields
}

func NewClient(host string) (Client, error)

func (c Client) Add(topic string) error

func (c Client) Del(topic string) error

func (c Client) List() ([]string, error)

func (c Client) Ping() error

func (c Client) Pub(topic string, data string) error
```