package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

var (
	listenAddress = flag.String(
		"web.listen-address", ":7654",
		"Address to listen on for requests.",
	)
	timeout = flag.Duration(
		"timeout", 5*time.Second,
		"Timeout for trying to query the livestatus socket.",
	)
	socket = flag.String(
		"socket-path", "/var/cache/naemon/live",
		"Path for Livestatus UNIX socket.",
	)
)

type Comment struct {
	ID                 int    `json:"id"`
	Author             string `json:"author"`
	Comment            string `json:"comment"`
	EntryTime          int    `json:"entry_time"`
	EntryType          int    `json:"entry_type"`
	ExpireTime         int    `json:"expire_time"`
	Expires            bool   `json:"expires"`
	Type               int    `json:"type"`
	HostName           string `json:"host_name"`
	ServiceDescription string `json:"service_description"`
}

func (c *Comment) UnmarshalJSON(b []byte) (err error) {
	var tmp []interface{}
	err = json.Unmarshal(b, &tmp)
	if err != nil {
		log.Printf("Error: Unmarshalling JSON for Comment object with error: %s", err)
		return err
	}

	c.ID = int(tmp[0].(float64))
	c.Author = tmp[1].(string)
	c.Comment = tmp[2].(string)
	c.EntryTime = int(tmp[3].(float64))
	c.EntryType = int(tmp[4].(float64))
	c.ExpireTime = int(tmp[5].(float64))
	c.Expires = tmp[6].(float64) != 0
	c.Type = int(tmp[7].(float64))
	c.HostName = tmp[8].(string)
	c.ServiceDescription = tmp[9].(string)

	return nil
}

type Contact struct {
	ID                          int    `json:"id"`
	Name                        string `json:"name"`
	Alias                       string `json:"alias"`
	Email                       string `json:"email"`
	Pager                       string `json:"pager"`
	HostNotificationPeriod      string `json:"host_notification_period"`
	HostNotificationsEnabled    bool   `json:"host_notifications_enabled"`
	ServiceNotificationPeriod   string `json:"service_notification_period"`
	ServiceNotificationsEnabled bool   `json:"service_notifications_enabled"`
}

func (c *Contact) UnmarshalJSON(b []byte) (err error) {
	var tmp []interface{}
	err = json.Unmarshal(b, &tmp)
	if err != nil {
		log.Printf("Error: Unmarshalling JSON for Contact object with error: %s", err)
		return err
	}

	c.ID = int(tmp[0].(float64))
	c.Name = tmp[1].(string)
	c.Alias = tmp[2].(string)
	c.Email = tmp[3].(string)
	c.Pager = tmp[4].(string)
	c.HostNotificationPeriod = tmp[5].(string)
	c.HostNotificationsEnabled = tmp[6].(float64) != 0
	c.ServiceNotificationPeriod = tmp[7].(string)
	c.ServiceNotificationsEnabled = tmp[8].(float64) != 0

	return nil
}

type Downtime struct {
	ID                 int    `json:"id"`
	Author             string `json:"author"`
	Comment            string `json:"comment"`
	Duration           int    `json:"duration"`
	StartTime          int    `json:"start_time"`
	EndTime            int    `json:"end_time"`
	EntryTime          int    `json:"entry_time"`
	Fixed              bool   `json:"fixed"`
	Type               int    `json:"type"`
	HostName           string `json:"host_name"`
	ServiceDescription string `json:"service_description"`
}

func (d *Downtime) UnmarshalJSON(b []byte) (err error) {
	var tmp []interface{}
	err = json.Unmarshal(b, &tmp)
	if err != nil {
		log.Printf("Error: Unmarshalling JSON for Downtime object with error: %s", err)
		return err
	}

	d.ID = int(tmp[0].(float64))
	d.Author = tmp[1].(string)
	d.Comment = tmp[2].(string)
	d.Duration = int(tmp[3].(float64))
	d.StartTime = int(tmp[4].(float64))
	d.EndTime = int(tmp[5].(float64))
	d.EntryTime = int(tmp[6].(float64))
	d.Fixed = tmp[7].(float64) != 0
	d.Type = int(tmp[8].(float64))
	d.HostName = tmp[9].(string)
	d.ServiceDescription = tmp[9].(string)

	return nil
}

