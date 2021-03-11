package models

type ErrorResponse struct {
	Error string `json:"error"`
}

type PingResponse struct {
	Message string `json:"msg" example:"pong"`
}

type ContactResponse struct {
	ID        int               `json:"id" example:"1"`
	FirstName string            `json:"firstName" example:"John"`
	LastName  string            `json:"lastName" example:"Doe"`
	Addresses []AddressResponse `json:"addresses"`
	CreatedAt int64             `json:"createdAt" example:"1554441489907"`
	UpdatedAt int64             `json:"updatedAt" example:"1554441489907"`
}

type AddressResponse struct {
	ID            int     `json:"id" example:"1"`
	Line1         string  `json:"line1" example:"1600 Pennsylvania Ave."`
	Line2         *string `json:"line2,omitempty" example:"Ste. 1234"`
	City          string  `json:"city" example:"Washington"`
	StateProvince string  `json:"stateProvince" example:"DC"`
	PostalCode    string  `json:"postalCode" example:"20006"`
	CreatedAt     int64   `json:"createdAt" example:"1554441489907"`
	UpdatedAt     int64   `json:"updatedAt" example:"1554441489907"`
}

type ProjectResponse struct {
	ID            string             `json:"id" example = "kXavp95Lej8x6GQgAym4ONQZw2J3WE"`
	Name          string             `json:"name" example:"New Project"`
	CreatedAt     int64              `json:"createdAt" example:"1554441489907"`
	UpdatedAt     int64              `json:"updatedAt" example:"1554441489907"`
	CreatedBy     string             `json:"createdby" example:auth0|5c89fc609686f13c3cd189c9`
	Role          string             `json:"role" example:"admin"`
	PropertyCount int                `json:"propertyCount" example:"5"`
	Property      []PropertyResponse `json:"property"`
	ImageUrl      string             `json:"image_url" example:"john"`
}

type RecentProjectResponse struct {
	ID                 string `json:"id" example = "kXavp95Lej8x6GQgAym4ONQZw2J3WE"`
	ProjectName        string `json:"name" example:"New Project"`
	ProjectDescription string `json:"description,omitempty" example:"Ste. 1234"`
	CreatedAt          int64  `json:"createdAt" example:"1554441489907"`
	UpdatedAt          int64  `json:"updatedAt" example:"1554441489907"`
}

type DetailResponse struct {
	Name        string `json:"name" example:"New Project"`
	Description string `json:"description,omitempty" example:"Ste. 1234"`
	CreatedAt   int64  `json:"createdAt" example:"1554441489907"`
	UpdatedAt   int64  `json:"updatedAt" example:"1554441489907"`
}

// A Response struct to map the Entire Response
type Response []struct {
	Id       int `json:"id"`
	ParsedId int `json:"parcelId"`
}

type ScoutRedResponse []struct {
	ID       int    `json:"id"`
	ParcelID int    `json:"parcelId"`
	Apn      string `json:"apn"`
	ApnFmt   string `json:"apnFmt"`
	Full     string `json:"full"`
	Street   struct {
		Number         int         `json:"number"`
		NumberFraction interface{} `json:"numberFraction"`
		PreDirection   interface{} `json:"preDirection"`
		Name           string      `json:"name"`
		Suffix         interface{} `json:"suffix"`
		PostDirection  interface{} `json:"postDirection"`
	} `json:"street"`
	Unit         interface{} `json:"unit"`
	Postal       string      `json:"postal"`
	Jurisdiction string      `json:"jurisdiction"`
	State        string      `json:"state"`
	Country      string      `json:"country"`
	Created      int64       `json:"created"`
	Updated      int64       `json:"updated"`
}

type RecentActivityResponse struct {
	Name      string `json:"name" example:"New Project"`
	Activity  string `json:"activity,omitempty" example:"Ste. 1234"`
	CreatedAt int64  `json:"createdAt" example:"1554441489907"`
	UpdatedAt int64  `json:"updatedAt" example:"1554441489907"`
}

