package templates

// GetResetPasswordErrorHTML generates an HTML page for invalid password reset links
func GetResetPasswordErrorHTML(forgotPasswordURL string) (string, error) {
	html := `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Password Reset Error</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            margin: 0;
            padding: 0;
            display: flex;
            justify-content: center;
            align-items: center;
            height: 100vh;
            background-color: #f5f5f5;
        }
        .container {
            text-align: center;
            padding: 2rem;
            background-color: white;
            border-radius: 8px;
            box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
            max-width: 500px;
        }
        h1 {
            color: #e74c3c;
        }
        p {
            margin: 1rem 0;
            line-height: 1.5;
            color: #333;
        }
        .btn {
            display: inline-block;
            background-color: #3498db;
            color: white;
            padding: 0.75rem 1.5rem;
            text-decoration: none;
            border-radius: 4px;
            margin-top: 1rem;
            transition: background-color 0.3s;
        }
        .btn:hover {
            background-color: #2980b9;
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>Password Reset Link Invalid</h1>
        <p>The password reset link you clicked is invalid or has expired.</p>
        <p>Please request a new password reset link.</p>
        <a href="` + forgotPasswordURL + `" class="btn">Go to Forgot Password</a>
    </div>
</body>
</html>
`
	return html, nil
}
