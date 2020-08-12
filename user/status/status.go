package status

//Status user status type
type Status int

var StatusUnkown = Status(0)

//StatusNormal user status normal
var StatusNormal = Status(1)

//StatusBanned user status banned
var StatusBanned = Status(-1)

//StatusLabelNormal status label for normal
var StatusLabelNormal = "normal"

//StatusLabelBanned status label for banned
var StatusLabelBanned = "banned"

//Service user status service interface
type Service interface {
	//IsAvailable check is status available
	IsAvailable(Status) (bool, error)
	//Label get status label
	//Empty string will be returned if status invalid
	Label(Status) (string, error)
}

//PlainService plain service type
type PlainService struct {
	//AvailableStatuses available statuses map
	AvailableStatuses map[Status]bool
	//StatusLabels status labels map
	StatusLabels map[Status]string
}

//IsAvailable check is status available
func (p *PlainService) IsAvailable(s Status) (bool, error) {
	return p.AvailableStatuses[s], nil
}

//Label get status label
func (p *PlainService) Label(s Status) (string, error) {
	return p.StatusLabels[s], nil
}

//NewService create new service
func NewService() *PlainService {
	return &PlainService{
		AvailableStatuses: map[Status]bool{},
		StatusLabels:      map[Status]string{},
	}
}

//NormalOrBannedService service only contains status normal and banned.
var NormalOrBannedService = &PlainService{
	AvailableStatuses: map[Status]bool{StatusNormal: true},
	StatusLabels:      map[Status]string{StatusNormal: StatusLabelNormal, StatusBanned: StatusLabelBanned},
}
