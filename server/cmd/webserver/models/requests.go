package models

type ContactRequest struct {
	FirstName string `json:"firstName" example:"John"`
	LastName  string `json:"lastName" example:"Doe"`
}

type AddressRequest struct {
	Line1         string  `json:"line1" example:"1600 Pennsylvania Ave."`
	Line2         *string `json:"line2,omitempty" example:"Ste. 1234"`
	City          string  `json:"city" example:"Washington"`
	StateProvince string  `json:"stateProvince" example:"DC"`
	PostalCode    string  `json:"postalCode" example:"20006"`
}

type PropertyRequest struct {
	Search string `json:"search" example:"search"`
}

type ProjectCreateRequest struct {
	Projectname string `json:"projectname" example:"John"`
}

type ProjectNameUpdateRequest struct {
	Projectname string `json:"projectname" example:"John"`
}

type InviteCollaboratorRequest struct {
	ProjectID int    `json:"projectID" example:19`
	InvitedBy string `json:"userID" example:"AUTH-sdhushduhsdu"`
	Name      string `json:"name" example:"John"`
	Email     string `json:"email" example:"John@vicesoftware.com"`
	Role      int    `json:"role" example:1`
	LastName  string `json:"lastname" example:"Miles"`
}

type SendEmailToSupportRequest struct {
	Email    string `json:"email" example:"sanjeev.sharma@vicesoftware.com"`
	Address  string `json:"address" example:2`
	Username string `json:"username" example:"Sanjeev"`
}

type AddPropertyGroupRequest struct {
	GroupID       string `json:"groupID" example:"tuighj-578uj"`
	ParcelID      int    `json:"parcelID" example:9`
	City          string `json:"city" example:"New Jersey""`
	StateProvince string `json:"stateProvince" example:"CA"`
	PostalCode    int    `json:"postalCode" example:"122002"`
	Country       string `json:"country" example:"USA"`
	PropertyInfo  string `json:"propertyInfo" example:"USA"`
	Jurisdiction  string `json:"jurisdiction" example:"USA"`
	StreetName    string `json:"StreetName" example:"New York Street"`
	StreetNumber  string `json:"streetNumber" example:"New York Street 44"`
	StreetSuffix  string `json:"streetSuffix" example:"AF"`
	Name          string `json:"name" example:"New York Street 44 NY"`
}

type NoteUpdateRequest struct {
	Note string `json:"note" example:"A quick brown fox"`
}

type PropertyUpdateRequest struct {
	ProjectID   int    `json:"projectID" example:"2"`
	ProjectName string `json:"projectName" example:"Project Name"`
}

type UserUpdateRequest struct {
	Name              string `json:"name" example:"Cody Miles"`
	Email             string `json:"email" example:"CodyMiles@vicesftware.com"`
	IsActive          bool   `json:"isActive" example:True`
	CreatedAt         string `json:"createdAt" example:"1554441489907"`
	UpdatedAt         string `json:"updatedAt" example:"1554441489907"`
	ImageUrl          string `json:"image_url" example:"https://s.gravatar.com/avatar/0a601b86076eb9fbf7d78da60a59ee8a?s=480&r=pg&d=https%3A%2F%2Fcdn.auth0.com%2Favatars%2Fsi.png"`
	AcceptedTermsDate string `json:"updatedAt" example:"1554441489907"`
	doIntro           bool   `json:"doIntro" example:True`
}

type UserRequest struct {
	Id                string `json:"id" example:"auth0|ad8b4fd5a359e62e34safas"`
	Name              string `json:"name" example:"Cody Miles"`
	Email             string `json:"email" example:"CodyMiles@vicesftware.com"`
	IsActive          bool   `json:"isActive" example:True`
	ImageUrl          string `json:"image_url" example:"https://s.gravatar.com/avatar/0a601b86076eb9fbf7d78da60a59ee8a?s=480&r=pg&d=https%3A%2F%2Fcdn.auth0.com%2Favatars%2Fsi.png"`
	CreatedAt         string `json:"createdAt" example:"1554441489907"`
	UpdatedAt         string `json:"updatedAt" example:"1554441489907"`
	AcceptedTermsDate string `json:"updatedAt" example:"1554441489907"`
	IsIntro           bool   `json:"isIntro" example:True`
}

