package main

import (
	"fmt"

	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/html; charset=utf-8")
		fmt.Fprintf(w, "Testing with <a href=\"https://github.com/rakyll/hey\">hey</a>" +
			"<ul>" +
			"<li>Test call from host to <a href='/localhost'>http://localhost:%s</a> (5s, local network)</li>" +
			"<li>Test call from host to <a href='/google-frontend/%s'>https://%s</a> (5s, through google frontend)</li>" +
			"</ul>",
				port(), strings.TrimSuffix(r.Host, ".a.run.app"), r.Host)
	})

	http.HandleFunc("/google-frontend/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/plain; charset=utf-8")

		host := strings.TrimPrefix(r.RequestURI, "/google-frontend/")
		if host == "" {
			host = "https://" + r.Host
		} else {
			host = "https://" + host + ".a.run.app/"
		}

		out, err := hey(host)
		if err != nil {
			http.Error(w,
				fmt.Sprintf("command failed with %s, %s\n", out, err),
				http.StatusInternalServerError )
			return
		}

		fmt.Fprintf(w, out)
	})

	http.HandleFunc("/localhost", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/plain; charset=utf-8")

		out, err := hey(fmt.Sprintf("http://localhost:%s/", port()))
		if err != nil {
			http.Error(w,
				fmt.Sprintf("command failed with %s, %s\n", out, err),
				http.StatusInternalServerError )
			return
		}

		fmt.Fprintf(w, out)
	})

	log.Fatal(http.ListenAndServe(":"+port(), nil))
}

func args(uri string) []string {
	return []string{"-z", "5s", "-c", "2", uri}
}

func hey(uri string) (string, error){
	args :=  args(uri)

	cmd := exec.Command("hey", args...)

	out, err := cmd.CombinedOutput()
	if err != nil {
		return string(out), err
	}
	result := fmt.Sprintf(
		"Running: hey %s\n\n" +
			"%s",
		strings.Join(args, " "),
		string(out))
	return result, nil
}

func port() string {
	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}
	return port
}