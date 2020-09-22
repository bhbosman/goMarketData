package register

import (
	"context"
	"encoding/json"
	"go.uber.org/fx"
	"golang.org/x/oauth2"
	"net/http"
	"os"
	"path/filepath"
)

//func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
//	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
//	fmt.Printf("Go to the following link in your browser then type the "+
//		"authorization code: \n%v\n", authURL)
//
//	var authCode string
//	if _, err := fmt.Scan(&authCode); err != nil {
//		log.Fatalf("Unable to read authorization code: %v", err)
//	}
//
//	tok, err := config.Exchange(context.TODO(), authCode)
//	if err != nil {
//		log.Fatalf("Unable to retrieve token from web: %v", err)
//	}
//	return tok
//}

// Retrieves a token from a local file.

type UserInformation struct {
	Oauth_token *oauth2.Token `json:"oauth_token"`
}

func UserInformationFrom(file string) (*UserInformation, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &UserInformation{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
//func saveToken(path string, token *oauth2.Token) {
//	fmt.Printf("Saving credential file to: %s\n", path)
//	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
//	if err != nil {
//		log.Fatalf("Unable to cache oauth token: %v", err)
//	}
//	defer f.Close()
//	json.NewEncoder(f).Encode(token)
//}

type NamedHttpClient struct {
	Client *http.Client
	Name   string
}

func ProvideMarketPricesClients() fx.Option {
	getClient := func(config *oauth2.Config, userInformation *UserInformation) (*http.Client, error) {
		return config.Client(context.Background(), userInformation.Oauth_token), nil
	}
	var options []fx.Option
	homePath, _ := os.UserHomeDir()
	path := filepath.Join(homePath, ".MarketPrices", "clients")
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		fn := path
		option := fx.Provide(
			fx.Annotated{
				Group: "MarketPricesClients",
				Target: func(params struct {
					fx.In
					MarketPricesConfig *oauth2.Config `name:"MarketPricesApplication"`
				}) (*NamedHttpClient, error) {
					userInformation, err := UserInformationFrom(fn)
					client, err := getClient(params.MarketPricesConfig, userInformation)
					if err != nil {
						return nil, err
					}
					return &NamedHttpClient{
						Client: client,
						Name:   fn,
					}, nil
				}})
		options = append(options, option)
		return nil
	})
	if err != nil {
		return fx.Error(err)
	}
	return fx.Options(options...)
}