type Host struct {
	ID                         int      `json:"id"`
	Name                       string   `json:"name"`
	Alias                      string   `json:"alias"`
	Acknowledged               bool     `json:"acknowledged"`
	Address                    string   `json:"address"`
	CheckPeriod                string   `json:"check_period"`
	CheckSource                string   `json:"check_source"`
	ChecksEnabled              bool     `json:"checks_enabled"`
	Comments                   []int    `json:"comments"`
	Contacts                   []string `json:"contacts"`
	Downtimes                  []int    `json:"downtimes"`
	EventHandler               string   `json:"event_handler"`
	EventHandlerEnabled        bool     `json:"event_handler_enabled"`
	ExecutionTime              int      `json:"execution_time"`
	FlapDetectionEnabled       bool     `json:"flap_detection_enabled"`
	Groups                     []string `json:"groups"`
	HardState                  int      `json:"hard_state"`
	HasBeenChecked             bool     `json:"has_been_checked"`
	InCheckPeriod              bool     `json:"in_check_period"`
	InNotificationPeriod       bool     `json:"in_notification_period"`
	IsFlapping                 bool     `json:"is_flapping"`
	LastCheck                  int      `json:"last_check"`
	LastNotification           int      `json:"last_notification"`
	LastStateChange            int      `json:"last_state_change"`
	LastTimeDown               int      `json:"last_time_down"`
	LastTimeUnreachable        int      `json:"last_time_unreachable"`
	LastTimeUp                 int      `json:"last_time_up"`
	Latency                    int      `json:"latency"`
	NextCheck                  int      `json:"next_check"`
	NextNotification           int      `json:"next_notification"`
	NotificationPeriod         string   `json:"notification_period"`
	NotificationsEnabled       bool     `json:"notifications_enabled"`
	NumberServices             int      `json:"number_of_services"`
	NumberServicesHardCritical int      `json:"number_of_services_hard_critical"`
	NumberServicesHardOK       int      `json:"number_of_services_hard_ok"`
	NumberServicesHardUnknown  int      `json:"number_of_services_hard_unknown"`
	NumberServicesHardWarning  int      `json:"number_of_services_hard_warning"`
	NumberServicesPending      int      `json:"number_of_services_pending"`
	State                      int      `json:"state"`
	StateType                  int      `json:"state_type"`
	Services                   []string `json:"services"`
}