type UserResponse struct {
	ID                string `json:"id" example:"1"`
	Name              string `json:"name" example:"1600 Pennsylvania Ave."`
	Email             string `json:"email" example:"Washington"`
	IsActive          bool   `json:"isActive" example:True`
	ImageUrl          string `json:"image_url" example:"Washington"`
	CreatedAt         int64  `json:"createdAt" example:"1554441489907"`
	UpdatedAt         int64  `json:"updatedAt" example:"1554441489907"`
 }

type RoleResponse struct {
	ID        int    `json:"id" example:"1"`
	Role      string `json:"role" example:"Manager"`
	CreatedAt int64  `json:"createdAt" example:"1554441489907"`
	UpdatedAt int64  `json:"updatedAt" example:"1554441489907"`
}

type ProjectUserResponse struct {
	ID                 int            `json:"id" example:"1"`
	ProjectName        string         `json:"name" example:"New Project"`
	ProjectDescription string         `json:"description,omitempty" example:"Ste. 1234"`
	CreatedAt          int64          `json:"createdAt" example:"1554441489907"`
	UpdatedAt          int64          `json:"updatedAt" example:"1554441489907"`
	Role               []RoleResponse `json:"role"`
}

type InviteCollaboratorResponse struct {
	ID                         string                `json:"id" example:WkvLw35q21ylPAZPGMJepzEQZ4XRbB`
	ProjectID                  string                `json:"projectID" example:WkvLw35q21ylPAZPGMJepzEQZ4XRbB`
	InvitedBy                  string                `json:"invitedBy" example:"stringID"`
	IsInviteAccepted           int                   `json:"isInviteAccepted" example:0`
	Name                       string                `json:"name" example:"John"`
	Email                      string                `json:"email" example:"John@vicesoftware.com"`
	Role                       int                   `json:"role" example:1`
	Invitation_sent_at         int64                 `json:"invitation_sent_at" example:"1554441489907"`
	Invitation_accepted_at     int64                 `json:"invitation_accepted_at" example:"1554441489907"`
	LastName                   string                `json:"lastname" example:"Miles"`
	Group                      []InviteGroupResponse `json:"groups" example:"[]"`
	HasAccessToCompleteProject bool                  `json:"hasAccessToCompleteProject" example:"true"`
}

type PropertyResponse struct {
	ID            string         `json:"id" example:"WkvLw35q21ylPAZPGMJepzEQZ4XRbB"`
	ParcelID      int            `json:"parcelID" example:"111001"`
	ProjectID     string         `json:"projectID" example:"WkvLw35q21ylPAZPGMJepzEQZ4XRbB"`
	ProjectName   string         `json:"projectName" example:"John"`
	PropertyInfo  string         `json:"propertyInfo" example:"John,john,john,john"`
	City          string         `json:"city" example:"Washington"`
	Jurisdiction  string         `json:"jurisdiction" example:"jurisdiction"`
	StreetName    string         `json:"streetName" example:"Washington Street"`
	StreetNumber  string         `json:"streetNumber" example:"Washington Street 44"`
	StreetSuffix  string         `json:"streetSuffix" example:"DR"`
	StateProvince string         `json:"stateProvince" example:"DC"`
	PostalCode    int            `json:"postalCode" example:"20006"`
	CreatedAt     int64          `json:"createdAt" example:"1554441489907"`
	UpdatedAt     int64          `json:"updatedAt" example:"1554441489907"`
	CreatedBy     string         `json:"createdBy" example:"1554441489907"`
	UpdatedBy     string         `json:"updatedBy" example:"1554441489907"`
	Country       string         `json:"country" example:"USA"`
	Note          []NoteResponse `json:"note"`
	Name          string         `json:"name" example:"New York Street 44 NY"`
}