type UpdatePropertyRequest struct {
	ParcelID      int    `json:"parcelID" example:9`
	City          string `json:"city" example:"New Jersey""`
	StateProvince string `json:"stateProvince" example:"CA"`
	PostalCode    int    `json:"postalCode" example:"122002"`
	Country       string `json:"country" example:"USA"`
	PropertyInfo  string `json:"propertyInfo" example:"USA"`
}

type DuplicateProjectRequest struct {
	ID int `json:"projectID" example:"2"`
}

type InviteUpdateRequest struct {
	IsInviteAccepted int `json:"isInviteAccepted" example:0`
}

type ProjectUserRoleRequest struct {
	userID    string `json:"userID" example:"auth|8834bhbdgx7733"`
	ProjectId int    `json:"projectId" example:1`
	Role      int    `json:"role" example:1`
}

type SubmitFeedbackRequest struct {
	UserId           string `json:"UserId" example:"auth|sdufyguebuhruehrue"`
	Description      string `json:"description" example:"A quick brown fox jumps over a little lazy dog"`
	FeedbackOverview string `json:"feedbackOverview" example:"Subject"`
	FeedbackCategory string `json:"feedbackCategory" example:"Category"`
}

type CompanyRequest struct {
	CompanyName         string `json:"companyName" example:"ashish.yadav@vicesoftware.com"`
	PhoneNumber         string `json:"phoneNumber" example:"auth|sdufyguebuhruehrue"`
	BillingContactEmail string `json:"email" example:"test@me.com"`
	BillingContactName  string `json:"email" example:"John Tester"`
	Address             string `json:"address" example:2`
	City                string `json:"city" example:"Gurgaon"`
	State               string `json:"state" example:"Haryana"`
	TimeZone            string `json:"timeZone" example:"IST"`
	Title               string `json:"title" example:"Software"`
	ZipCode             string `json:"zipCode" example:"122001"`
	CreatedBy           string `json:"createdBy" example:"AuthIISiosd"`
	CreatedAt           string `json:"createdAt" example:"89382984334"`
}

type CodeNoteRequest struct {
	CodeID          int    `json:"codeID" example:9`
	Description     string `json:"note" example:"notes"`
	HighLightedText string `json:"highLightedText" example:"notes"`
}

type CodeNoteUpdateRequest struct {
	Note string `json:"note" example:"A quick brown fox"`
}

type CodeNoteCommentRequest struct {
	NoteID      int    `json:"noteID" example:9`
	Description string `json:"note" example:"notes"`
}

type CodeRequest struct {
	TopicSummary string `json:"topicSummary" example:"Washington"`
	Category     string `json:"category" example:"1554441489907"`
	Number       string `json:"number" example:"$762387"`
	CreatedBy    string `json:"createdBy" example:"AuthIISiosd"`
	CreatedAt    string `json:"createdAt" example:"89382984334"`
}
type GroupRequest struct {
	//ID        string `json:"id" example:"3f06af63-a93c-11e4"`
	Name      string `json:"name" example:"John"`
	CreatedBy string `json:"createdBy" example:"1554441489907"`
	// UpdatedBy string `json:"updatedBy" example:"1554441489907"`
	CreatedAt int64 `json:"createdAt" example:"1554441489907"`
	//UpdatedAt int64 `json:"updatedAt" example:"1554441489907"`
	IsActive bool `json:"isActive" example:"true"`
}
type CodeCollectionAssociationRequest struct {
	ID               string `json:"id"  example:"9aa90bac-e5de-11e9-8529-97b20f595a0e"`
	CodeID           int    `json:"codeId"  example:"1"`
	CodeCollectionID string `json:"codeCollectionID" example:"9aa90bac-e5de-11e9-8529-97b20f595a0e"`
}
type CodeCollectionsRequest struct {
	ID        string `json:"id" example:"3f06af63-a93c-11e4"`
	Name      string `json:"name" example:"John"`
	CreatedBy string `json:"createdBy" example:"1554441489907"`
	UpdatedBy string `json:"updatedBy" example:"1554441489907"`
	CreatedAt int64  `json:"createdAt" example:"1554441489907"`
	UpdatedAt int64  `json:"updatedAt" example:"1554441489907"`
	IsActive  bool   `json:"isActive" example:"true"`
}

