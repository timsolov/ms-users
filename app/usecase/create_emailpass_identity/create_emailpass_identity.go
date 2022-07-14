package create_emailpass_identity

import (
	"context"
	"fmt"
	"ms-users/app/common/event"
	"ms-users/app/common/jsonschema"
	"ms-users/app/common/password"
	"ms-users/app/domain"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/tidwall/gjson"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	confirmEmailTemplateName = "confirm-registration-email-pass"
)

// Repository describes repository contract
type Repository interface {
	// CreateUserAggregate creates new ident record with user record.
	CreateUserAggregate(ctx context.Context, ua *domain.UserAggregate, confirm *domain.Confirm, events []event.Event) error
}

// Params describes parameters
type Params struct {
	Email          string
	EmailConfirmed bool
	Password       string
	Profile        []byte // json
}

// UseCase describes usecase
type UseCase struct {
	repo           Repository
	baseURL        string
	fromEmail      string
	fromName       string
	confirmLife    time.Duration
	jsonSchema     jsonschema.Schema
	jsonSchemaName string
}

func New(
	repo Repository,
	baseURL,
	fromEmail,
	fromName string,
	confirmLife time.Duration,
	jsonSchema jsonschema.Schema,
	jsonSchemaName string,
) *UseCase {
	uc := &UseCase{
		repo:           repo,
		baseURL:        baseURL,
		fromEmail:      fromEmail,
		fromName:       fromName,
		confirmLife:    confirmLife,
		jsonSchema:     jsonSchema,
		jsonSchemaName: jsonSchemaName,
	}

	return uc
}

func (uc *UseCase) Do(ctx context.Context, cmd *Params) (profileID uuid.UUID, err error) {
	// create encrypted password
	encryptedPass, err := password.Encrypt(cmd.Password)
	if err != nil {
		return uuid.Nil, errors.Wrap(err, "encrypting password by bcrypt")
	}

	err = uc.ValidateProfile(ctx, cmd.Profile)
	if err != nil {
		return uuid.Nil, err // 400
	}

	// fill the profile

	// TODO: when we'll have different identities we should check does identity
	//       already exist and attach new identity to existing user
	// new user record
	user := domain.User{
		UserID:  uuid.New(),
		View:    uc.jsonSchemaName,
		Profile: cmd.Profile,
	}

	// new identity record
	ident := domain.Ident{
		UserID:         user.UserID,
		Ident:          cmd.Email,
		IdentConfirmed: false,
		Kind:           domain.EmailPassIdent,
		Password:       encryptedPass,
	}

	// fill user aggregate
	ua := domain.UserAggregate{
		User:   user,
		Idents: []domain.Ident{ident},
	}

	confirmRecord, err := uc.PrepareConfirmEmailRecord(cmd.Email)
	if err != nil {
		return uuid.Nil, errors.Wrap(err, "prepare confirm email record") // 500
	}

	// encode to json object which encoded by base64
	code, err := confirmRecord.ToBase64JSON()
	if err != nil {
		err = errors.Wrap(err, "encode confirm struct to base64") // 500
		return
	}

	confirmEmailEvent, err := uc.PrepareConfirmEmailEvent(cmd.Email, "en", code, cmd.Profile)
	if err != nil {
		return uuid.Nil, errors.Wrap(err, "prepare confirm email event") // 500
	}

	err = uc.repo.CreateUserAggregate(ctx, &ua, &confirmRecord, []event.Event{confirmEmailEvent})
	if err != nil {
		if errors.Cause(err) == domain.ErrNotUnique {
			err = domain.ErrIdentityDuplicated // 400
		}
		return uuid.Nil, errors.Wrap(err, "create user record in db") // 500
	}

	return user.UserID, nil
}

// ValidateProfile validates profile object by jsonschema and build gRPC status error with bad request.
func (uc *UseCase) ValidateProfile(ctx context.Context, profile []byte) error {
	errs, err := uc.jsonSchema.ValidateBytes(ctx, profile)
	if err != nil {
		return errors.Wrap(err, "jsonschema validation") // 500
	}
	// if there's no errors
	if len(errs) == 0 { // 200
		return nil
	}

	br := &errdetails.BadRequest{}

	for _, e := range errs {
		v := &errdetails.BadRequest_FieldViolation{
			Field:       e.PropertyPath,
			Description: e.Message,
		}

		br.FieldViolations = append(br.FieldViolations, v)
	}

	st := status.New(codes.InvalidArgument, "invalid request")
	std, err := st.WithDetails(br)
	if err != nil {
		return st.Err() // 400
	}

	return std.Err() // 400
}

// PrepareConfirmEmailRecord creates instance of Confirm struct for confirming email.
func (uc *UseCase) PrepareConfirmEmailRecord(email string) (confirm domain.Confirm, err error) {
	// prepare variables for email sending
	confirmPassword := uuid.New().String()
	confirm, err = domain.NewConfirm(
		domain.EmailConfirmKind,
		confirmPassword,
		uc.confirmLife,
		map[string]string{ /* vars */
			"email": email,
		},
	)

	return
}

// PrepareConfirmEmailEvent creates instance of Event for sending confirmation email.
func (uc *UseCase) PrepareConfirmEmailEvent(email, lang, code string, profile []byte) (confirmEmail event.Event, err error) {
	// prepare url
	url := fmt.Sprintf("%s/confirm/%s", uc.baseURL, code)

	vars := map[string]any{
		"url": url,
	}

	// fill vars with all profile variables
	gjson.ParseBytes(profile).ForEach(func(key, value gjson.Result) bool {
		vars[key.String()] = value.Value()
		return true
	})

	// prepare event email.SendTemplate
	toEmail := email
	confirmEmail, err = event.EmailSendTemplate(
		confirmEmailTemplateName,
		lang, // TODO: en language should be user's language not constant
		uc.fromEmail,
		uc.fromName,
		toEmail,
		vars,
	)
	if err != nil {
		err = errors.Wrap(err, "prepare event for sending email-pass confirmation")
		return
	}

	return confirmEmail, err
}