func (h *Host) UnmarshalJSON(b []byte) (err error) {
	var tmp []interface{}
	err = json.Unmarshal(b, &tmp)
	if err != nil {
		log.Printf("Error: Unmarshalling JSON for Host object with error: %s", err)
		return err
	}

	h.ID = int(tmp[0].(float64))
	h.Name = tmp[1].(string)
	h.Alias = tmp[2].(string)
	h.Acknowledged = tmp[3].(float64) != 0
	h.Address = tmp[4].(string)
	h.CheckPeriod = tmp[5].(string)
	h.CheckSource = tmp[6].(string)
	h.ChecksEnabled = tmp[7].(float64) != 0
	h.Comments = make([]int, len(tmp[8].([]interface{})))
	for i := range tmp[8].([]interface{}) {
		h.Comments[i] = int(tmp[8].([]interface{})[i].(float64))
	}
	h.Contacts = make([]string, len(tmp[9].([]interface{})))
	for i := range tmp[9].([]interface{}) {
		h.Contacts[i] = tmp[9].([]interface{})[i].(string)
	}
	h.Downtimes = make([]int, len(tmp[10].([]interface{})))
	for i := range tmp[10].([]interface{}) {
		h.Downtimes[i] = int(tmp[10].([]interface{})[i].(float64))
	}
	h.EventHandler = tmp[11].(string)
	h.EventHandlerEnabled = tmp[12].(float64) != 0
	h.ExecutionTime = int(tmp[13].(float64))
	h.FlapDetectionEnabled = tmp[14].(float64) != 0
	h.Groups = make([]string, len(tmp[15].([]interface{})))
	for i := range tmp[15].([]interface{}) {
		h.Groups[i] = tmp[15].([]interface{})[i].(string)
	}
	h.HardState = int(tmp[16].(float64))
	h.HasBeenChecked = tmp[17].(float64) != 0
	h.InCheckPeriod = tmp[18].(float64) != 0
	h.InNotificationPeriod = tmp[19].(float64) != 0
	h.IsFlapping = tmp[20].(float64) != 0
	h.LastCheck = int(tmp[21].(float64))
	h.LastNotification = int(tmp[22].(float64))
	h.LastStateChange = int(tmp[23].(float64))
	h.LastTimeDown = int(tmp[24].(float64))
	h.LastTimeUnreachable = int(tmp[25].(float64))
	h.LastTimeUp = int(tmp[26].(float64))
	h.Latency = int(tmp[27].(float64))
	h.NextCheck = int(tmp[28].(float64))
	h.NextNotification = int(tmp[29].(float64))
	h.NotificationPeriod = tmp[30].(string)
	h.NotificationsEnabled = tmp[31].(float64) != 0
	h.NumberServices = int(tmp[32].(float64))
	h.NumberServicesHardCritical = int(tmp[33].(float64))
	h.NumberServicesHardOK = int(tmp[34].(float64))
	h.NumberServicesHardUnknown = int(tmp[35].(float64))
	h.NumberServicesHardWarning = int(tmp[36].(float64))
	h.NumberServicesPending = int(tmp[37].(float64))
	h.State = int(tmp[38].(float64))
	h.StateType = int(tmp[39].(float64))
	h.Services = make([]string, len(tmp[40].([]interface{})))
	for i := range tmp[40].([]interface{}) {
		h.Services[i] = tmp[40].([]interface{})[i].(string)
	}

	return nil
}

type Service struct {
	ID                   int      `json:"id"`
	Acknowledged         bool     `json:"acknowledged"`
	CheckPeriod          string   `json:"check_period"`
	CheckSource          string   `json:"check_source"`
	CheckType            int      `json:"check_type"`
	ChecksEnabled        bool     `json:"checks_enabled"`
	Comments             []int    `json:"comments"`
	Contacts             []string `json:"contacts"`
	Description          string   `json:"description"`
	Downtimes            []int    `json:"downtimes"`
	EventHandler         string   `json:"event_handler"`
	EventHandlerEnabled  bool     `json:"event_handler_enabled"`
	ExecutionTime        int      `json:"execution_time"`
	FlapDetectionEnabled bool     `json:"flap_detection_enabled"`
	Groups               []string `json:"groups"`
	HasBeenChecked       bool     `json:"has_been_checked"`
	InCheckPeriod        bool     `json:"in_check_period"`
	InNotificationPeriod bool     `json:"in_notification_period"`
	IsFlapping           bool     `json:"is_flapping"`
	LastCheck            int      `json:"last_check"`
	LastNotification     int      `json:"last_notification"`
	LastStateChange      int      `json:"last_state_change"`
	LastTimeCritical     int      `json:"last_time_critical"`
	LastTimeOK           int      `json:"last_time_ok"`
	LastTimeUnknown      int      `json:"last_time_unknown"`
	LastTimeWarning      int      `json:"last_time_warning"`
	Latency              int      `json:"latency"`
	NextCheck            int      `json:"next_check"`
	NextNotification     int      `json:"next_notification"`
	NotificationPeriod   string   `json:"notification_period"`
	NotificationsEnabled bool     `json:"notifications_enabled"`
	State                int      `json:"state"`
	StateType            int      `json:"state_type"`
	HostName             string   `json:"host"`
}

