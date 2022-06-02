# Reset password process init for email-pass identity

This usecase is useful when we've to init password reset process for email-password identity.
For email-pass identity should be provided email and if identity with related
email exists confirmation code for recovery process will be sent to that email.
In that email will be stored link with comfirm_id (i) and verifycation code (v) inside
query parameters of the link. The link should leads the user to the web page where he will 
see input for a `new password`.

## Dependencies
- pgq queue `outbox` for sending events to message broker.
- email service for sending email confirmation.

  Email service should listen for subject `email.SendTemplate` with format:
  ```json
  {
    "template": "reset-email-password-init",
    "language": "en", // i18n 
    "vars": { // marshaled json object with variables such as subject, receiver, send, sender_name and others to use in template
      "sender": "test@example.org", // sender email
      "sender_name": "Sender Name", // sender name
      "receiver": "new_user@example.org", // email of the new identity
      "receiver_name": "New User", // full name of the new identity
      // url - is link to the Web page shows field "new password"
      // query i - should be confirm_id
      // query v - should be password for confirm record in db
      // above vars utilizes in next request to set new password for identity
      "url": "http://example.org/reset-password?i=Jjquws&v=kiquywha"
    }
  }
  ```

## Configuration
- FROM_EMAIL - from this email will be sent confirmation email;
- FROM_NAME  - from this name will be sent confirmation email;
- WEB_BASE_URL - base url of Front-End;
- CONFIRM_LIFE - duration for confirm life;
