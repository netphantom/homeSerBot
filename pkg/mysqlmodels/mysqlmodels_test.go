package mysqlmodels

import (
	"errors"
	"testing"

	tb "gopkg.in/tucnak/telebot.v2"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *DbModel {
	t.Helper()

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}

	err = db.SetupJoinTable(&User{}, "Subscription", &UserProcess{})
	if err != nil {
		t.Fatalf("failed to setup join table: %v", err)
	}

	err = db.AutoMigrate(&User{}, &Process{}, &Notification{})
	if err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}

	return &DbModel{Db: db}
}

func makeUser(username string) *User {
	return &User{
		User: tb.User{
			ID:        0,
			FirstName: "Test",
			LastName:  "User",
			Username:  username,
		},
		Allowed:      false,
		Password:     nil,
		Subscription: nil,
	}
}

func TestRegisterUser(t *testing.T) {
	m := setupTestDB(t)

	user := makeUser("testuser")
	err := m.RegisterUser(user)
	if err != nil {
		t.Fatalf("RegisterUser failed: %v", err)
	}
	if user.Id == 0 {
		t.Fatal("expected user ID to be set after registration")
	}
}

func TestRegisterUser_DuplicateUsername(t *testing.T) {
	m := setupTestDB(t)

	user1 := makeUser("dupe")
	err := m.RegisterUser(user1)
	if err != nil {
		t.Fatalf("first RegisterUser failed: %v", err)
	}

	user2 := makeUser("dupe")
	err = m.RegisterUser(user2)
	if err != nil {
		t.Fatalf("second RegisterUser should succeed (no unique constraint on username): %v", err)
	}
	if user2.Id == 0 {
		t.Fatal("expected second user to get an ID")
	}
}

func TestUserByUsername_Found(t *testing.T) {
	m := setupTestDB(t)

	original := makeUser("findme")
	err := m.RegisterUser(original)
	if err != nil {
		t.Fatalf("RegisterUser failed: %v", err)
	}

	found := m.UserByUsername("findme")
	if found == nil {
		t.Fatal("expected to find user, got nil")
	}
	if found.Username != "findme" {
		t.Fatalf("expected username 'findme', got '%s'", found.Username)
	}
}

func TestUserByUsername_NotFound(t *testing.T) {
	m := setupTestDB(t)

	found := m.UserByUsername("nonexistent")
	if found != nil {
		t.Fatal("expected nil for nonexistent user")
	}
}

func TestVerifyId_Found(t *testing.T) {
	m := setupTestDB(t)

	original := makeUser("verifyid")
	err := m.RegisterUser(original)
	if err != nil {
		t.Fatalf("RegisterUser failed: %v", err)
	}

	found := m.VerifyId(original.Id)
	if found == nil {
		t.Fatal("expected to find user, got nil")
	}
	if found.Username != "verifyid" {
		t.Fatalf("expected username 'verifyid', got '%s'", found.Username)
	}
}

func TestVerifyId_NotFound(t *testing.T) {
	m := setupTestDB(t)

	found := m.VerifyId(999)
	if found != nil {
		t.Fatal("expected nil for nonexistent ID")
	}
}

func TestAuthenticate_NoPassword(t *testing.T) {
	m := setupTestDB(t)

	user := makeUser("nopass")
	err := m.RegisterUser(user)
	if err != nil {
		t.Fatalf("RegisterUser failed: %v", err)
	}

	id, err := m.Authenticate("nopass", "")
	if err != nil {
		t.Fatalf("Authenticate failed: %v", err)
	}
	if id != int(user.Id) {
		t.Fatalf("expected id %d, got %d", user.Id, id)
	}
}

func TestAuthenticate_WithPassword(t *testing.T) {
	m := setupTestDB(t)

	user := makeUser("withpass")
	err := m.RegisterUser(user)
	if err != nil {
		t.Fatalf("RegisterUser failed: %v", err)
	}

	err = m.ChangePsw("newpassword", "", int(user.Id))
	if err != nil {
		t.Fatalf("ChangePsw failed: %v", err)
	}

	id, err := m.Authenticate("withpass", "newpassword")
	if err != nil {
		t.Fatalf("Authenticate failed: %v", err)
	}
	if id != int(user.Id) {
		t.Fatalf("expected id %d, got %d", user.Id, id)
	}
}