type PropertyGroupResponse struct {
	ID            int            `json:"id" example:"1"`
	ParcelID      int            `json:"parcelID" example:"111001"`
	GroupID       string         `json:"groupID" example:"12345-68-tihgbhj"`
	ProjectName   string         `json:"projectName" example:"John"`
	PropertyInfo  string         `json:"propertyInfo" example:"John,john,john,john"`
	City          string         `json:"city" example:"Washington"`
	Jurisdiction  string         `json:"jurisdiction" example:"jurisdiction"`
	StreetName    string         `json:"streetName" example:"Washington Street"`
	StreetNumber  string         `json:"streetNumber" example:"Washington Street 44"`
	StreetSuffix  string         `json:"streetSuffix" example:"DR"`
	StateProvince string         `json:"stateProvince" example:"DC"`
	PostalCode    int            `json:"postalCode" example:"20006"`
	CreatedAt     int64          `json:"createdAt" example:"1554441489907"`
	UpdatedAt     int64          `json:"updatedAt" example:"1554441489907"`
	CreatedBy     string         `json:"createdBy" example:"1554441489907"`
	UpdatedBy     string         `json:"updatedBy" example:"1554441489907"`
	Country       string         `json:"country" example:"USA"`
	Note          []NoteResponse `json:"note"`
	Name          string         `json:"name" example:"New York Street 44 NY"`
}

type NoteResponse struct {
	ID          string `json:"id" example:"WkvLw35q21ylPAZPGMJepzEQZ4XRbB"`
	Description string `json:"description" example:"Washington"`
	CreatedAt   int64  `json:"createdAt" example:"1554441489907"`
	UpdatedAt   int64  `json:"updatedAt" example:"1554441489907"`
	CreatedBy   string `json:"createdBy" example:"1554441489907"`
	UpdatedBy   string `json:"updatedBy" example:"1554441489907"`
}

type CompanyResponse struct {
	ID          string    `json:"id" example:"1"`
	CompanyName string `json:"companyName" example:"Innostax"`
	PhoneNumber string `json:"phoneNumber" example:"9813553355"`
	BillingContactEmail       string `json:"email" example:"ashish.yadav@vicesoftware.com"`
	Address     string `json:"address" example:"Gurgaon"`
	City        string `json:"city" example:"Gurgaon"`
	State       string `json:"state" example:"Haryana"`
	ZipCode     string `json:"zipCode" example:"122001"`
	TimeZone    string `json:"timeZone" example:"ITC"`
	Title       string `json:"title" example:"Software"`
	CreatedBy   string `json:"createdBy" example:"1554441489907"`
	CreatedAt   int64  `json:"createdAt" example:"1554441489907"`
}

type SearchCodesResponse struct {
	ID           int    `json:"id" example:"1"`
	TopicSummary string `json:"topicSummary" example:"Washington"`
	Category     string `json:"category" example:"1554441489907"`
	Number       string `json:"number" example:"$762387"`
}

type RolesResponse struct {
	ID        int    `json:"id" example:"1"`
	Role      string `json:"role" example:"Admin"`
	CreatedAt int64  `json:"createdAt" example:"1554441489907"`
	UpdatedAt int64  `json:"updatedAt" example:"1554441489907"`
}

type CodesResponse struct {
	ID           string `json:"id" example:"WkvLw35q21ylPAZPGMJepzEQZ4XRbB"`
	TopicSummary string `json:"topicSummary" example:"Washington"`
	Category     string `json:"category" example:"1554441489907"`
	Number       string `json:"number" example:"$762387"`
	CreatedBy    string `json:"createdBy" example:"1554441489907"`
	CreatedAt    int64  `json:"createdAt" example:"1554441489907"`
	ProjectID    string `json:"ProjectID" example:"WkvLw35q21ylPAZPGMJepzEQZ4XRbB"`
}

type CodeNoteResponse struct {
	ID          string `json:"id" example:"WkvLw35q21ylPAZPGMJepzEQZ4XRbB"`
	Description string `json:"description" example:"Washington"`
	Metdata     string `json:"metadata" example:"Washington"`
	ThreadCount int    `json:"threadCount" example:1`
	CreatedAt   int64  `json:"createdAt" example:"1554441489907"`
	UpdatedAt   int64  `json:"updatedAt" example:"1554441489907"`
	CreatedBy   string `json:"createdBy" example:"1554441489907"`
	UpdatedBy   string `json:"updatedBy" example:"1554441489907"`
}

