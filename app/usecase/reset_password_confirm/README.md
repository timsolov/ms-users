# Reset password confirm for email-pass identity

This usecase is useful for set new password for email-password identity.
The end-point waits for:
- confirm_id - confirm id from `i` query parameter from email link from reset_password_init usecase;
- verifycation - confirm password from `v` query parameter from email link from reset_password_init usecase;
- new_password - new password for identity.