func (s *Service) UnmarshalJSON(b []byte) (err error) {
	var tmp []interface{}
	err = json.Unmarshal(b, &tmp)
	if err != nil {
		log.Printf("Error: Unmarshalling JSON for Service object with error: %s", err)
		return err
	}

	s.ID = int(tmp[0].(float64))
	s.Acknowledged = tmp[1].(float64) != 0
	s.CheckPeriod = tmp[2].(string)
	s.CheckSource = tmp[3].(string)
	s.CheckType = int(tmp[4].(float64))
	s.ChecksEnabled = tmp[5].(float64) != 0
	s.Comments = make([]int, len(tmp[6].([]interface{})))
	for i := range tmp[6].([]interface{}) {
		s.Comments[i] = int(tmp[6].([]interface{})[i].(float64))
	}
	s.Contacts = make([]string, len(tmp[7].([]interface{})))
	for i := range tmp[7].([]interface{}) {
		s.Contacts[i] = tmp[7].([]interface{})[i].(string)
	}
	s.Description = tmp[8].(string)
	s.Downtimes = make([]int, len(tmp[9].([]interface{})))
	for i := range tmp[9].([]interface{}) {
		s.Downtimes[i] = int(tmp[9].([]interface{})[i].(float64))
	}
	s.EventHandler = tmp[10].(string)
	s.EventHandlerEnabled = tmp[11].(float64) != 0
	s.ExecutionTime = int(tmp[12].(float64))
	s.FlapDetectionEnabled = tmp[13].(float64) != 0
	s.Groups = make([]string, len(tmp[14].([]interface{})))
	for i := range tmp[14].([]interface{}) {
		s.Groups[i] = tmp[14].([]interface{})[i].(string)
	}
	s.HasBeenChecked = tmp[15].(float64) != 0
	s.InCheckPeriod = tmp[16].(float64) != 0
	s.InNotificationPeriod = tmp[17].(float64) != 0
	s.IsFlapping = tmp[18].(float64) != 0
	s.LastCheck = int(tmp[19].(float64))
	s.LastNotification = int(tmp[20].(float64))
	s.LastStateChange = int(tmp[21].(float64))
	s.LastTimeCritical = int(tmp[22].(float64))
	s.LastTimeOK = int(tmp[23].(float64))
	s.LastTimeUnknown = int(tmp[24].(float64))
	s.LastTimeWarning = int(tmp[25].(float64))
	s.Latency = int(tmp[26].(float64))
	s.NextCheck = int(tmp[27].(float64))
	s.NextNotification = int(tmp[28].(float64))
	s.NotificationPeriod = tmp[29].(string)
	s.NotificationsEnabled = tmp[30].(float64) != 0
	s.State = int(tmp[31].(float64))
	s.StateType = int(tmp[32].(float64))
	s.HostName = tmp[33].(string)

	return nil
}

type Status struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func query(q string) (io.ReadCloser, error) {
	f, err := net.DialTimeout("unix", *socket, *timeout)
	if err != nil {
		log.Printf("Error: Could not connect to unix socket: %s", err)
		return nil, err
	}
	if err := f.SetDeadline(time.Now().Add(*timeout)); err != nil {
		f.Close()
		log.Printf("Error: Connection to unix socket timed out: %s", err)
		return nil, err
	}

	_, err = io.WriteString(f, q+"\nOutputFormat: json\n\n")
	if err != nil {
		log.Printf("Error: Could not query Livestatus: %s", err)
		f.Close()
		return nil, err
	}

	return f, nil
}

func returnJson(w http.ResponseWriter, j interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(j)
}

func returnCode404(w http.ResponseWriter, message string) {
	s := Status{Code: 404, Message: message}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(s)
}

func returnCode500(w http.ResponseWriter, message string) {
	s := Status{Code: 500, Message: message}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(s)
}

func getComments(w http.ResponseWriter, r *http.Request) {
	var comments []Comment

	raw, err := query("GET comments\nColumns:id author comment entry_time entry_type expire_time expires type host_name service_description")
	if err != nil {
		returnCode500(w, "Could not query Livestatus")
		return
	}
	defer raw.Close()

	err = json.NewDecoder(raw).Decode(&comments)
	if err != nil {
		log.Printf("Error: Decoding JSON from Livestatus: %s", err)
		returnCode500(w, "Could not decode response from Livestatus")
		return
	}
	returnJson(w, comments)
}