type GroupResponse struct {
	ID        string `json:"id" example:"9aa90bac-e5de-11e9-8529-97b20f595a0e"`
	Name      string `json:"name" example:"John"`
	Role      int    `json:"role" example:1`
	ProjectID string `json:"projectID" example:WkvLw35q21ylPAZPGMJepzEQZ4XRbB`
	CreatedBy string `json:"createdBy" example:"1554441489907"`
	UpdatedBy string `json:"updatedBy" example:"1554441489907"`
	CreatedAt int64  `json:"createdAt" example:"1554441489907"`
	UpdatedAt int64  `json:"updatedAt" example:"1554441489907"`
	IsActive  bool   `json:"isActive" example:"true"`
}

type CodeCollectionAssociationResponse struct {
	ID               string `json:"id"  example:"9aa90bac-e5de-11e9-8529-97b20f595a0e"`
	CodeID           string `json:"codeId"  example:"WkvLw35q21ylPAZPGMJepzEQZ4XRbB"`
	CodeCollectionID string `json:"codeCollectionID" example:"9aa90bac-e5de-11e9-8529-97b20f595a0e"`
}

type CodeCollectionsResponse struct {
	ID        string `json:"id" example:"9aa90bac-e5de-11e9-8529-97b20f595a0e"`
	Name      string `json:"name" example:"John"`
	CreatedBy string `json:"createdBy" example:"1554441489907"`
	UpdatedBy string `json:"updatedBy" example:"1554441489907"`
	CreatedAt int64  `json:"createdAt" example:"1554441489907"`
	UpdatedAt int64  `json:"updatedAt" example:"1554441489907"`
	IsActive  bool   `json:"isActive" example:"true"`
}

type ProjectGroupsResponse struct {
	ID        string `json:"id"  example:"9aa90bac-e5de-11e9-8529-97b20f595a0e"`
	ProjectID string `json:"projectId"  example:"WkvLw35q21ylPAZPGMJepzEQZ4XRbB"`
	GroupID   string `json:"groupId" example:"9aa90bac-e5de-11e9-8529-97b20f595a0e"`
}

type CodeGroupAssociationResponse struct {
	ID      string `json:"id"  example:"9aa90bac-e5de-11e9-8529-97b20f595a0e"`
	CodeID  string `json:"codeId"  example:"WkvLw35q21ylPAZPGMJepzEQZ4XRbB"`
	GroupID string `json:"groupId" example:"9aa90bac-e5de-11e9-8529-97b20f595a0e"`
}

type PropertyGroupAssociationResponse struct {
	ID         string `json:"id"  example:"9aa90bac-e5de-11e9-8529-97b20f595a0e"`
	PropertyID int    `json:"propertyId"  example:"1"`
	GroupID    string `json:"groupId" example:"9aa90bac-e5de-11e9-8529-97b20f595a0e"`
}

type FilesResponse struct {
	ID        string `json:"id" example:"WkvLw35q21ylPAZPGMJepzEQZ4XRbB"`
	Url       string `json:"url" example:"http://www.google.com"`
	CreatedBy string `json:"createdBy" example:"1554441489907"`
	UpdatedBy string `json:"updatedBy" example:"1554441489907"`
	CreatedAt int64  `json:"createdAt" example:"1554441489907"`
	UpdatedAt int64  `json:"updatedAt" example:"1554441489907"`
	IsActive  bool   `json:"isActive" example:"true"`
}

type FilesGroupAssociationResponse struct {
	ID      string `json:"id"  example:"9aa90bac-e5de-11e9-8529-97b20f595a0e"`
	FileID  int    `json:"fileId"  example:"1"`
	GroupID string `json:"groupId" example:"9aa90bac-e5de-11e9-8529-97b20f595a0e"`
}
type GroupCodeCollectionAssociationResponse struct {
	ID               string `json:"id"  example:"9aa90bac-e5de-11e9-8529-97b20f595a0e"`
	GroupID          string `json:"groupId" example:"9aa90bac-e5de-11e9-8529-97b20f595a0e"`
	CodeCollectionID string `json:"codeCollectionId" example:"9aa90bac-e5de-11e9-8529-97b20f595a0e"`
}
type GroupUserAssociationResponse struct {
	ID      string `json:"id"  example:"9aa90bac-e5de-11e9-8529-97b20f595a0e"`
	GroupID string `json:"groupId" example:"9aa90bac-e5de-11e9-8529-97b20f595a0e"`
	UserID  string `json:"userId" example:"9aa90bac-e5de-11e9-8529-97b20f595a0e"`
}