func TestAuthenticate_WrongPassword(t *testing.T) {
	m := setupTestDB(t)

	user := makeUser("wrongpass")
	err := m.RegisterUser(user)
	if err != nil {
		t.Fatalf("RegisterUser failed: %v", err)
	}

	err = m.ChangePsw("correct", "", int(user.Id))
	if err != nil {
		t.Fatalf("ChangePsw failed: %v", err)
	}

	_, err = m.Authenticate("wrongpass", "wrong")
	if !errors.Is(err, ErrInvalidCredentials) {
		t.Fatalf("expected ErrInvalidCredentials, got %v", err)
	}
}

func TestAuthenticate_NonexistentUser(t *testing.T) {
	m := setupTestDB(t)

	_, err := m.Authenticate("ghost", "x")
	if !errors.Is(err, ErrInvalidCredentials) {
		t.Fatalf("expected ErrInvalidCredentials, got %v", err)
	}
}

func TestChangePsw_FirstTime(t *testing.T) {
	m := setupTestDB(t)

	user := makeUser("firstpsw")
	err := m.RegisterUser(user)
	if err != nil {
		t.Fatalf("RegisterUser failed: %v", err)
	}

	err = m.ChangePsw("hunter2", "", int(user.Id))
	if err != nil {
		t.Fatalf("ChangePsw failed: %v", err)
	}

	id, err := m.Authenticate("firstpsw", "hunter2")
	if err != nil {
		t.Fatalf("Authenticate after ChangePsw failed: %v", err)
	}
	if id != int(user.Id) {
		t.Fatalf("expected id %d, got %d", user.Id, id)
	}
}

func TestChangePsw_WithCurrentPassword(t *testing.T) {
	m := setupTestDB(t)

	user := makeUser("changepsw")
	err := m.RegisterUser(user)
	if err != nil {
		t.Fatalf("RegisterUser failed: %v", err)
	}

	err = m.ChangePsw("oldpass", "", int(user.Id))
	if err != nil {
		t.Fatalf("first ChangePsw failed: %v", err)
	}

	err = m.ChangePsw("newpass", "oldpass", int(user.Id))
	if err != nil {
		t.Fatalf("second ChangePsw failed: %v", err)
	}

	_, err = m.Authenticate("changepsw", "oldpass")
	if !errors.Is(err, ErrInvalidCredentials) {
		t.Fatal("old password should no longer work")
	}

	id, err := m.Authenticate("changepsw", "newpass")
	if err != nil {
		t.Fatalf("Authenticate with new password failed: %v", err)
	}
	if id != int(user.Id) {
		t.Fatalf("expected id %d, got %d", user.Id, id)
	}
}

func TestChangePsw_WrongCurrentPassword(t *testing.T) {
	m := setupTestDB(t)

	user := makeUser("wrongcurrent")
	err := m.RegisterUser(user)
	if err != nil {
		t.Fatalf("RegisterUser failed: %v", err)
	}

	err = m.ChangePsw("realpass", "", int(user.Id))
	if err != nil {
		t.Fatalf("first ChangePsw failed: %v", err)
	}

	err = m.ChangePsw("newpass", "wrong", int(user.Id))
	if !errors.Is(err, ErrInvalidCredentials) {
		t.Fatalf("expected ErrInvalidCredentials, got %v", err)
	}
}

func TestListNewUsers(t *testing.T) {
	m := setupTestDB(t)

	users, err := m.ListNewUsers()
	if err != nil {
		t.Fatalf("ListNewUsers failed: %v", err)
	}
	if len(users) != 0 {
		t.Fatalf("expected 0 new users, got %d", len(users))
	}

	u1 := makeUser("user1")
	u2 := makeUser("user2")
	_ = m.RegisterUser(u1)
	_ = m.RegisterUser(u2)

	users, err = m.ListNewUsers()
	if err != nil {
		t.Fatalf("ListNewUsers failed: %v", err)
	}
	if len(users) != 2 {
		t.Fatalf("expected 2 new users, got %d", len(users))
	}
}

