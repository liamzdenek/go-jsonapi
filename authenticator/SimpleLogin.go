package authenticator;

/*
import(
    . ".."
    "hash"
    "crypto/sha256"
    "math/rand"
    "time"
    "strconv"
    "errors"
);

type SimpleLogin struct {
    UserResource, SessionResource Resource
    SaltGenFunc func() string
    SessionIdGenFunc func() string
    CookieName string
    Hash hash.Hash
}

type SimpleLoginSession struct {
    Id int `meddler:"id,pk"`
    SessionId string `jsonapi:"id" meddler:"session_id"`
    UserId int `json:"-"`
    Created time.Time `json:"created"`
    Updated time.Time `json:"updated"`
}

func NewSimpleLogin(user_resource, session_resource Resource) *SimpleLogin {
    rand_str := func() string {
        var length = 10;
        var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
        out := make([]rune, length);
        for i := 0; i < length; i++ {
            out[i] = letters[rand.Intn(len(letters))];
        }
        return string(out)
    }
    return &SimpleLogin{
        UserResource: user_resource,
        SessionResource: session_resource,
        Hash: sha256.New(),
        CookieName: "go-session-id",
        SaltGenFunc: rand_str,
        SessionIdGenFunc: rand_str,
    }
}

func(sl *SimpleLogin) GetUserId(r *Request) *string {
    session_id, err := r.HttpRequest.Cookie(sl.CookieName);
    r.API.Logger.Debugf("Got login cookie: %#v\n", session_id);
    if err != nil {
        return nil;
    }
    record, err := sl.SessionResource.FindOne(r, RequestParams{}, session_id.Value);
    if err != nil {
        panic(err);
    }
    if record == nil {
        return nil;
    }
    attr := record.Attributes.(*SimpleLoginSession);
    attr.Updated = time.Now();
    sl.SessionResource.Update(r, record);
    id := strconv.Itoa(attr.UserId);
    return &id;
}

func(sl *SimpleLogin) FindDefault(r *Request, rp RequestParams) ([]*Record, error) {
    panic(TODO());
}

func(sl *SimpleLogin) FindOne(r *Request, rp RequestParams, id string) (*Record, error) {
    panic(TODO());
}

func(sl *SimpleLogin) FindMany(r *Request, rp RequestParams, ids []string) ([]*Record, error) {
    panic(Unimplemented());
}

func(sl *SimpleLogin) FindManyByField(r *Request, rp RequestParams, field, value string) ([]*Record, error) {
    panic(Unimplemented());
}

func(sl *SimpleLogin) Delete(r *Request, id string) error {
    if id != "self" {
        panic(NewResponderBaseErrors(403, errors.New("This endpoint only supports id=\"self\"")));
    }
    session_id, err := r.HttpRequest.Cookie(sl.CookieName);
    if err != nil {
        panic(NewResponderBaseErrors(401, errors.New("You must be logged in to log out")));
    }

    err = sl.SessionResource.Delete(r, session_id.Value);
    return err;
}

func(sl *SimpleLogin) ParseJSON(r *Request, src *Record, raw []byte) (*Record, error) {
    panic(TODO());
}

func(sl *SimpleLogin) Create(r *Request, record *Record) (status RecordCreatedStatus, err error) {
    panic(TODO());
}

func(sl *SimpleLogin) Update(r *Request, record *Record) error {
    panic(TODO());
}
*/
