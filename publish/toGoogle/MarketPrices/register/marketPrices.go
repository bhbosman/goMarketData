package register

import (
	"go.uber.org/fx"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

func ProvideMarketPriceGoogleApplication() fx.Option {
	return fx.Options(
		fx.Provide(
			fx.Annotated{
				Name: "MarketPricesApplication",
				Target: func() (*oauth2.Config, error) {
					homePath, _ := os.UserHomeDir()
					b, err := ioutil.ReadFile(filepath.Join(homePath, ".MarketPrices", "credentials.json"))
					if err != nil {
						log.Fatalf("Unable to read client secret file: %v", err)
					}
					// If modifying these scopes, delete your previously saved token.json.
					return google.ConfigFromJSON(
						b,
						"https://www.googleapis.com/auth/spreadsheets.readonly",
						"https://www.googleapis.com/auth/drive.file")
				},
			}),
	)
}