type GroupUpdateRequest struct {
	Name string `json:"name" example:"John"`
}

type ProjectGroupsRequest struct {
	ID        string `json:"id"  example:"9aa90bac-e5de-11e9-8529-97b20f595a0e"`
	ProjectID int    `json:"projectId"  example:"1"`
	GroupID   string `json:"groupId" example:"9aa90bac-e5de-11e9-8529-97b20f595a0e"`
}

type CodeCollectionsUpdateRequest struct {
	Name string `json:"name" example:"John"`
}
type CodeGroupAssociationRequest struct {
	ID      string `json:"id"  example:"9aa90bac-e5de-11e9-8529-97b20f595a0e"`
	CodeID  int    `json:"codeId"  example:"1"`
	GroupID string `json:"groupId" example:"9aa90bac-e5de-11e9-8529-97b20f595a0e"`
}
type PropertyGroupAssociationRequest struct {
	ID         string `json:"id"  example:"9aa90bac-e5de-11e9-8529-97b20f595a0e"`
	PropertyID int    `json:"propertyId"  example:"1"`
	GroupID    string `json:"groupId" example:"9aa90bac-e5de-11e9-8529-97b20f595a0e"`
}

type FilesRequest struct {
	ID        int    `json:"id" example:"9"`
	Url       string `json:"url" example:"http://www.google.com"`
	CreatedBy string `json:"createdBy" example:"1554441489907"`
	UpdatedBy string `json:"updatedBy" example:"1554441489907"`
	CreatedAt int64  `json:"createdAt" example:"1554441489907"`
	UpdatedAt int64  `json:"updatedAt" example:"1554441489907"`
	IsActive  bool   `json:"isActive" example:"true"`
}

type FilesUpdateRequest struct {
	Url string `json:"url" example:"http://www.google.com"`
}

type FilesGroupAssociationRequest struct {
	ID      string `json:"id"  example:"9aa90bac-e5de-11e9-8529-97b20f595a0e"`
	FileID  int    `json:"fileId"  example:"1"`
	GroupID string `json:"groupId" example:"9aa90bac-e5de-11e9-8529-97b20f595a0e"`
}

type GroupCodeCollectionAssociationRequest struct {
	ID               string `json:"id"  example:"9aa90bac-e5de-11e9-8529-97b20f595a0e"`
	GroupID          string `json:"groupId" example:"9aa90bac-e5de-11e9-8529-97b20f595a0e"`
	CodeCollectionID string `json:"codeCollectionId" example:"9aa90bac-e5de-11e9-8529-97b20f595a0e"`
}

type GroupUserAssociationRequest struct {
	ID      string `json:"id"  example:"9aa90bac-e5de-11e9-8529-97b20f595a0e"`
	GroupID string `json:"groupId" example:"9aa90bac-e5de-11e9-8529-97b20f595a0e"`
	UserID  string `json:"userId" example:"9aa90bac-e5de-11e9-8529-97b20f595a0e"`
}

type SignUpAuth0Request struct {
	Email    string `json:"email" example:"ashish.yadav218@gmail.com"`
	Password string `json: "password" example: "Ashish@1234"`
}

type ResetPasswordLink struct {
	Email string `json:"email" example:"ashish.yadav@vicesoftware.com"`
}

type ActiveCollaboratorRequest struct {
	UserID    string `json:"userID" example:"google-oauth2|100700139019068209339"`
	ProjectID int    `json:"projectID" example:1`
	RoleID    int    `json:"roleID" example:1`
	InviteID  int    `json:"inviteID" example:1`
}

type InvitePersonToCollaborateOnProject struct {
	Email           string                                   `json:"email" example:"ashish.yadav@vicesoftware.com"`
	InviteBy        string                                   `json:"invitedBy" example:"auth|sdufyguebuhruehrue"`
	Body            string                                   `json:"body" example:"A quick brown fox jumps over a little lazy dog"`
	ProjectId       int                                      `json:"projectId" example:2`
	Role            int                                      `json:"role" example:2`
	InviteId        string                                   `json:"inviteId" example:2"`
	Personalization []DemoInvitePersonToCollaborateOnProject `json: "personalization"`
}

