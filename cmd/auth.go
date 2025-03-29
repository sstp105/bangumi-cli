package cmd

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"runtime"

	"github.com/spf13/cobra"
	"github.com/sstp105/bangumi-cli/internal/bangumi"
	"github.com/sstp105/bangumi-cli/internal/config"
	"github.com/sstp105/bangumi-cli/internal/ospath"
)

var authCmd = &cobra.Command{
	Use:   "auth",
	Run: func(cmd *cobra.Command, args []string) {
		authHandler()
	},
}

func init() {
	rootCmd.AddCommand(authCmd)
}

var stop = make(chan os.Signal, 1)

func authHandler() {
	clientID := config.AppConfig.BangumiClientID
	clientSecret := config.AppConfig.BangumiClientSecret
	port := config.AppConfig.Port
	redirectURI := "http://localhost:8765" 
	authURL := fmt.Sprintf("https://bgm.tv/oauth/authorize?client_id=%s&response_type=code&redirect_uri=%s", clientID, redirectURI)

	fmt.Println(clientID, clientSecret, port)

	fmt.Println("Opening URL in the browser for bangumi.tv authentication...")
	if err := openBrowser(authURL); err != nil {
		log.Fatalf("failed to open URL in browser:%s\n", err)
	}

	signal.Notify(stop, os.Interrupt)
	go func() {
		serve()
	}()
	<-stop 

	fmt.Println("Callback processed, shutting down the server...")
}

func serve() {
	http.HandleFunc("/", listenCallback)

	port := config.AppConfig.Port
	addr := ":" + port
	fmt.Printf("Listening for callback at localhost%s\n", addr)

	server := &http.Server{
		Addr: addr,
	}
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Error starting the HTTP server: %s", err)
	}
}

func listenCallback(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	code := q.Get("code")
	if code == "" {
		http.Error(w, "Missing code parameter in callback", http.StatusBadRequest)
		return
	}

	fmt.Printf("Authorization code received: %s\n", code)

	client := bangumi.NewOAuthClient()
	redirectURI := fmt.Sprintf("http://localhost:%s", config.AppConfig.Port)
	creds, err := client.GetAccessToken(
		config.AppConfig.BangumiClientID,
		config.AppConfig.BangumiClientSecret,
		code,
		redirectURI,
	)
	if err != nil {
		fmt.Printf("error:%s", err)
		http.Error(w, "Failed to request access token", http.StatusInternalServerError)
	}

	// save the creds locally
	if err := ospath.SaveJSONConfigFile(ospath.BangumiCredentialFile, creds); err != nil {
		fmt.Printf("error saving cred file:%s", err)
	}

	w.Write([]byte("You have successfully authenticated. You can close this window now."))

	stop <- os.Interrupt
}

var osURLCommands = map[string][]string{
	"linux":   {"xdg-open"},
	"windows": {"rundll32", "url.dll,FileProtocolHandler"},
	"darwin":   {"open"},
}

func openBrowser(url string) error {
	goos := runtime.GOOS
	v, exist := osURLCommands[goos]
	if !exist {
		return fmt.Errorf("unsupported os platform %s", goos)
	}

	command := v[0]
	args := append(v[1:], url)

	return exec.Command(command, args...).Start()
}