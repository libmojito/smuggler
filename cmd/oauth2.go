package cmd

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"runtime"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
)

const callBackHost = ":8080"
const callBackPath = "/callback"

var shutdownChan = make(chan bool)

func oauthConf() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     viper.GetString("smuggler.oauth2.clientID"),
		ClientSecret: viper.GetString("smuggler.oauth2.clientSecret"),
		Scopes:       []string{viper.GetString("smuggler.oauth2.scopes")},
		Endpoint: oauth2.Endpoint{
			AuthURL:  viper.GetString("smuggler.oauth2.endPoint.auth"),
			TokenURL: viper.GetString("smuggler.oauth2.endPoint.token"),
		},
		RedirectURL: callbackURL(),
	}
}

func callbackURL() string {
	return fmt.Sprintf("http://localhost%s%s", callBackHost, callBackPath)
}

func callbackHandler(c *gin.Context) {
	code := c.Query("code")

	ctx := context.Background()
	tok, err := oauthConf().Exchange(ctx, code)
	if err != nil {
		log.Fatal(err)
	}

	c.String(http.StatusOK, tok.AccessToken)

	shutdownChan <- true
}

func openURL(url string) error {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	}

	return err
}

func startServer() *http.Server {

	router := gin.Default()
	router.GET(callBackPath, callbackHandler)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	return srv

}

func NewOauth2Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "oauth2",
		Short: "oauth2 client",
		Long:  `oauth2 client`,
		Run: func(cmd *cobra.Command, args []string) {

			url := oauthConf().AuthCodeURL("state")
			fmt.Printf("Visit the URL for the auth dialog: %v\n", url)
			openURL(url)

			srv := startServer()

			<-shutdownChan
			fmt.Println("shutting down server...")
			srv.Shutdown(context.Background())
		},
	}

	clientIDFlag := "clientID"
	cmd.Flags().String(clientIDFlag, "", "The clientID to use")
	viper.BindPFlag("smuggler.oauth2.clientID", cmd.Flags().Lookup(clientIDFlag))

	clientSecretFlag := "clientSecret"
	cmd.Flags().String(clientSecretFlag, "", "The clientSecret to use")
	viper.BindPFlag("smuggler.oauth2.clientSecret", cmd.Flags().Lookup(clientSecretFlag))

	scopesFlag := "scopes"
	cmd.Flags().String(scopesFlag, "", "The scope")
	viper.BindPFlag("smuggler.oauth2.scopes", cmd.Flags().Lookup(scopesFlag))

	authEPFlag := "auth-url"
	cmd.Flags().String(authEPFlag, "", "The auth end point")
	viper.BindPFlag("smuggler.oauth2.endPoint.auth", cmd.Flags().Lookup(authEPFlag))

	tokenEPFlag := "token-url"
	cmd.Flags().String(tokenEPFlag, "", "The token end point")
	viper.BindPFlag("smuggler.oauth2.endPoint.token", cmd.Flags().Lookup(tokenEPFlag))

	return cmd
}