func getComment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var comments []Comment

	raw, err := query(fmt.Sprintf("GET comments\nFilter: id = %s\nColumns:id author comment entry_time entry_type expire_time expires type host_name service_description", vars["id"]))
	if err != nil {
		returnCode500(w, "Could not query Livestatus")
		return
	}
	defer raw.Close()

	err = json.NewDecoder(raw).Decode(&comments)
	if err != nil {
		log.Printf("Error: Decoding JSON from Livestatus: %s", err)
		returnCode500(w, "Could not decode response from Livestatus")
		return
	}
	if len(comments) > 0 {
		returnJson(w, comments[0])
	} else {
		returnCode404(w, "Comment not found")
	}
}

func getContacts(w http.ResponseWriter, r *http.Request) {
	var contacts []Contact

	raw, err := query("GET contacts\nColumns:id name alias email pager host_notification_period host_notifications_enabled service_notification_period service_notifications_enabled")
	if err != nil {
		returnCode500(w, "Could not query Livestatus")
		return
	}
	defer raw.Close()

	err = json.NewDecoder(raw).Decode(&contacts)
	if err != nil {
		log.Printf("Error: Decoding JSON from Livestatus: %s", err)
		returnCode500(w, "Could not decode response from Livestatus")
		return
	}
	returnJson(w, contacts)
}

func getContact(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var contacts []Contact

	raw, err := query(fmt.Sprintf("GET contacts\nFilter: name = %s\nColumns:id name alias email pager host_notification_period host_notifications_enabled service_notification_period service_notifications_enabled", vars["name"]))
	if err != nil {
		returnCode500(w, "Could not query Livestatus")
		return
	}
	defer raw.Close()

	err = json.NewDecoder(raw).Decode(&contacts)
	if err != nil {
		log.Printf("Error: Decoding JSON from Livestatus: %s", err)
		returnCode500(w, "Could not decode response from Livestatus")
		return
	}
	if len(contacts) > 0 {
		returnJson(w, contacts[0])
	} else {
		returnCode404(w, "Contact not found")
	}
}

func getDowntimes(w http.ResponseWriter, r *http.Request) {
	var downtimes []Downtime

	raw, err := query("GET downtimes\nColumns:id author comment duration start_time end_time entry_time fixed type host_name service_description")
	if err != nil {
		log.Printf("Error: %s", err)
		returnCode500(w, "Could not query Livestatus")
		return
	}
	defer raw.Close()

	err = json.NewDecoder(raw).Decode(&downtimes)
	if err != nil {
		log.Printf("Error: Decoding JSON from Livestatus: %s", err)
		returnCode500(w, "Could not decode response from Livestatus")
		return
	}
	returnJson(w, downtimes)
}

func getDowntime(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var downtimes []Downtime

	raw, err := query(fmt.Sprintf("GET downtimes\nFilter: id = %s\nColumns:id author comment duration start_time end_time entry_time fixed type host_name service_description", vars["id"]))
	if err != nil {
		returnCode500(w, "Could not query Livestatus")
		return
	}
	defer raw.Close()

	err = json.NewDecoder(raw).Decode(&downtimes)
	if err != nil {
		log.Printf("Error: Decoding JSON from Livestatus: %s", err)
		returnCode500(w, "Could not decode response from Livestatus")
		return
	}
	if len(downtimes) > 0 {
		returnJson(w, downtimes[0])
	} else {
		returnCode404(w, "Downtime not found")
	}
}

func getHosts(w http.ResponseWriter, r *http.Request) {
	var hosts []Host

	raw, err := query("GET hosts\nColumns:id name alias acknowledged address check_period check_source checks_enabled comments contacts downtimes event_handler event_handler_enabled execution_time flap_detection_enabled groups hard_state has_been_checked in_check_period in_notification_period is_flapping last_check last_notification last_state_change last_time_down last_time_unreachable last_time_up latency next_check next_notification notification_period notifications_enabled num_services num_services_hard_crit num_services_hard_ok num_services_hard_unknown num_services_hard_warn num_services_pending state state_type services")
	if err != nil {
		returnCode500(w, "Could not query Livestatus")
		return
	}
	defer raw.Close()

	err = json.NewDecoder(raw).Decode(&hosts)
	if err != nil {
		log.Printf("Error: Decoding JSON from Livestatus: %s", err)
		returnCode500(w, "Could not decode response from Livestatus")
		return
	}
	returnJson(w, hosts)
}