func TestAllowUser(t *testing.T) {
	m := setupTestDB(t)

	user := makeUser("toallow")
	err := m.RegisterUser(user)
	if err != nil {
		t.Fatalf("RegisterUser failed: %v", err)
	}

	err = m.AllowUser("toallow")
	if err != nil {
		t.Fatalf("AllowUser failed: %v", err)
	}

	found := m.UserByUsername("toallow")
	if found == nil {
		t.Fatal("expected to find user")
	}
	if !found.Allowed {
		t.Fatal("expected user to be allowed")
	}
}

func TestAllowUser_Nonexistent(t *testing.T) {
	m := setupTestDB(t)

	err := m.AllowUser("ghost")
	if !errors.Is(err, ErrNoRecord) {
		t.Fatalf("expected ErrNoRecord, got %v", err)
	}
}

func TestNotAllowUser(t *testing.T) {
	m := setupTestDB(t)

	user := makeUser("toremove")
	err := m.RegisterUser(user)
	if err != nil {
		t.Fatalf("RegisterUser failed: %v", err)
	}

	err = m.NotAllowUser("toremove")
	if err != nil {
		t.Fatalf("NotAllowUser failed: %v", err)
	}

	found := m.UserByUsername("toremove")
	if found != nil {
		t.Fatal("expected user to be deleted")
	}
}

func TestNotAllowUser_Nonexistent(t *testing.T) {
	m := setupTestDB(t)

	err := m.NotAllowUser("ghost")
	if !errors.Is(err, ErrNoRecord) {
		t.Fatalf("expected ErrNoRecord, got %v", err)
	}
}

func TestListAllUsers(t *testing.T) {
	m := setupTestDB(t)

	all, err := m.ListAllUsers()
	if err != nil {
		t.Fatalf("ListAllUsers failed: %v", err)
	}
	if len(all) != 0 {
		t.Fatalf("expected 0 users, got %d", len(all))
	}

	_ = m.RegisterUser(makeUser("u1"))
	_ = m.RegisterUser(makeUser("u2"))
	_ = m.RegisterUser(makeUser("u3"))

	all, err = m.ListAllUsers()
	if err != nil {
		t.Fatalf("ListAllUsers failed: %v", err)
	}
	if len(all) != 3 {
		t.Fatalf("expected 3 users, got %d", len(all))
	}
}

func TestAddAndListProcesses(t *testing.T) {
	m := setupTestDB(t)

	m.AddProcess("nginx.service", "Nginx web server")
	m.AddProcess("postgresql.service", "PostgreSQL database")

	list, err := m.ProcessList()
	if err != nil {
		t.Fatalf("ProcessList failed: %v", err)
	}
	if len(list) != 2 {
		t.Fatalf("expected 2 processes, got %d", len(list))
	}
}

func TestGetProcessInfo(t *testing.T) {
	m := setupTestDB(t)

	m.AddProcess("test.service", "Test process")

	list, _ := m.ProcessList()
	pid := int(list[0].ID)

	p, err := m.GetProcessInfo(pid)
	if err != nil {
		t.Fatalf("GetProcessInfo failed: %v", err)
	}
	if p.Name != "test.service" {
		t.Fatalf("expected 'test.service', got '%s'", p.Name)
	}
	if p.Description != "Test process" {
		t.Fatalf("expected 'Test process', got '%s'", p.Description)
	}
}

func TestGetProcessInfo_NotFound(t *testing.T) {
	m := setupTestDB(t)

	_, err := m.GetProcessInfo(999)
	if err == nil {
		t.Fatal("expected error for nonexistent process")
	}
}

