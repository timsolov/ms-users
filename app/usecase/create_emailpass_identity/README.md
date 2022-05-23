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
    "template": "email-pass-confirm",
    "language": "en", // i18n 
    "vars": { // marshaled json object with variables such as subject, receiver, send, sender_name and others to use in template
      "sender": "test@example.org", // sender email
      "sender_name": "Sender Name", // sender name
      "receiver": "new_user@example.org", // email of the new identity
      "receiver_name": "New User", // full name of the new identity
      "url": "http://example.org/confirm/JAJSUDfy18" // confirmation url
    }
  }
  ```

## Configuration
- FROM_EMAIL - from this email will be sent confirmation email;
- FROM_NAME  - from this name will be sent confirmation email;
- CONFIRM_LIFE - duration for confirm life;
