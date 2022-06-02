# Retry to send confirmation email for new user

This usecase is useful when user wants to resend confirmation email on registration process.

## TODO
- Rate limit for new requests

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
      "receiver_name": "New User", // full name of the new identity
      "url": "http://api.example.org/confirm/JAJSUDfy18" // confirmation url
    }
  }
  ```

## Configuration
- FROM_EMAIL - from this email will be sent confirmation email;
- FROM_NAME  - from this name will be sent confirmation email;
- API_BASE_URL - base url to the service;
- CONFIRM_LIFE - duration for confirm life;