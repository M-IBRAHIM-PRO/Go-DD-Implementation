package membersvc

// RegisterMemberInput — what comes IN (from HTTP handler)
type RegisterMemberInput struct {
    Name  string
    Email string
}

// MemberOutput — what goes OUT (to HTTP handler)
// Notice: no domain types leak out. The app layer translates.
type MemberOutput struct {
    ID       string
    Name     string
    Email    string
    Status   string
    JoinedAt string
}