func getHost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var hosts []Host

	raw, err := query(fmt.Sprintf("GET hosts\nFilter: name = %s\nColumns:id name alias acknowledged address check_period check_source checks_enabled comments contacts downtimes event_handler event_handler_enabled execution_time flap_detection_enabled groups hard_state has_been_checked in_check_period in_notification_period is_flapping last_check last_notification last_state_change last_time_down last_time_unreachable last_time_up latency next_check next_notification notification_period notifications_enabled num_services num_services_hard_crit num_services_hard_ok num_services_hard_unknown num_services_hard_warn num_services_pending state state_type services", vars["name"]))
	if err != nil {
		returnCode500(w, "Could not query Livestatus")
		return
	}
	defer raw.Close()

	err = json.NewDecoder(raw).Decode(&hosts)
	if err != nil {
		log.Printf("Error: Decoding JSON from Livestatus: %s", err)
		returnCode500(w, "Could not decode response from Livestatus")
		return
	}
	if len(hosts) > 0 {
		returnJson(w, hosts[0])
	} else {
		returnCode404(w, "Host not found")
	}
}

func getServices(w http.ResponseWriter, r *http.Request) {
	var services []Service

	raw, err := query("GET services\nColumns:id acknowledged check_period check_source check_type checks_enabled comments contacts description downtimes event_handler event_handler_enabled execution_time flap_detection_enabled groups has_been_checked in_check_period in_notification_period is_flapping last_check last_notification last_state_change last_time_critical last_time_ok last_time_unknown last_time_warning latency next_check next_notification notification_period notifications_enabled state state_type host_name")
	if err != nil {
		returnCode500(w, "Could not query Livestatus")
		return
	}
	defer raw.Close()

	err = json.NewDecoder(raw).Decode(&services)
	if err != nil {
		log.Printf("Error: Decoding JSON from Livestatus: %s", err)
		returnCode500(w, "Could not decode response from Livestatus")
		return
	}
	returnJson(w, services)
}

func getService(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var services []Service

	raw, err := query(fmt.Sprintf("GET services\nFilter: host_name = %s\nFilter: description = %s\nColumns:id acknowledged check_period check_source check_type checks_enabled comments contacts description downtimes event_handler event_handler_enabled execution_time flap_detection_enabled groups has_been_checked in_check_period in_notification_period is_flapping last_check last_notification last_state_change last_time_critical last_time_ok last_time_unknown last_time_warning latency next_check next_notification notification_period notifications_enabled state state_type host_name", vars["host_name"], vars["name"]))
	if err != nil {
		returnCode500(w, "Could not query Livestatus")
		return
	}
	defer raw.Close()

	err = json.NewDecoder(raw).Decode(&services)
	if err != nil {
		log.Printf("Error: Decoding JSON from Livestatus: %s", err)
		returnCode500(w, "Could not decode response from Livestatus")
		return
	}
	if len(services) > 0 {
		returnJson(w, services[0])
	} else {
		returnCode404(w, "Service not found")
	}
}

func main() {
	flag.Parse()
	log.Printf("Starting up...")

	router := mux.NewRouter()
	router.HandleFunc("/comments", getComments)
	router.HandleFunc("/comments/{id:[0-9]+}", getComment)
	router.HandleFunc("/contacts", getContacts)
	router.HandleFunc("/contacts/{name}", getContact)
	router.HandleFunc("/downtimes", getDowntimes)
	router.HandleFunc("/downtimes/{id:[0-9]+}", getDowntime)
	router.HandleFunc("/hosts", getHosts)
	router.HandleFunc("/hosts/{name}", getHost)
	router.HandleFunc("/services", getServices)
	router.HandleFunc("/hosts/{host_name}/services/{name}", getService)
	log.Fatal(http.ListenAndServe(*listenAddress, router))
}
