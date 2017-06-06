package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"
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
	ID         int    `json:"id"`
	Author     string `json:"author"`
	Comment    string `json:"comment"`
	EntryTime  int    `json:"entry_time"`
	EntryType  int    `json:"entry_type"`
	ExpireTime int    `json:"expire_time"`
	Expires    bool   `json:"expires"`
	Type       int    `json:"type"`
}

func (c *Comment) UnmarshalJSON(b []byte) (err error) {
	var tmp []interface{}
	err = json.Unmarshal(b, &tmp)
	if err != nil {
		log.Fatal(err)
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
		log.Fatal(err)
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
	ID        int    `json:"id"`
	Author    string `json:"author"`
	Comment   string `json:"comment"`
	Duration  int    `json:"duration"`
	StartTime int    `json:"start_time"`
	EndTime   int    `json:"end_time"`
	EntryTime int    `json:"entry_time"`
	Fixed     bool   `json:"fixed"`
	Type      int    `json:"type"`
}

func (d *Downtime) UnmarshalJSON(b []byte) (err error) {
	var tmp []interface{}
	err = json.Unmarshal(b, &tmp)
	if err != nil {
		log.Fatal(err)
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
}

func (h *Host) UnmarshalJSON(b []byte) (err error) {
	var tmp []interface{}
	err = json.Unmarshal(b, &tmp)
	if err != nil {
		log.Fatal(err)
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
	HostID               int      `json:"host_id"`
}

func (s *Service) UnmarshalJSON(b []byte) (err error) {
	var tmp []interface{}
	err = json.Unmarshal(b, &tmp)
	if err != nil {
		log.Fatal(err)
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
	s.HostID = int(tmp[33].(float64))

	return nil
}

func parseID(s string) (id int, err error) {
	p := strings.LastIndex(s, "/")
	if p < 0 {
		return 0, errors.New("Invalid URL")
	}

	id, err = strconv.Atoi(s[p+1:])
	if err != nil {
		return 0, errors.New("Invalid URL")
	}
	return id, nil
}

func query(q string) (io.ReadCloser, error) {
	f, err := net.DialTimeout("unix", *socket, *timeout)
	if err != nil {
		return nil, err
	}
	if err := f.SetDeadline(time.Now().Add(*timeout)); err != nil {
		f.Close()
		return nil, err
	}

	//cmd := "GET contacts\nColumns:id name alias email pager host_notification_period host_notifications_enabled service_notification_period service_notifications_enabled\nOutputFormat: json\n\n"
	_, err = io.WriteString(f, q+"\nOutputFormat: json\n\n")
	if err != nil {
		f.Close()
		return nil, err
	}

	return f, nil
}

func returnAllComments(w http.ResponseWriter, r *http.Request) {
	var comments []Comment
	raw, err := query("GET comments\nColumns:id author comment entry_time entry_type expire_time expires type")
	if err != nil {
		log.Fatal(err)
	}
	defer raw.Close()

	err = json.NewDecoder(raw).Decode(&comments)
	json.NewEncoder(w).Encode(comments)
}

func returnComment(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r.URL.Path)
	if err != nil {
		log.Fatal(err)
	}

	var comments []Comment
	raw, err := query(fmt.Sprintf("GET comments\nFilter: id = %d\nColumns:id author comment entry_time entry_type expire_time expires type", id))
	if err != nil {
		log.Fatal(err)
	}
	defer raw.Close()

	err = json.NewDecoder(raw).Decode(&comments)
	if len(comments) > 0 {
		json.NewEncoder(w).Encode(comments[0])
	} else {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404 - Comment not found"))
	}
}

func returnAllContacts(w http.ResponseWriter, r *http.Request) {
	var contacts []Contact
	raw, err := query("GET contacts\nColumns:id name alias email pager host_notification_period host_notifications_enabled service_notification_period service_notifications_enabled")
	if err != nil {
		log.Fatal(err)
	}
	defer raw.Close()

	err = json.NewDecoder(raw).Decode(&contacts)
	json.NewEncoder(w).Encode(contacts)
}

func returnContact(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r.URL.Path)
	if err != nil {
		log.Fatal(err)
	}

	var contacts []Contact
	raw, err := query(fmt.Sprintf("GET contacts\nFilter: id = %d\nColumns:id name alias email pager host_notification_period host_notifications_enabled service_notification_period service_notifications_enabled", id))
	if err != nil {
		log.Fatal(err)
	}
	defer raw.Close()

	err = json.NewDecoder(raw).Decode(&contacts)
	if len(contacts) > 0 {
		json.NewEncoder(w).Encode(contacts[0])
	} else {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404 - Contact not found"))
	}
}

func returnAllDowntimes(w http.ResponseWriter, r *http.Request) {
	var downtimes []Downtime
	raw, err := query("GET downtimes\nColumns:id author comment duration start_time end_time entry_time fixed type")
	if err != nil {
		log.Fatal(err)
	}
	defer raw.Close()

	err = json.NewDecoder(raw).Decode(&downtimes)
	json.NewEncoder(w).Encode(downtimes)
}

func returnDowntime(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r.URL.Path)
	if err != nil {
		log.Fatal(err)
	}

	var downtimes []Downtime
	raw, err := query(fmt.Sprintf("GET downtimes\nFilter: id = %d\nColumns:id author comment duration start_time end_time entry_time fixed type", id))
	if err != nil {
		log.Fatal(err)
	}
	defer raw.Close()

	err = json.NewDecoder(raw).Decode(&downtimes)
	if len(downtimes) > 0 {
		json.NewEncoder(w).Encode(downtimes[0])
	} else {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404 - Downtime not found"))
	}
}

func returnAllHosts(w http.ResponseWriter, r *http.Request) {
	var hosts []Host
	raw, err := query("GET hosts\nColumns:id name alias acknowledged address check_period check_source checks_enabled comments contacts downtimes event_handler event_handler_enabled execution_time flap_detection_enabled groups hard_state has_been_checked in_check_period in_notification_period is_flapping last_check last_notification last_state_change last_time_down last_time_unreachable last_time_up latency next_check next_notification notification_period notifications_enabled num_services num_services_hard_crit num_services_hard_ok num_services_hard_unknown num_services_hard_warn num_services_pending state state_type")
	if err != nil {
		log.Fatal(err)
	}
	defer raw.Close()

	err = json.NewDecoder(raw).Decode(&hosts)
	json.NewEncoder(w).Encode(hosts)
}

func returnHost(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r.URL.Path)
	if err != nil {
		log.Fatal(err)
	}

	var hosts []Host
	raw, err := query(fmt.Sprintf("GET hosts\nFilter: id = %d\nColumns:id name alias acknowledged address check_period check_source checks_enabled comments contacts downtimes event_handler event_handler_enabled execution_time flap_detection_enabled groups hard_state has_been_checked in_check_period in_notification_period is_flapping last_check last_notification last_state_change last_time_down last_time_unreachable last_time_up latency next_check next_notification notification_period notifications_enabled num_services num_services_hard_crit num_services_hard_ok num_services_hard_unknown num_services_hard_warn num_services_pending state state_type", id))
	if err != nil {
		log.Fatal(err)
	}
	defer raw.Close()

	err = json.NewDecoder(raw).Decode(&hosts)
	if len(hosts) > 0 {
		json.NewEncoder(w).Encode(hosts[0])
	} else {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404 - Host not found"))
	}
}