func TestUpdateDescription(t *testing.T) {
	m := setupTestDB(t)

	m.AddProcess("updatable.service", "Old description")
	list, _ := m.ProcessList()
	pid := int(list[0].ID)

	err := m.UpdateDescription(pid, "New description")
	if err != nil {
		t.Fatalf("UpdateDescription failed: %v", err)
	}

	p, _ := m.GetProcessInfo(pid)
	if p.Description != "New description" {
		t.Fatalf("expected 'New description', got '%s'", p.Description)
	}
}

func TestDeleteProcess(t *testing.T) {
	m := setupTestDB(t)

	m.AddProcess("delete.service", "To be deleted")
	list, _ := m.ProcessList()
	pid := int(list[0].ID)

	err := m.DeleteProcess(pid)
	if err != nil {
		t.Fatalf("DeleteProcess failed: %v", err)
	}

	_, err = m.GetProcessInfo(pid)
	if err == nil {
		t.Fatal("expected error after deletion")
	}
}

func TestSubscribeAndUnsubscribe(t *testing.T) {
	m := setupTestDB(t)

	err := m.RegisterUser(makeUser("subuser"))
	if err != nil {
		t.Fatalf("RegisterUser failed: %v", err)
	}
	user := m.UserByUsername("subuser")
	if user == nil {
		t.Fatal("expected to find user after registration")
	}

	m.AddProcess("proc1.service", "Process 1")
	m.AddProcess("proc2.service", "Process 2")
	list, _ := m.ProcessList()
	pid1 := int(list[0].ID)
	pid2 := int(list[1].ID)

	sub, err := m.SubscribeToProcess(user, pid1)
	if err != nil {
		t.Fatalf("SubscribeToProcess failed: %v", err)
	}
	if sub.Name != "proc1.service" {
		t.Fatalf("expected 'proc1.service', got '%s'", sub.Name)
	}

	_, err = m.SubscribeToProcess(user, pid2)
	if err != nil {
		t.Fatalf("second SubscribeToProcess failed: %v", err)
	}

	subs := m.ListSubscribed(user)
	if len(subs) != 2 {
		t.Fatalf("expected 2 subscriptions, got %d", len(subs))
	}

	err = m.UnsubscribeToProcess(user, pid1)
	if err != nil {
		t.Fatalf("UnsubscribeToProcess failed: %v", err)
	}

	subs = m.ListSubscribed(user)
	if len(subs) != 1 {
		t.Fatalf("expected 1 subscription after unsub, got %d", len(subs))
	}
	if subs[0].Name != "proc2.service" {
		t.Fatalf("expected remaining subscription 'proc2.service', got '%s'", subs[0].Name)
	}
}

func TestSubscribeToProcess_NotFound(t *testing.T) {
	m := setupTestDB(t)

	user := makeUser("failuser")
	err := m.RegisterUser(user)
	if err != nil {
		t.Fatalf("RegisterUser failed: %v", err)
	}

	_, err = m.SubscribeToProcess(user, 999)
	if err == nil {
		t.Fatal("expected error subscribing to nonexistent process")
	}
}

func TestListSubscribed_Empty(t *testing.T) {
	m := setupTestDB(t)

	err := m.RegisterUser(makeUser("nosubs"))
	if err != nil {
		t.Fatalf("RegisterUser failed: %v", err)
	}
	user := m.UserByUsername("nosubs")
	if user == nil {
		t.Fatal("expected to find user after registration")
	}

	subs := m.ListSubscribed(user)
	if subs == nil {
		t.Fatal("expected empty slice, not nil, for user with no subscriptions")
	}
	if len(subs) != 0 {
		t.Fatalf("expected 0 subscriptions, got %d", len(subs))
	}
}

func TestAddNotification(t *testing.T) {
	m := setupTestDB(t)

	user := makeUser("notifuser")
	err := m.RegisterUser(user)
	if err != nil {
		t.Fatalf("RegisterUser failed: %v", err)
	}

	m.AddProcess("testproc.service", "Test")
	list, _ := m.ProcessList()
	pid := int(list[0].ID)

	n := &Notification{
		UserID:    int(user.Id),
		ProcessID: pid,
		Active:    "active",
		Process:   "0",
		Sent:      false,
	}
	m.AddNotification(n)

	notifs := m.UserProcessNotification(user)
	if len(notifs) != 1 {
		t.Fatalf("expected 1 notification, got %d", len(notifs))
	}
	if notifs[0].Active != "active" {
		t.Fatalf("expected 'active', got '%s'", notifs[0].Active)
	}
}

