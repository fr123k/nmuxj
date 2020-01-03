package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

type VM struct {
	ID          string `json:"id"`
	Group       string `json:"group"`
	ProvisionID string `json:"provisionid"`
}

type Group struct {
	Group       string `json:"group"`
	Total       int64  `json:"total"`
	Count       int64  `json:"count"`
	ProvisionID string `json:"provisionid"`
	Vms         []VM   `json:"vm"`
}

func (group Group) trigger(r *http.Request) Group {
	if group.Count >= group.Total {
		delete(provisions[group.ProvisionID], group.Group)
		var jHook JenkinsHook
		json.NewDecoder(r.Body).Decode(&jHook)
		startJenkinsJob(cfg, jHook)
	}
	return group
}

type JenkinsHook struct {
	Context      string            `json:"ctx"`
	Job          string            `json:"job"`
	JobParameter map[string]string `json:"jobParams"`
}

type Configuration struct {
	JenkinsServer    string
	JenkinsUser      string
	JenkinsToken     string
	ApplicationToken string
}

var cfg Configuration

var provisions map[string]map[string]Group

func inc(i *int64) int64 { *i++; return *i }

func prettyPrint(i interface{}) string {
	s, _ := json.MarshalIndent(i, "", "\t")
	return string(s)
}

func createGroup(vmID string, group string, provisionID string, total int64) Group {
	vmGroup := Group{
		Group:       group,
		Total:       total,
		Count:       1,
		ProvisionID: provisionID,
		Vms: []VM{
			VM{
				ID:          vmID,
				Group:       group,
				ProvisionID: provisionID,
			},
		},
	}
	return vmGroup
}

func registerVM(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	provisionID := params["provisionID"]
	group := params["group"]
	total, err := strconv.ParseInt(params["total"], 10, 64)
	if err != nil {
		panic(err)
		return
	}

	fmt.Printf("Received call %s\n", r.URL.Path)
	vmID := params["vmID"]

	if groups, found := provisions[provisionID]; found {
		if vmGroup, found := groups[group]; found {
			found := false
			for _, vm := range vmGroup.Vms {
				if strings.Compare(vm.ID, vmID) == 0 {
					found = true
					break
				}
			}
			if found == false {
				vm := VM{
					ID:          vmID,
					Group:       group,
					ProvisionID: provisionID,
				}
				vmGroup.Vms = append(vmGroup.Vms, vm)
				inc(&vmGroup.Count)
				groups[group] = vmGroup
			}
			json.NewEncoder(w).Encode(vmGroup)
			vmGroup.trigger(r)
		} else {
			groups[group] = createGroup(vmID, group, provisionID, total).trigger(r)
		}
	} else {
		groups := make(map[string]Group)
		vmGroup := createGroup(vmID, group, provisionID, total).trigger(r)
		groups[group] = vmGroup
		provisions[provisionID] = groups
		json.NewEncoder(w).Encode(vmGroup)
		return
	}
	return
}

func startJenkinsJob(cfg Configuration, jHook JenkinsHook) {
	var buffer bytes.Buffer
	for k, v := range jHook.JobParameter {
		buffer.WriteString(fmt.Sprintf("%s=%s&", k, v))
	}

	fmt.Printf("Body %s\n", prettyPrint(jHook))
	var buildPath = "buildWithParameters"

	if jHook.JobParameter == nil {
		buildPath = "build"
	}
	//TODO make the jenkin job url als configurable
	url := fmt.Sprintf("%s/job/%s/%s?%scause=%s", cfg.JenkinsServer, jHook.Job, buildPath, buffer.String(), "cloud-init")
	fmt.Printf("Call jenkins %s\n", url)
	client := &http.Client{}
	req, _ := http.NewRequest("POST", url, nil)
	req.SetBasicAuth(cfg.JenkinsUser, cfg.JenkinsToken)
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error calling Jenkins '%s'\n", err)
		return
	}
	bodyText, err := ioutil.ReadAll(resp.Body)
	log.Printf("response %s", string(bodyText))
}

//Initialize the Configuration struct by reading the values for it from the environment variables.
func Init() Configuration {
	return Configuration{
		JenkinsServer:    os.Getenv("JENKINS_SERVER"),
		JenkinsUser:      os.Getenv("JENKINS_USER"),
		JenkinsToken:     os.Getenv("JENKINS_TOKEN"),
		ApplicationToken: os.Getenv("APP_TOKEN"),
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return defaultValue
	}
	return value
}

// Middleware function, which will be called for each request
func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")

		splitToken := strings.Split(auth, " ")
		if len(splitToken) != 2 {
			log.Printf("Authorization header malformed '%s'\n", auth)
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		authType := strings.TrimSpace(splitToken[0])
		if authType != "Bearer" {
			log.Printf("Only Authorization header type 'Bearer' is supported not '%s'\n", authType)
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		b64Token := strings.TrimSpace(splitToken[1])
		token, err := base64.StdEncoding.DecodeString(b64Token)
		if err != nil {
			log.Printf("Authorization header token base64 malformed '%s'\n", b64Token)
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		if string(token) == cfg.ApplicationToken {
			next.ServeHTTP(w, r)
		} else {
			log.Printf("Authentication failed secret mismatch \n")
			http.Error(w, "Forbidden", http.StatusForbidden)
		}
		return
	})
}

func main() {
	provisions = make(map[string]map[string]Group)
	router := mux.NewRouter()

	//TODO define timeout for consider when counter less then total the missing vms as dead and trigger an alert
	cfg = Init()

	router.HandleFunc("/provision/{provisionID}/group/{group}/total/{total:[0-9]+}/vm/{vmID}", registerVM).
		Methods("POST")

	if len(cfg.ApplicationToken) > 0 {
		router.Use(Middleware)
	}

	log.Fatal(http.ListenAndServe(":"+getEnv("PORT", "8080"), router))
}