func returnAllServices(w http.ResponseWriter, r *http.Request) {
	var services []Service
	raw, err := query("GET services\nColumns:id acknowledged check_period check_source check_type checks_enabled comments contacts description downtimes event_handler event_handler_enabled execution_time flap_detection_enabled groups has_been_checked in_check_period in_notification_period is_flapping last_check last_notification last_state_change last_time_critical last_time_ok last_time_unknown last_time_warning latency next_check next_notification notification_period notifications_enabled state state_type host_id")
	if err != nil {
		log.Fatal(err)
	}
	defer raw.Close()

	err = json.NewDecoder(raw).Decode(&services)
	json.NewEncoder(w).Encode(services)
}

func returnService(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r.URL.Path)
	if err != nil {
		log.Fatal(err)
	}

	var services []Service
	raw, err := query(fmt.Sprintf("GET services\nFilter: id = %d\nColumns:id acknowledged check_period check_source check_type checks_enabled comments contacts description downtimes event_handler event_handler_enabled execution_time flap_detection_enabled groups has_been_checked in_check_period in_notification_period is_flapping last_check last_notification last_state_change last_time_critical last_time_ok last_time_unknown last_time_warning latency next_check next_notification notification_period notifications_enabled state state_type host_id", id))
	if err != nil {
		log.Fatal(err)
	}
	defer raw.Close()

	err = json.NewDecoder(raw).Decode(&services)
	if len(services) > 0 {
		json.NewEncoder(w).Encode(services[0])
	} else {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404 - Service not found"))
	}
}

func handleRequests() {
	http.HandleFunc("/comments", returnAllComments)
	http.HandleFunc("/comments/", returnComment)
	http.HandleFunc("/contacts", returnAllContacts)
	http.HandleFunc("/contacts/", returnContact)
	http.HandleFunc("/downtimes", returnAllDowntimes)
	http.HandleFunc("/downtimes/", returnDowntime)
	http.HandleFunc("/hosts", returnAllHosts)
	http.HandleFunc("/hosts/", returnHost)
	http.HandleFunc("/services", returnAllServices)
	http.HandleFunc("/services/", returnService)
	log.Fatal(http.ListenAndServe(*listenAddress, nil))
}

func main() {
	handleRequests()
}