type DemoInvitePersonToCollaborateOnProject struct {
	Personalization string `json:"email" example:"ashish.yadav@vicesoftware.com"`
	InviteBy        string `json:"invitedBy" example:"auth|sdufyguebuhruehrue"`
	Body            string `json:"body" example:"A quick brown fox jumps over a little lazy dog"`
	ProjectId       int    `json:"projectId" example:2`
	Role            int    `json:"role" example:2`
	InviteId        string `json:"inviteId" example:2"`
}

type Payload struct {
	From             From               `json:"from"`
	Personalizations []Personalizations `json:"personalizations"`
	TemplateID       string             `json:"template_id"`
}

type From struct {
	Email string `json:"email"`
}
type To struct {
	Email string `json:"email"`
}
type Items struct {
	Text  string `json:"text"`
	Image string `json:"image"`
	Price string `json:"price"`
}
type DynamicTemplateData struct {
	Header     string `json:"header"`
	Text       string `json:"text"`
	C2a_link   string `json:"c2a_link"`
	C2a_button string `json:"c2a_button"`
}
type Personalizations struct {
	To                  []To                `json:"to"`
	DynamicTemplateData DynamicTemplateData `json:"dynamic_template_data"`
}

type UserEmailRequest struct {
	Email string `json:"email" example:"cody@vicesoftware.com"`
}

type DuplicateRequest struct {
	IsDuplicateProperty          bool   `json:"isDuplicateProperty" example:"true"`
	IsDuplicateNotesFromProperty bool   `json:"isDuplicateNotesFromProperty" example:"true"`
	IsDuplicateCodes             bool   `json:"isDuplicateCodes" example:"true"`
	IsDuplicateNotesFromCodes    bool   `json:"isDuplicateNotesFromCodes" example:"true"`
	IsAllSelect                  bool   `json:"isAllSelect" example:"true"`
	Name                         string `json:"name" example:"true"`
}

type InviteGroupRequest struct {
	GroupID   string `json:"groupID" example:"546879-gb68-764"`
	InvitedBy string `json:"userID" example:"AUTH-sdhushduhsdu"`
	Name      string `json:"name" example:"John"`
	Email     string `json:"email" example:"John@vicesoftware.com"`
	Role      int    `json:"role" example:1`
	LastName  string `json:"lastname" example:"Miles"`
}

type CodeLibraryRequest struct {
	ID        string `json:"id" example:"3f06af63-a93c-11e4"`
	Name      string `json:"name" example:"John"`
	CreatedBy string `json:"createdBy" example:"1554441489907"`
	UpdatedBy string `json:"updatedBy" example:"1554441489907"`
	CreatedAt int64  `json:"createdAt" example:"1554441489907"`
	UpdatedAt int64  `json:"updatedAt" example:"1554441489907"`
	IsActive  bool   `json:"isActive" example:"true"`
}

type CompanyDetailsRequest struct {
	CompanyName         string `json: "companyName" example:"3f06af63-a93c-1dskg"`
	PhoneNumber         string `json: "phoneNumber" example:"445645612851"`
	Address             string `json: "streetAddress" example:"3428"`
	City                string `json: "city" example:"San Diago"`
	State               string `json: "state" example:"California"`
	BillingContactName  string `json: "billingContactName" example:"3f06af63-a93c-1dskg"`
	BillingContactEmail string `json: "billingContactEmail" example:"3f06af63-a93c-1dskg"`
	ZipCode             string `json: "zipCode" example:"3f06af63-a93c-1dskg"`
	CreatedBy           string `json:"createdBy" example:"1554441489907"`
	UpdatedBy           string `json:"updatedBy" example:"1554441489907"`
	IsActive            bool   `json:"isActive" example:"true"`
}

type DoIntroRequest struct {
	IsIntro bool `json:"isIntro" example:true`
}

type UserRolesRequest struct {
	InviteID int    `json:"inviteId" example:70`
	GroupID  string `json:"groupId" example:"546879-gb68-764"`
	RoleID   int    `json:"roleID" example:"2"`
}

type InviteResendRequest struct {
	InviteID int    `json:"inviteID" example:19`
	GroupID  string `json:"groupID" example:"19"`
}
