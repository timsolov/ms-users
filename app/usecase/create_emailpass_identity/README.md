# Create Email-Password Identity

This usecase is useful when we've to create new email-password identity with user profile.
It's always creates both a identity email-password and user profile.
So if you need to add new identity to existing profile you've to use another usecase.

## Dependencies
- pgq queue `outbox` for sending events to message broker.
- email service for sending email confirmation.
  
  Email service should listen for subject `email.SendTemplate` with format:
  ```json
  {
    "template": "confirm-registration-email-pass",
    "language": "en", // i18n 
    "vars": { // marshaled json object with variables such as subject, receiver, send, sender_name and others to use in template
      "sender": "test@example.org", // sender email
      "sender_name": "Sender Name", // sender name
      "receiver": "new_user@example.org", // email of the new identity
      "url": "http://api.example.org/confirm/JAJSUDfy18" // confirmation url
      // and whole profile object also will be included here
    }
  }
  ```
  If you'd like to use first_name, last_name from profile in email, you've to use template
  placeholders from profile object (e.g. `{{.first_name}}` if user provided his `first_name`). Also you've to provoide for possibility of empty `first_name`, `last_name` when they're not reqired in jsonschema.

## Configuration
- FROM_EMAIL - from this email will be sent confirmation email;
- FROM_NAME  - from this name will be sent confirmation email;
- API_BASE_URL - base url to the service;
- CONFIRM_LIFE - duration for confirm life;
- PROFILE_JSONSCHEMA_PATH - path to file with jsonschema for validate provided profile object (example: settings/_profile.json);