type PropertyForGroupResponse struct {
	ID            string `json:"id" example:"WkvLw35q21ylPAZPGMJepzEQZ4XRbB"`
	ParcelID      int    `json:"parcelID" example:"111001"`
	PropertyInfo  string `json:"propertyInfo" example:"John,john,john,john"`
	City          string `json:"city" example:"Washington"`
	Jurisdiction  string `json:"jurisdiction" example:"jurisdiction"`
	StreetName    string `json:"streetName" example:"Washington Street"`
	StreetNumber  string `json:"streetNumber" example:"Washington Street 44"`
	StreetSuffix  string `json:"streetSuffix" example:"DR"`
	StateProvince string `json:"stateProvince" example:"DC"`
	PostalCode    int    `json:"postalCode" example:"20006"`
	Country       string `json:"country" example:"USA"`
	CreatedAt     int64  `json:"createdAt" example:"1554441489907"`
	UpdatedAt     int64  `json:"updatedAt" example:"1554441489907"`
	CreatedBy     string `json:"createdBy" example:"1554441489907"`
	UpdatedBy     string `json:"updatedBy" example:"1554441489907"`
	Name          string `json:"name" example:"New York Street 44 NY"`
}

type GroupPropertyCodesFilesResponse struct {
	ID              string                    `json:"id" example:"9aa90bac-e5de-11e9-8529-97b20f595a0e"`
	Name            string                    `json:"name" example:"John"`
	ProjectName     string                    `json:"projectName" example:"Demo Project"`
	ProjectID       string                    `json:"projectID" example:WkvLw35q21ylPAZPGMJepzEQZ4XRbB`
	ProjectRole     string                    `json:"projectRole" example:"admin"`
	CreatedBy       string                    `json:"createdBy" example:"155441489907"`
	UpdatedBy       string                    `json:"updatedBy" example:"1554441489907"`
	CreatedAt       int64                     `json:"createdAt" example:"1554441489907"`
	UpdatedAt       int64                     `json:"updatedAt" example:"1554441489907"`
	IsActive        bool                      `json:"isActive" example:"true"`
	Properties      PropertyResponse          `json:"properties" example:"true"`
	CodeCollections []CodeCollectionsResponse `json:"codeCollection"`
	Files           []FilesResponse           `json:"files"`
}

type UserAuth0Response struct {
	Id            string `json:"id" example:"5d8085e34c09310dc47280a9"`
	EmailVerified string `json:"email_verified" example:"false"`
	Email         string `json:"email" example:"ashish.yadav218@gmail.com"`
}

type Resp struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type ActiveCollaboratorResponse struct {
	UserID    string `json:"userID" example:"google-oauth2|100700139019068209339"`
	ProjectID int    `json:"projectID" example:1`
	RoleID    int    `json:"roleID" example:1`
	InviteID  int    `json:"inviteID" example:1`
}

type GroupCodesFilesRespone struct {
	ID              string                    `json:"id" example:"9aa90bac-e5de-11e9-8529-97b20f595a0e"`
	Name            string                    `json:"name" example:"John"`
	CreatedBy       string                    `json:"createdBy" example:"1554441489907"`
	UpdatedBy       string                    `json:"updatedBy" example:"1554441489907"`
	CreatedAt       int64                     `json:"createdAt" example:"1554441489907"`
	UpdatedAt       int64                     `json:"updatedAt" example:"1554441489907"`
	IsActive        bool                      `json:"isActive" example:"true"`
	CodeCollections []CodeCollectionsResponse `json:"codeCollection"`
	Files           []FilesResponse           `json:"files"`
}

type InviteGroupResponse struct {
	ID               string `json:"id" example:"1"`
	GroupID          string `json:"groupID" example:"546879-gb68-764"`
	IsInviteAccepted int    `json:"isInviteAccepted" example:0`
	RoleID           int    `json:"role" example:1`
	InviteID         int    `json:"inviteID" example:12`
	GroupName        string `json:"groupName" example:"12"`
}