func TestAddNotification_Dedup(t *testing.T) {
	m := setupTestDB(t)

	user := makeUser("dedupuser")
	err := m.RegisterUser(user)
	if err != nil {
		t.Fatalf("RegisterUser failed: %v", err)
	}

	m.AddProcess("dedup.service", "Dedup test")
	list, _ := m.ProcessList()
	pid := int(list[0].ID)

	n1 := &Notification{
		UserID:    int(user.Id),
		ProcessID: pid,
		Active:    "active",
		Process:   "0",
		Sent:      false,
	}
	m.AddNotification(n1)

	n2 := &Notification{
		UserID:    int(user.Id),
		ProcessID: pid,
		Active:    "active",
		Process:   "0",
		Sent:      false,
	}
	m.AddNotification(n2)

	notifs := m.UserProcessNotification(user)
	if len(notifs) != 1 {
		t.Fatalf("expected 1 notification after dedup, got %d", len(notifs))
	}
}

func TestAddNotification_NoDedupOnDifferentStatus(t *testing.T) {
	m := setupTestDB(t)

	user := makeUser("nodup")
	err := m.RegisterUser(user)
	if err != nil {
		t.Fatalf("RegisterUser failed: %v", err)
	}

	m.AddProcess("status.service", "Status test")
	list, _ := m.ProcessList()
	pid := int(list[0].ID)

	n1 := &Notification{
		UserID:    int(user.Id),
		ProcessID: pid,
		Active:    "active",
		Process:   "0",
		Sent:      false,
	}
	m.AddNotification(n1)

	n2 := &Notification{
		UserID:    int(user.Id),
		ProcessID: pid,
		Active:    "inactive",
		Process:   "1",
		Sent:      false,
	}
	m.AddNotification(n2)

	notifs := m.UserProcessNotification(user)
	if len(notifs) != 2 {
		t.Fatalf("expected 2 notifications (different status), got %d", len(notifs))
	}
}

func TestMarkAsSent(t *testing.T) {
	m := setupTestDB(t)

	user := makeUser("markuser")
	err := m.RegisterUser(user)
	if err != nil {
		t.Fatalf("RegisterUser failed: %v", err)
	}

	m.AddProcess("mark.service", "Mark test")
	list, _ := m.ProcessList()
	pid := int(list[0].ID)

	n := &Notification{
		UserID:    int(user.Id),
		ProcessID: pid,
		Active:    "active",
		Process:   "0",
		Sent:      false,
	}
	m.AddNotification(n)

	notifs := m.UserProcessNotification(user)
	if len(notifs) != 1 {
		t.Fatalf("expected 1 notification before marking, got %d", len(notifs))
	}

	m.MarkAsSent(&notifs[0])

	notifs = m.UserProcessNotification(user)
	if len(notifs) != 0 {
		t.Fatalf("expected 0 notifications after marking sent, got %d", len(notifs))
	}
}

func TestRemoveNotification(t *testing.T) {
	m := setupTestDB(t)

	user := makeUser("remuser")
	err := m.RegisterUser(user)
	if err != nil {
		t.Fatalf("RegisterUser failed: %v", err)
	}

	m.AddProcess("rem.service", "Remove test")
	list, _ := m.ProcessList()
	pid := int(list[0].ID)

	n := &Notification{
		UserID:    int(user.Id),
		ProcessID: pid,
		Active:    "active",
		Process:   "0",
		Sent:      false,
	}
	m.AddNotification(n)

	notifs := m.UserProcessNotification(user)
	m.RemoveNotification(&notifs[0])

	notifs = m.UserProcessNotification(user)
	if len(notifs) != 0 {
		t.Fatalf("expected 0 notifications after removal, got %d", len(notifs))
	}
}
