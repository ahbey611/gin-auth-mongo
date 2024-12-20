package mail

// Verification code email template
var EmailRegisterLinkTemplate string = `<h1>Email Confirmation</h1>
<h2>Hello %s</h2>
<p>Thank you for registering. You can complete your registration by clicking the link below:</p>
<a href="%s">%s</a>
<p>This link will expire in <strong>%d minutes</strong>.</p>
<p>Expired time: %s</p>
<p>This email is auto generated, please do not reply to this email.</p>
<p>If you did not request this email, please ignore it.</p>
`

var PasswordResetLinkTemplate string = `<h1>Password Reset</h1>
<h2>Hello %s</h2>
<p>You can reset your password by clicking the link below:</p>
<a href="%s">%s</a>
<p>This link will expire in <strong>%d minutes</strong>.</p>
<p>Expired time: %s</p>
<p>This email is auto generated, please do not reply to this email.</p>
<p>If you did not request this email, please ignore it.</p>`

var EmailRegisterCodeTemplate string = `<h1>Email Confirmation</h1>
<h2>Hello %s</h2>
<p>Thank you for registering. Here is your verification code:</p>
<h2>%s</h2>
<p>This code will expire in <strong>%d minutes</strong>.</p>
<p>Expired time: %s</p>
<p>This email is auto generated, please do not reply to this email.</p>
<p>If you did not request this email, please ignore it.</p>
`

var PasswordResetCodeTemplate string = `<h1>Password Reset</h1>
<h2>Hello %s</h2>
<p>Here is your password reset verification code:</p>
<h2>%s</h2>
<p>This code will expire in <strong>%d minutes</strong>.</p>
<p>Expired time: %s</p>
<p>This email is auto generated, please do not reply to this email.</p>
<p>If you did not request this email, please ignore it.</p>`