type CodeLibraryResponse struct {
	ID        string `json:"id" example:"9aa90bac-e5de-11e9-8529-97b20f595a0e"`
	Name      string `json:"name" example:"John"`
	CreatedBy string `json:"createdBy" example:"1554441489907"`
	UpdatedBy string `json:"updatedBy" example:"1554441489907"`
	CreatedAt int64  `json:"createdAt" example:"1554441489907"`
	UpdatedAt int64  `json:"updatedAt" example:"1554441489907"`
	IsActive  bool   `json:"isActive" example:"true"`
}

type InviteGroupUserResponse struct {
	ID                   int             `json:"id" example:"1"`
	Group                []GroupResponse `json:"groupName" example:"Test group"`
	FirstName            string          `json:"firstname" example:"John"`
	LastName             string          `json:"lastname" example:"Miles"`
	IsInvitationAccepted int             `json:"isInviteAccepted" example:0`
	Email                string          `json:"email" example:"John@vicesoftware.com"`
}

type CompanyDetailsResponse struct {
	ID                  int    `json:"id" example:"9aa90bac-e5de-11e9-8529-97b20f595a0e"`
	CompanyName         string `json:"companyName" example:"Company Test Pvt Ltd"`
	PhoneNumber         string `json:"phoneNumber" example:"9999999999"`
	Address             string `json:"address" example:"34289 talbot street"`
	City                string `json:"city" example:"San Diago"`
	State               string `json:"state" example:"California"`
	BillingContactName  string `json:"billingContactName" example:"Aditya"`
	BillingContactEmail string `json:"billingContactEmail" example:"aditya.vashistha@vicesoftware.com"`
	ZipCode             string `json:"zipCode" example:"122001"`
	CreatedAt           int64  `json:"createdAt" example:"15544414899074"`
	CreatedBy           string `json:"createdBy" example:"15544414899074"`
	PlanExpDate         int64  `json:"planExpDate" example:"15544414899074"`
	UpdatedAt           int64  `json:"updatedAt" example:"15544414899074"`
	UpdatedBy           string `json:"updatedBy" example:"15544414899074"`
}

type InviteUserResponse struct {
	ID       string `json:"id" example:"9aa90bac-e5de-11e9-8529-97b20f595a0e"`
	InviteID int    `json:"inviteId" example:80`
	GroupID  string `json:"groupId" example:"9aa90bac-e5de-11e9-8529-97b20f595a0e"`
	// Name                   string		`json:"name" example:"Aditya"`
	// Email                  string		`json:"email" example:"aditya.vashistha@vicesoftware.com"`
	RoleID int `json:"role" example:2`
	// LastName               string		`json:"lastName" example:"Vashistha"`
}

type InviterDetailResponse struct {
	ID          string                `json:"id" example:"1"`
	Name        string                `json:"groupID" example:"546879-gb68-764"`
	Image       string                `json:"isInviteAccepted" example:"0"`
	ProjectID   int                   `json:"role" example:1`
	ProjectName string                `json:"inviteID" example:"12"`
	GroupName   string                `json:"groupName" example:"12"`
	Group       []InviteGroupResponse `json:"groups" example:"[]"`
}

type InvitationGroupResponse struct {
	ID                         string          `json:"id" example:"9aa90bac-e5de-11e9-8529-97b20f595a0e"`
	ProjectID                  string          `json:"projectid" example:"9aa90bac-e5de-11e9-8529-97b20f595a0e"`
	InvitedBy                  string          `json:"invitedBy" example:"stringID"`
	IsInviteAccepted           int             `json:"isInviteAccepted" example:0`
	Name                       string          `json:"name" example:"John"`
	Email                      string          `json:"email" example:"John@vicesoftware.com"`
	Role                       int             `json:"role" example:1`
	Invitation_sent_at         int64           `json:"invitation_sent_at" example:"1554441489907"`
	Invitation_accepted_at     int64           `json:"invitation_accepted_at" example:"1554441489907"`
	LastName                   string          `json:"lastname" example:"Miles"`
	Group                      []GroupResponse `json:"groups" example:"[]"`
	HasAccessToCompleteProject bool            `json:"hasAccessToCompleteProject" example:"true"`
}
